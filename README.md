# omada_exporter
![docker-publish](https://github.com/charlie-haley/omada_exporter/actions/workflows/docker-publish.yml/badge.svg)
<p align="center" style="text-align: center">
    <img src="./docs/images/logo-dark-mode.svg#gh-dark-mode-only" width="70%"><br/>
    <img src="./docs/images/logo-light-mode.svg#gh-light-mode-only" width="70%"><br/>
    Prometheus Exporter for TP-Link Omada Controller SDN. <br/>
</p>

## Installation

__I *highly* recommend you create a new user in the Omada SDN that has the `Viewer` role and use that to authenticate instead of your primary admin user.__

### 🐋 Docker
```bash
docker run -d \
    -p 9202:9202 \
    -e OMADA_HOST='https://192.168.1.20' \
    -e OMADA_USER='exporter' \
    -e OMADA_PASS='mypassword' \
    -e OMADA_SITE='Default' \
    chhaley/omada_exporter
```

### ☸️ Helm
```bash
helm repo add charlie-haley http://charts.charliehaley.dev
helm repo update
helm install omada-exporter charlie-haley/omada-exporter \
    --set omada.host=https://192.1.1.20 \ 
    --set omada.username=exporter \
    --set omada.password=mypassword \
    --set omada.site=Default \
    -n monitoring
```

If you want to use the ServiceMonitor (which is enabled by default) you'll need to have [prometheus-operator](https://github.com/prometheus-operator/prometheus-operator) deployed to your cluster, see [values](https://github.com/charlie-haley/private-charts/blob/main/charts/omada-exporter/values.yaml) to disable it if you'd like use ingress instead.

[You can find the chart repo here](https://github.com/charlie-haley/private-charts), if you'd like to contribute.

### 🖥️ Command Line
```
NAME:
   omada-exporter - Prometheus Exporter for TP-Link Omada Controller SDN.

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   development

AUTHOR:
   Charlie Haley <charlie-haley@users.noreply.github.com>

COMMANDS:
   version, v  prints the current version.
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value      The hostname of the Omada Controller, including protocol. [$OMADA_HOST]
   --username value  Username of the Omada user you'd like to use to fetch metrics. [$OMADA_USER]
   --password value  Password for your Omada user. [$OMADA_PASS]
   --port value      Port on which to expose the Prometheus metrics. (default: "9202") [$OMADA_PORT]
   --site value      Password for your Omada user. [$OMADA_SITE]
   --interval value  Interval between scrapes, in seconds. (default: 5) [$OMADA_SCRAPE_INTERVAL]
   --insecure        Whether to skip verifying the SSL certificate on the controller. (default: false) [$OMADA_INSECURE]
   --help, -h        show help (default: false)
   --version, -v     print the version (default: false)
```

## ⚙️ Configuration

### Environment Variables
Variable                 | Purpose
-------------------------|-----------------------------------
OMADA_HOST               | The hostname of the Omada Controller, including protocol.
OMADA_USER               | Username of the Omada user you'd like to use to fetch metrics.
OMADA_PASS               | Password for your Omada user.
OMADA_SITE               | Site you'd like to get metrics from. (default: "Default")
OMADA_PORT               | Port on which to expose the Prometheus metrics. (default: 9202)
OMADA_INSECURE           | Whether to skip verifying the SSL certificate on the controller. (default: false)
OMADA_SCRAPE_INTERVAL    | Interval between scrapes, in seconds. (default: 5)

### Helm
```
# values.yaml
omada:
    host: "https://192.1.1.20" # The hostname of the Omada Controller, including protocol.
    username: "exporter"       # Username of the Omada user you'd like to use to fetch metrics.
    password: "mypassword"     # Password for your Omada user.
    site: "Default"            # Site you'd like to get metrics from. (default: "Default")
    insecure: false            # Whether to skip verifying the SSL certificate on the controller. (default: false)
    scrape_interval: 5         # Interval between scrapes, in seconds. (default: 5)
```

## 📊 Metrics
Name                               | Description                                                 | Labels
-----------------------------------|-------------------------------------------------------------|---------------------------------------------------
omada_uptime_seconds               | Uptime of the device.                                       | device, model, version, ip, mac, site, device_type
omada_cpu_percentage               | Percentage of device CPU used.                              | device, model, version, ip, mac, site, device_type
omada_mem_percentage               | Percentage of device Memory used.                           | device, model, version, ip, mac, site, device_type
omada_device_need_upgrade          | A boolean on whether the device needs an upgrade.           | device, model, version, ip, mac, site, device_type
omada_tx_rate                      | The tx rate of the device.                                  | device, model, version, ip, mac, site, device_type
omada_rx_rate                      | The rx rate of the device.                                  | device, model, version, ip, mac, site, device_type
omada_poe_remain_watts             | The remaining amount of PoE power for the device in watts.  | device, model, version, ip, mac, site, device_type
omada_download_activity_bytes      | The current download activity for the LAN client in bytes.  | client, vendor, switch_port, vlan_id, ip, mac, site
omada_wlan_download_activity_bytes | The current download activity for the WLAN client in bytes. | client, vendor, ip, mac, ap_name, site, ssid, wifi_mode
omada_client_signal_dbm            | The signal level for the wireless client in dBm.            | device, model, version, ip, mac, site, device_type
omada_port_power_watts             | The current PoE usage of the port in watts.                 | device, model, version, ip, mac, site, device_type
omada_port_link_status             | A boolean representing the link status of the port.         | device, model, version, ip, mac, site, device_type
omada_port_link_speed_mbps         | Port link speed in mbps. This is the capability of the connection, not the active throughput. | device, model, version, ip, mac, site, device_type
