package handle

import (
	"encoding/json"
	"net/http"

	"github.com/NinoMatskepladze/wallet/errors"
	"github.com/NinoMatskepladze/wallet/models"
	"github.com/NinoMatskepladze/wallet/responder"
	"github.com/NinoMatskepladze/wallet/service"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ServiceRoutes for wallet routes
type ServiceRoutes struct {
	service service.ServiceIface
	res     responder.Iface
}

// NewServiceRoutes returns new wallet
func NewServiceRoutes(service service.ServiceIface, res responder.Iface) *ServiceRoutes {
	return &ServiceRoutes{
		service: service,
		res:     res,
	}
}

// CreateWallet controller function
func (r *ServiceRoutes) CreateWallet(w http.ResponseWriter, req *http.Request) {
	wallet, err := r.service.CreateWallet(req.Context())
	if err != nil {
		r.res.Error(req.Context(), w, err)
		return
	}
	r.res.JSON(w, http.StatusCreated, &models.CreateWalletResponse{
		Wallet: wallet,
	})
}

// UpdateBalance controller function for wallet balance increase/decreases
func (r *ServiceRoutes) UpdateBalance(w http.ResponseWriter, req *http.Request) {
	walletID := chi.URLParam(req, "wallet_id")
	updateWalletReq := &models.UpdateWalletRequest{}

	err := json.NewDecoder(req.Body).Decode(updateWalletReq)
	if err != nil {
		r.res.Error(req.Context(), w, &errors.ValidationError{
			HttpError: errors.HttpError{
				HttpCode: 400,
				Message:  err.Error(),
			},
		})
		return
	}

	// Validations for UpdateWalletRequest struct
	err = validate.Struct(updateWalletReq)
	if err != nil {
		r.res.Error(req.Context(), w, &errors.ValidationError{
			HttpError: errors.HttpError{HttpCode: 400, Message: err.Error()},
		})
		return
	}

	err = r.service.UpdateBalance(req.Context(), walletID, *updateWalletReq)

	if err != nil {
		r.res.Error(req.Context(), w, err)
		return
	}
	r.res.Empty(w)
}

// GetWallet for getting wallet information
func (r *ServiceRoutes) GetWallet(w http.ResponseWriter, req *http.Request) {
	walletID := chi.URLParam(req, "wallet_id")

	wallet, err := r.service.GetWallet(req.Context(), walletID)
	if err != nil {
		r.res.Error(req.Context(), w, err)
		return
	}
	r.res.JSON(w, http.StatusOK, wallet)
}
