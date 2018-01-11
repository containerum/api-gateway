package manage

import (
	"encoding/json"
	"net/http"
	"time"

	"git.containerum.net/ch/api-gateway/pkg/store"
	logdec "git.containerum.net/ch/api-gateway/pkg/utils/logger"
	"git.containerum.net/ch/api-gateway/pkg/utils/snake"

	log "github.com/Sirupsen/logrus"
)

//ManagerHandlers implements manage handlers
type ManagerHandlers interface {
	//Listeners
	GetAllRouter() http.HandlerFunc
	GetRouter() http.HandlerFunc
	CreateRouter() http.HandlerFunc
	UpdateRouter() http.HandlerFunc
	DeleteRouter() http.HandlerFunc
	//Groups
	GetAllGroup() http.HandlerFunc
	// createGroup(router *Router) http.HandlerFunc
}

type manage struct {
	st *store.Store
}

type errAnswer struct {
	Errors []string
}

var (
	errLogger = log.WithField("Time", time.Now().Format(time.RFC822))
)

//NewManager return managers handlers
func NewManager(st *store.Store) ManagerHandlers {
	return manage{st}
}

//WriteAnswer render answer with all headers and json content
func WriteAnswer(status int, answerObject interface{}, errs *[]error, reqName string, w *http.ResponseWriter) error {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("X-Request-Name", snake.StrToSnake(reqName))

	log.WithField("Errors", errs).Debug("Error list")

	var answerBytes []byte
	var err error
	//Make good answer
	if status >= http.StatusOK && status <= http.StatusNoContent && errs == nil {
		answerBytes, err = json.Marshal(answerObject)
		if err != nil {
			status = http.StatusInternalServerError
			logdec.DecorateLoggerWithRuntimeContext(errLogger).WithError(err).Error("Unable make answer")
		}
	}
	//Make error answer
	if (status < http.StatusOK && status > http.StatusNoContent) || errs != nil {
		answer := errAnswer{}
		for _, e := range *errs {
			answer.Errors = append(answer.Errors, e.Error())
		}
		answerBytes, err = json.Marshal(answer)
		log.WithField("Answer", string(answerBytes)).Debug("Error answer")
		if err != nil {
			status = http.StatusInternalServerError
			logdec.DecorateLoggerWithRuntimeContext(errLogger).WithError(err).Error("Unable make error answer")
		}
	}
	//If answer is NULL, write empty string
	if string(answerBytes) == "null" {
		answerBytes = []byte("")
	}

	(*w).WriteHeader(status)
	(*w).Write(answerBytes)
	return nil
}
