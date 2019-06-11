package middleware

import (
	"github.com/adverax/echo"
	"net/http"
	"net/http/httptest"
	"time"
)

// Cached is middleware for cache whole html page.
func PageCache(
	e *echo.Echo,
	class string,
	dependencies func(c echo.Context) (map[string]string, error),
	duration time.Duration,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var ctx echo.Context
			if e.Complex {
				ctx = echo.RequestContext(r)
			} else {
				ctx = e.AcquireContext()
				defer e.ReleaseContext(ctx)
				ctx.Reset(r, w)
			}

			deps, err := dependencies(ctx)
			if err != nil {
				ctx.Error(err)
				return
			}

			var content []byte
			err = e.Cacher.FetchData(
				class,
				deps,
				&content,
				func() (interface{}, error) {
					rec := httptest.NewRecorder()
					next.ServeHTTP(rec, r)

					for k, v := range rec.Header() {
						w.Header()[k] = v
					}

					w.WriteHeader(rec.Code)
					return rec.Body.Bytes(), nil
				},
				duration,
			)

			if err != nil {
				ctx.Error(err)
				return
			}

			_, _ = w.Write(content)
		}

		return http.HandlerFunc(fn)
	}
}
