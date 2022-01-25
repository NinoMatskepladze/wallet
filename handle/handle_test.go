package handle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NinoMatskepladze/wallet/models"
	"github.com/NinoMatskepladze/wallet/responder"
	"github.com/NinoMatskepladze/wallet/test"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap"
)

func TestHandler_createWallet(t *testing.T) {
	l := zap.NewExample().Sugar()

	cases := []struct {
		test.HTTPTestCase
		want *models.CreateWalletResponse
	}{
		{
			test.HTTPTestCase{
				Name: "should succeed",
				ReqFn: func() *http.Request {
					ctx := context.Background()
					r, err := http.NewRequestWithContext(ctx, "POST", "/wallets", nil)
					require.Nil(t, err)
					return r
				},
				WantCode: http.StatusCreated,
			},
			&models.CreateWalletResponse{
				Wallet: models.Wallet{
					ID:      "",
					Balance: 0,
				},
			},
		},
	}
	for _, c := range cases {
		svc := test.FakeService{}
		res := responder.NewResponder(l)
		cont := NewServiceRoutes(svc, res)
		handler := chi.NewRouter()

		t.Run(c.Name, func(t *testing.T) {
			handler.Post("/wallets", cont.CreateWallet)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, c.HTTPTestCase.ReqFn())
			got := rr.Result()

			require.Equal(t, c.HTTPTestCase.WantCode, got.StatusCode)

			body, err := io.ReadAll(got.Body)
			require.Nil(t, err)
			require.Nil(t, got.Body.Close())

			if c.want != nil {
				resp := new(models.CreateWalletResponse)
				require.Nil(t, json.Unmarshal(body, resp))
			}
		})
	}
}

func TestHandler_updateBalance(t *testing.T) {
	l := zap.NewExample().Sugar()
	fakeID := "FakeID"
	balance := 100

	cases := []struct {
		test.HTTPTestCase
		want *models.CreateWalletResponse
	}{
		{
			test.HTTPTestCase{
				Name: "should succeed",
				ReqFn: func() *http.Request {
					mr := &models.UpdateWalletRequest{
						Amount: balance,
					}
					body, err := json.Marshal(mr)
					require.Nil(t, err)
					ctx := context.Background()
					r, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("/wallets/%s", fakeID), bytes.NewReader(body))
					require.Nil(t, err)
					return r
				},
				WantCode: http.StatusNoContent,
			},
			&models.CreateWalletResponse{
				Wallet: models.Wallet{
					ID:      models.WalletID(fakeID),
					Balance: int64(balance),
				},
			},
		},
		{
			test.HTTPTestCase{
				Name: "should fail",
				ReqFn: func() *http.Request {
					body := []byte(`{amount: "200"}`)
					ctx := context.Background()
					r, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("/wallets/%s", fakeID), bytes.NewReader(body))
					require.Nil(t, err)
					return r
				},
				WantCode: http.StatusBadRequest,
			},
			nil,
		},
	}
	for _, c := range cases {
		svc := test.FakeService{}
		res := responder.NewResponder(l)
		cont := NewServiceRoutes(svc, res)
		handler := chi.NewRouter()

		t.Run(c.Name, func(t *testing.T) {
			handler.Post("/wallets/{wallet_id}", cont.UpdateBalance)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, c.HTTPTestCase.ReqFn())
			got := rr.Result()

			require.Equal(t, c.HTTPTestCase.WantCode, got.StatusCode)

		})
	}
}

func TestHandler_getWallet(t *testing.T) {
	l := zap.NewExample().Sugar()
	fakeID := "FakeID"
	balance := 100

	cases := []struct {
		test.HTTPTestCase
		want *models.GetWalletResponse
	}{
		{
			test.HTTPTestCase{
				Name: "should succeed",
				ReqFn: func() *http.Request {
					ctx := context.Background()
					r, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("/wallets/%s", fakeID), nil)
					require.Nil(t, err)
					return r
				},
				WantCode: http.StatusOK,
			},
			&models.GetWalletResponse{
				Wallet: models.Wallet{
					ID:      models.WalletID(fakeID),
					Balance: int64(balance),
				},
			},
		},
	}
	for _, c := range cases {
		svc := test.FakeService{}
		res := responder.NewResponder(l)
		cont := NewServiceRoutes(svc, res)
		handler := chi.NewRouter()

		t.Run(c.Name, func(t *testing.T) {
			handler.Get("/wallets/{wallet_id}", cont.GetWallet)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, c.HTTPTestCase.ReqFn())
			got := rr.Result()

			require.Equal(t, c.HTTPTestCase.WantCode, got.StatusCode)
			body, err := io.ReadAll(got.Body)
			require.Nil(t, err)

			if c.want != nil {
				resp := new(models.GetWalletResponse)
				require.Nil(t, json.Unmarshal(body, resp))

				require.Equal(t, c.want, resp)
			}
		})
	}
}
