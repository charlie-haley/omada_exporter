# omada_exporter
![docker-publish](https://github.com/charlie-haley/omada_exporter/actions/workflows/docker-publish.yml/badge.svg)
<p align="center" style="text-align: center">
    <img src="./docs/images/logo-dark-mode.svg#gh-dark-mode-only" width="70%"><br/>
    <img src="./docs/images/logo-light-mode.svg#gh-light-mode-only" width="70%"><br/>
    Prometheus Exporter for TP-Link Omada Controller SDN. <br/>
</p>

### 📈 Dashboard

There's a [default dashboard in this repo](./dashboard.json), which would be a good starting point for visualising your metrics. [The dashboard is also available on Grafana.com.](https://grafana.com/grafana/dashboards/16343)

<p align="center" style="text-align: center">
    <img src="./docs/images/dashboard-dark-mode.png#gh-dark-mode-only" width="70%"><br/>
    <img src="./docs/images/dashboard-light-mode.png#gh-light-mode-only" width="70%"><br/>
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

__There's also a GHCR mirror available if you'd prefer to not use Docker Hub. `ghcr.io/charlie-haley/omada_exporter`__

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
[You can download the latest binary release here.](https://github.com/charlie-haley/omada_exporter/releases/latest)
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
   --host value                 The hostname of the Omada Controller, including protocol. [$OMADA_HOST]
   --username value             Username of the Omada user you'd like to use to fetch metrics. [$OMADA_USER]
   --password value             Password for your Omada user. [$OMADA_PASS]
   --port value                 Port on which to expose the Prometheus metrics. (default: "9202") [$OMADA_PORT]
   --site value                 Omada site to scrape metrics from. (default: "Default") [$OMADA_SITE]
   --log-level value            Application log level. (default: "error") [$LOG_LEVEL]
   --timeout value              Timeout when making requests to the Omada Controller. (default: 15) [$OMADA_REQUEST_TIMEOUT]
   --insecure                   Whether to skip verifying the SSL certificate on the controller. (default: false) [$OMADA_INSECURE]
   --disable-go-collector       Disable Go collector metrics. (default: true) [$OMADA_DISABLE_GO_COLLECTOR]
   --disable-process-collector  Disable process collector metrics. (default: true) [$OMADA_DISABLE_PROCESS_COLLECTOR]
   --help, -h                   show help (default: false)
   --version, -v                print the version (default: false)
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
OMADA_REQUEST_TIMEOUT    | Timeout when making requests to the Omada Controller. (default: 15)
OMADA_DISABLE_GO_COLLECTOR | Disable Go collector metrics. (default: true)
OMADA_DISABLE_PROCESS_COLLECTOR | Disable process collector metrics. (default: true)
LOG_LEVEL                       | Application log level. (default: "error")

### Helm
```
# values.yaml
omada:
    host: "https://192.1.1.20" # The hostname of the Omada Controller, including protocol.
    username: "exporter"       # Username of the Omada user you'd like to use to fetch metrics.
    password: "mypassword"     # Password for your Omada user.
    site: "Default"            # Site you'd like to get metrics from. (default: "Default")
    insecure: false            # Whether to skip verifying the SSL certificate on the controller. (default: false)
```

## 📊 Metrics
| Name | Description | Labels |
|--|--|--|
| omada_client_download_activity_bytes | The current download activity for the client in bytes. | client vendor ip mac host_name site site_id connection_mode wifi_mode ap_name ssid switch_port vlan_id |
| omada_client_signal_pct | The signal quality for the wireless client in percent. | client vendor ip mac host_name site site_id connection_mode wifi_mode ap_name ssid |
| omada_client_snr_dbm | The signal to noise ration for the wireless client in dBm. | client vendor ip mac host_name site site_id connection_mode wifi_mode ap_name ssid |
| omada_client_rssi_dbm | The RSSI for the wireless client in dBm. | client vendor ip mac host_name site site_id connection_mode wifi_mode ap_name ssid |
| omada_client_traffic_down_bytes | Total bytes received by wireless client. | client vendor ip mac host_name site site_id connection_mode wifi_mode ap_name ssid |
| omada_client_traffic_up_bytes | Total bytes sent by wireless client. | client vendor ip mac host_name site site_id connection_mode wifi_mode ap_name ssid |
| omada_client_tx_rate | TX rate of wireless client. | client vendor ip mac host_name site site_id connection_mode wifi_mode ap_name ssid |
| omada_client_rx_rate | RX rate of wireless client. | client vendor ip mac host_name site site_id connection_mode wifi_mode ap_name ssid |
| omada_client_connected_total | Total number of connected clients. | site site_id connection_mode wifi_mode |
| omada_controller_uptime_seconds | Uptime of the controller. | controller_name model controller_version firmware_version mac site site_id |
| omada_controller_storage_used_bytes | Storage used on the controller. | storage_name controller_name model controller_version firmware_version mac site site_id |
| omada_controller_storage_available_bytes | Total storage available for the controller. | storage_name controller_name model controller_version firmware_version mac site site_id |
| omada_device_uptime_seconds | Uptime of the device. | device model version ip mac site site_id device_type |
| omada_device_uptime_seconds | Uptime of the device. | device model version ip mac site site_id device_type |
| omada_device_cpu_percentage | Percentage of device CPU used. | device model version ip mac site site_id device_type |
| omada_device_mem_percentage | Percentage of device Memory used. | device model version ip mac site site_id device_type |
| omada_device_need_upgrade | A boolean on whether the device needs an upgrade. | device model version ip mac site site_id device_type |
| omada_device_tx_rate | The tx rate of the device. | device model version ip mac site site_id device_type |
| omada_device_rx_rate | The rx rate of the device. | device model version ip mac site site_id device_type |
| omada_device_poe_remain_watts | The remaining amount of PoE power for the device in watts. | device model version ip mac site site_id device_type |
| omada_device_download | Device download traffic. | device model version ip mac site site_id device_type |
| omada_device_upload | Device upload traffic. | device model version ip mac site site_id device_type |
| omada_port_power_watts | The current PoE usage of the port in watts. | device device_mac client vendor switch_port name switch_mac switch_id vlan_id profile site site_id |
| omada_port_link_status | A boolean representing the link status of the port. | device device_mac client vendor switch_port name switch_mac switch_id vlan_id profile site site_id |
| omada_port_link_speed_mbps | Port link speed in mbps. This is the capability of the connection, not the active throughput. | device device_mac client vendor switch_port name switch_mac switch_id vlan_id profile site site_id |
| omada_port_link_rx | Bytes recieved on a port. | device device_mac client vendor switch_port name switch_mac switch_id vlan_id profile site site_id |
| omada_port_link_tx | Bytes transmitted on a port. | device device_mac client vendor switch_port name switch_mac switch_id vlan_id profile site site_id |
