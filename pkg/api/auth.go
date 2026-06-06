package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// isAuthError returns true if the error code indicates an authentication
// or session problem that should trigger a re-login rather than a fatal error.
func isAuthError(code int) bool {
	switch code {
	case -1200, // not authenticated (classic)
		-1201,  // token expired
		-30109, // invalid token
		-40011, // token mismatch
		-40012, // CSRF token mismatch
		-1:     // generic error (often auth-related on newer firmware)
		return true
	}
	return false
}

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
	if err != nil {
		return false, err
	}

	// Treat any authentication-related error code as "not logged in"
	// so that the caller can re-attempt login.
	if isAuthError(loginstatus.ErrorCode) {
		return false, nil
	}
	if loginstatus.ErrorCode != 0 {
		return false, fmt.Errorf("invalid error code %d returned from loginStatus API. Response Body: %s", loginstatus.ErrorCode, string(body))
	}

	return loginstatus.Result.Login, nil
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

	// Check for API-level errors first
	var apiResp struct {
		ErrorCode int    `json:"errorCode"`
		Msg       string `json:"msg"`
	}
	if err := json.Unmarshal(body, &apiResp); err == nil && apiResp.ErrorCode != 0 {
		return fmt.Errorf("login failed with error %d: %s", apiResp.ErrorCode, apiResp.Msg)
	}

	logindata := loginResponse{}
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
