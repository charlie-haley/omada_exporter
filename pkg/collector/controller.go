package collector

import (
	"github.com/charlie-haley/omada_exporter/pkg/api"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/rs/zerolog/log"
)

type controllerCollector struct {
	omadaControllerUptimeSeconds         *prometheus.Desc
	omadaControllerStorageUsedBytes      *prometheus.Desc
	omadaControllerStorageAvailableBytes *prometheus.Desc
	client                               *api.Client
}

func (c *controllerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.omadaControllerUptimeSeconds
	ch <- c.omadaControllerStorageUsedBytes
	ch <- c.omadaControllerStorageAvailableBytes
}

func (c *controllerCollector) Collect(ch chan<- prometheus.Metric) {
	client := c.client
	config := c.client.Config

	site := config.Site
	controller, err := client.GetController()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get controller")
	}

	ch <- prometheus.MustNewConstMetric(c.omadaControllerUptimeSeconds, prometheus.GaugeValue, controller.Uptime/1000,
		controller.Name, controller.Model, controller.ControllerVersion, controller.ControllerVersion, controller.MacAddress, site, client.SiteId)

	for _, s := range controller.Storage {
		ch <- prometheus.MustNewConstMetric(c.omadaControllerStorageUsedBytes, prometheus.GaugeValue, s.Used*1000000000,
			s.Name, controller.Name, controller.Model, controller.ControllerVersion, controller.ControllerVersion, controller.MacAddress, site, client.SiteId)

		ch <- prometheus.MustNewConstMetric(c.omadaControllerStorageAvailableBytes, prometheus.GaugeValue, s.Total*100000000,
			s.Name, controller.Name, controller.Model, controller.ControllerVersion, controller.ControllerVersion, controller.MacAddress, site, client.SiteId)
	}

}

func NewControllerCollector(c *api.Client) *controllerCollector {
	return &controllerCollector{
		omadaControllerUptimeSeconds: prometheus.NewDesc("omada_controller_uptime_seconds",
			"Uptime of the controller.",
			[]string{"controller_name", "model", "controller_version", "firmware_version", "mac", "site", "site_id"},
			nil,
		),
		omadaControllerStorageUsedBytes: prometheus.NewDesc("omada_controller_storage_used_bytes",
			"Storage used on the controller.",
			[]string{"storage_name", "controller_name", "model", "controller_version", "firmware_version", "mac", "site", "site_id"},
			nil,
		),
		omadaControllerStorageAvailableBytes: prometheus.NewDesc("omada_controller_storage_available_bytes",
			"Total storage available for the controller.",
			[]string{"storage_name", "controller_name", "model", "controller_version", "firmware_version", "mac", "site", "site_id"},
			nil,
		),
		client: c,
	}
}
