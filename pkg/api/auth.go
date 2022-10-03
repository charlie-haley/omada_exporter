package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) IsLoggedIn() (bool, error) {
	loginstatus := loginStatus{}

	url := fmt.Sprintf("%s/%s/api/v2/loginStatus", c.Config.String("host"), c.omadaCID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	setHeaders(req, c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(body, &loginstatus)
	if loginstatus.ErrorCode == -1200 {
		return false, nil
	}
	if loginstatus.ErrorCode != 0 {
		return false, fmt.Errorf("invalid error code returned from API. Response Body: %s", string(body))
	}

	return loginstatus.Result.Login, err
}

// one of the "quirks" of the omada API - it requires a CID to be part of the path
// fetching this from the path after redirect seems like the best way
func (c *Client) getCid() (string, error) {
	host := c.Config.String("host")
	req, err := http.NewRequest("GET", host, nil)
	if err != nil {
		return "", err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	location := res.Request.URL.Path

	// remove `/login` path from url
	location = strings.Replace(location, "/login", "", -1)

	// trim first rune from string - leading `/`
	return location[1:], nil
}

func (c *Client) Login() error {
	logindata := loginResponse{}

	url := fmt.Sprintf("%s/%s/api/v2/login", c.Config.String("host"), c.omadaCID)
	jsonStr := []byte(fmt.Sprintf(`{"username":"%s","password":"%s"}`, c.Config.String("username"), c.Config.String("password")))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	setHeaders(req, "")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &logindata)
	if err != nil {
		return err
	}

	c.token = logindata.Result.Token
	return nil
}

type loginResponse struct {
	Result loginResult `json:"result"`
}
type loginResult struct {
	Token string `json:"token"`
}
type loginStatus struct {
	ErrorCode int            `json:"errorCode"`
	Result    loggedInResult `json:"result"`
}
type loggedInResult struct {
	Login bool `json:"login"`
}
