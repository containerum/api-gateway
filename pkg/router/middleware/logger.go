package middleware

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	clickhouse "git.containerum.net/ch/api-gateway/pkg/utils/clickhouselog"

	log "github.com/Sirupsen/logrus"
	"github.com/cactus/go-statsd-client/statsd"

	b64 "encoding/base64"
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

			userID := w.Header().Get("X-User-ID")
			if userID == "" {
				userID = "unknow"
			}

			//TODO ADD secret to env
			//Write Log to Clickhouse
			logKeyStr := fmt.Sprintf("%s+%s", w.Header().Get("X-Request-ID"), "Secret")
			logKey := sha256.Sum256([]byte(logKeyStr))
			logRecord := clickhouse.LogRecord{
				Method:          r.Method,
				RequestTime:     time.Now(),
				RequestSize:     uint(r.ContentLength),
				ResponseSize:    uint(lw.BytesWritten()),
				User:            userID,
				Path:            r.RequestURI,
				Latency:         latency,
				ID:              w.Header().Get("X-Request-ID"),
				Status:          uint(lw.Status()),
				Upstream:        w.Header().Get("X-Upstream"),
				UserAgent:       r.UserAgent(),
				Fingerprint:     w.Header().Get("X-User-Fingerprint"),
				RequestHeaders:  makeBase64Headers(r.Header, logKey[:]),
				RequestBody:     makebase64Body(reqBody, logKey[:]),
				ResponseHeaders: makeBase64Headers(lw.Header(), logKey[:]),
				ResponseBody:    makebase64Body(lw.Bytes(), logKey[:]),
				GatewayID:       w.Header().Get("X-Gateway-ID"),
			}
			clickLogs.WriteLog(logRecord)

			//Write log after
			if log.GetLevel() == log.InfoLevel {
				log.WithFields(log.Fields{
					"Method":       r.Method,
					"Path":         r.RequestURI,
					"Latency":      fmt.Sprintf("%v", latency),
					"Status":       lw.Status(),
					"RequestID":    w.Header().Get("X-Request-ID"),
					"ResponseSize": lw.BytesWritten(),
					"RequestSize":  r.ContentLength,
				}).Info("Request")
			} else {
				log.WithFields(log.Fields{
					"Method":       r.Method,
					"Path":         r.RequestURI,
					"Latency":      fmt.Sprintf("%v", latency),
					"Status":       lw.Status(),
					"RequestID":    w.Header().Get("X-Request-ID"),
					"ResponseSize": lw.BytesWritten(),
					"RequestSize":  r.ContentLength,
					"User":         userID,
					"Upstream":     w.Header().Get("X-Upstream"),
					"UserAgent":    r.UserAgent(),
					"Fingerprint":  w.Header().Get("X-User-Fingerprint"),
					// "RequestHeaders":  string(r.Header),
					"RequestBody": string(reqBody),
					// "ResponseHeaders": string(lw.Header()),
					"ResponseBody": string(lw.Bytes()),
					"GatewayID":    w.Header().Get("X-Gateway-ID"),
				}).Info("Request")
			}
		})
	}
}

func makeBase64Headers(headers http.Header, key []byte) string {
	headers64, err := json.Marshal(headers)
	if err != nil {
		return err.Error()
	}
	return makebase64Body(headers64, key)
}

func makebase64Body(body []byte, key []byte) string {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	cipherdata := make([]byte, aes.BlockSize+len(body))
	iv := cipherdata[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherdata[aes.BlockSize:], body)

	// convert to base64
	return b64.URLEncoding.EncodeToString(cipherdata)
}
