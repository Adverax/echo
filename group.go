package echo

/*
import (
	"net/http"
	"path"
)

// Abstract VUX interface
type Mux interface {
	Pre(middleware ...MiddlewareFunc)
	Use(middleware ...MiddlewareFunc)
	CONNECT(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	DELETE(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	GET(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	HEAD(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	OPTIONS(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	PATCH(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	POST(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	PUT(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	TRACE(path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	FORM(path string, h HandlerFunc, m ...MiddlewareFunc) []*Route
	Any(path string, h HandlerFunc, m ...MiddlewareFunc) []*Route
	Match(methods []string, path string, h HandlerFunc, m ...MiddlewareFunc) []*Route
	Static(prefix, root string) *Route
	File(path, file string, m ...MiddlewareFunc) *Route
	Add(method, path string, h HandlerFunc, m ...MiddlewareFunc) *Route
	Group(prefix string, m ...MiddlewareFunc) Mux
	Union(fn func(mux Mux), m ...MiddlewareFunc)
	Route(prefix string, fn func(mux Mux), m ...MiddlewareFunc)
}

// Group is a set of sub-routes for a specified route. It can be used for inner
// routes that share a common middleware or functionality that should be separate
// from the parent echo instance while still inheriting from it.
type group struct {
	prefix     string
	middleware []MiddlewareFunc
	echo       *Echo
}

// Pre implements `Echo#pre()` for sub-routes within the Group.
func (g *group) Pre(middleware ...MiddlewareFunc) {
	g.middleware = append(middleware, g.middleware...)
	// Allow all requests to reach the group as they might get dropped if router
	// doesn't find a match, making none of the group middleware process.
	for _, p := range []string{"", "/*"} {
		g.echo.Any(path.Clean(g.prefix+p), func(c Context) error {
			return NotFoundHandler(c)
		}, g.middleware...)
	}
}

// Use implements `Echo#Use()` for sub-routes within the Group.
func (g *group) Use(middleware ...MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
	// Allow all requests to reach the group as they might get dropped if router
	// doesn't find a match, making none of the group middleware process.
	for _, p := range []string{"", "/*"} {
		g.echo.Any(path.Clean(g.prefix+p), func(c Context) error {
			return NotFoundHandler(c)
		}, g.middleware...)
	}
}

// CONNECT implements `Echo#CONNECT()` for sub-routes within the Group.
func (g *group) CONNECT(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodConnect, path, h, m...)
}

// DELETE implements `Echo#DELETE()` for sub-routes within the Group.
func (g *group) DELETE(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodDelete, path, h, m...)
}

// GET implements `Echo#GET()` for sub-routes within the Group.
func (g *group) GET(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodGet, path, h, m...)
}

// HEAD implements `Echo#HEAD()` for sub-routes within the Group.
func (g *group) HEAD(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodHead, path, h, m...)
}

// OPTIONS implements `Echo#OPTIONS()` for sub-routes within the Group.
func (g *group) OPTIONS(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodOptions, path, h, m...)
}

// PATCH implements `Echo#PATCH()` for sub-routes within the Group.
func (g *group) PATCH(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodPatch, path, h, m...)
}

// POST implements `Echo#POST()` for sub-routes within the Group.
func (g *group) POST(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodPost, path, h, m...)
}

// PUT implements `Echo#PUT()` for sub-routes within the Group.
func (g *group) PUT(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodPut, path, h, m...)
}

// TRACE implements `Echo#TRACE()` for sub-routes within the Group.
func (g *group) TRACE(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return g.Add(http.MethodTrace, path, h, m...)
}

// FORM implements `Echo#Form()` for sub-routes within the Group.
func (g *group) FORM(path string, h HandlerFunc, m ...MiddlewareFunc) []*Route {
	return []*Route{
		g.Add(http.MethodGet, path, h, m...),
		g.Add(http.MethodPost, path, h, m...),
	}
}

// Any implements `Echo#Any()` for sub-routes within the Group.
func (g *group) Any(path string, h HandlerFunc, m ...MiddlewareFunc) []*Route {
	routes := make([]*Route, len(methods))
	for i, method := range methods {
		routes[i] = g.Add(method, path, h, m...)
	}
	return routes
}

// Match implements `Echo#Match()` for sub-routes within the Group.
func (g *group) Match(methods []string, path string, h HandlerFunc, m ...MiddlewareFunc) []*Route {
	routes := make([]*Route, len(methods))
	for i, method := range methods {
		routes[i] = g.Add(method, path, h, m...)
	}
	return routes
}

// Group creates a new sub-group with prefix and optional sub-group-level middleware.
func (g *group) Group(prefix string, m ...MiddlewareFunc) Mux {
	ms := make([]MiddlewareFunc, 0, len(g.middleware)+len(m))
	ms = append(ms, g.middleware...)
	ms = append(ms, m...)
	return g.echo.Group(g.prefix+prefix, ms...)
}

// Union creates a new sub-group with prefix and optional sub-group-level middleware.
// After that, routine calls custom function with this group.
func (g *group) Union(fn func(mux Mux), m ...MiddlewareFunc) {
	fn(g.Group("", m...))
	return
}

// Route creates a new sub-group with prefix and optional sub-group-level middleware.
// After that, routine calls custom function with this group.
func (g *group) Route(prefix string, fn func(mux Mux), m ...MiddlewareFunc) {
	fn(g.Group(prefix, m...))
	return
}

// Static implements `Echo#Static()` for sub-routes within the Group.
func (g *group) Static(prefix, root string) *Route {
	return static(g, prefix, root)
}

// File implements `Echo#File()` for sub-routes within the Group.
func (g *group) File(path, file string, m ...MiddlewareFunc) *Route {
	return g.echo.File(g.prefix+path, file, m...)
}

// Add implements `Echo#Add()` for sub-routes within the Group.
func (g *group) Add(method, path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	// Combine into a new slice to avoid accidentally passing the same slice for
	// multiple routes, which would lead to later add() calls overwriting the
	// middleware from earlier calls.
	ms := make([]MiddlewareFunc, 0, len(g.middleware)+len(m))
	ms = append(ms, g.middleware...)
	ms = append(ms, m...)
	return g.echo.Add(method, g.prefix+path, h, ms...)
}
*/
