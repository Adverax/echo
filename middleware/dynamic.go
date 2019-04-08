package middleware

import (
	"github.com/adverax/echo"
	"net/http"
)

// WithDynamic is a middleware that added supports dynamic content in the COMPLEX mode.
// NOTE. Do not forget to switch the flag e.Complex = true.
func WithDynamic(
	e *echo.Echo,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			e.Dynamic(w, r, next)
		}

		return http.HandlerFunc(fn)
	}
}
