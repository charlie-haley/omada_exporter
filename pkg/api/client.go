package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (c *Client) GetClients() ([]client, error) {
	clientdata := clientResponse{}

	loggedIn, err := c.IsLoggedIn()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if !loggedIn {
		log.Info(fmt.Sprintf("Not logged in, logging in with user: %s...", c.Config.String("username")))
		err := c.Login()
		if err != nil || c.token == "" {
			log.Error(fmt.Sprintf("Failed to login: %s", err))
			return clientdata.Result.Data, err
		}
	}

	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/clients", c.Config.String("host"), c.omadaCID, c.Config.String("site"))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("currentPage", "1")
	q.Add("currentPageSize", "10000")
	q.Add("filters.active", "true")
	req.URL.RawQuery = q.Encode()

	setHeaders(req, c.token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = json.Unmarshal(body, &clientdata)

	return clientdata.Result.Data, err
}

type clientResponse struct {
	Result data `json:"result"`
}
type data struct {
	Data []client `json:"data"`
}
type client struct {
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
