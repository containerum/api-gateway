package manage

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

var (
	server *httptest.Server
)

func init() {
	r := chi.NewRouter()
	r.Get("/manage/route", getTest())

	server = httptest.NewServer(r)
}

func getTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}
}

func TestGetRouter(t *testing.T) {
	request, _ := http.NewRequest("GET", server.URL+"/manage/route", nil)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(res.StatusCode)
}
