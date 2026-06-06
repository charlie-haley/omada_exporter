package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/charlie-haley/omada_exporter/pkg/config"
	log "github.com/rs/zerolog/log"
)

// apiResponse is a generic Omada API response with the result as raw JSON.
type apiResponse struct {
	ErrorCode int             `json:"errorCode"`
	Msg       string          `json:"msg"`
	Result    json.RawMessage `json:"result"`
}

// checkResponse parses an API response body, checks for API-level errors,
// and returns the raw result payload.
func checkResponse(body []byte) (json.RawMessage, error) {
	var resp apiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("API error %d: %s", resp.ErrorCode, resp.Msg)
	}
	return resp.Result, nil
}

// parseListResult unmarshals a list result from an API response body.
// It handles both direct array format ({"result": [...]}) and the paginated
// format used by newer Omada controllers ({"result": {"data": [...]}}).
func parseListResult[T any](body []byte) ([]T, error) {
	raw, err := checkResponse(body)
	if err != nil {
		return nil, err
	}

	// Try paginated format first: {"data": [...], "totalRows": N, ...}
	var paginated struct {
		Data []T `json:"data"`
	}
	if err := json.Unmarshal(raw, &paginated); err == nil && paginated.Data != nil {
		return paginated.Data, nil
	}

	// Fallback to direct array format: [...]
	var direct []T
	if err := json.Unmarshal(raw, &direct); err != nil {
		return nil, fmt.Errorf("failed to parse list result: %w", err)
	}
	return direct, nil
}

type Client struct {
	Config     *config.Config
	httpClient *http.Client
	token      string
	omadaCID   string
	SiteId     string
}

func setuphttpClient(insecure bool, timeout int) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to init cookiejar")
	}
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	client := &http.Client{Transport: t, Timeout: time.Duration(timeout) * time.Second, Jar: jar}

	if insecure {
		t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return client, nil
}

func Configure(c *config.Config) (*Client, error) {
	httpClient, err := setuphttpClient(c.Insecure, c.Timeout)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Config:     c,
		httpClient: httpClient,
	}
	cid, err := client.getCid()
	if err != nil {
		return nil, err
	}
	client.omadaCID = cid

	sid, err := client.getSiteId(c.Site)
	if err != nil {
		return nil, err
	}
	client.SiteId = *sid

	return client, nil
}

func (c *Client) makeRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "omada_exporter")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Omada-Request-Source", "web-local")

	if c.token != "" {
		req.Header.Set("Csrf-Token", c.token)
	}

	return c.httpClient.Do(req)
}

// makeLoggedInRequest ensures we are logged in before making the request.
// Instead of pre-checking login status on every call, it logs in once on first
// use and only re-authenticates when a request returns an auth error.
//
// Omada session cookies and CSRF tokens expire after a period of inactivity.
// When that happens the controller answers with HTTP 200 but an auth errorCode
// in the JSON body, so we inspect the payload (not the status code), re-login
// once, and replay the request. Without this the exporter would keep sending an
// expired token forever and silently stop producing metrics until restarted.
func (c *Client) makeLoggedInRequest(req *http.Request) (*http.Response, error) {
	// Buffer the request body up front so the request can be replayed after a
	// re-login (http.Request bodies are single-use once sent).
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return nil, err
		}
	}

	// Login once if we don't have a token yet.
	if c.token == "" {
		log.Info().Msg(fmt.Sprintf("not logged in, logging in with user: %s", c.Config.Username))
		if err := c.Login(); err != nil {
			return nil, fmt.Errorf("login failed: %w", err)
		}
	}

	resp, body, authErr, err := c.doDetectingAuth(req, bodyBytes)
	if err != nil {
		return nil, err
	}

	// Session expired or token rejected: re-authenticate once and replay.
	if authErr {
		log.Info().Msg("session expired or token rejected, re-authenticating")
		if err := c.Login(); err != nil {
			return nil, fmt.Errorf("re-login failed: %w", err)
		}
		resp, body, _, err = c.doDetectingAuth(req, bodyBytes)
		if err != nil {
			return nil, err
		}
	}

	// Restore the consumed body so callers can read the response normally.
	resp.Body = io.NopCloser(bytes.NewReader(body))
	return resp, nil
}

// doDetectingAuth sends req (restoring its body from bodyBytes first), reads the
// full response body, and reports whether the controller returned an auth error
// code. The response body is consumed and returned separately so the caller can
// decide whether to retry before handing it back.
func (c *Client) doDetectingAuth(req *http.Request, bodyBytes []byte) (*http.Response, []byte, bool, error) {
	if bodyBytes != nil {
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		req.ContentLength = int64(len(bodyBytes))
	}

	resp, err := c.makeRequest(req)
	if err != nil {
		return nil, nil, false, err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, nil, false, err
	}

	var probe struct {
		ErrorCode int `json:"errorCode"`
	}
	authErr := json.Unmarshal(body, &probe) == nil && isAuthError(probe.ErrorCode)

	return resp, body, authErr, nil
}
