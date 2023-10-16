package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/rs/zerolog/log"
)

func (c *Client) GetController() (*Controller, error) {
	url := fmt.Sprintf("%s/%s/api/v2/maintenance/controllerStatus?", c.Config.Host, c.omadaCID)
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
	log.Debug().Bytes("data", body).Msg("Received data from controllerStatus endpoint")

	controllerData := controllerResponse{}
	err = json.Unmarshal(body, &controllerData)

	return &controllerData.Result, err
}

type controllerResponse struct {
	Result Controller `json:"result"`
}
type Controller struct {
	Name              string       `json:"name"`
	MacAddress        string       `json:"macAddress"`
	FirmwareVersion   string       `json:"firmwareVersion"`
	ControllerVersion string       `json:"controllerVersion"`
	Model             string       `json:"model"`
	Uptime            float64      `json:"upTime"`
	Storage           []hwcStorage `json:"hwcStorage"`
}
type hwcStorage struct {
	Name  string  `json:"name"`
	Total float64 `json:"totalStorage"`
	Used  float64 `json:"usedStorage"`
}
