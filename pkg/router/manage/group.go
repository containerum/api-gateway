package manage

import (
	"errors"
	"net/http"
)

const (
	getAllGroupMethod = "Get all group"
)

var (
	//ErrUnableFindGroups - error when unable get groups from db
	ErrUnableFindGroups = errors.New("Unable to find groups")
)

//GetAllGroup return listeners list
func (m manage) GetAllGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groups, err := (*m.st).GetGroupList(nil)
		if err != nil {
			WriteAnswer(http.StatusBadRequest, getAllGroupMethod, &w, nil, ErrUnableFindGroups)
			return
		}
		WriteAnswer(http.StatusOK, getAllGroupMethod, &w, groups)
	}
}
