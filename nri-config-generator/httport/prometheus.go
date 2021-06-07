package httport

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/newrelic/infra-integrations-sdk/v4/log"
)

const (
	healthCheckEndpoint = "http://%s:%v/metrics"
)

var (
	httpPrometheusExporterClient *http.Client
	exporterURL                  string
)

func SetPrometheusExporterPort(hostname string, port int) {
	exporterURL = fmt.Sprintf(healthCheckEndpoint, hostname, port)
	tr := &http.Transport{
		MaxIdleConnsPerHost: 1024,
		TLSHandshakeTimeout: 0 * time.Second,
	}
	httpPrometheusExporterClient = &http.Client{Transport: tr}
}

func IsPrometheusExporterRunning() bool {
	log.Debug("check if prometheus exporter is running on %s", exporterURL)
	resp, err := httpPrometheusExporterClient.Get(exporterURL)
	if err != nil {
		log.Error("error while checking the prometheus exporter 'health check': %s",err.Error())
		return false
	}
	io.Copy(ioutil.Discard, resp.Body) // <= NOTE
	defer resp.Body.Close()
	return true
}
