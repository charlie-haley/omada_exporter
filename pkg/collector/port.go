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
	omadaPortLinkRx        *prometheus.Desc
	omadaPortLinkTx        *prometheus.Desc
	client                 *api.Client
}

func (c *portCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.omadaPortPowerWatts
	ch <- c.omadaPortLinkStatus
	ch <- c.omadaPortLinkSpeedMbps
	ch <- c.omadaPortLinkRx
	ch <- c.omadaPortLinkTx
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
			var cHostName, cVendor, cVlanID string
			linkSpeed := getPortByLinkSpeed(p.PortStatus.LinkSpeed)

			portClient, err := client.GetClientByPort(device.Mac, p.Port)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get client by port")
			}

			port := fmt.Sprintf("%.0f", p.Port)
			if portClient != nil {
				cHostName = portClient.HostName
				cVendor = portClient.Vendor
				cVlanID = fmt.Sprintf("%.0f", portClient.VlanId)
			}

			labels := []string{device.Name, device.Mac, cHostName, cVendor, port, p.Name, p.SwitchMac, p.SwitchId, cVlanID, p.ProfileName, site, client.SiteId}

			ch <- prometheus.MustNewConstMetric(c.omadaPortPowerWatts, prometheus.GaugeValue, p.PortStatus.PoePower, labels...)
			ch <- prometheus.MustNewConstMetric(c.omadaPortLinkStatus, prometheus.GaugeValue, p.PortStatus.LinkStatus, labels...)
			ch <- prometheus.MustNewConstMetric(c.omadaPortLinkSpeedMbps, prometheus.GaugeValue, linkSpeed, labels...)
			ch <- prometheus.MustNewConstMetric(c.omadaPortLinkRx, prometheus.CounterValue, p.PortStatus.Rx, labels...)
			ch <- prometheus.MustNewConstMetric(c.omadaPortLinkTx, prometheus.CounterValue, p.PortStatus.Tx, labels...)
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
	labels := []string{"device", "device_mac", "client", "vendor", "switch_port", "name", "switch_mac", "switch_id", "vlan_id", "profile", "site", "site_id"}

	return &portCollector{
		omadaPortPowerWatts: prometheus.NewDesc("omada_port_power_watts",
			"The current PoE usage of the port in watts.",
			labels,
			nil,
		),
		omadaPortLinkStatus: prometheus.NewDesc("omada_port_link_status",
			"A boolean representing the link status of the port.",
			labels,
			nil,
		),
		omadaPortLinkSpeedMbps: prometheus.NewDesc("omada_port_link_speed_mbps",
			"Port link speed in mbps. This is the capability of the connection, not the active throughput.",
			labels,
			nil,
		),
		omadaPortLinkRx: prometheus.NewDesc("omada_port_link_rx",
			"Bytes recieved on a port.",
			labels,
			nil,
		),
		omadaPortLinkTx: prometheus.NewDesc("omada_port_link_tx",
			"Bytes transmitted on a port.",
			labels,
			nil,
		),
		client: c,
	}
}
