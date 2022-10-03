package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (c *Client) GetDevices() ([]Device, error) {
	loggedIn, err := c.IsLoggedIn()
	if err != nil {
		return nil, err
	}
	if !loggedIn {
		log.Info(fmt.Errorf("not logged in, logging in with user: %s", c.Config.String("username")))
		err := c.Login()
		if err != nil || c.token == "" {
			log.Error(fmt.Errorf("failed to login: %s", err))
			return nil, err
		}
	}

	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/devices", c.Config.String("host"), c.omadaCID, c.SiteId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	setHeaders(req, c.token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	devicedata := deviceResponse{}
	err = json.Unmarshal(body, &devicedata)

	var ports []Port
	for i, d := range devicedata.Result {
		if d.Type == "switch" {
			switchPorts, err := c.GetPorts(d.Mac)
			if err != nil {
				return nil, fmt.Errorf("failed to get ports: %s", err)
			}
			ports = append(ports, switchPorts...)
			device := &devicedata.Result[i]
			device.Ports = ports
		}
	}

	return devicedata.Result, err
}

type deviceResponse struct {
	Result []Device `json:"result"`
}
type Device struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Mac         string  `json:"mac"`
	Model       string  `json:"model"`
	Version     string  `json:"version"`
	Ip          string  `json:"ip"`
	CpuUtil     float64 `json:"cpuUtil"`
	MemUtil     float64 `json:"memUtil"`
	Uptime      float64 `json:"uptimeLong"`
	NeedUpgrade bool    `json:"needUpgrade"`
	TxRate      float64 `json:"txRate"`
	RxRate      float64 `json:"rxRate"`
	PoeRemain   float64 `json:"poeRemain"`
	Ports       []Port  `json:"ports"`
}
