package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// there's no nice way of fetching the site ID from the `Viewer` role
// calling the user endpoint seems to return a list of sites for the user
func (c *Client) getSiteId(name string) (*string, error) {
	url := fmt.Sprintf("%s/%s/api/v2/users/current", c.Config.Host, c.omadaCID)
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
