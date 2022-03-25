package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// there's no nice way of fetching the site ID from the `Viewer` role
// calling the timeout endpoint seems to return a siteId(?!)
func (c *Client) getSiteId(name string) (*string, error) {
	loggedIn, err := c.IsLoggedIn()
	if err != nil {
		return nil, err
	}
	if !loggedIn {
		log.Info(fmt.Errorf("not logged in, logging in with user: %s", c.Config.String("username")))
		err := c.Login()
		if err != nil || c.token == "" {
			return nil, fmt.Errorf("failed to login: %s", err)
		}
	}

	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/setting/firewall/timeout", c.Config.String("host"), c.omadaCID, name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	setHeaders(req, c.token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	sites := timeoutResponse{}
	err = json.Unmarshal(body, &sites)
	if err != nil {
		return nil, err
	}

	return &sites.Result.Id, nil
}

type timeoutResponse struct {
	Result timeout `json:"result"`
}
type timeout struct {
	Id string `json:"siteId"`
}
