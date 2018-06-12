package main 

import (
	"net/http"
	"context"
//	"log"
//	"fmt"

)

func AddContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	lo("val1")
		token := r.Header.Get("Authorization")
//		lo("token: " + token)
		var ok bool
		var claims interface{}
		if len(token) > 7 {
//			lo("checking")
			claims, ok = ValidateToken(token[7:])
		} else {
			ok = false
		}
		if ok {
//			log.Println("valid token" + fmt.Sprint(claims["id"]))
			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			unAuthHTTPReturn(w, r)
//			next.ServeHTTP(w, r)
		}
	})
}