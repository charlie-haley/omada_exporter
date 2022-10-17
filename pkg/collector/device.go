package collector

import (
	"github.com/charlie-haley/omada_exporter/pkg/api"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/rs/zerolog/log"
)

type deviceCollector struct {
	omadaDeviceUptimeSeconds  *prometheus.Desc
	omadaDeviceCpuPercentage  *prometheus.Desc
	omadaDeviceMemPercentage  *prometheus.Desc
	omadaDeviceNeedUpgrade    *prometheus.Desc
	omadaDeviceTxRate         *prometheus.Desc
	omadaDeviceRxRate         *prometheus.Desc
	omadaDevicePoeRemainWatts *prometheus.Desc
	client                    *api.Client
}

func (c *deviceCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.omadaDeviceUptimeSeconds
	ch <- c.omadaDeviceUptimeSeconds
	ch <- c.omadaDeviceCpuPercentage
	ch <- c.omadaDeviceMemPercentage
	ch <- c.omadaDeviceNeedUpgrade
	ch <- c.omadaDeviceTxRate
	ch <- c.omadaDeviceRxRate
	ch <- c.omadaDevicePoeRemainWatts
}

func (c *deviceCollector) Collect(ch chan<- prometheus.Metric) {
	client := c.client
	config := c.client.Config

	site := config.Site
	devices, err := client.GetDevices()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get devices")
	}

	for _, item := range devices {
		needUpgrade := float64(0)
		if item.NeedUpgrade {
			needUpgrade = 1
		}
		ch <- prometheus.MustNewConstMetric(c.omadaDeviceUptimeSeconds, prometheus.GaugeValue, item.Uptime,
			item.Name, item.Model, item.Version, item.Ip, item.Mac, site, client.SiteId, item.Type)

		ch <- prometheus.MustNewConstMetric(c.omadaDeviceCpuPercentage, prometheus.GaugeValue, item.CpuUtil,
			item.Name, item.Model, item.Version, item.Ip, item.Mac, site, client.SiteId, item.Type)

		ch <- prometheus.MustNewConstMetric(c.omadaDeviceMemPercentage, prometheus.GaugeValue, item.MemUtil,
			item.Name, item.Model, item.Version, item.Ip, item.Mac, site, client.SiteId, item.Type)

		ch <- prometheus.MustNewConstMetric(c.omadaDeviceNeedUpgrade, prometheus.GaugeValue, needUpgrade,
			item.Name, item.Model, item.Version, item.Ip, item.Mac, site, client.SiteId, item.Type)

		if item.Type == "ap" {
			ch <- prometheus.MustNewConstMetric(c.omadaDeviceTxRate, prometheus.GaugeValue, item.TxRate,
				item.Name, item.Model, item.Version, item.Ip, item.Mac, site, client.SiteId, item.Type)

			ch <- prometheus.MustNewConstMetric(c.omadaDeviceRxRate, prometheus.GaugeValue, item.RxRate,
				item.Name, item.Model, item.Version, item.Ip, item.Mac, site, client.SiteId, item.Type)
		}
		if item.Type == "switch" {
			ch <- prometheus.MustNewConstMetric(c.omadaDevicePoeRemainWatts, prometheus.GaugeValue, item.PoeRemain,
				item.Name, item.Model, item.Version, item.Ip, item.Mac, site, client.SiteId, item.Type)
		}
	}
}

func NewDeviceCollector(c *api.Client) *deviceCollector {
	return &deviceCollector{
		omadaDeviceUptimeSeconds: prometheus.NewDesc("omada_device_uptime_seconds",
			"Uptime of the device.",
			[]string{"device", "model", "version", "ip", "mac", "site", "site_id", "device_type"},
			nil,
		),
		omadaDeviceCpuPercentage: prometheus.NewDesc("omada_device_cpu_percentage",
			"Percentage of device CPU used.",
			[]string{"device", "model", "version", "ip", "mac", "site", "site_id", "device_type"},
			nil,
		),
		omadaDeviceMemPercentage: prometheus.NewDesc("omada_device_mem_percentage",
			"Percentage of device Memory used.",
			[]string{"device", "model", "version", "ip", "mac", "site", "site_id", "device_type"},
			nil,
		),
		omadaDeviceNeedUpgrade: prometheus.NewDesc("omada_device_need_upgrade",
			"A boolean on whether the device needs an upgrade.",
			[]string{"device", "model", "version", "ip", "mac", "site", "site_id", "device_type"},
			nil,
		),
		omadaDeviceTxRate: prometheus.NewDesc("omada_device_tx_rate",
			"The tx rate of the device.",
			[]string{"device", "model", "version", "ip", "mac", "site", "site_id", "device_type"},
			nil,
		),
		omadaDeviceRxRate: prometheus.NewDesc("omada_device_rx_rate",
			"The rx rate of the device.",
			[]string{"device", "model", "version", "ip", "mac", "site", "site_id", "device_type"},
			nil,
		),
		omadaDevicePoeRemainWatts: prometheus.NewDesc("omada_device_poe_remain_watts",
			"The remaining amount of PoE power for the device in watts.",
			[]string{"device", "model", "version", "ip", "mac", "site", "site_id", "device_type"},
			nil,
		),
		client: c,
	}
}
