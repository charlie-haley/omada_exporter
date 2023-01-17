package collector

import (
	"fmt"

	"github.com/charlie-haley/omada_exporter/pkg/api"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/rs/zerolog/log"
)

type portCollector struct {
	omadaPortPowerWatts    *prometheus.Desc
	omadaPortLinkStatus    *prometheus.Desc
	omadaPortLinkSpeedMbps *prometheus.Desc
	client                 *api.Client
}

func (c *portCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.omadaPortPowerWatts
	ch <- c.omadaPortLinkStatus
	ch <- c.omadaPortLinkSpeedMbps
}

func (c *portCollector) Collect(ch chan<- prometheus.Metric) {
	client := c.client
	config := c.client.Config

	site := config.Site
	devices, err := client.GetDevices()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get devices")
		return
	}

	for _, device := range devices {
		for _, p := range device.Ports {
			linkSpeed := float64(0)
			if p.PortStatus.LinkSpeed == 0 {
				linkSpeed = 0
			}
			if p.PortStatus.LinkSpeed == 1 {
				linkSpeed = 10
			}
			if p.PortStatus.LinkSpeed == 2 {
				linkSpeed = 100
			}
			if p.PortStatus.LinkSpeed == 3 {
				linkSpeed = 1000
			}

			portClient, err := client.GetClientByPort(device.Mac, p.Port)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get client by port")
			}

			port := fmt.Sprintf("%.0f", p.Port)
			if portClient != nil {
				vlanId := fmt.Sprintf("%.0f", portClient.VlanId)

				ch <- prometheus.MustNewConstMetric(c.omadaPortPowerWatts, prometheus.GaugeValue, p.PortStatus.PoePower,
					portClient.HostName, portClient.Vendor, port, p.SwitchMac, p.SwitchId, vlanId, p.ProfileName, site, client.SiteId)

				ch <- prometheus.MustNewConstMetric(c.omadaPortLinkStatus, prometheus.GaugeValue, p.PortStatus.LinkStatus,
					portClient.HostName, portClient.Vendor, port, p.SwitchMac, p.SwitchId, vlanId, p.ProfileName, site, client.SiteId)

				ch <- prometheus.MustNewConstMetric(c.omadaPortLinkSpeedMbps, prometheus.GaugeValue, linkSpeed,
					portClient.HostName, portClient.Vendor, port, p.SwitchMac, p.SwitchId, vlanId, p.ProfileName, site, client.SiteId)
			} else {

				ch <- prometheus.MustNewConstMetric(c.omadaPortPowerWatts, prometheus.GaugeValue, p.PortStatus.PoePower,
					"", "", port, p.SwitchMac, p.SwitchId, "", p.ProfileName, site, client.SiteId)

				ch <- prometheus.MustNewConstMetric(c.omadaPortLinkStatus, prometheus.GaugeValue, p.PortStatus.LinkStatus,
					"", "", port, p.SwitchMac, p.SwitchId, "", p.ProfileName, site, client.SiteId)

				ch <- prometheus.MustNewConstMetric(c.omadaPortLinkSpeedMbps, prometheus.GaugeValue, linkSpeed,
					"", "", port, p.SwitchMac, p.SwitchId, "", p.ProfileName, site, client.SiteId)
			}
		}
	}
}

func NewPortCollector(c *api.Client) *portCollector {
	return &portCollector{
		omadaPortPowerWatts: prometheus.NewDesc("omada_port_power_watts",
			"The current PoE usage of the port in watts.",
			[]string{"client", "vendor", "switch_port", "switch_mac", "switch_id", "vlan_id", "profile", "site", "site_id"},
			nil,
		),
		omadaPortLinkStatus: prometheus.NewDesc("omada_port_link_status",
			"A boolean representing the link status of the port.",
			[]string{"client", "vendor", "switch_port", "switch_mac", "switch_id", "vlan_id", "profile", "site", "site_id"},
			nil,
		),
		omadaPortLinkSpeedMbps: prometheus.NewDesc("omada_port_link_speed_mbps",
			"Port link speed in mbps. This is the capability of the connection, not the active throughput.",
			[]string{"client", "vendor", "switch_port", "switch_mac", "switch_id", "vlan_id", "profile", "site", "site_id"},
			nil,
		),
		client: c,
	}
}
