package responder

import (
	"context"
	"net/http"
)

type Iface interface {
	JSON(w http.ResponseWriter, statusCode int, obj interface{})
	Error(ctx context.Context, w http.ResponseWriter, err error)
	Empty(w http.ResponseWriter)
}
