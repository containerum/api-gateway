package router

import (
	"net/http"
)

func noRouteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No route"))
	}
}

func rootRouteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello! I am API Gateway."))
	}
}
