package http

import (
	"context"
	"encoding/json"
	httpError "github.com/xesina/go-kit-realworld-example-app/http/error"
	"github.com/xesina/go-kit-realworld-example-app/user"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, response interface{}, code int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(response)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(user.Response); ok && e.Err != nil {
		httpError.EncodeError(ctx, e.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}