package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (c *Client) GetPorts(switchMac string) ([]Port, error) {
	loggedIn, err := c.IsLoggedIn()
	if err != nil {
		log.Error(err)
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

	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/switches/%s/ports", c.Config.String("host"), c.omadaCID, c.Config.String("site"), switchMac)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

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
}
