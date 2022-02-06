package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (c *Client) GetController() (*Controller, error) {
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

	url := fmt.Sprintf("%s/%s/api/v2/maintenance/controllerStatus?", c.Config.String("host"), c.omadaCID)
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
