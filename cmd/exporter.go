package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/charlie-haley/omada_exporter/pkg/omada"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func Run() {
	go omada.UpdateStatus()
	exporterPort := os.Getenv("OMADA_EXPORTER_PORT")
	if exporterPort == "" {
		exporterPort = "9202"
	}

	log.Info(fmt.Sprintf("Listening on :%s", exporterPort))
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%s", exporterPort), nil)
}
