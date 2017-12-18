package clickhouselog

import (
	"encoding/json"
	"net"
	"time"
)

type LogClient struct {
	con *net.UDPConn
}

type LogRecord struct {
	Method       string        `json:"method"`
	RequestTime  time.Time     `json:"request_time"`
	RequestSize  uint          `json:"request_size"`
	ResponseSize uint          `json:"response_size"`
	User         string        `json:"user"`
	Path         string        `json:"path"`
	Latency      time.Duration `json:"latency"`
	ID           string        `json:"id"`

	Status          uint   `json:"status"`
	Upstream        string `json:"upstream"`
	UserAgent       string `json:"user_agent"`
	Fingerprint     string `json:"fingerprint"`
	RequestHeaders  string `json:"request_headers"`
	RequestBody     string `json:"request_body"`
	ResponseHeaders string `json:"response_headers"`
	ResponseBody    string `json:"response_body"`
	GatewayID       string `json:"gateway_id"`
}

func OpenConenction(addr string) (*LogClient, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	con, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}

	return &LogClient{con: con}, nil
}

func (lc *LogClient) WriteLog(lr LogRecord) error {
	js, err := json.Marshal(lr)
	if err != nil {
		return err
	}
	_, err = lc.con.Write(js)
	return err
}
