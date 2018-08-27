package server

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/containerum/cherry"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/gin-gonic/gin"
)

func statusWithError(serviceName string, err error) model.ServiceStatus {
	return model.ServiceStatus{
		Name:     serviceName,
		StatusOK: false,
		Details: map[string]string{
			"error": err.Error(),
		},
	}
}

func (s *Server) getServiceStatus(backend *url.URL) model.ServiceStatus {
	resp, err := http.Get(s.addPrefixToBackendURL(backend).String())
	if err != nil {
		return statusWithError(backend.Hostname(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var returnedErr cherry.Err
		if err := json.NewDecoder(resp.Body).Decode(&returnedErr); err != nil {
			return statusWithError(backend.Hostname(), err)
		}
		return statusWithError(backend.Hostname(), &returnedErr)
	}

	var returnedStatus model.ServiceStatus
	if err := json.NewDecoder(resp.Body).Decode(&returnedStatus); err != nil {
		return statusWithError(backend.Hostname(), err)
	}
	return returnedStatus
}

func (s *Server) myStatus() model.ServiceStatus {
	return model.ServiceStatus{
		Name:     "api-gateway",
		Version:  s.options.Version,
		StatusOK: true,
	}
}

func (s *Server) healthCheckHandler(ctx *gin.Context) {
	var serviceStatuses []model.ServiceStatus
	for _, serviceURL := range s.options.Config.HealthCheck.URLs {
		serviceStatuses = append(serviceStatuses, s.getServiceStatus(serviceURL.URL))
	}
	serviceStatuses = append(serviceStatuses, s.myStatus())
	ctx.JSON(http.StatusOK, serviceStatuses)
}
