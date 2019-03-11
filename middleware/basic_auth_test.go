package middleware

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/adverax/echo"
)

func TestBasicAuth(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	f := func(u, p string, c echo.Context) (bool, error) {
		if u == "joe" && p == "secret" {
			return true, nil
		}
		return false, nil
	}
	h := BasicAuth(f)(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	assert := assert.New(t)

	// Valid credentials
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	auth := basic + " " + base64.StdEncoding.EncodeToString([]byte("joe:secret"))
	req.Header.Set(echo.HeaderAuthorization, auth)
	assert.NoError(h(c))

	h = BasicAuthWithConfig(BasicAuthConfig{
		Skipper:   nil,
		Validator: f,
		Realm:     "someRealm",
	})(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	// Valid credentials
	res = httptest.NewRecorder()
	c = e.NewContext(req, res)
	auth = basic + " " + base64.StdEncoding.EncodeToString([]byte("joe:secret"))
	req.Header.Set(echo.HeaderAuthorization, auth)
	assert.NoError(h(c))

	// Case-insensitive header scheme
	res = httptest.NewRecorder()
	c = e.NewContext(req, res)
	auth = strings.ToUpper(basic) + " " + base64.StdEncoding.EncodeToString([]byte("joe:secret"))
	req.Header.Set(echo.HeaderAuthorization, auth)
	assert.NoError(h(c))

	// Invalid credentials
	res = httptest.NewRecorder()
	c = e.NewContext(req, res)
	auth = basic + " " + base64.StdEncoding.EncodeToString([]byte("joe:invalid-password"))
	req.Header.Set(echo.HeaderAuthorization, auth)
	he := h(c).(*echo.HTTPError)
	assert.Equal(http.StatusUnauthorized, he.Code)
	assert.Equal(basic+` realm="someRealm"`, res.Header().Get(echo.HeaderWWWAuthenticate))

	// Missing Authorization header
	res = httptest.NewRecorder()
	c = e.NewContext(req, res)
	req.Header.Del(echo.HeaderAuthorization)
	he = h(c).(*echo.HTTPError)
	assert.Equal(http.StatusUnauthorized, he.Code)

	// Invalid Authorization header
	res = httptest.NewRecorder()
	c = e.NewContext(req, res)
	auth = base64.StdEncoding.EncodeToString([]byte("invalid"))
	req.Header.Set(echo.HeaderAuthorization, auth)
	he = h(c).(*echo.HTTPError)
	assert.Equal(http.StatusUnauthorized, he.Code)
}
