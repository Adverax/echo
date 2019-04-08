package echo

import (
	stdContext "context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/adverax/echo/log"
	"github.com/go-chi/chi"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const (
	defaultMemory = 32 << 20 // 32 MB
	indexPage     = "index.html"
	defaultIndent = "  "
	ContextKey    = contextType(1)
)

type contextType int

// Context represents the context of the current HTTP request. It holds request and
// response objects, path, path parameters, data and registered handler.
type Context interface {
	stdContext.Context

	// Request returns `*http.Request`.
	Request() *http.Request

	// SetRequest sets `*http.Request`.
	SetRequest(r *http.Request)

	// Response returns `*Response`.
	Response() *Response

	// IsTLS returns true if HTTP connection is TLS otherwise false.
	IsTLS() bool

	// IsWebSocket returns true if HTTP connection is WebSocket otherwise false.
	IsWebSocket() bool

	// Scheme returns the HTTP protocol scheme, `http` or `https`.
	Scheme() string

	// RealIP returns the client's network address based on `X-Forwarded-For`
	// or `X-Real-IP` request header.
	RealIP() string

	// Param returns path parameter by name.
	Param(name string) string
	ParamString(name string, defaults string) string
	ParamInt(name string, defaults int) int
	ParamInt8(name string, defaults int8) int8
	ParamInt16(name string, defaults int16) int16
	ParamInt32(name string, defaults int32) int32
	ParamInt64(name string, defaults int64) int64
	ParamUint(name string, defaults uint) uint
	ParamUint8(name string, defaults uint8) uint8
	ParamUint16(name string, defaults uint16) uint16
	ParamUint32(name string, defaults uint32) uint32
	ParamUint64(name string, defaults uint64) uint64
	ParamFloat32(name string, defaults float32) float32
	ParamFloat64(name string, defaults float64) float64
	ParamBoolean(name string, defaults bool) bool

	// ParamNames returns path parameter names.
	ParamNames() []string

	// SetParamNames sets path parameter names.
	SetParamNames(names ...string)

	// ParamValues returns path parameter values.
	ParamValues() []string

	// SetParamValues sets path parameter values.
	SetParamValues(values ...string)

	// QueryParam returns the query param for the provided name.
	QueryParam(name string) string

	// QueryParams returns the query parameters as `url.Values`.
	QueryParams() url.Values

	// QueryString returns the URL query string.
	QueryString() string

	// FormValue returns the form field value for the provided name.
	FormValue(name string) string

	// FormParams returns the form parameters as `url.Values`.
	FormParams() (url.Values, error)

	// FormFile returns the multipart form file for the provided name.
	FormFile(name string) (*multipart.FileHeader, error)

	// MultipartForm returns the multipart form.
	MultipartForm() (*multipart.Form, error)

	// Cookie returns the named cookie provided in the request.
	Cookie(name string) (*http.Cookie, error)

	// SetCookie adds a `Set-Cookie` header in HTTP response.
	SetCookie(cookie *http.Cookie)

	// Cookies returns the HTTP cookies sent with the request.
	Cookies() []*http.Cookie

	// Get retrieves data from the context.
	Get(key interface{}) interface{}

	// Set saves data in the context.
	Set(key interface{}, val interface{})

	// Create new context with new value.
	WithValue(key interface{}, val interface{}) Context

	// HTML sends an HTTP response with status code.
	HTML(code int, html string) error

	// HTMLBlob sends an HTTP blob response with status code.
	HTMLBlob(code int, b []byte) error

	// String sends a string response with status code.
	String(code int, s string) error

	// JSON sends a JSON response with status code.
	JSON(code int, i interface{}) error

	// JSONPretty sends a pretty-print JSON with status code.
	JSONPretty(code int, i interface{}, indent string) error

	// JSONBlob sends a JSON blob response with status code.
	JSONBlob(code int, b []byte) error

	// JSONP sends a JSONP response with status code. It uses `callback` to construct
	// the JSONP payload.
	JSONP(code int, callback string, i interface{}) error

	// JSONPBlob sends a JSONP blob response with status code. It uses `callback`
	// to construct the JSONP payload.
	JSONPBlob(code int, callback string, b []byte) error

	// XML sends an XML response with status code.
	XML(code int, i interface{}) error

	// XMLPretty sends a pretty-print XML with status code.
	XMLPretty(code int, i interface{}, indent string) error

	// XMLBlob sends an XML blob response with status code.
	XMLBlob(code int, b []byte) error

	// Blob sends a blob response with status code and content type.
	Blob(code int, contentType string, b []byte) error

	// Stream sends a streaming response with status code and content type.
	Stream(code int, contentType string, r io.Reader) error

	// Template sends a HTML response with status code,
	Template(code int, t Template, data interface{}) (err error)

	// File sends a response with the content of the file.
	File(file string) error

	// Attachment sends a response as attachment, prompting client to save the
	// file.
	Attachment(file string, name string) error

	// Inline sends a response as inline, opening the file in the browser.
	Inline(file string, name string) error

	// NoContent sends a response with no body and a status code.
	NoContent(code int) error

	// Redirect redirects the request to a provided URL with status code.
	Redirect(code int, url string) error

	// Refresh redirects the request to the corrent URL with status code.
	Refresh(code int) error

	// Revert redirects the request to a prev (referrer URL) address with status code.
	Revert(code int) error

	// HeaderNoCache writes the header for disable any caching
	HeaderNoCache()

	// Error invokes the registered HTTP error handler. Generally used by middleware.
	Error(err error)

	// Handler returns the matched handler by router.
	Handler() HandlerFunc

	// SetHandler sets the matched handler by router.
	SetHandler(h HandlerFunc)

	// Logger returns the `Logger` instance.
	Logger() log.Logger

	// Session returns the `Session` instance.
	Session() Session

	// SetSession sets the `Session` instance.
	SetSession(session Session)

	// Echo returns the `Echo` instance.
	Echo() *Echo

	// Reset resets the context after request completes. It must be called along
	// with `Echo#AcquireContext()` and `Echo#ReleaseContext()`.
	// See `Echo#ServeHTTP()`
	Reset(r *http.Request, w http.ResponseWriter)

	// Get active locale
	Locale() Locale

	// Set active locale
	SetLocale(locale Locale)

	// Add flash message
	AddFlash(class FlashClass, message interface{}) error
}

type context struct {
	stdContext.Context
	request  *http.Request
	response *Response
	path     string
	query    url.Values
	handler  HandlerFunc
	store    map[interface{}]interface{}
	echo     *Echo
	lock     sync.RWMutex
	locale   Locale
	session  Session
}

func (c *context) writeContentType(value string) {
	header := c.Response().Header()
	if header.Get(HeaderContentType) == "" {
		header.Set(HeaderContentType, value)
	}
}

func (c *context) Request() *http.Request {
	return c.request
}

func (c *context) SetRequest(r *http.Request) {
	c.request = r
}

func (c *context) Response() *Response {
	return c.response
}

func (c *context) IsTLS() bool {
	return c.request.TLS != nil
}

func (c *context) IsWebSocket() bool {
	upgrade := c.request.Header.Get(HeaderUpgrade)
	return strings.ToLower(upgrade) == "websocket"
}

func (c *context) Scheme() string {
	// Can't use `r.Request.URL.Scheme`
	// See: https://groups.google.com/forum/#!topic/golang-nuts/pMUkBlQBDF0
	if c.IsTLS() {
		return "https"
	}
	if scheme := c.request.Header.Get(HeaderXForwardedProto); scheme != "" {
		return scheme
	}
	if scheme := c.request.Header.Get(HeaderXForwardedProtocol); scheme != "" {
		return scheme
	}
	if ssl := c.request.Header.Get(HeaderXForwardedSsl); ssl == "on" {
		return "https"
	}
	if scheme := c.request.Header.Get(HeaderXUrlScheme); scheme != "" {
		return scheme
	}
	return "http"
}

func (c *context) RealIP() string {
	if ip := c.request.Header.Get(HeaderXForwardedFor); ip != "" {
		return strings.Split(ip, ", ")[0]
	}
	if ip := c.request.Header.Get(HeaderXRealIP); ip != "" {
		return ip
	}
	ra, _, _ := net.SplitHostPort(c.request.RemoteAddr)
	return ra
}

func (c *context) Param(name string) string {
	return chi.URLParamFromCtx(c, name)
}

func (c *context) ParamString(name string, defaults string) string {
	v := c.Param(name)
	if v == "" {
		return defaults
	}
	return v
}

func (c *context) ParamInt(name string, defaults int) int {
	s := c.Param(name)
	if s == "" {
		return defaults
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaults
	}
	return int(v)
}

func (c *context) ParamInt8(name string, defaults int8) int8 {
	s := c.Param(name)
	if s == "" {
		return defaults
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaults
	}
	return int8(v)
}

func (c *context) ParamInt16(name string, defaults int16) int16 {
	s := c.Param(name)
	if s == "" {
		return defaults
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaults
	}
	return int16(v)
}

func (c *context) ParamInt32(name string, defaults int32) int32 {
	s := c.Param(name)
	if s == "" {
		return defaults
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaults
	}
	return int32(v)
}

func (c *context) ParamInt64(name string, defaults int64) int64 {
	s := c.Param(name)
	if s == "" {
		return defaults
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaults
	}
	return int64(v)
}

func (c *context) ParamUint(name string, defaults uint) uint {
	s := c.Param(name)
	if s == "" {
		return defaults
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return defaults
	}
	return uint(v)
}

func (c *context) ParamUint8(name string, defaults uint8) uint8 {
	s := c.Param(name)
	if s == "" {
		return defaults
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return defaults
	}
	return uint8(v)
}

func (c *context) ParamUint16(name string, defaults uint16) uint16 {
	s := c.Param(name)
	if s == "" {
		return defaults
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return defaults
	}
	return uint16(v)
}

func (c *context) ParamUint32(name string, defaults uint32) uint32 {
	s := c.Param(name)
	if s == "" {
		return defaults
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return defaults
	}
	return uint32(v)
}

func (c *context) ParamUint64(name string, defaults uint64) uint64 {
	s := c.Param(name)
	if s == "" {
		return defaults
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return defaults
	}
	return uint64(v)
}

func (c *context) ParamFloat32(name string, defaults float32) float32 {
	s := c.Param(name)
	if s == "" {
		return defaults
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaults
	}
	return float32(v)
}

func (c *context) ParamFloat64(name string, defaults float64) float64 {
	s := c.Param(name)
	if s == "" {
		return defaults
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaults
	}
	return v
}

func (c *context) ParamBoolean(name string, defaults bool) bool {
	s := strings.ToLower(c.Param(name))
	if s == "" {
		return defaults
	}
	if v, has := boolMap[s]; has {
		return v
	}
	return defaults
}

func (c *context) ParamNames() []string {
	cc := chi.RouteContext(c)
	return cc.URLParams.Keys
}

func (c *context) SetParamNames(names ...string) {
	cc := chi.RouteContext(c)
	cc.URLParams.Keys = names
}

func (c *context) ParamValues() []string {
	cc := chi.RouteContext(c)
	return cc.URLParams.Values[:len(cc.URLParams.Keys)]
}

func (c *context) SetParamValues(values ...string) {
	cc := chi.RouteContext(c)
	cc.URLParams.Values = values
}

func (c *context) QueryParam(name string) string {
	if c.query == nil {
		c.query = c.request.URL.Query()
	}
	return c.query.Get(name)
}

func (c *context) QueryParams() url.Values {
	if c.query == nil {
		c.query = c.request.URL.Query()
	}
	return c.query
}

func (c *context) QueryString() string {
	return c.request.URL.RawQuery
}

func (c *context) FormValue(name string) string {
	return c.request.FormValue(name)
}

func (c *context) FormParams() (url.Values, error) {
	if strings.HasPrefix(c.request.Header.Get(HeaderContentType), MIMEMultipartForm) {
		if err := c.request.ParseMultipartForm(defaultMemory); err != nil {
			return nil, err
		}
	} else {
		if err := c.request.ParseForm(); err != nil {
			return nil, err
		}
	}
	return c.request.Form, nil
}

func (c *context) FormFile(name string) (*multipart.FileHeader, error) {
	_, fh, err := c.request.FormFile(name)
	return fh, err
}

func (c *context) MultipartForm() (*multipart.Form, error) {
	err := c.request.ParseMultipartForm(defaultMemory)
	return c.request.MultipartForm, err
}

func (c *context) Cookie(name string) (*http.Cookie, error) {
	return c.request.Cookie(name)
}

func (c *context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Response(), cookie)
}

func (c *context) Cookies() []*http.Cookie {
	return c.request.Cookies()
}

func (c *context) WithValue(key interface{}, val interface{}) Context {
	return &valueCtx{
		Context: c,
		key:     key,
		val:     val,
	}
}

func (c *context) Get(key interface{}) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.store[key]
}

func (c *context) Set(key interface{}, val interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.store == nil {
		c.store = make(map[interface{}]interface{})
	}
	c.store[key] = val
}

func (c *context) HTML(code int, html string) (err error) {
	return c.HTMLBlob(code, []byte(html))
}

func (c *context) HTMLBlob(code int, b []byte) (err error) {
	return c.Blob(code, MIMETextHTMLCharsetUTF8, b)
}

func (c *context) String(code int, s string) (err error) {
	return c.Blob(code, MIMETextPlainCharsetUTF8, []byte(s))
}

func (c *context) jsonPBlob(code int, callback string, i interface{}) (err error) {
	enc := json.NewEncoder(c.response)
	_, pretty := c.QueryParams()["pretty"]
	if c.echo.Debug || pretty {
		enc.SetIndent("", "  ")
	}
	c.writeContentType(MIMEApplicationJavaScriptCharsetUTF8)
	c.response.WriteHeader(code)
	if _, err = c.response.Write([]byte(callback + "(")); err != nil {
		return
	}
	if err = enc.Encode(i); err != nil {
		return
	}
	if _, err = c.response.Write([]byte(");")); err != nil {
		return
	}
	return
}

func (c *context) json(code int, i interface{}, indent string) error {
	enc := json.NewEncoder(c.response)
	if indent != "" {
		enc.SetIndent("", indent)
	}
	c.writeContentType(MIMEApplicationJSONCharsetUTF8)
	c.response.WriteHeader(code)
	return enc.Encode(i)
}

func (c *context) JSON(code int, i interface{}) (err error) {
	indent := ""
	if _, pretty := c.QueryParams()["pretty"]; c.echo.Debug || pretty {
		indent = defaultIndent
	}
	return c.json(code, i, indent)
}

func (c *context) JSONPretty(code int, i interface{}, indent string) (err error) {
	return c.json(code, i, indent)
}

func (c *context) JSONBlob(code int, b []byte) (err error) {
	return c.Blob(code, MIMEApplicationJSONCharsetUTF8, b)
}

func (c *context) JSONP(code int, callback string, i interface{}) (err error) {
	return c.jsonPBlob(code, callback, i)
}

func (c *context) JSONPBlob(code int, callback string, b []byte) (err error) {
	c.writeContentType(MIMEApplicationJavaScriptCharsetUTF8)
	c.response.WriteHeader(code)
	if _, err = c.response.Write([]byte(callback + "(")); err != nil {
		return
	}
	if _, err = c.response.Write(b); err != nil {
		return
	}
	_, err = c.response.Write([]byte(");"))
	return
}

func (c *context) xml(code int, i interface{}, indent string) (err error) {
	c.writeContentType(MIMEApplicationXMLCharsetUTF8)
	c.response.WriteHeader(code)
	enc := xml.NewEncoder(c.response)
	if indent != "" {
		enc.Indent("", indent)
	}
	if _, err = c.response.Write([]byte(xml.Header)); err != nil {
		return
	}
	return enc.Encode(i)
}

func (c *context) XML(code int, i interface{}) (err error) {
	indent := ""
	if _, pretty := c.QueryParams()["pretty"]; c.echo.Debug || pretty {
		indent = defaultIndent
	}
	return c.xml(code, i, indent)
}

func (c *context) XMLPretty(code int, i interface{}, indent string) (err error) {
	return c.xml(code, i, indent)
}

func (c *context) XMLBlob(code int, b []byte) (err error) {
	c.writeContentType(MIMEApplicationXMLCharsetUTF8)
	c.response.WriteHeader(code)
	if _, err = c.response.Write([]byte(xml.Header)); err != nil {
		return
	}
	_, err = c.response.Write(b)
	return
}

func (c *context) Blob(code int, contentType string, b []byte) (err error) {
	c.writeContentType(contentType)
	c.response.WriteHeader(code)
	_, err = c.response.Write(b)
	return
}

func (c *context) Stream(code int, contentType string, r io.Reader) (err error) {
	c.writeContentType(contentType)
	c.response.WriteHeader(code)
	_, err = io.Copy(c.response, r)
	return
}

func (c *context) Template(code int, t Template, data interface{}) (err error) {
	c.writeContentType(MIMETextHTMLCharsetUTF8)
	c.response.WriteHeader(code)
	return t.Execute(c.response.Writer, data)
}

func (c *context) File(file string) (err error) {
	f, err := os.Open(file)
	if err != nil {
		return NotFoundHandler(c)
	}
	defer f.Close()

	fi, _ := f.Stat()
	if fi.IsDir() {
		file = filepath.Join(file, indexPage)
		f, err = os.Open(file)
		if err != nil {
			return NotFoundHandler(c)
		}
		defer f.Close()
		if fi, err = f.Stat(); err != nil {
			return
		}
	}
	http.ServeContent(c.Response(), c.Request(), fi.Name(), fi.ModTime(), f)
	return
}

func (c *context) Attachment(file, name string) error {
	return c.contentDisposition(file, name, "attachment")
}

func (c *context) Inline(file, name string) error {
	return c.contentDisposition(file, name, "inline")
}

func (c *context) contentDisposition(file, name, dispositionType string) error {
	c.response.Header().Set(HeaderContentDisposition, fmt.Sprintf("%s; filename=%q", dispositionType, name))
	return c.File(file)
}

func (c *context) NoContent(code int) error {
	c.response.WriteHeader(code)
	return nil
}

func (c *context) Redirect(code int, url string) error {
	if code < 300 || code > 308 {
		return ErrInvalidRedirectCode
	}
	if c.session != nil {
		if err := c.session.Save(c); err != nil {
			return err
		}
	}
	url = c.echo.UrlLinker.Expand(c, url)
	c.response.Header().Set(HeaderLocation, url)
	c.response.WriteHeader(code)
	return nil
}

func (c *context) Refresh(code int) error {
	return c.Redirect(code, c.request.URL.String())
}

func (c *context) Revert(code int) error {
	return c.Redirect(code, c.request.Referer())
}

func (c *context) HeaderNoCache() {
	c.response.Writer.Header().Set("Cache-Control", "no-cache")
}

func (c *context) Error(err error) {
	c.echo.HTTPErrorHandler(c, err)
}

func (c *context) Echo() *Echo {
	return c.echo
}

func (c *context) Handler() HandlerFunc {
	return c.handler
}

func (c *context) SetHandler(h HandlerFunc) {
	c.handler = h
}

func (c *context) Logger() log.Logger {
	return c.echo.Logger
}

func (c *context) Session() Session {
	return c.session
}

func (c *context) SetSession(session Session) {
	c.session = session
}

func (c *context) Reset(r *http.Request, w http.ResponseWriter) {
	c.Context, c.request = makeContext(r)
	c.response.reset(w)
	c.query = nil
	c.handler = NotFoundHandler
	c.store = nil
	c.path = ""
	c.session = nil
	// NOTE: Don't reset because it has to have length c.echo.maxParam at all times
}

func (c *context) Locale() Locale {
	return c.locale
}

func (c *context) SetLocale(locale Locale) {
	c.locale = locale
}

func (c *context) AddFlash(class FlashClass, message interface{}) error {
	msg, err := RenderWidget(c, message)
	if err != nil {
		return err
	}

	c.session.AddFlash(class, msg)

	return nil
}

func (c *context) Value(key interface{}) interface{} {
	if key == ContextKey {
		return Context(c)
	}

	return c.Context.Value(key)
}

var boolMap = map[string]bool{
	"yes":   true,
	"no":    false,
	"true":  true,
	"false": false,
	"0":     false,
	"1":     true,
	"on":    true,
	"off":   false,
}

// A valueCtx carries a key-value pair. It implements Value for that key and
// delegates all other calls to the embedded Context.
type valueCtx struct {
	Context
	key, val interface{}
}

func (c *valueCtx) Value(key interface{}) interface{} {
	if c.key == key {
		return c.val
	}
	return c.Context.Value(key)
}

func ContextFromRequest(r *http.Request) Context {
	return r.Context().Value(ContextKey).(Context)
}
