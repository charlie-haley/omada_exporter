package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// there's no nice way of fetching the site ID from the `Viewer` role
// calling the user endpoint seems to return a list of sites for the user
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

	url := fmt.Sprintf("%s/%s/api/v2/users/current", c.Config.String("host"), c.omadaCID)
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	user := userResponse{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	for _, s := range user.Result.Privilege.Sites {
		if s.Key == name {
			return &s.Value, nil
		}
	}

	return nil, fmt.Errorf("failed to find site with name %s", name)
}

type userResponse struct {
	Result user `json:"result"`
}

type user struct {
	Privilege privilege `json:"privilege"`
}

type privilege struct {
	Sites []site `json:"sites"`
}

type site struct {
	Key   string `json:"name"`
	Value string `json:"key"`
}
