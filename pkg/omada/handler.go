package omada

import (
	"fmt"

	"github.com/charlie-haley/omada_exporter/pkg/api"
)

func Scrape(c *api.Client) error {
	devices, err := c.GetDevices()
	if err != nil {
		return fmt.Errorf("failed to get devices: %s", err)
	}
	clients, err := c.GetClients()
	if err != nil {
		return fmt.Errorf("failed to get client: %s", err)
	}
	var ports []api.Port
	// ensure slice is empty
	ports = nil

	for _, item := range devices {
		needUpgrade := float64(0)
		if item.NeedUpgrade {
			needUpgrade = 1
		}
		omada_uptime_seconds.WithLabelValues(item.Name, item.Model, item.Version, item.Ip, item.Mac, c.Config.String("site"), item.Type).Set(item.Uptime)
		omada_cpu_percentage.WithLabelValues(item.Name, item.Model, item.Version, item.Ip, item.Mac, c.Config.String("site"), item.Type).Set(item.CpuUtil)
		omada_mem_percentage.WithLabelValues(item.Name, item.Model, item.Version, item.Ip, item.Mac, c.Config.String("site"), item.Type).Set(item.MemUtil)
		omada_device_need_upgrade.WithLabelValues(item.Name, item.Model, item.Version, item.Ip, item.Mac, c.Config.String("site"), item.Type).Set(needUpgrade)
		if item.Type == "ap" {
			omada_tx_rate.WithLabelValues(item.Name, item.Model, item.Version, item.Ip, item.Mac, c.Config.String("site"), item.Type).Set(item.TxRate)
			omada_rx_rate.WithLabelValues(item.Name, item.Model, item.Version, item.Ip, item.Mac, c.Config.String("site"), item.Type).Set(item.RxRate)
		}
		if item.Type == "switch" {
			omada_poe_remain_watts.WithLabelValues(item.Name, item.Model, item.Version, item.Ip, item.Mac, c.Config.String("site"), item.Type).Set(item.PoeRemain)
			switchPorts, err := c.GetPorts(item.Mac)
			if err != nil {
				return fmt.Errorf("failed to get ports: %s", err)
			}
			ports = append(ports, switchPorts...)
		}
	}
	for _, item := range clients {
		vlanId := fmt.Sprintf("%.0f", item.VlanId)
		port := fmt.Sprintf("%.0f", item.Port)
		if item.Wireless {
			wifiMode := fmt.Sprintf("%.0f", item.WifiMode)
			omada_download_activity_bytes_wlan.WithLabelValues(item.HostName, item.Vendor, item.Ip, item.Mac, item.ApName, c.Config.String("site"), item.Ssid, wifiMode).Set(item.Activity)
			omada_client_signal_dbm.WithLabelValues(item.HostName, item.Vendor, item.Ip, item.Mac, item.ApName, c.Config.String("site"), item.Ssid, wifiMode).Set(item.SignalLevel)
		}
		if item.Wireless {
			for _, p := range ports {
				if p.SwitchMac == item.SwitchMac && p.Port == item.Port {
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
					omada_port_power_watts.WithLabelValues(item.HostName, item.Vendor, port, p.SwitchId, p.SwitchMac, vlanId, p.ProfileName, c.Config.String("site")).Set(p.PortStatus.PoePower)
					omada_port_link_status.WithLabelValues(item.HostName, item.Vendor, port, p.SwitchId, p.SwitchMac, vlanId, p.ProfileName, c.Config.String("site")).Set(p.PortStatus.LinkStatus)
					omada_port_link_speed_mbps.WithLabelValues(item.HostName, item.Vendor, p.SwitchId, p.SwitchMac, port, vlanId, p.ProfileName, c.Config.String("site")).Set(linkSpeed)
				}
			}
			omada_download_activity_bytes.WithLabelValues(item.HostName, item.Vendor, port, vlanId, item.Ip, item.Mac, c.Config.String("site")).Set(item.Activity)
		}
	}
	return nil
}
