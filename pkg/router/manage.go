package router

import (
	"net/http"

	"github.com/go-chi/chi"
)

//CreateManageRouter return manage handlers
func CreateManageRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getAllRouter)
	r.Post("/", addRouter)
	r.Get("/{id}", getRouter)
	r.Put("/{id}", updateRouter)
	r.Delete("/{id}", removeRouter)
	return r
}

func getAllRouter(w http.ResponseWriter, r *http.Request) {
	rs, err := (*st).GetRoutesList()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	for _, rr := range *rs {
		w.Write([]byte(rr.ID))
	}
}

func getRouter(w http.ResponseWriter, r *http.Request) {
	rs, err := (*st).GetRouter(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rs.ID))
}

func addRouter(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create route"))
}

func updateRouter(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update route"))
}

func removeRouter(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Remove route"))
}
