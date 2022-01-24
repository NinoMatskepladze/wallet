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
	"github.com/go-kit/kit/log"
)

const (
	DefaultErrMsg = "internal server error - our team is already looking into it"
)

type Responder struct {
	logger log.Logger
}

func NewResponder(log log.Logger) *Responder {
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
			r.logger.Log("broken pipe - could not encode response\n%v", err)
			return
		}

		r.logger.Log("could not write error response to response writer\n%v", err)
		r.Error(context.Background(), w, err)
	}
}

func (r *Responder) Error(ctx context.Context, w http.ResponseWriter, err error) {
	s := &struct {
		ID      string    `json:"id"`
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
		r.logger.Log(fmt.Errorf("server error [%s] at [%s]: %+v", s.ID, s.Date, err))
	}

	r.JSON(w, s.Code, s)
}

func (r *Responder) Empty(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
