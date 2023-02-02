package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/rs/zerolog/log"
)

// gets clients by switch mac address
func (c *Client) GetClientByPort(switchMac string, port float64) (*NetworkClient, error) {
	clients, err := c.getClientsWithFilters(true, switchMac)
	if err != nil {
		return nil, err
	}
	for _, client := range clients {
		if client.Port == port {
			return &client, nil
		}
	}
	return nil, nil
}

// gets all clients
func (c *Client) GetClients() ([]NetworkClient, error) {
	client, err := c.getClientsWithFilters(false, "")
	if err != nil {
		return nil, err
	}

	return client, nil
}

// gets clients by filters in omada - currentl supports SwitchMac
func (c *Client) getClientsWithFilters(filtersEnabled bool, mac string) ([]NetworkClient, error) {
	loggedIn, err := c.IsLoggedIn()
	if err != nil {
		return nil, err
	}
	if !loggedIn {
		log.Info().Msg(fmt.Sprintf("not logged in, logging in with user: %s", c.Config.Username))
		err := c.Login()
		if err != nil || c.token == "" {
			log.Error().Err(err).Msg("Failed to login")
			return nil, err
		}
	}

	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/clients", c.Config.Host, c.omadaCID, c.SiteId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("currentPage", "1")
	q.Add("currentPageSize", "10000")
	q.Add("filters.active", "true")
	if filtersEnabled {
		q.Add("filters.switchMac=", mac)
	}

	req.URL.RawQuery = q.Encode()

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
	log.Debug().Bytes("data", body).Msg("Received data from clients endpoint")

	clientdata := clientResponse{}
	err = json.Unmarshal(body, &clientdata)

	return clientdata.Result.Data, err
}

type clientResponse struct {
	Result data `json:"result"`
}
type data struct {
	Data []NetworkClient `json:"data"`
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
	WifiMode    float64 `json:"wifiMode"`
	Ssid        string  `json:"ssid"`
}
