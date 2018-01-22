package grpcutils

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

var GRPCToHTTPCode = map[codes.Code]int{
	codes.OK:                 http.StatusOK,
	codes.Canceled:           444, // no response
	codes.Unknown:            http.StatusInternalServerError,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.DeadlineExceeded:   http.StatusRequestTimeout,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.ResourceExhausted:  http.StatusInsufficientStorage, // maybe better solution
	codes.FailedPrecondition: http.StatusPreconditionFailed,
	codes.Aborted:            http.StatusGatewayTimeout,
	codes.OutOfRange:         http.StatusRequestedRangeNotSatisfiable,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.DataLoss:           http.StatusInternalServerError,
	codes.Unauthenticated:    http.StatusUnauthorized, // authentication and authorization are not the same :)
}
