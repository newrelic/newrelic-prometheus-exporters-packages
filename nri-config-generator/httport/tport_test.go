package httport

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_findAvailablePort(t *testing.T) {
	port, err := findAvailablePort()
	assert.Nil(t, err)
	assert.NotEmpty(t, port)

}

func launchServerOn(port string) *http.Server {
	server := &http.Server{Addr: fmt.Sprintf(":%s", port)}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Warn(err.Error())
		}
	}()
	time.Sleep(2*time.Second)
	return server
}

func Test_isPortAvailable(t *testing.T) {
	port, err := findAvailablePort()
	assert.Nil(t, err)
	portStr := fmt.Sprintf("%v", port)
	server := launchServerOn(portStr)
	assert.False(t, isPortAvailable(portStr))
	assert.Nil(t, server.Shutdown(context.Background()))
	assert.True(t, isPortAvailable(portStr))
}

func Test_isConnectionRefusedError(t *testing.T) {
	assert.False(t, isConnectionRefusedError(errors.New("test error")))
	assert.True(t, isConnectionRefusedError(&net.OpError{Op: "read", Err: fmt.Errorf("test error")}))
	assert.False(t, isConnectionRefusedError(&net.OpError{Op: "dial", Err: fmt.Errorf("test error")}))
	assert.False(t, isConnectionRefusedError(&net.OpError{Op: "write", Err: fmt.Errorf("test error")}))
}

func Test_GetAvailablePort(t *testing.T) {
	availablePort, err := findAvailablePort()
	assert.Nil(t, err)
	portStr := fmt.Sprintf("%v", availablePort)
	port, err := GetAvailablePort(portStr)
	assert.Nil(t, err)
	assert.Equal(t, availablePort, port)
	port, err = GetAvailablePort("9999")
	assert.Nil(t, err)
	assert.Equal(t, availablePort, port)
	server := launchServerOn(portStr)
	port, err = GetAvailablePort("9999")
	assert.Nil(t, err)
	assert.NotEqual(t, availablePort, port)
	assert.Nil(t, server.Shutdown(context.Background()))

}

func Test_isPortAvailable2(t *testing.T) {
	port, err := findAvailablePort()
	assert.Nil(t, err)
	port2, err := findAvailablePort()
	assert.NotEqual(t, port, port2)
	port3, err := findAvailablePort()
	assert.NotEqual(t, port, port3)
}
