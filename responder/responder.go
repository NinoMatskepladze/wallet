package responder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"syscall"
	"time"

	customErrors "github.com/NinoMatskepladze/wallet/errors"
	"go.uber.org/zap"
)

const (
	DefaultErrMsg = "internal server error - our team is already looking into it"
)

type Responder struct {
	logger *zap.SugaredLogger
}

func NewResponder(log *zap.SugaredLogger) *Responder {
	return &Responder{
		logger: log,
	}
}

func (r *Responder) JSON(w http.ResponseWriter, statusCode int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		if errors.Is(err, syscall.EPIPE) {
			r.logger.Error("broken pipe - could not encode response\n%v", err)
			return
		}
		r.logger.Error("could not write error response to response writer\n%v", err)
		r.Error(context.Background(), w, err)
	}
}

func (r *Responder) Error(ctx context.Context, w http.ResponseWriter, err error) {
	s := &struct {
		Code    int       `json:"code"`
		Message string    `json:"message"`
		Date    time.Time `json:"date"`
	}{
		Date: time.Now(),
	}

	switch t := err.(type) {
	case *customErrors.ValidationError:
		s.Code = http.StatusBadRequest
		s.Message = t.Error()
	default:
		s.Code = http.StatusInternalServerError
		s.Message = DefaultErrMsg
		r.logger.Error(fmt.Errorf("server error at [%s]: %+v", s.Date, err))
	}

	r.JSON(w, s.Code, s)
}

func (r *Responder) Empty(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
