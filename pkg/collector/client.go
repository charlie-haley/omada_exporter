package collector

import (
	"fmt"

	"github.com/charlie-haley/omada_exporter/pkg/api"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/rs/zerolog/log"
)

type clientCollector struct {
	omadaClientDownloadActivityBytes *prometheus.Desc
	omadaClientSignalDbm             *prometheus.Desc
	omadaClientConnectedTotal        *prometheus.Desc
	client                           *api.Client
}

func (c *clientCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.omadaClientDownloadActivityBytes
	ch <- c.omadaClientSignalDbm
	ch <- c.omadaClientConnectedTotal
}

func (c *clientCollector) Collect(ch chan<- prometheus.Metric) {
	client := c.client
	config := c.client.Config

	site := config.Site
	clients, err := client.GetClients()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get clients")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.omadaClientConnectedTotal, prometheus.GaugeValue, float64(len(clients)),
		site, client.SiteId)

	for _, item := range clients {
		vlanId := fmt.Sprintf("%.0f", item.VlanId)
		port := fmt.Sprintf("%.0f", item.Port)
		if item.Wireless {
			wifiMode := fmt.Sprintf("%.0f", item.WifiMode)
			ch <- prometheus.MustNewConstMetric(c.omadaClientDownloadActivityBytes, prometheus.GaugeValue, item.Activity,
				item.HostName, item.Vendor, "", "", item.Ip, item.Mac, site, client.SiteId, item.ApName, item.Ssid, wifiMode)

			ch <- prometheus.MustNewConstMetric(c.omadaClientSignalDbm, prometheus.GaugeValue, item.SignalLevel,
				item.HostName, item.Vendor, item.Ip, item.Mac, item.ApName, site, client.SiteId, item.Ssid, wifiMode)
		}
		if !item.Wireless {
			ch <- prometheus.MustNewConstMetric(c.omadaClientDownloadActivityBytes, prometheus.GaugeValue, item.Activity,
				item.HostName, item.Vendor, port, vlanId, item.Ip, item.Mac, site, client.SiteId, "", "", "")
		}
	}
}

func NewClientCollector(c *api.Client) *clientCollector {
	return &clientCollector{
		omadaClientDownloadActivityBytes: prometheus.NewDesc("omada_client_download_activity_bytes",
			"The current download activity for the client in bytes.",
			[]string{"client", "vendor", "switch_port", "vlan_id", "ip", "mac", "site", "site_id", "ap_name", "ssid", "wifi_mode"},
			nil,
		),

		omadaClientSignalDbm: prometheus.NewDesc("omada_client_signal_dbm",
			"The signal level for the wireless client in dBm.",
			[]string{"client", "vendor", "ip", "mac", "ap_name", "site", "site_id", "ssid", "wifi_mode"},
			nil,
		),

		omadaClientConnectedTotal: prometheus.NewDesc("omada_client_connected_total",
			"Total number of connected clients.",
			[]string{"site", "site_id"},
			nil,
		),
		client: c,
	}
}
