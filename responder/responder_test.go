package responder

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	customErrors "github.com/NinoMatskepladze/wallet/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"github.com/stretchr/testify/require"
)

func TestResponder_Error(t *testing.T) {
	cases := []struct {
		name        string
		err         error
		wantMessage string
		wantInLog   string
		wantCode    int
	}{
		{
			name:        "should return validation error",
			err:         customErrors.NewValidationError("this is a validation error"),
			wantMessage: "this is a validation error",
			wantInLog:   "",
			wantCode:    400,
		},
		{
			name:        "should return not found error",
			err:         customErrors.NewNotFoundError("this is a not found error"),
			wantMessage: "this is a not found error",
			wantInLog:   "",
			wantCode:    404,
		},
		{
			name:        "should return internal service error",
			err:         errors.New("this is an internal service error that we might want to see in logs"),
			wantMessage: DefaultErrMsg,
			wantInLog:   "this is an internal service error that we might want to see in logs",
			wantCode:    500,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			core, logs := observer.New(zap.DebugLevel)
			logger := zap.New(core).Sugar()
			defer logger.Sync() // nolint:errcheck

			w := &httptest.ResponseRecorder{
				Body: &bytes.Buffer{},
			}
			res := NewResponder(logger)
			res.Error(context.Background(), w, c.err)

			body := &struct {
				ID      string    `json:"id"`
				Code    int       `json:"code"`
				Message string    `json:"message"`
				Date    time.Time `json:"date"`
			}{}
			err := json.NewDecoder(w.Body).Decode(body)
			require.Nil(t, err)

			if c.wantMessage != "" {
				require.Equal(t, c.wantMessage, body.Message)
			}

			if c.wantInLog != "" {
				require.Equal(t, 1, len(logs.All()))
				l := logs.All()[0]

				wantLevel := zap.ErrorLevel
				gotLevel := l.Level
				require.Equal(t, wantLevel, gotLevel)

				// make sure the error message is in the log
				require.True(t, strings.Contains(l.Message, c.wantInLog))

				// make sure the random response id is included in the log
				require.True(t, strings.Contains(l.Message, body.ID))
			}

			gotCode := w.Code
			require.Equal(t, c.wantCode, gotCode)
		})
	}
}
