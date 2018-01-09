// +build !stub

package manage

import "net/http"

//GetAllRouter return listeners list
func (m manage) GetAllRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("Not stub")
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
