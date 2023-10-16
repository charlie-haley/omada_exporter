package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/rs/zerolog/log"
)

func (c *Client) GetPorts(switchMac string) ([]Port, error) {
	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/switches/%s/ports", c.Config.Host, c.omadaCID, c.SiteId, switchMac)
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
	log.Debug().Bytes("data", body).Msg("Received data from ports endpoint")

	portdata := portResponse{}
	err = json.Unmarshal(body, &portdata)

	return portdata.Result, err
}

type portResponse struct {
	Result []Port `json:"result"`
}
type Port struct {
	Id          string     `json:"id"`
	SwitchId    string     `json:"switchId"`
	SwitchMac   string     `json:"switchMac"`
	Name        string     `json:"name"`
	PortStatus  portStatus `json:"portStatus"`
	Port        float64    `json:"port"`
	ProfileName string     `json:"profileName"`
}
type portStatus struct {
	Port       float64 `json:"id"`
	LinkStatus float64 `json:"linkStatus"`
	LinkSpeed  float64 `json:"linkSpeed"`
	PoePower   float64 `json:"poePower"`
	Poe        bool    `json:"poe"`
	Rx         float64 `json:"rx"`
	Tx         float64 `json:"tx"`
}
