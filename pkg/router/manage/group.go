package manage

import (
	"errors"
	"net/http"

	"git.containerum.net/ch/api-gateway/pkg/model"
)

var (
	//ErrUnableFindGroups - error when unable get groups from db
	ErrUnableFindGroups = errors.New("Unable to find groups")
)

//GetAllGroup return listeners list
func (m manage) GetAllGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqName := "Get all group"
		groups, err := (*m.st).GetGroupList(&model.Group{})
		if err != nil {
			WriteAnswer(http.StatusBadRequest, nil, &[]error{ErrUnableFindGroups}, reqName, &w)
			return
		}
		WriteAnswer(http.StatusOK, groups, nil, reqName, &w)
	}
}
