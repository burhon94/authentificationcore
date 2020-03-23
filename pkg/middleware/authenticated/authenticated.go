package authenticated

import (
	"context"
	"net/http"
)

func Authenticated(Auth func(ctx context.Context) bool) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			if !Auth(request.Context()) {
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			next(writer, request)
		}
	}
}

