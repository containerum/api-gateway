package clickhouselog

import (
	"net"
	"time"
)

type LogClient struct {
	*net.UDPConn
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
}

func OpenConenction(addr string) (*LogClient, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}
	con, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}
	return &LogClient{UDPConn: con}, nil
}

func (lc *LogClient) WriteLog(lr LogRecord) error {

	return nil
}
