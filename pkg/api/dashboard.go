package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/rs/zerolog/log"
)

func (c *Client) GetSiteOverview() (*SiteOverview, error) {
	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/dashboard/overviewDiagram", c.Config.Host, c.omadaCID, c.SiteId)
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
	log.Debug().Bytes("data", body).Msg("Received data from overviewDiagram endpoint")

	raw, err := checkResponse(body)
	if err != nil {
		return nil, fmt.Errorf("failed to get site overview: %w", err)
	}

	var overview SiteOverview
	err = json.Unmarshal(raw, &overview)
	if err != nil {
		return nil, fmt.Errorf("failed to parse site overview: %w", err)
	}

	return &overview, nil
}

type SiteOverview struct {
	TotalApNum          float64 `json:"totalApNum"`
	ConnectedApNum      float64 `json:"connectedApNum"`
	DisconnectedApNum   float64 `json:"disconnectedApNum"`
	IsolatedApNum       float64 `json:"isolatedApNum"`
	TotalSwitchNum      float64 `json:"totalSwitchNum"`
	ConnectedSwitchNum  float64 `json:"connectedSwitchNum"`
	DisconnectedSwitchNum float64 `json:"disconnectedSwitchNum"`
	TotalGatewayNum     float64 `json:"totalGatewayNum"`
	ConnectedGatewayNum float64 `json:"connectedGatewayNum"`
	DisconnectedGatewayNum float64 `json:"disconnectedGatewayNum"`
	TotalPorts          float64 `json:"totalPorts"`
	AvailablePorts      float64 `json:"availablePorts"`
	PowerConsumption    float64 `json:"powerConsumption"`
	TotalClientNum      float64 `json:"totalClientNum"`
	WiredClientNum      float64 `json:"wiredClientNum"`
	WirelessClientNum   float64 `json:"wirelessClientNum"`
	GuestNum            float64 `json:"guestNum"`
}
