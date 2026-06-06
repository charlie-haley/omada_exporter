package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/rs/zerolog/log"
)

// GetClients fetches all active clients. It tries the OpenAPI v2 endpoint first
// (POST with JSON body, same as the web UI), then falls back to the legacy api/v2
// GET endpoint, and finally to the insight endpoint.
func (c *Client) GetClients() ([]NetworkClient, error) {
	// Try OpenAPI v2 endpoint (works for all clients on v6.x, same as web UI)
	clients, err := c.getClientsOpenAPI()
	if err == nil {
		return clients, nil
	}
	log.Debug().Err(err).Msg("OpenAPI clients endpoint failed, trying legacy endpoint")

	// Fallback: try legacy api/v2 GET endpoint (works on older controllers)
	clients, err = c.getClientsLegacy()
	if err == nil {
		return clients, nil
	}
	log.Debug().Err(err).Msg("Legacy clients endpoint failed, trying insight fallback")

	// Final fallback: insight endpoint (limited data but always works)
	clients, insightErr := c.getClientsFromInsight()
	if insightErr != nil {
		return nil, fmt.Errorf("all client endpoints failed: openapi, legacy, insight (%v)", insightErr)
	}

	return clients, nil
}

// getClientsOpenAPI fetches clients via the OpenAPI v2 POST endpoint.
// This is the same endpoint the Omada web UI uses and returns full data
// for both wired and wireless clients, even with Viewer role on v6.x.
func (c *Client) getClientsOpenAPI() ([]NetworkClient, error) {
	url := fmt.Sprintf("%s/openapi/v2/%s/sites/%s/clients", c.Config.Host, c.omadaCID, c.SiteId)

	filters := map[string]interface{}{"active": true}

	reqBody := map[string]interface{}{
		"filters":  filters,
		"page":     1,
		"pageSize": 1000,
		"scope":    1,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Omada-Request-Source", "web-local")

	resp, err := c.makeLoggedInRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Debug().Bytes("data", body).Msg("Received data from OpenAPI clients endpoint")

	clients, err := parseListResult[NetworkClient](body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI clients: %w", err)
	}

	log.Info().Int("clients", len(clients)).Msg("Fetched clients via OpenAPI v2 endpoint")
	return clients, nil
}

// getClientsLegacy fetches clients via the legacy api/v2 GET endpoint.
// This works on older Omada controllers but may fail on v6.x for Viewer role.
func (c *Client) getClientsLegacy() ([]NetworkClient, error) {
	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/clients", c.Config.Host, c.omadaCID, c.SiteId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("currentPage", "1")
	q.Add("currentPageSize", "10000")
	q.Add("filters.active", "true")

	req.URL.RawQuery = q.Encode()

	resp, err := c.makeLoggedInRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Debug().Bytes("data", body).Msg("Received data from legacy clients endpoint")

	clients, err := parseListResult[NetworkClient](body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse legacy clients: %w", err)
	}

	return clients, nil
}

// getClientsFromInsight fetches clients from the insight endpoint and filters
// for recently active ones. This endpoint is available to Viewer roles on v6.x
// controllers where the standard clients endpoint returns "General error."
func (c *Client) getClientsFromInsight() ([]NetworkClient, error) {
	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/insight/clients", c.Config.Host, c.omadaCID, c.SiteId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("currentPage", "1")
	q.Add("currentPageSize", "10000")
	req.URL.RawQuery = q.Encode()

	resp, err := c.makeLoggedInRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Debug().Bytes("data", body).Msg("Received data from insight/clients endpoint")

	raw, err := checkResponse(body)
	if err != nil {
		return nil, err
	}

	type insightClient struct {
		Name     string  `json:"name"`
		Mac      string  `json:"mac"`
		Wireless bool    `json:"wireless"`
		Download float64 `json:"download"`
		Upload   float64 `json:"upload"`
		LastSeen int64   `json:"lastSeen"`
		VlanId   float64 `json:"vid"`
	}

	// Parse paginated result
	var paginated struct {
		Data []insightClient `json:"data"`
	}
	if err := json.Unmarshal(raw, &paginated); err != nil {
		return nil, fmt.Errorf("failed to parse insight clients: %w", err)
	}

	// Filter for active clients (seen in the last 10 minutes)
	cutoff := time.Now().Add(-10 * time.Minute).UnixMilli()
	var clients []NetworkClient
	for _, ic := range paginated.Data {
		if ic.LastSeen < cutoff {
			continue
		}
		clients = append(clients, NetworkClient{
			Name:        ic.Name,
			Mac:         ic.Mac,
			Wireless:    ic.Wireless,
			TrafficDown: ic.Download,
			TrafficUp:   ic.Upload,
			VlanId:      ic.VlanId,
		})
	}

	log.Info().Int("active", len(clients)).Int("total", len(paginated.Data)).
		Msg("Using insight/clients fallback (some metrics unavailable)")

	return clients, nil
}

type NetworkClient struct {
	Name        string  `json:"name"`
	HostName    string  `json:"hostName"`
	Mac         string  `json:"mac"`
	Port        float64 `json:"port"`
	Ip          string  `json:"ip"`
	VlanId      float64 `json:"vid"`
	ApName      string  `json:"apName"`
	Wireless    bool    `json:"wireless"`
	SwitchMac   string  `json:"switchMac"`
	Vendor      string  `json:"vendor"`
	Activity    float64 `json:"activity"`
	SignalLevel float64 `json:"signalLevel"`
	SignalNoise float64 `json:"snr"`
	WifiMode    float64 `json:"wifiMode"`
	Ssid        string  `json:"ssid"`
	Rssi        float64 `json:"rssi"`
	TrafficDown float64 `json:"trafficDown"`
	TrafficUp   float64 `json:"trafficUp"`
	RxRate      float64 `json:"rxRate"`
	TxRate      float64 `json:"txRate"`
}
