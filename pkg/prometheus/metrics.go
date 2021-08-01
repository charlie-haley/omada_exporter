package metrics


import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	Omada_uptime_seconds = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omada_uptime_seconds",
		Help: "Uptime of the device.",
	},
	[]string{"device", "model", "version", "ip", "mac", "site", "device_type"})
)

var (
	Omada_cpu_percentage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omada_cpu_percentage",
		Help: "Percentage of device CPU used.",
	},
	[]string{"device", "model", "version", "ip", "mac", "site", "device_type"})
)

var (
	Omada_mem_percentage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omada_mem_percentage",
		Help: "Percentage of device Memory used.",
	},
	[]string{"device", "model", "version", "ip", "mac", "site", "device_type"})
)

var (
	Omada_device_need_upgrade = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omada_device_need_upgrade",
		Help: "A boolean on whether the device needs an upgrade.",
	},
	[]string{"device", "model", "version", "ip", "mac", "site", "device_type"})
)

var (
	Omada_tx_rate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omada_tx_rate",
		Help: "The tx rate of the device.",
	},
	[]string{"device", "model", "version", "ip", "mac", "site", "device_type"})
)

var (
	Omada_rx_rate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omada_rx_rate",
		Help: "The rx rate of the device.",
	},
	[]string{"device", "model", "version", "ip", "mac", "site", "device_type"})
)

var (
	Omada_poe_remain_watts = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omada_poe_remain_watts",
		Help: "The remaining amount of PoE power for the device in watts.",
	},
	[]string{"device", "model", "version", "ip", "mac", "site", "device_type"})
)

var (
	Omada_download_activity_bytes = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omada_download_activity_bytes",
		Help: "The current download activity for the LAN client in bytes.",
	},
	[]string{"client", "vendor", "switch_port", "vlan_id", "ip", "mac", "site"})
)

var (
	Omada_download_activity_bytes_wlan = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omada_wlan_download_activity_bytes",
		Help: "The current download activity for the WLAN client in bytes.",
	},
	[]string{"client", "vendor", "ip", "mac", "ap_name", "site", "ssid", "wifi_mode"})
)

var (
	Omada_client_signal_dbm = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omada_client_signal_dbm",
		Help: "The signal level for the wireless client in dBm.",
	},
	[]string{"client", "vendor", "ip", "mac", "ap_name", "site", "ssid", "wifi_mode"})
)

var (
	Omada_port_power_watts = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omada_port_power_watts",
		Help: "The current PoE usage of the port in watts.",
	},
	[]string{"client", "vendor", "switch_port","switch_mac", "switch_id", "vlan_id", "profile", "site"})
)

var (
	Omada_port_link_status = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omada_port_link_status",
		Help: "A boolean representing the link status of the port.",
	},
	[]string{"client", "vendor", "switch_port", "switch_mac", "switch_id", "vlan_id", "profile", "site"})
)

var (
	Omada_port_link_speed_mbps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omada_port_link_speed_mbps",
		Help: "Port link speed in mbps. This is the capability of the connection, not the active throughput.",
	},
	[]string{"client", "vendor", "switch_port", "switch_mac", "switch_id", "vlan_id", "profile", "site"})
)

