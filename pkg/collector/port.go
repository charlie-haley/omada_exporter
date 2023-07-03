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
		// The Omada exporter sometimes returns duplicate ports. e.g an 8 port switch will return 16 ports with identical ports
		// this causes issues with Prometheus as it tries to register duplicate metrics. A bit of hacky fix, but here we remove
		// duplicate ports to prevent this error.
		ports := removeDuplicates(device.Ports)
		for _, p := range ports {
			linkSpeed := getPortByLinkSpeed(p.PortStatus.LinkSpeed)

			portClient, err := client.GetClientByPort(device.Mac, p.Port)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get client by port")
			}

			port := fmt.Sprintf("%.0f", p.Port)
			if portClient != nil {
				vlanId := fmt.Sprintf("%.0f", portClient.VlanId)

				ch <- prometheus.MustNewConstMetric(c.omadaPortPowerWatts, prometheus.GaugeValue, p.PortStatus.PoePower,
					device.Name, device.Mac, portClient.HostName, portClient.Vendor, port, p.SwitchMac, p.SwitchId, vlanId, p.ProfileName, site, client.SiteId)

				ch <- prometheus.MustNewConstMetric(c.omadaPortLinkStatus, prometheus.GaugeValue, p.PortStatus.LinkStatus,
					device.Name, device.Mac, portClient.HostName, portClient.Vendor, port, p.SwitchMac, p.SwitchId, vlanId, p.ProfileName, site, client.SiteId)

				ch <- prometheus.MustNewConstMetric(c.omadaPortLinkSpeedMbps, prometheus.GaugeValue, linkSpeed,
					device.Name, device.Mac, portClient.HostName, portClient.Vendor, port, p.SwitchMac, p.SwitchId, vlanId, p.ProfileName, site, client.SiteId)
			} else {
				ch <- prometheus.MustNewConstMetric(c.omadaPortPowerWatts, prometheus.GaugeValue, p.PortStatus.PoePower,
					device.Name, device.Mac, "", "", port, p.SwitchMac, p.SwitchId, "", p.ProfileName, site, client.SiteId)

				ch <- prometheus.MustNewConstMetric(c.omadaPortLinkStatus, prometheus.GaugeValue, p.PortStatus.LinkStatus,
					device.Name, device.Mac, "", "", port, p.SwitchMac, p.SwitchId, "", p.ProfileName, site, client.SiteId)

				ch <- prometheus.MustNewConstMetric(c.omadaPortLinkSpeedMbps, prometheus.GaugeValue, linkSpeed,
					device.Name, device.Mac, "", "", port, p.SwitchMac, p.SwitchId, "", p.ProfileName, site, client.SiteId)
			}
		}
	}
}

func getPortByLinkSpeed(ls float64) float64 {
	switch ls {
	case 0:
		return 0
	case 1:
		return 10
	case 2:
		return 100
	case 3:
		return 1000
	case 4:
		return 2500
	case 5:
		return 10000
	}
	return 0
}

func removeDuplicates(s []api.Port) []api.Port {
	// create map to track found items
	found := map[api.Port]bool{}
	res := []api.Port{}

	for v := range s {
		if found[s[v]] {
			// skip adding to new array if it exists
			continue
		}
		// add to new array, mark as found
		found[s[v]] = true
		res = append(res, s[v])
	}
	return res
}

func NewPortCollector(c *api.Client) *portCollector {
	return &portCollector{
		omadaPortPowerWatts: prometheus.NewDesc("omada_port_power_watts",
			"The current PoE usage of the port in watts.",
			[]string{"device", "device_mac", "client", "vendor", "switch_port", "switch_mac", "switch_id", "vlan_id", "profile", "site", "site_id"},
			nil,
		),
		omadaPortLinkStatus: prometheus.NewDesc("omada_port_link_status",
			"A boolean representing the link status of the port.",
			[]string{"device", "device_mac", "client", "vendor", "switch_port", "switch_mac", "switch_id", "vlan_id", "profile", "site", "site_id"},
			nil,
		),
		omadaPortLinkSpeedMbps: prometheus.NewDesc("omada_port_link_speed_mbps",
			"Port link speed in mbps. This is the capability of the connection, not the active throughput.",
			[]string{"device", "device_mac", "client", "vendor", "switch_port", "switch_mac", "switch_id", "vlan_id", "profile", "site", "site_id"},
			nil,
		),
		client: c,
	}
}
