package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/adverax/echo"
)

func TestRequestID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	rid := RequestIDWithConfig(RequestIDConfig{})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := rid(handler)
	h(c)
	assert.Len(t, rec.Header().Get(echo.HeaderXRequestID), 32)

	// Custom generator
	rid = RequestIDWithConfig(RequestIDConfig{
		Generator: func() string { return "customGenerator" },
	})
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	h = rid(handler)
	h(c)
	assert.Equal(t, rec.Header().Get(echo.HeaderXRequestID), "customGenerator")
}
