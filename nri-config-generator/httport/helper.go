package httport

import (
	"net"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/newrelic/infra-integrations-sdk/v4/log"
)

const tcp = "tcp"

func GetAvailablePort(defPort string) (int, error) {
	port, err := strconv.Atoi(defPort)
	if err == nil && isPortAvailable(defPort) {
		return port, nil
	}
	if err != nil {
		log.Warn(err.Error())
	}
	return findAvailablePort()
}

func findAvailablePort() (int, error) {
	addr, err := net.ResolveTCPAddr(tcp, "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP(tcp, addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func isPortAvailable(port string) bool {
	conn, err := net.DialTimeout(tcp, net.JoinHostPort("", port), 1*time.Second)
	if err != nil {
		if isConnectionRefusedError(err) {
			return true
		}
	}
	if conn != nil {
		conn.Close()
		return false
	}
	return true
}

func isConnectionRefusedError(err error) bool {
	if netError, ok := err.(net.Error); ok && netError.Timeout() {
		return false
	}
	switch t := err.(type) {
	case *net.OpError:
		if t.Op == "dial" {
			return false
		} else if t.Op == "read" {
			return true
		}
	case *os.SyscallError:
		if (*t).Err == syscall.ECONNREFUSED {
			return true
		}
	}
	return false
}
