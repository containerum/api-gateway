package router

import (
	"net/http"

	"git.containerum.net/ch/api-gateway/pkg/router/manage"

	"github.com/go-chi/chi"
)

//CreateManageRouter return manage handlers
func CreateManageRouter(router *Router) http.Handler {
	r := chi.NewRouter()
	// Router handlers
	m := manage.NewManager(router.store)
	/* Listeners */
	r.Get("/route", m.GetAllRouter())
	r.Post("/route", m.CreateRouter())
	r.Get("/route/{id}", m.GetRouter())
	r.Put("/route/{id}", m.UpdateRouter())
	r.Delete("/route/{id}", m.DeleteRouter())
	/* Groups */
	r.Get("/group", m.GetAllGroup())

	return r
}
