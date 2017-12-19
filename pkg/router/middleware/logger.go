package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	clickhouse "git.containerum.net/ch/api-gateway/pkg/utils/clickhouselog"

	log "github.com/Sirupsen/logrus"
	"github.com/cactus/go-statsd-client/statsd"
)

type LoggerResponseWritter interface {
	http.ResponseWriter
	Status() int
	BytesWritten() int
	Bytes() []byte
}

type loggerWritter struct {
	http.ResponseWriter
	wroteHeader bool
	code        int
	bytes       int
	bytesArr    []byte
	headers     []string
}

func NewLoggerResponseWritter(w http.ResponseWriter) LoggerResponseWritter {
	return &loggerWritter{ResponseWriter: w}
}

func (lw *loggerWritter) WriteHeader(code int) {
	if !lw.wroteHeader {
		lw.code = code
		lw.wroteHeader = true
		lw.ResponseWriter.WriteHeader(code)
	}
}

func (lw *loggerWritter) Write(buf []byte) (int, error) {
	lw.WriteHeader(http.StatusOK)
	n, err := lw.ResponseWriter.Write(buf)
	lw.bytes += n
	lw.bytesArr = append(lw.bytesArr, buf...)
	return n, err
}

func (lw *loggerWritter) Status() int {
	return lw.code
}

func (lw *loggerWritter) BytesWritten() int {
	return lw.bytes
}

func (lw *loggerWritter) Bytes() []byte {
	return lw.bytesArr
}

func Logger(stats *statsd.Statter, clickLogs *clickhouse.LogClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lw := NewLoggerResponseWritter(w)
			next.ServeHTTP(lw, r)
			latency := time.Now().Sub(start)

			//Set status in Statsd
			if stats != nil {
				statusCall := fmt.Sprintf("call.status.%v", lw.Status())
				methodCall := fmt.Sprintf("call.method.%v", r.Method)
				reqName := w.Header().Get("X-Request-Name")
				if reqName == "" {
					reqName = "system"
				}
				statusCallNamed := fmt.Sprintf("call.route.%v.status.%v", reqName, lw.Status())
				s := *stats
				s.Inc("call.status.all", 1, 1.0)
				s.Inc(statusCall, 1, 1.0)
				s.Inc(methodCall, 1, 1.0)
				s.Inc(statusCallNamed, 1, 1.0)
			}

			var reqBody []byte
			_, err := r.Body.Read(reqBody)
			if err != nil {
				// TODO
			}

			headersRequest, err := json.Marshal(r.Header)
			if err != nil {
				//TODO
			}

			headersResponse, err := json.Marshal(w.Header())
			if err != nil {
				// TODO
			}

			userId := w.Header().Get("X-User-ID")
			if userId == "" {
				userId = "unknow"
			}

			//Write Log to Clickhouse
			clickLogs.WriteLog(clickhouse.LogRecord{
				Method:          r.Method,
				RequestTime:     time.Now(),
				RequestSize:     uint(r.ContentLength),
				ResponseSize:    uint(lw.BytesWritten()),
				User:            userId,
				Path:            r.RequestURI,
				Latency:         latency,
				ID:              w.Header().Get("X-Request-ID"),
				Status:          uint(lw.Status()),
				Upstream:        w.Header().Get("X-Upstream"),
				UserAgent:       r.UserAgent(),
				Fingerprint:     w.Header().Get("X-User-Fingerprint"),
				RequestHeaders:  string(headersRequest),
				RequestBody:     string(reqBody),
				ResponseHeaders: string(headersResponse),
				ResponseBody:    string(lw.Bytes()),
				GatewayID:       w.Header().Get("X-Gateway-ID"),
			})

			//Write log after
			log.WithFields(log.Fields{
				"Method":       r.Method,
				"Path":         r.RequestURI,
				"Latency":      fmt.Sprintf("%v", latency),
				"Status":       lw.Status(),
				"RequestID":    w.Header().Get("X-Request-ID"),
				"ResponseSize": lw.BytesWritten(),
				"RequestSize":  r.ContentLength,
			}).Info("Request")
		})
	}
}
