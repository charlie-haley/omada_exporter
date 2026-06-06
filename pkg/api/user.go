package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/rs/zerolog/log"
)

// getSiteId resolves a site name to its internal ID.
// It tries the user privileges endpoint first, then falls back to the sites
// listing endpoint for users with allSite privileges or different API versions.
func (c *Client) getSiteId(name string) (*string, error) {
	// Try user privileges endpoint first
	sid, err := c.getSiteIdFromUser(name)
	if err == nil {
		return sid, nil
	}
	log.Debug().Err(err).Msg("Failed to get site ID from user endpoint, trying sites list")

	// Fallback: list all sites
	sid, err = c.getSiteIdFromSitesList(name)
	if err == nil {
		return sid, nil
	}

	return nil, fmt.Errorf("failed to find site with name %q: tried user privileges and sites list endpoints", name)
}

func (c *Client) getSiteIdFromUser(name string) (*string, error) {
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

	raw, err := checkResponse(body)
	if err != nil {
		return nil, err
	}

	user := user{}
	err = json.Unmarshal(raw, &user)
	if err != nil {
		return nil, err
	}

	for _, s := range user.Privilege.Sites {
		if s.Key == name {
			return &s.Value, nil
		}
	}

	return nil, fmt.Errorf("site %q not found in user privileges (%d sites listed)", name, len(user.Privilege.Sites))
}

// getSiteIdFromSitesList fetches the site ID from the sites listing endpoint.
// This works for admin users and newer Omada controller versions where
// the user privileges endpoint may not include the site list.
func (c *Client) getSiteIdFromSitesList(name string) (*string, error) {
	url := fmt.Sprintf("%s/%s/api/v2/sites", c.Config.Host, c.omadaCID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("currentPage", "1")
	q.Add("currentPageSize", "1000")
	req.URL.RawQuery = q.Encode()

	resp, err := c.makeLoggedInRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	raw, err := checkResponse(body)
	if err != nil {
		return nil, err
	}

	// The sites endpoint may return paginated or direct results
	type siteEntry struct {
		Name   string `json:"name"`
		SiteId string `json:"id"`
		Key    string `json:"key"`
	}

	// Try paginated format
	var paginated struct {
		Data []siteEntry `json:"data"`
	}
	if err := json.Unmarshal(raw, &paginated); err == nil && paginated.Data != nil {
		for _, s := range paginated.Data {
			if s.Name == name {
				id := s.SiteId
				if id == "" {
					id = s.Key
				}
				return &id, nil
			}
		}
		return nil, fmt.Errorf("site %q not found in sites list (%d sites)", name, len(paginated.Data))
	}

	// Try direct array format
	var sites []siteEntry
	if err := json.Unmarshal(raw, &sites); err != nil {
		return nil, fmt.Errorf("failed to parse sites list: %w", err)
	}

	for _, s := range sites {
		if s.Name == name {
			id := s.SiteId
			if id == "" {
				id = s.Key
			}
			return &id, nil
		}
	}

	return nil, fmt.Errorf("site %q not found in sites list (%d sites)", name, len(sites))
}

type user struct {
	Privilege privilege `json:"privilege"`
}

type privilege struct {
	AllSite bool   `json:"allSite"`
	Sites   []site `json:"sites"`
}

type site struct {
	Key   string `json:"name"`
	Value string `json:"key"`
}
