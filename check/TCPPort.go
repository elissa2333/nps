package check

import (
	"net"
	"time"
)

func TCPPort(ip, port string, in2 int) bool {
	tou := time.Duration(in2)
	address := ip + ":" + port
	conn, err := net.DialTimeout("tcp", address, tou*time.Second)
	if err != nil {
		return false
	} else {
		conn.Close()
		return true
	}
}
