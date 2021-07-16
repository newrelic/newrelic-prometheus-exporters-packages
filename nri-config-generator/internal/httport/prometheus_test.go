package httport

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_IsPrometheusExporterRunning(t *testing.T) {
	port, err := findAvailablePort()
	assert.Nil(t, err)
	SetPrometheusExporterPort("localhost", port)
	assert.False(t, IsPrometheusExporterRunning())
	assert.False(t, IsPrometheusExporterRunning())
	server := &http.Server{Addr: fmt.Sprintf(":%v", port)}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Warn(err.Error())
		}
	}()
	assert.True(t, IsPrometheusExporterRunning())
	defer func() {
		if err := server.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()
}
