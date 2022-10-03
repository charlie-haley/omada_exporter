package omada

import (
	"fmt"

	"github.com/charlie-haley/omada_exporter/pkg/api"
)

func Scrape(c *api.Client) error {
	site := c.Config.String("site")

	devices, err := c.GetDevices()
	if err != nil {
		return fmt.Errorf("failed to get devices: %s", err)
	}
	setDeviceMetrics(devices, site, c.SiteId)

	clients, err := c.GetClients()
	if err != nil {
		return fmt.Errorf("failed to get clients: %s", err)
	}
	setClientMetrics(clients, site, c.SiteId)

	for _, device := range devices {
		err := setPortMetricsByDevice(c, device, site, c.SiteId)
		if err != nil {
			return fmt.Errorf("failed to set port metrics: %s", err)
		}
	}

	controller, err := c.GetController()
	if err != nil {
		return fmt.Errorf("failed to get controller: %s", err)
	}
	setControllerMetrics(controller)

	return nil
}

// set prometheus metrics for controller
func setControllerMetrics(controller *api.Controller) {
	omada_controller_uptime_seconds.WithLabelValues(controller.Name, controller.Model, controller.ControllerVersion, controller.ControllerVersion, controller.MacAddress).Set(controller.Uptime / 1000)
	for _, s := range controller.Storage {
		omada_controller_storage_used_bytes.WithLabelValues(s.Name, controller.Name, controller.Model, controller.ControllerVersion, controller.ControllerVersion, controller.MacAddress).Set(s.Used * 1000000000)
		omada_controller_storage_available_bytes.WithLabelValues(s.Name, controller.Name, controller.Model, controller.ControllerVersion, controller.ControllerVersion, controller.MacAddress).Set(s.Total * 1000000000)
	}
}

// set prometheus metrics for devices
func setDeviceMetrics(devices []api.Device, site string, siteId string) {
	for _, item := range devices {
		needUpgrade := float64(0)
		if item.NeedUpgrade {
			needUpgrade = 1
		}
		omada_device_uptime_seconds.WithLabelValues(item.Name, item.Model, item.Version, item.Ip, item.Mac, site, siteId, item.Type).Set(item.Uptime)
		omada_device_cpu_percentage.WithLabelValues(item.Name, item.Model, item.Version, item.Ip, item.Mac, site, siteId, item.Type).Set(item.CpuUtil)
		omada_device_mem_percentage.WithLabelValues(item.Name, item.Model, item.Version, item.Ip, item.Mac, site, siteId, item.Type).Set(item.MemUtil)
		omada_device_need_upgrade.WithLabelValues(item.Name, item.Model, item.Version, item.Ip, item.Mac, site, siteId, item.Type).Set(needUpgrade)
		if item.Type == "ap" {
			omada_device_tx_rate.WithLabelValues(item.Name, item.Model, item.Version, item.Ip, item.Mac, site, siteId, item.Type).Set(item.TxRate)
			omada_device_rx_rate.WithLabelValues(item.Name, item.Model, item.Version, item.Ip, item.Mac, site, siteId, item.Type).Set(item.RxRate)
		}
		if item.Type == "switch" {
			omada_device_poe_remain_watts.WithLabelValues(item.Name, item.Model, item.Version, item.Ip, item.Mac, site, siteId, item.Type).Set(item.PoeRemain)
		}
	}
}

// set prometheus metrics for ports
func setPortMetricsByDevice(c *api.Client, device api.Device, site string, siteId string) error {
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

		client, err := c.GetClientByPort(device.Mac, p.Port)
		if err != nil {
			return fmt.Errorf("failed to get clients: %s", err)
		}

		port := fmt.Sprintf("%.0f", p.Port)
		if client != nil {
			vlanId := fmt.Sprintf("%.0f", client.VlanId)

			omada_port_power_watts.WithLabelValues(client.HostName, client.Vendor, port, p.SwitchMac, p.SwitchId, vlanId, p.ProfileName, site, siteId).Set(p.PortStatus.PoePower)
			omada_port_link_status.WithLabelValues(client.HostName, client.Vendor, port, p.SwitchMac, p.SwitchId, vlanId, p.ProfileName, site, siteId).Set(p.PortStatus.LinkStatus)
			omada_port_link_speed_mbps.WithLabelValues(client.HostName, client.Vendor, port, p.SwitchMac, p.SwitchId, vlanId, p.ProfileName, site, siteId).Set(linkSpeed)
		} else {
			omada_port_power_watts.WithLabelValues("", "", port, p.SwitchMac, p.SwitchId, "", p.ProfileName, site, siteId).Set(p.PortStatus.PoePower)
			omada_port_link_status.WithLabelValues("", "", port, p.SwitchMac, p.SwitchId, "", p.ProfileName, site, siteId).Set(p.PortStatus.LinkStatus)
			omada_port_link_speed_mbps.WithLabelValues("", "", port, p.SwitchMac, p.SwitchId, "", p.ProfileName, site, siteId).Set(linkSpeed)
		}

	}
	return nil
}

// set prometheus metrics for clients
func setClientMetrics(clients []api.NetworkClient, site string, siteId string) {
	omada_client_connected_total.WithLabelValues(site, siteId).Set(float64(len(clients)))

	for _, item := range clients {
		vlanId := fmt.Sprintf("%.0f", item.VlanId)
		port := fmt.Sprintf("%.0f", item.Port)
		if item.Wireless {
			wifiMode := fmt.Sprintf("%.0f", item.WifiMode)
			omada_client_download_activity_bytes.WithLabelValues(item.HostName, item.Vendor, "", "", item.Ip, item.Mac, site, siteId, item.ApName, item.Ssid, wifiMode).Set(item.Activity)
			omada_client_signal_dbm.WithLabelValues(item.HostName, item.Vendor, item.Ip, item.Mac, item.ApName, site, siteId, item.Ssid, wifiMode).Set(item.SignalLevel)
		}
		//
		if !item.Wireless {
			omada_client_download_activity_bytes.WithLabelValues(item.HostName, item.Vendor, port, vlanId, item.Ip, item.Mac, site, siteId, "", "", "").Set(item.Activity)
		}
	}
}
