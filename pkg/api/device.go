package api

import (
	"fmt"
	"io"
	"net/http"

	log "github.com/rs/zerolog/log"
)

func (c *Client) GetDevices() ([]Device, error) {
	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/devices", c.Config.Host, c.omadaCID, c.SiteId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.makeLoggedInRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Debug().Bytes("data", body).Msg("Received data from devices endpoint")

	devices, err := parseListResult[Device](body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse devices: %w", err)
	}

	return devices, nil
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
	Download    int64   `json:"download"`
	Upload      int64   `json:"upload"`
}
