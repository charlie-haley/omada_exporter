# omada_exporter
![docker-publish](https://github.com/charlie-haley/omada_exporter/actions/workflows/docker-publish.yml/badge.svg)

Prometheus Exporter for TP-Link Omada Controller SDN.

### Tested Devices
- TL-SG3428MP
- EAP245 Access Point

## Installation
The exporter listens on port `9202` by default. It can be overidden with the `OMADA_EXPORTER_PORT` environment variable.

### TP-Link Omada User
I *highly* recommend you create a new user in the Omada SDN that has the `Viewer` role and use that to authenticate instead of your primary admin user.

### Docker
```
docker run -d \
    --network host \
    -e OMADA_HOST='https://192.168.1.20' \
    -e OMADA_USER='exporter' \
    -e OMADA_PASS='mypassword' \
    -e OMADA_SITE='Default' \
    chhaley/omada_exporter
```

### Helm
```
helm repo add charlie-haley http://charts.charliehaley.dev
helm repo update
helm install omada-exporter charlie-haley/omada-exporter \
    --set omada.ip=192.1.1.20 \ 
    --set omada.username=exporter \
    --set omada.password=mypassword \
    --set omada.site=Default \
    -n monitoring
```

If you want to use the ServiceMonitor (which is enabled by default) you'll need to have [prometheus-operator](https://github.com/prometheus-operator/prometheus-operator) deployed to your cluster, see [values](https://github.com/charlie-haley/private-charts/blob/main/charts/omada-exporter/values.yaml) to disable it if you'd like use ingress instead.

[You can find the chart repo here](https://github.com/charlie-haley/private-charts), if you'd like to contribute.

## Configuration

### Environment Variables
Variable                 | Purpose
-------------------------|-----------------------------------
OMADA_HOST               | Host of the Omada Controller SDN
OMADA_USER               | Username for the Omada user.
OMADA_PASS               | Password for the Omada user.
OMADA_SITE               | Site you'd like to get metrics from.
OMADA_EXPORTER_PORT      | Port the exporter should run on, default 9202
OMADA_INSECURE           | Whether to skip verifying the SSL certificate on the controller, defaults to false.

### Helm
```
# values.yaml
omada:
    host: 192.1.1.20       #Host of the Omada Controller SDN
    username: exporter     #Username for the Omada user.
    password: mypassword   #Host of the Omada Controller SDN
    site: Default          #Site you'd like to get metrics from.
    insecure: false        #Whether to skip verifying the SSL certificate on the controller, defaults to false.
```

## Metrics
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