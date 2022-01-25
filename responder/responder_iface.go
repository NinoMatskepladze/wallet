package responder

import (
	"context"
	"net/http"
)

// Iface for responder methods, returning responses regarding its response with body, error or just an empty one
type Iface interface {
	JSON(w http.ResponseWriter, statusCode int, obj interface{})
	Error(ctx context.Context, w http.ResponseWriter, err error)
	Empty(w http.ResponseWriter)
}
