package check

import (
	"net"
	"strconv"
	"time"
)

func TCPPort(host string, port uint16, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(int(port))), timeout)
	if err != nil {
		return false
	} else {
		conn.Close()
		return true
	}
}
