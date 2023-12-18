package collector

import (
	"fmt"

	"github.com/charlie-haley/omada_exporter/pkg/api"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/rs/zerolog/log"
)

type clientCollector struct {
	omadaClientDownloadActivityBytes *prometheus.Desc
	omadaClientSignalPct             *prometheus.Desc
	omadaClientRssiDbm               *prometheus.Desc
	omadaClientTrafficDown           *prometheus.Desc
	omadaClientTrafficUp             *prometheus.Desc
	omadaClientTxRate                *prometheus.Desc
	omadaClientRxRate                *prometheus.Desc
	omadaClientConnectedTotal        *prometheus.Desc
	client                           *api.Client
}

func (c *clientCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.omadaClientDownloadActivityBytes
	ch <- c.omadaClientSignalPct
	ch <- c.omadaClientRssiDbm
	ch <- c.omadaClientTrafficDown
	ch <- c.omadaClientTrafficUp
	ch <- c.omadaClientTxRate
	ch <- c.omadaClientRxRate
	ch <- c.omadaClientConnectedTotal
}

func FormatWifiMode(wifiMode int) string {
	mapping := map[int]string{
		0: "802.11a",
		1: "802.11b",
		2: "802.11g",
		3: "802.11na",
		4: "802.11ng",
		5: "802.11ac",
		6: "802.11axa",
		7: "802.11axg",
	}
	formatted, ok := mapping[wifiMode]
	if !ok {
		return ""
	}
	return formatted
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

	totals := map[string]int{}

	for _, item := range clients {
		vlanId := fmt.Sprintf("%.0f", item.VlanId)
		port := fmt.Sprintf("%.0f", item.Port)
		if item.Wireless {
			wifiMode := FormatWifiMode(int(item.WifiMode))

			CollectWirelessMetrics := func(desc *prometheus.Desc, valueType prometheus.ValueType, value float64) {
				ch <- prometheus.MustNewConstMetric(desc, valueType, value,
					item.HostName, item.Vendor, item.Ip, item.Mac, site, client.SiteId, "wireless", wifiMode, item.ApName, item.Ssid)
			}
			CollectWirelessMetrics(c.omadaClientSignalPct, prometheus.GaugeValue, item.SignalLevel)
			CollectWirelessMetrics(c.omadaClientRssiDbm, prometheus.GaugeValue, item.Rssi)
			CollectWirelessMetrics(c.omadaClientTrafficDown, prometheus.CounterValue, item.TrafficDown)
			CollectWirelessMetrics(c.omadaClientTrafficUp, prometheus.CounterValue, item.TrafficUp)
			CollectWirelessMetrics(c.omadaClientTxRate, prometheus.GaugeValue, item.TxRate)
			CollectWirelessMetrics(c.omadaClientRxRate, prometheus.GaugeValue, item.RxRate)

			totals[wifiMode] += 1
			ch <- prometheus.MustNewConstMetric(c.omadaClientDownloadActivityBytes, prometheus.GaugeValue, item.Activity,
				item.HostName, item.Vendor, item.Ip, item.Mac, site, client.SiteId, "wireless", wifiMode, item.ApName, item.Ssid, "", "")
		}
		if !item.Wireless {
			totals["wired"] += 1
			ch <- prometheus.MustNewConstMetric(c.omadaClientDownloadActivityBytes, prometheus.GaugeValue, item.Activity,
				item.HostName, item.Vendor, item.Ip, item.Mac, site, client.SiteId, "wired", "", "", "", port, vlanId)
		}
	}

	for connectionModeFmt, v := range totals {
		if connectionModeFmt == "wired" {
			ch <- prometheus.MustNewConstMetric(c.omadaClientConnectedTotal, prometheus.GaugeValue, float64(v),
				site, client.SiteId, "wired", "")
		} else {
			ch <- prometheus.MustNewConstMetric(c.omadaClientConnectedTotal, prometheus.GaugeValue, float64(v),
				site, client.SiteId, "wireless", connectionModeFmt)
		}
	}
}

func NewClientCollector(c *api.Client) *clientCollector {
	client_labels := []string{"client", "vendor", "ip", "mac", "site", "site_id", "connection_mode", "wifi_mode", "ap_name", "ssid"}
	wired_client_labels := append(client_labels, "switch_port", "vlan_id")

	return &clientCollector{
		omadaClientDownloadActivityBytes: prometheus.NewDesc("omada_client_download_activity_bytes",
			"The current download activity for the client in bytes.",
			wired_client_labels,
			nil,
		),

		omadaClientSignalPct: prometheus.NewDesc("omada_client_signal_pct",
			"The signal quality for the wireless client in percent.",
			client_labels,
			nil,
		),

		omadaClientRssiDbm: prometheus.NewDesc("omada_client_rssi_dbm",
			"The RSSI for the wireless client in dBm.",
			client_labels,
			nil,
		),

		omadaClientTrafficDown: prometheus.NewDesc("omada_client_traffic_down_bytes",
			"Total bytes received by wireless client.",
			client_labels,
			nil,
		),

		omadaClientTrafficUp: prometheus.NewDesc("omada_client_traffic_up_bytes",
			"Total bytes sent by wireless client.",
			client_labels,
			nil,
		),

		omadaClientTxRate: prometheus.NewDesc("omada_client_tx_rate",
			"TX rate of wireless client.",
			client_labels,
			nil,
		),

		omadaClientRxRate: prometheus.NewDesc("omada_client_rx_rate",
			"RX rate of wireless client.",
			client_labels,
			nil,
		),

		omadaClientConnectedTotal: prometheus.NewDesc("omada_client_connected_total",
			"Total number of connected clients.",
			[]string{"site", "site_id", "connection_mode", "wifi_mode"},
			nil,
		),

		client: c,
	}
}
