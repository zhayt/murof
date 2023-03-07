package controller

import (
	"fmt"
	"log"
	"net/http"
)

func (c *Controller) switchController(next []http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			next[0].ServeHTTP(w, r)
		case http.MethodPost:
			next[1].ServeHTTP(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	}
}

func (c *Controller) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("%s - %s %s", r.Proto, r.Method, r.URL.String()))
		next.ServeHTTP(w, r)
	})
}
