package httport

import (
	"net"
	"os"
	"strconv"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

const tcp = "tcp"

func GetAvailablePort(configPort string) (int, error) {
	if configPort != "" {
		port, err := strconv.Atoi(configPort)
		if err == nil && isPortAvailable(configPort) {
			return port, nil
		}
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
	defer func() {
		if err := l.Close(); err != nil {
			log.Warn(err)
		}
	}()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func isPortAvailable(port string) bool {
	conn, err := net.DialTimeout(tcp, net.JoinHostPort("", port), 1*time.Second)
	if conn != nil {
		if err := conn.Close(); err != nil {
			log.Warn(err)
		}
		return false
	}
	return !isConnectionRefusedError(err)
}

func isConnectionRefusedError(err error) bool {
	switch t := err.(type) {
	case *net.OpError:
		return t.Op == "read"
	case *os.SyscallError:
		if (*t).Err == syscall.ECONNREFUSED {
			return true
		}
	}
	return false
}
