package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) IsLoggedIn() (bool, error) {
	loginstatus := loginStatus{}

	url := fmt.Sprintf("%s/%s/api/v2/loginStatus", c.Config.Host, c.omadaCID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	res, err := c.makeRequest(req)
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
func (c *Client) getCid() (string, error) {
	url := fmt.Sprintf("%s/api/info", c.Config.Host)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	res, err := c.makeRequest(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var infoResponse struct {
		ErrorCode int    `json:"errorCode"`
		Msg       string `json:"msg"`
		Result    struct {
			OmadaCID string `json:"omadacId"`
		}
	}
	err = json.NewDecoder(res.Body).Decode(&infoResponse)
	if err != nil {
		return "", err
	}

	if infoResponse.Result.OmadaCID == "" {
		return "", fmt.Errorf("no CID found in response")
	}

	return infoResponse.Result.OmadaCID, nil
}

func (c *Client) Login() error {
	logindata := loginResponse{}

	url := fmt.Sprintf("%s/%s/api/v2/login", c.Config.Host, c.omadaCID)
	jsonStr := []byte(fmt.Sprintf(`{"username":"%s","password":"%s"}`, c.Config.Username, c.Config.Password))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	res, err := c.makeRequest(req)
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
