package middlewares

import (
	"context"
	"github.com/authentication-app-server/helpers"
	"net/http"
)

func TokenValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		customClaims, err := helpers.ValidateToken(r.Header.Get("Authorization"))
		if err != nil {
			errorMap := helpers.ErrorsMap{
				Success: false,
				Message: map[string]string{
					"unauthorized": "You have not access to this area.",
				},
			}
			_, output := helpers.CustomResponse(nil, errorMap)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(output)
			return
		}

		ctx := context.WithValue(r.Context(), "ID", customClaims.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
