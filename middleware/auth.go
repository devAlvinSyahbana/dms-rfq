package middlewares

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/devAlvinSyahbana/golang-rfq/service"
)

type authString string

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth == "" {
			next.ServeHTTP(w, r)
			return
		}

		bearer := "Bearer "
		auth = auth[len(bearer):]
		validate, err := service.JwtValidate(context.Background(), auth)
		if err != nil || !validate.Valid {
			jData, _ := json.Marshal(map[string]string{"Error": "Invalid Auth token"})
			w.Header().Set("Content-Type", "application/json")
			w.Write(jData)
			return
		}

		ctx := context.WithValue(r.Context(), authString("auth"), validate.Claims)
		// w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5173")
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func CtxValue(ctx context.Context) *service.JwtCustomClaim {
	raw, _ := ctx.Value(authString("auth")).(*service.JwtCustomClaim)
	return raw
}
func CtxValueRaw(ctx context.Context) string {
	raw, _ := ctx.Value(authString("auth")).(string)
	return raw
}
