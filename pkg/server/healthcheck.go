package server

import (
	"encoding/json"
	"net/http"
	"net/url"

	"git.containerum.net/ch/api-gateway/pkg/server/middleware"
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

func extractXHeaders(sourceHeaders http.Header) http.Header {
	ret := make(http.Header)
	for name, value := range sourceHeaders {
		if middleware.XHeaderRegexp.MatchString(name) {
			ret[name] = value
		}
	}
	return ret
}

func (s *Server) getServiceStatus(backend *url.URL, requestHeaders http.Header) model.ServiceStatus {
	req := http.Request{
		Method: http.MethodGet,
		Header: extractXHeaders(requestHeaders),
		URL:    s.addPrefixToBackendURL(backend),
	}
	resp, err := http.DefaultClient.Do(&req)
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
		serviceStatuses = append(serviceStatuses, s.getServiceStatus(serviceURL.URL, ctx.Request.Header))
	}
	serviceStatuses = append(serviceStatuses, s.myStatus())
	ctx.JSON(http.StatusOK, serviceStatuses)
}
