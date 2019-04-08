package echo

import (
	stdContext "context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	user struct {
		ID   int    `json:"id" xml:"id" form:"id" query:"id"`
		Name string `json:"name" xml:"name" form:"name" query:"name"`
	}
)

const (
	userJSON                    = `{"id":1,"name":"Jon Snow"}`
	userXML                     = `<user><id>1</id><name>Jon Snow</name></user>`
	userForm                    = `id=1&name=Jon Snow`
	invalidContent              = "invalid content"
	userJSONInvalidType         = `{"id":"1","name":"Jon Snow"}`
	userXMLConvertNumberError   = `<user><id>Number one</id><name>Jon Snow</name></user>`
	userXMLUnsupportedTypeError = `<user><>Number one</><name>Jon Snow</name></user>`
)

const userJSONPretty = `{
  "id": 1,
  "name": "Jon Snow"
}`

const userXMLPretty = `<user>
  <id>1</id>
  <name>Jon Snow</name>
</user>`

func TestEcho(t *testing.T) {
	e := New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Router
	assert.NotNil(t, e.Router())

	// DefaultHTTPErrorHandler
	e.DefaultHTTPErrorHandler(c, errors.New("error"))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestEchoWrapHandler(t *testing.T) {
	e := New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	}))
	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "test", rec.Body.String())
	}
}

func TestEchoConnect(t *testing.T) {
	e := New()
	testMethod(t, http.MethodConnect, "/", e)
}

func TestEchoDelete(t *testing.T) {
	e := New()
	testMethod(t, http.MethodDelete, "/", e)
}

func TestEchoGet(t *testing.T) {
	e := New()
	testMethod(t, http.MethodGet, "/", e)
}

func TestEchoHead(t *testing.T) {
	e := New()
	testMethod(t, http.MethodHead, "/", e)
}

func TestEchoOptions(t *testing.T) {
	e := New()
	testMethod(t, http.MethodOptions, "/", e)
}

func TestEchoPatch(t *testing.T) {
	e := New()
	testMethod(t, http.MethodPatch, "/", e)
}

func TestEchoPost(t *testing.T) {
	e := New()
	testMethod(t, http.MethodPost, "/", e)
}

func TestEchoPut(t *testing.T) {
	e := New()
	testMethod(t, http.MethodPut, "/", e)
}

func TestEchoTrace(t *testing.T) {
	e := New()
	testMethod(t, http.MethodTrace, "/", e)
}

func TestEchoForm(t *testing.T) {
	e := New()
	e.FORM(e.Router(), "/", func(c Context) error {
		return c.String(http.StatusOK, "Any")
	})
}

func TestEchoNotFound(t *testing.T) {
	e := New()
	req := httptest.NewRequest(http.MethodGet, "/files", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestEchoMethodNotAllowed(t *testing.T) {
	e := New()
	e.GET(e.router, "/", func(c Context) error {
		return c.String(http.StatusOK, "Echo!")
	})
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestEchoContext(t *testing.T) {
	e := New()
	c := e.AcquireContext()
	assert.IsType(t, new(context), c)
	e.ReleaseContext(c)
}

func TestEchoStart(t *testing.T) {
	e := New()
	go func() {
		assert.NoError(t, e.Start(":0"))
	}()
	time.Sleep(200 * time.Millisecond)
}

func TestEchoStartTLS(t *testing.T) {
	e := New()
	go func() {
		err := e.StartTLS(":0", "_fixture/certs/cert.pem", "_fixture/certs/key.pem")
		// Prevent the test to fail after closing the servers
		if err != http.ErrServerClosed {
			assert.NoError(t, err)
		}
	}()
	time.Sleep(200 * time.Millisecond)

	e.Close()
}

func TestEchoStartTLSByteString(t *testing.T) {
	cert, err := ioutil.ReadFile("_fixture/certs/cert.pem")
	require.NoError(t, err)
	key, err := ioutil.ReadFile("_fixture/certs/key.pem")
	require.NoError(t, err)

	testCases := []struct {
		cert        interface{}
		key         interface{}
		expectedErr error
		name        string
	}{
		{
			cert:        "_fixture/certs/cert.pem",
			key:         "_fixture/certs/key.pem",
			expectedErr: nil,
			name:        `ValidCertAndKeyFilePath`,
		},
		{
			cert:        cert,
			key:         key,
			expectedErr: nil,
			name:        `ValidCertAndKeyByteString`,
		},
		{
			cert:        cert,
			key:         1,
			expectedErr: ErrInvalidCertOrKeyType,
			name:        `InvalidKeyType`,
		},
		{
			cert:        0,
			key:         key,
			expectedErr: ErrInvalidCertOrKeyType,
			name:        `InvalidCertType`,
		},
		{
			cert:        0,
			key:         1,
			expectedErr: ErrInvalidCertOrKeyType,
			name:        `InvalidCertAndKeyTypes`,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			e := New()
			e.HideBanner = true

			go func() {
				err := e.StartTLS(":0", test.cert, test.key)
				if test.expectedErr != nil {
					require.EqualError(t, err, test.expectedErr.Error())
				} else if err != http.ErrServerClosed { // Prevent the test to fail after closing the servers
					require.NoError(t, err)
				}
			}()
			time.Sleep(200 * time.Millisecond)

			require.NoError(t, e.Close())
		})
	}
}

func TestEchoStartAutoTLS(t *testing.T) {
	e := New()
	errChan := make(chan error, 0)

	go func() {
		errChan <- e.StartAutoTLS(":0")
	}()
	time.Sleep(200 * time.Millisecond)

	select {
	case err := <-errChan:
		assert.NoError(t, err)
	default:
		assert.NoError(t, e.Close())
	}
}

func testMethod(t *testing.T, method, path string, e *Echo) {
	r := reflect.ValueOf(e.router)
	p := reflect.ValueOf(path)
	h := reflect.ValueOf(func(c Context) error {
		return c.String(http.StatusOK, method)
	})
	i := interface{}(e)
	reflect.ValueOf(i).MethodByName(method).Call([]reflect.Value{r, p, h})
	_, body := request(method, path, e)
	assert.Equal(t, method, body)
}

func request(method, path string, e *Echo) (int, string) {
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func TestHTTPError(t *testing.T) {
	err := NewHTTPError(http.StatusBadRequest, map[string]interface{}{
		"code": 12,
	})
	assert.Equal(t, "code=400, message=map[code:12]", err.Error())
}

func TestEchoClose(t *testing.T) {
	e := New()
	errCh := make(chan error)

	go func() {
		errCh <- e.Start(":0")
	}()

	time.Sleep(200 * time.Millisecond)

	if err := e.Close(); err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, e.Close())

	err := <-errCh
	assert.Equal(t, err.Error(), "http: Server closed")
}

func TestEchoShutdown(t *testing.T) {
	e := New()
	errCh := make(chan error)

	go func() {
		errCh <- e.Start(":0")
	}()

	time.Sleep(200 * time.Millisecond)

	if err := e.Close(); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := stdContext.WithTimeout(stdContext.Background(), 10*time.Second)
	defer cancel()
	assert.NoError(t, e.Shutdown(ctx))

	err := <-errCh
	assert.Equal(t, err.Error(), "http: Server closed")
}
