package router

import (
	"net/http"
	"os"
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

// TODO: make checker
func statusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, err := os.Hostname()
		if err != nil {
			// TODO: return error
			panic(err)
		}
		w.Write([]byte(name))
	}
}
