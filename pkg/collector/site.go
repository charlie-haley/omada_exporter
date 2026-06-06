package collector

import (
	"github.com/charlie-haley/omada_exporter/pkg/api"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/rs/zerolog/log"
)

type siteCollector struct {
	omadaSiteConnectedAp          *prometheus.Desc
	omadaSiteDisconnectedAp       *prometheus.Desc
	omadaSiteIsolatedAp           *prometheus.Desc
	omadaSiteConnectedSwitch      *prometheus.Desc
	omadaSiteDisconnectedSwitch   *prometheus.Desc
	omadaSiteConnectedGateway     *prometheus.Desc
	omadaSiteDisconnectedGateway  *prometheus.Desc
	omadaSiteTotalPorts           *prometheus.Desc
	omadaSiteAvailablePorts       *prometheus.Desc
	omadaSitePowerConsumptionWatts *prometheus.Desc
	omadaSiteWiredClients         *prometheus.Desc
	omadaSiteWirelessClients      *prometheus.Desc
	omadaSiteGuestClients         *prometheus.Desc
	client                        *api.Client
}

func (c *siteCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.omadaSiteConnectedAp
	ch <- c.omadaSiteDisconnectedAp
	ch <- c.omadaSiteIsolatedAp
	ch <- c.omadaSiteConnectedSwitch
	ch <- c.omadaSiteDisconnectedSwitch
	ch <- c.omadaSiteConnectedGateway
	ch <- c.omadaSiteDisconnectedGateway
	ch <- c.omadaSiteTotalPorts
	ch <- c.omadaSiteAvailablePorts
	ch <- c.omadaSitePowerConsumptionWatts
	ch <- c.omadaSiteWiredClients
	ch <- c.omadaSiteWirelessClients
	ch <- c.omadaSiteGuestClients
}

func (c *siteCollector) Collect(ch chan<- prometheus.Metric) {
	client := c.client
	config := c.client.Config

	site := config.Site
	overview, err := client.GetSiteOverview()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get site overview")
		return
	}

	labels := []string{site, client.SiteId}

	ch <- prometheus.MustNewConstMetric(c.omadaSiteConnectedAp, prometheus.GaugeValue, overview.ConnectedApNum, labels...)
	ch <- prometheus.MustNewConstMetric(c.omadaSiteDisconnectedAp, prometheus.GaugeValue, overview.DisconnectedApNum, labels...)
	ch <- prometheus.MustNewConstMetric(c.omadaSiteIsolatedAp, prometheus.GaugeValue, overview.IsolatedApNum, labels...)
	ch <- prometheus.MustNewConstMetric(c.omadaSiteConnectedSwitch, prometheus.GaugeValue, overview.ConnectedSwitchNum, labels...)
	ch <- prometheus.MustNewConstMetric(c.omadaSiteDisconnectedSwitch, prometheus.GaugeValue, overview.DisconnectedSwitchNum, labels...)
	ch <- prometheus.MustNewConstMetric(c.omadaSiteConnectedGateway, prometheus.GaugeValue, overview.ConnectedGatewayNum, labels...)
	ch <- prometheus.MustNewConstMetric(c.omadaSiteDisconnectedGateway, prometheus.GaugeValue, overview.DisconnectedGatewayNum, labels...)
	ch <- prometheus.MustNewConstMetric(c.omadaSiteTotalPorts, prometheus.GaugeValue, overview.TotalPorts, labels...)
	ch <- prometheus.MustNewConstMetric(c.omadaSiteAvailablePorts, prometheus.GaugeValue, overview.AvailablePorts, labels...)
	ch <- prometheus.MustNewConstMetric(c.omadaSitePowerConsumptionWatts, prometheus.GaugeValue, overview.PowerConsumption, labels...)
	ch <- prometheus.MustNewConstMetric(c.omadaSiteWiredClients, prometheus.GaugeValue, overview.WiredClientNum, labels...)
	ch <- prometheus.MustNewConstMetric(c.omadaSiteWirelessClients, prometheus.GaugeValue, overview.WirelessClientNum, labels...)
	ch <- prometheus.MustNewConstMetric(c.omadaSiteGuestClients, prometheus.GaugeValue, overview.GuestNum, labels...)
}

func NewSiteCollector(c *api.Client) *siteCollector {
	labels := []string{"site", "site_id"}

	return &siteCollector{
		omadaSiteConnectedAp: prometheus.NewDesc("omada_site_connected_ap_num",
			"Number of connected access points.", labels, nil),
		omadaSiteDisconnectedAp: prometheus.NewDesc("omada_site_disconnected_ap_num",
			"Number of disconnected access points.", labels, nil),
		omadaSiteIsolatedAp: prometheus.NewDesc("omada_site_isolated_ap_num",
			"Number of isolated access points.", labels, nil),
		omadaSiteConnectedSwitch: prometheus.NewDesc("omada_site_connected_switch_num",
			"Number of connected switches.", labels, nil),
		omadaSiteDisconnectedSwitch: prometheus.NewDesc("omada_site_disconnected_switch_num",
			"Number of disconnected switches.", labels, nil),
		omadaSiteConnectedGateway: prometheus.NewDesc("omada_site_connected_gateway_num",
			"Number of connected gateways.", labels, nil),
		omadaSiteDisconnectedGateway: prometheus.NewDesc("omada_site_disconnected_gateway_num",
			"Number of disconnected gateways.", labels, nil),
		omadaSiteTotalPorts: prometheus.NewDesc("omada_site_total_ports",
			"Total number of switch ports.", labels, nil),
		omadaSiteAvailablePorts: prometheus.NewDesc("omada_site_available_ports",
			"Number of available (unused) switch ports.", labels, nil),
		omadaSitePowerConsumptionWatts: prometheus.NewDesc("omada_site_power_consumption_watts",
			"Total PoE power consumption in watts.", labels, nil),
		omadaSiteWiredClients: prometheus.NewDesc("omada_site_wired_clients",
			"Number of connected wired clients.", labels, nil),
		omadaSiteWirelessClients: prometheus.NewDesc("omada_site_wireless_clients",
			"Number of connected wireless clients.", labels, nil),
		omadaSiteGuestClients: prometheus.NewDesc("omada_site_guest_clients",
			"Number of connected guest clients.", labels, nil),
		client: c,
	}
}
