// +build stub

package manage

import (
	"net/http"

	"git.containerum.net/ch/api-gateway/pkg/model"
)

var (
	l1 = model.Listener{
		DefaultModel: model.DefaultModel{
			ID: "11472afa-48ff-4583-b2e4-74119993f22a",
		},
		Name:   "Router N1",
		OAuth:  newBool(false),
		Active: newBool(true),
	}
	l2 = model.Listener{
		DefaultModel: model.DefaultModel{
			ID: "fdbcf79a-c254-48ec-b1df-60a0984fab5e",
		},
		Name:   "Router N1",
		OAuth:  new(bool),
		Active: new(bool),
	}
	l3 = model.Listener{
		DefaultModel: model.DefaultModel{
			ID: "126a743d-add6-4de9-971d-6691f316d530",
		},
		Name:   "Router N1",
		OAuth:  new(bool),
		Active: new(bool),
	}
	l4 = model.Listener{
		DefaultModel: model.DefaultModel{
			ID: "9fd10b6a-f58c-4d3f-90b9-0990ee80f684",
		},
		Name:   "Router N1",
		OAuth:  new(bool),
		Active: new(bool),
	}
)

//GetAllRouter return listeners list
func (m manage) GetAllRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqName := "Get all router stub"
		listeners := append([]model.Listener{},
			l1, l2, l3, l4,
		)
		WriteAnswer(http.StatusOK, &listeners, nil, reqName, &w)
	}
}

//GetRouter return listeners by id
func (m manage) GetRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

//CreateRouter return listeners list
func (m manage) CreateRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

//UpdateRouter return listeners list
func (m manage) UpdateRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

//DeleteRouter return listeners list
func (m manage) DeleteRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func newBool(b bool) *bool {
	res := new(bool)
	*res = b
	return res
}
