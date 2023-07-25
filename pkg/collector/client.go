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
	omadaClientRssiDbm				 *prometheus.Desc
	omadaClientTrafficDown			 *prometheus.Desc
	omadaClientTrafficUp			 *prometheus.Desc
	omadaClientTxRate				 *prometheus.Desc
	omadaClientRxRate				 *prometheus.Desc

	omadaClientConnectedTotal        *prometheus.Desc
	omadaClientWirelessTotal         *prometheus.Desc
	omadaClientWireless2g            *prometheus.Desc
	omadaClientWireless5g            *prometheus.Desc
	omadaClientWireless6g            *prometheus.Desc
	omadaClientWiredTotal         	 *prometheus.Desc
	client                           *api.Client
}

func (c *clientCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.omadaClientDownloadActivityBytes
	ch <- c.omadaClientSignalDbm
	ch <- c.omadaClientRssiDbm		
	ch <- c.omadaClientTrafficDown	
	ch <- c.omadaClientTrafficUp	
	ch <- c.omadaClientTxRate		
	ch <- c.omadaClientRxRate	
	ch <- c.omadaClientConnectedTotal
	ch <- c.omadaClientWirelessTotal
	ch <- c.omadaClientWireless2g
	ch <- c.omadaClientWireless5g
	ch <- c.omadaClientWireless6g
	ch <- c.omadaClientWiredTotal
}

func (c *clientCollector) Collect(ch chan<- prometheus.Metric) {
	client := c.client
	config := c.client.Config

	site := config.Site
	clients, stats, err := client.GetClients()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get clients")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.omadaClientConnectedTotal, prometheus.GaugeValue, stats.Total,
		site, client.SiteId)
	ch <- prometheus.MustNewConstMetric(c.omadaClientWirelessTotal, prometheus.GaugeValue, stats.Wireless,
		site, client.SiteId)
	ch <- prometheus.MustNewConstMetric(c.omadaClientWireless2g, prometheus.GaugeValue, stats.Num2g,
		site, client.SiteId)
	ch <- prometheus.MustNewConstMetric(c.omadaClientWireless5g, prometheus.GaugeValue, stats.Num5g,
		site, client.SiteId)
	ch <- prometheus.MustNewConstMetric(c.omadaClientWireless6g, prometheus.GaugeValue, stats.Num6g,
		site, client.SiteId)
	ch <- prometheus.MustNewConstMetric(c.omadaClientWiredTotal, prometheus.GaugeValue, stats.Wired,
		site, client.SiteId)

	for _, item := range clients {
		vlanId := fmt.Sprintf("%.0f", item.VlanId)
		port := fmt.Sprintf("%.0f", item.Port)
		if item.Wireless {
			wifiMode := fmt.Sprintf("%.0f", item.WifiMode)
			ch <- prometheus.MustNewConstMetric(c.omadaClientDownloadActivityBytes, prometheus.GaugeValue, item.Activity,
				item.HostName, item.Vendor, "", "", item.Ip, item.Mac, site, client.SiteId, item.ApName, item.Ssid, wifiMode)

			ch <- prometheus.MustNewConstMetric(c.omadaClientSignalDbm, prometheus.GaugeValue, -item.SignalLevel,
				item.HostName, item.Vendor, item.Ip, item.Mac, item.ApName, site, client.SiteId, item.Ssid, wifiMode)

			ch <- prometheus.MustNewConstMetric(c.omadaClientRssiDbm, prometheus.GaugeValue, item.Rssi,
				item.HostName, item.Vendor, item.Ip, item.Mac, item.ApName, site, client.SiteId, item.Ssid, wifiMode)
			
			ch <- prometheus.MustNewConstMetric(c.omadaClientTrafficDown, prometheus.CounterValue, item.TrafficDown,
				item.HostName, item.Vendor, item.Ip, item.Mac, item.ApName, site, client.SiteId, item.Ssid, wifiMode)
			
			ch <- prometheus.MustNewConstMetric(c.omadaClientTrafficUp, prometheus.CounterValue, item.TrafficUp,
				item.HostName, item.Vendor, item.Ip, item.Mac, item.ApName, site, client.SiteId, item.Ssid, wifiMode)
			
			ch <- prometheus.MustNewConstMetric(c.omadaClientTxRate, prometheus.GaugeValue, item.TxRate,
				item.HostName, item.Vendor, item.Ip, item.Mac, item.ApName, site, client.SiteId, item.Ssid, wifiMode)

			ch <- prometheus.MustNewConstMetric(c.omadaClientRxRate, prometheus.GaugeValue, item.RxRate,
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
			"The noise level for the wireless client in dBm.",
			[]string{"client", "vendor", "ip", "mac", "ap_name", "site", "site_id", "ssid", "wifi_mode"},
			nil,
		),

		omadaClientRssiDbm: prometheus.NewDesc("omada_client_rssi_dbm",
			"The RSSI for the wireless client in dBm.",
			[]string{"client", "vendor", "ip", "mac", "ap_name", "site", "site_id", "ssid", "wifi_mode"},
			nil,
		),

		omadaClientTrafficDown: prometheus.NewDesc("omada_client_traffic_down",
			"Total bytes received by wireless client.",
			[]string{"client", "vendor", "ip", "mac", "ap_name", "site", "site_id", "ssid", "wifi_mode"},
			nil,
		),

		omadaClientTrafficUp: prometheus.NewDesc("omada_client_traffic_up",
			"Total bytes sent by wireless client.",
			[]string{"client", "vendor", "ip", "mac", "ap_name", "site", "site_id", "ssid", "wifi_mode"},
			nil,
		),

		omadaClientTxRate: prometheus.NewDesc("omada_client_tx_rate",
			"TX rate of wireless client.",
			[]string{"client", "vendor", "ip", "mac", "ap_name", "site", "site_id", "ssid", "wifi_mode"},
			nil,
		),

		omadaClientRxRate: prometheus.NewDesc("omada_client_rx_rate",
			"RX rate of wireless client.",
			[]string{"client", "vendor", "ip", "mac", "ap_name", "site", "site_id", "ssid", "wifi_mode"},
			nil,
		),

		omadaClientConnectedTotal: prometheus.NewDesc("omada_client_connected_total",
			"Total number of connected clients.",
			[]string{"site", "site_id"},
			nil,
		),

		omadaClientWirelessTotal: prometheus.NewDesc("omada_client_wireless_total",
			"Total number of connected wireless clients.",
			[]string{"site", "site_id"},
			nil,
		),

		omadaClientWireless2g: prometheus.NewDesc("omada_client_wireless_2g",
			"Total number of connected wireless 2.4GHz clients.",
			[]string{"site", "site_id"},
			nil,
		),

		omadaClientWireless5g: prometheus.NewDesc("omada_client_wireless_5g",
			"Total number of connected wireless 5GHz clients.",
			[]string{"site", "site_id"},
			nil,
		),

		omadaClientWireless6g: prometheus.NewDesc("omada_client_wireless_6g",
			"Total number of connected wireless 6GHz clients.",
			[]string{"site", "site_id"},
			nil,
		),
		
		omadaClientWiredTotal: prometheus.NewDesc("omada_client_wired_total",
			"Total number of connected wired clients.",
			[]string{"site", "site_id"},
			nil,
		),

		client: c,
	}
}
