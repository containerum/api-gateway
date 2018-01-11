package router

import (
	"encoding/json"
	"net/http"
	"time"

	logdec "git.containerum.net/ch/api-gateway/pkg/utils/logger"
	"git.containerum.net/ch/api-gateway/pkg/utils/snake"

	log "github.com/Sirupsen/logrus"
)

var (
	errLogger = log.WithField("Time", time.Now().Format(time.RFC822))
)

type errAnswer struct {
	Errors []string
}

func noRouteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No route"))
	}
}

func rootRouteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Containerum.io API Gateway"))
	}
}

//WriteAnswer render answer with all headers and json content
func WriteAnswer(status int, answerObject interface{}, errs *[]error, reqName string, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-Name", snake.StrToSnake(reqName))

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
	w.WriteHeader(status)
	w.Write(answerBytes)
	return nil
}
