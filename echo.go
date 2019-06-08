/*
Package echo implements high performance, minimalist Go web framework.

Example:

  package main

  import (
    "net/http"

    "github.com/adverax/echo"
    "github.com/adverax/echo/middleware"
  )

  // Handler
  func hello(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
  }

  func main() {
    // Echo instance
    e := echo.New()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Routes
    router := e.Router()
    e.GET(router, "/", hello)

    // Start server
    e.Logger.Error(e.Start(":1323"))
  }

*/
package echo

import (
	stdContext "context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/adverax/echo/cacher"
	"github.com/adverax/echo/sync/arbiter"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/acme/autocert"

	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/adverax/echo/cache"
	"github.com/adverax/echo/data"
	"github.com/adverax/echo/generic"
	"github.com/adverax/echo/log"
)

// Echo is the top-level framework instance.
type Echo struct {
	maxParam         *int
	router           Router
	notFoundHandler  HandlerFunc
	roots            sync.Pool
	Server           *http.Server
	TLSServer        *http.Server
	Listener         net.Listener
	TLSListener      net.Listener
	AutoTLSManager   autocert.Manager
	DisableHTTP2     bool
	Debug            bool
	HideBanner       bool
	HidePort         bool
	Complex          bool
	HTTPErrorHandler HTTPErrorHandler
	Logger           log.Logger
	Locale           Locale // Prototype
	UrlLinker        UrlLinker
	Cache            cache.Cache
	Cacher           cacher.Cacher
	Messages         MessageManager
	Resources        ResourceManager
	DataSets         DataSetManager
	Arbiter          arbiter.Arbiter
}

// NewContext returns a Context instance.
func (e *Echo) NewContext(r *http.Request, w http.ResponseWriter) Context {
	locale := generic.MakePointerTo(generic.CloneValue(e.Locale))
	ctx, r := makeContext(r)

	return &context{
		Context:  ctx,
		request:  r,
		response: NewResponse(w, e),
		store:    make(map[interface{}]interface{}),
		echo:     e,
		handler:  NotFoundHandler,
		locale:   locale.(Locale),
	}
}

// NewRouter return a new Router instance.
func (e *Echo) NewRouter() Router {
	return &router{
		echo:   e,
		Router: chi.NewRouter(),
	}
}

// Router returns router.
func (e *Echo) Router() Router {
	return e.router
}

// DefaultHTTPErrorHandler is the default HTTP error handler. It sends a JSON response
// with status code.
func (e *Echo) DefaultHTTPErrorHandler(c Context, err error) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)

	if err != ErrNotFound {
		e.Logger.Error(err)
	}

	if err == data.ErrNoMatch {
		err = ErrNotFound
	}

	if he, ok := err.(*HTTPError); ok {
		code = he.Code
		msg = he.Message
		if he.Internal != nil {
			err = fmt.Errorf("%v, %v", err, he.Internal)
		}
	} else if e.Debug {
		msg = err.Error()
	} else {
		msg = http.StatusText(code)
	}
	if _, ok := msg.(string); ok {
		msg = Map{"message": msg}
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, msg)
		}
		if err != nil {
			e.Logger.Error(err)
		}
	}
}

// AcquireContext returns an empty `Context` instance from the pool.
// You must return the context by calling `ReleaseContext()`.
func (e *Echo) AcquireContext() Context {
	return e.roots.Get().(Context)
}

// ReleaseContext returns the `Context` instance back to the pool.
// You must call it after `AcquireContext()`.
func (e *Echo) ReleaseContext(c Context) {
	e.roots.Put(c)
}

// ServeHTTP implements `http.Handler` interface, which serves HTTP requests.
func (e *Echo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e.Complex {
		e.router.ServeHTTP(w, r)
		return
	}

	// Acquire context
	c := e.roots.Get().(*context)
	c.Reset(r, w)

	// Attach context to the response
	r = r.WithContext(stdContext.WithValue(r.Context(), ContextKey, c))
	c.request = r

	// Execute chain
	e.router.ServeHTTP(w, r)

	// Release context
	e.roots.Put(c)
}

// Dynamic is internal method, that used for init context for handling dynamic HTTP requests in the COMPLEX mode.
func (e *Echo) Dynamic(w http.ResponseWriter, r *http.Request, next http.Handler) {
	// Acquire context
	c := e.roots.Get().(*context)
	defer e.roots.Put(c)

	c.Reset(r, w)

	// Attach context to the response
	r = r.WithContext(stdContext.WithValue(r.Context(), ContextKey, c))
	c.request = r

	// Execute chain
	next.ServeHTTP(w, r)
}

// Invoke http handler
func (e *Echo) dispatch(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := r.Context().Value(ContextKey).(*context)
		c.Context = r.Context()
		c.request = r.WithContext(c)

		err := handler(c)
		if err != nil {
			c.Error(err)
		}
	}
}

// Start starts an HTTP server.
func (e *Echo) Start(address string) error {
	e.Server.Addr = address
	return e.StartServer(e.Server)
}

// StartTLS starts an HTTPS server.
// If `certFile` or `keyFile` is `string` the values are treated as file paths.
// If `certFile` or `keyFile` is `[]byte` the values are treated as the certificate or key as-is.
func (e *Echo) StartTLS(address string, certFile, keyFile interface{}) (err error) {
	var cert []byte
	if cert, err = filepathOrContent(certFile); err != nil {
		return
	}

	var key []byte
	if key, err = filepathOrContent(keyFile); err != nil {
		return
	}

	s := e.TLSServer
	s.TLSConfig = new(tls.Config)
	s.TLSConfig.Certificates = make([]tls.Certificate, 1)
	if s.TLSConfig.Certificates[0], err = tls.X509KeyPair(cert, key); err != nil {
		return
	}

	return e.startTLS(address)
}

// StartAutoTLS starts an HTTPS server using certificates automatically installed from https://letsencrypt.org.
func (e *Echo) StartAutoTLS(address string) error {
	s := e.TLSServer
	s.TLSConfig = new(tls.Config)
	s.TLSConfig.GetCertificate = e.AutoTLSManager.GetCertificate
	return e.startTLS(address)
}

func (e *Echo) startTLS(address string) error {
	s := e.TLSServer
	s.Addr = address
	if !e.DisableHTTP2 {
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "h2")
	}
	return e.StartServer(e.TLSServer)
}

// StartServer starts a custom http server.
func (e *Echo) StartServer(s *http.Server) (err error) {
	// Setup
	s.Handler = e

	if !e.HideBanner {
		e.Logger.Info(fmt.Sprintf(banner, "v"+Version, website))
	}

	if s.TLSConfig == nil {
		if e.Listener == nil {
			e.Listener, err = newListener(s.Addr)
			if err != nil {
				return err
			}
		}
		if !e.HidePort {
			e.Logger.Info(fmt.Sprintf("⇨ http server started on %s\n", e.Listener.Addr()))
		}
		return s.Serve(e.Listener)
	}
	if e.TLSListener == nil {
		l, err := newListener(s.Addr)
		if err != nil {
			return err
		}
		e.TLSListener = tls.NewListener(l, s.TLSConfig)
	}
	if !e.HidePort {
		e.Logger.Info(fmt.Sprintf("⇨ https server started on %s\n", e.TLSListener.Addr()))
	}
	return s.Serve(e.TLSListener)
}

// Close immediately stops the server.
// It internally calls `http.Server#Close()`.
func (e *Echo) Close() error {
	if err := e.TLSServer.Close(); err != nil {
		return err
	}
	return e.Server.Close()
}

// Shutdown stops the server gracefully.
// It internally calls `http.Server#Shutdown()`.
func (e *Echo) Shutdown(ctx stdContext.Context) error {
	if err := e.TLSServer.Shutdown(ctx); err != nil {
		return err
	}
	return e.Server.Shutdown(ctx)
}

// HTTPError represents an error that occurred while handling a request.
type HTTPError struct {
	Code     int
	Message  interface{}
	Internal error // Stores the error returned by an external dependency
}

// Error makes it compatible with `error` interface.
func (he *HTTPError) Error() string {
	return fmt.Sprintf("code=%d, message=%v", he.Code, he.Message)
}

// SetInternal sets error to HTTPError.Internal
func (he *HTTPError) SetInternal(err error) *HTTPError {
	he.Internal = err
	return he
}

type Handler interface {
	ServeHTTP(ctx Context) error
}

// HandlerFunc defines a function to serve HTTP requests.
type HandlerFunc func(ctx Context) error

// HTTPErrorHandler is a centralized HTTP error handler.
type HTTPErrorHandler func(Context, error)

// Map defines a generic map of type `map[string]interface{}`.
type Map map[string]interface{}

func (fn HandlerFunc) ServeHTTP(ctx Context) error {
	return fn(ctx)
}

// HTTP methods
// NOTE: Deprecated, please use the stdlib constants directly instead.
const (
	CONNECT = http.MethodConnect
	DELETE  = http.MethodDelete
	GET     = http.MethodGet
	HEAD    = http.MethodHead
	OPTIONS = http.MethodOptions
	PATCH   = http.MethodPatch
	POST    = http.MethodPost
	// PROPFIND = "PROPFIND"
	PUT   = http.MethodPut
	TRACE = http.MethodTrace
)

// MIME types
const (
	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + charsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + charsetUTF8
	MIMETextXML                          = "text/xml"
	MIMETextXMLCharsetUTF8               = MIMETextXML + "; " + charsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + charsetUTF8
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
)

const (
	charsetUTF8 = "charset=UTF-8"
	// PROPFIND Method can be used on collection and property resources.
	PROPFIND = "PROPFIND"
)

// Headers
const (
	HeaderAccept              = "Accept"
	HeaderAcceptEncoding      = "Accept-Encoding"
	HeaderAllow               = "Allow"
	HeaderAuthorization       = "Authorization"
	HeaderContentDisposition  = "Content-Disposition"
	HeaderContentEncoding     = "Content-Encoding"
	HeaderContentLength       = "Content-Length"
	HeaderContentType         = "Content-Type"
	HeaderCookie              = "Cookie"
	HeaderSetCookie           = "Set-Cookie"
	HeaderIfModifiedSince     = "If-Modified-Since"
	HeaderLastModified        = "Last-Modified"
	HeaderLocation            = "Location"
	HeaderUpgrade             = "Upgrade"
	HeaderVary                = "Vary"
	HeaderWWWAuthenticate     = "WWW-Authenticate"
	HeaderXForwardedFor       = "X-Forwarded-For"
	HeaderXForwardedProto     = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       = "X-Forwarded-Ssl"
	HeaderXUrlScheme          = "X-Url-Scheme"
	HeaderXHTTPMethodOverride = "X-HTTP-Method-Override"
	HeaderXRealIP             = "X-Real-IP"
	HeaderXRequestID          = "X-Request-ID"
	HeaderXRequestedWith      = "X-Requested-With"
	HeaderServer              = "Server"
	HeaderOrigin              = "Origin"

	// Access control
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity         = "Strict-Transport-Security"
	HeaderXContentTypeOptions             = "X-Content-Type-Options"
	HeaderXXSSProtection                  = "X-XSS-Protection"
	HeaderXFrameOptions                   = "X-Frame-Options"
	HeaderContentSecurityPolicy           = "Content-Security-Policy"
	HeaderContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	HeaderXCSRFToken                      = "X-CSRF-Token"
)

const (
	// Version of Echo
	Version = "5.0.0"
	website = "https://github.com/adverax/echo"
	// http://patorjk.com/software/taag/#p=display&f=Small%20Slant&t=Echo
	banner = `
   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ %s
High performance, minimalist Go web framework
%s
____________________________________O/_______
                                    O\
`
)

// Errors
var (
	ErrUnsupportedMediaType        = NewHTTPError(http.StatusUnsupportedMediaType)
	ErrNotFound                    = NewHTTPError(http.StatusNotFound)
	ErrUnauthorized                = NewHTTPError(http.StatusUnauthorized)
	ErrForbidden                   = NewHTTPError(http.StatusForbidden)
	ErrMethodNotAllowed            = NewHTTPError(http.StatusMethodNotAllowed)
	ErrStatusRequestEntityTooLarge = NewHTTPError(http.StatusRequestEntityTooLarge)
	ErrTooManyRequests             = NewHTTPError(http.StatusTooManyRequests)
	ErrBadRequest                  = NewHTTPError(http.StatusBadRequest)
	ErrBadGateway                  = NewHTTPError(http.StatusBadGateway)
	ErrInternalServerError         = NewHTTPError(http.StatusInternalServerError)
	ErrRequestTimeout              = NewHTTPError(http.StatusRequestTimeout)
	ErrServiceUnavailable          = NewHTTPError(http.StatusServiceUnavailable)
	ErrValidatorNotRegistered      = errors.New("validator not registered")
	ErrRendererNotRegistered       = errors.New("renderer not registered")
	ErrInvalidRedirectCode         = errors.New("invalid redirect status code")
	ErrCookieNotFound              = errors.New("cookie not found")
	ErrInvalidCertOrKeyType        = errors.New("invalid cert or key type, must be string or []byte")
	ErrAbort                       = errors.New("abort")
	ErrModelSealed                 = errors.New("model is accepted")
)

// Error handlers
var (
	NotFoundHandler = func(c Context) error {
		return ErrNotFound
	}

	MethodNotAllowedHandler = func(c Context) error {
		return ErrMethodNotAllowed
	}
)

// New creates an instance of Echo.
func New() (e *Echo) {
	e = &Echo{
		Server:    new(http.Server),
		TLSServer: new(http.Server),
		Locale:    Defaults.Locale,
		UrlLinker: Defaults.UrlLinker,
		Cache:     Defaults.Cache,
		Cacher:    Defaults.Cacher,
		Arbiter:   Defaults.Arbiter,
		Logger:    log.NewDebug("\n"),
		DataSets:  Defaults.DataSets,
		AutoTLSManager: autocert.Manager{
			Prompt: autocert.AcceptTOS,
		},
		maxParam: new(int),
	}
	e.Server.Handler = e
	e.TLSServer.Handler = e
	e.HTTPErrorHandler = e.DefaultHTTPErrorHandler
	e.roots.New = func() interface{} {
		return e.NewContext(nil, nil)
	}
	e.router = e.NewRouter()
	return
}

// NewHTTPError creates a new HTTPError instance.
func NewHTTPError(code int, message ...interface{}) *HTTPError {
	he := &HTTPError{Code: code, Message: http.StatusText(code)}
	if len(message) > 0 {
		he.Message = message[0]
	}
	return he
}

// WrapHandler wraps `http.Handler` into `echo.HandlerFunc`.
func WrapHandler(h http.Handler) HandlerFunc {
	return func(c Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	if c, err = ln.AcceptTCP(); err != nil {
		return
	} else if err = c.(*net.TCPConn).SetKeepAlive(true); err != nil {
		return
	} else if err = c.(*net.TCPConn).SetKeepAlivePeriod(3 * time.Minute); err != nil {
		return
	}
	return
}

func filepathOrContent(fileOrContent interface{}) (content []byte, err error) {
	switch v := fileOrContent.(type) {
	case string:
		return ioutil.ReadFile(v)
	case []byte:
		return v, nil
	default:
		return nil, ErrInvalidCertOrKeyType
	}
}

func newListener(address string) (*tcpKeepAliveListener, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return &tcpKeepAliveListener{l.(*net.TCPListener)}, nil
}

func makeContext(r *http.Request) (ctx stdContext.Context, req *http.Request) {
	if r == nil {
		return stdContext.WithValue(
			stdContext.Background(),
			chi.RouteCtxKey,
			chi.NewRouteContext(),
		), nil
	}

	ctx = r.Context()
	rc := ctx.Value(chi.RouteCtxKey)
	if _, ok := rc.(*chi.Context); !ok {
		ctx = stdContext.WithValue(
			ctx,
			chi.RouteCtxKey,
			chi.NewRouteContext(),
		)
		return ctx, r.WithContext(ctx)
	}

	return ctx, r
}

type Router interface {
	http.Handler
	chi.Routes

	// Use appends one of more middlewares onto the Router stack.
	Use(middlewares ...func(http.Handler) http.Handler)

	// With adds inline middlewares for an endpoint handler.
	With(middlewares ...func(http.Handler) http.Handler) Router

	// Group adds a new inline-Router along the current routing
	// path, with a fresh middleware stack for the inline-Router.
	Group(fn func(r Router)) Router

	// Route mounts a sub-Router along a `pattern`` string.
	Route(pattern string, fn func(r Router)) Router

	// Mount attaches another http.Handler along ./pattern/*
	Mount(pattern string, h HandlerFunc)

	// Handle and HandleFunc adds routes for `pattern` that matches
	// all HTTP methods.
	Handle(pattern string, h HandlerFunc)

	// Method and MethodFunc adds routes for `pattern` that matches
	// the `method` HTTP method.
	Method(method, pattern string, h HandlerFunc)

	// HTTP-method routing along `pattern`
	Connect(pattern string, h HandlerFunc)
	Delete(pattern string, h HandlerFunc)
	Get(pattern string, h HandlerFunc)
	Head(pattern string, h HandlerFunc)
	Options(pattern string, h HandlerFunc)
	Patch(pattern string, h HandlerFunc)
	Post(pattern string, h HandlerFunc)
	Put(pattern string, h HandlerFunc)
	Trace(pattern string, h HandlerFunc)
	Form(pattern string, h HandlerFunc)

	// NotFound defines a handler to respond whenever a route could
	// not be found.
	NotFound(h HandlerFunc)

	// MethodNotAllowed defines a handler to respond whenever a method is
	// not allowed.
	MethodNotAllowed(h HandlerFunc)
}

type router struct {
	chi.Router
	echo *Echo
}

func (r *router) With(middlewares ...func(http.Handler) http.Handler) Router {
	rr := r.Router.With(middlewares...)
	return &router{echo: r.echo, Router: rr}
}

func (r *router) Group(fn func(Router)) (res Router) {
	r.Router.Group(func(rr chi.Router) {
		res := &router{Router: rr, echo: r.echo}
		fn(res)
	})
	return
}

func (r *router) Route(pattern string, fn func(r Router)) (res Router) {
	r.Router.Route(pattern, func(rr chi.Router) {
		res := &router{Router: rr, echo: r.echo}
		fn(res)
	})
	return
}

func (r *router) Method(method, pattern string, handler HandlerFunc) {
	r.Router.MethodFunc(method, pattern, r.echo.dispatch(handler))
}

func (r *router) Mount(pattern string, handler HandlerFunc) {
	r.Router.Mount(pattern, r.echo.dispatch(handler))
}

func (r *router) Handle(pattern string, handler HandlerFunc) {
	r.Router.Handle(pattern, r.echo.dispatch(handler))
}

func (r *router) Connect(pattern string, handler HandlerFunc) {
	r.Router.Connect(pattern, r.echo.dispatch(handler))
}

func (r *router) Delete(pattern string, handler HandlerFunc) {
	r.Router.Delete(pattern, r.echo.dispatch(handler))
}

func (r *router) Get(pattern string, handler HandlerFunc) {
	r.Router.Get(pattern, r.echo.dispatch(handler))
}

func (r *router) Head(pattern string, handler HandlerFunc) {
	r.Router.Head(pattern, r.echo.dispatch(handler))
}

func (r *router) Options(pattern string, handler HandlerFunc) {
	r.Router.Options(pattern, r.echo.dispatch(handler))
}

func (r *router) Patch(pattern string, handler HandlerFunc) {
	r.Router.Patch(pattern, r.echo.dispatch(handler))
}

func (r *router) Post(pattern string, handler HandlerFunc) {
	r.Router.Post(pattern, r.echo.dispatch(handler))
}

func (r *router) Put(pattern string, handler HandlerFunc) {
	r.Router.Put(pattern, r.echo.dispatch(handler))
}

func (r *router) Trace(pattern string, handler HandlerFunc) {
	r.Router.Trace(pattern, r.echo.dispatch(handler))
}

func (r *router) Form(pattern string, handler HandlerFunc) {
	r.Get(pattern, handler)
	r.Post(pattern, handler)
}

func (r *router) NotFound(handler HandlerFunc) {
	r.Router.NotFound(r.echo.dispatch(handler))
}

func (r *router) MethodNotAllowed(handler HandlerFunc) {
	r.Router.MethodNotAllowed(r.echo.dispatch(handler))
}
