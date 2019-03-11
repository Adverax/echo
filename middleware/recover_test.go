package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/adverax/echo"
	"github.com/adverax/echo/log"
)

func TestRecover(t *testing.T) {
	e := echo.New()
	buf := new(bytes.Buffer)
	e.Logger = log.New(buf, buf, buf, buf, "")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := Recover()(echo.HandlerFunc(func(c echo.Context) error {
		panic("test")
	}))
	h(c)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, buf.String(), "PANIC RECOVER")
}
