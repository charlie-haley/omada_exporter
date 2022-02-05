package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (c *Client) GetDevices() ([]device, error) {
	devicedata := deviceResponse{}

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
			return devicedata.Result, err
		}
	}

	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/devices", c.Config.String("host"), c.omadaCID, c.Config.String("site"))
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

	err = json.Unmarshal(body, &devicedata)

	return devicedata.Result, err
}

type deviceResponse struct {
	Result []device `json:"result"`
}
type device struct {
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
}
