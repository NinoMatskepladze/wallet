package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NinoMatskepladze/wallet/errors"
	"github.com/stretchr/testify/require"
)

// HTTPTestCase struct for testing purposes
type HTTPTestCase struct {
	Name       string
	ReqFn      func() *http.Request
	WantCode   int
	WantErrMsg string
}

// Run func runs http test cases with the help of testing package
func (c *HTTPTestCase) Run(t *testing.T, h http.Handler) func(t *testing.T) {
	return func(t *testing.T) {
		c.RunHTTPTestCase(t, h)
	}
}

// RunHTTPTestCase
func (c *HTTPTestCase) RunHTTPTestCase(t *testing.T, h http.Handler) []byte {
	t.Helper()
	// invoke req
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, c.ReqFn())
	resp := rr.Result()

	// validate response code
	require.Equal(t, c.WantCode, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	defer resp.Body.Close()
	t.Log("response", string(bodyBytes))

	if c.WantErrMsg != "" {
		require.NotEmpty(t, bodyBytes, "expected a non empty http response")
		gotErr := new(errors.HttpError)
		err := json.Unmarshal(bodyBytes, gotErr)
		require.Nil(t, err)
		require.Equal(t, c.WantErrMsg, gotErr.Message)
	}

	return bodyBytes
}
