package api

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/urfave/cli/v2"
)

type Client struct {
	Config     *cli.Context
	httpClient *http.Client
	token      string
	omadaCID   string
}

func setuphttpClient(insecure bool) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to init cookiejar")
	}
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	client := &http.Client{Transport: t, Timeout: time.Duration(10) * time.Second, Jar: jar}

	if insecure {
		t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return client, nil
}

func Configure(c *cli.Context) (*Client, error) {
	httpClient, err := setuphttpClient(c.Bool("insecure"))
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
	return client, nil
}

func setHeaders(r *http.Request, crsfToken string) {
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json; charset=UTF-8")
	r.Header.Add("X-Requested-With", "XMLHttpRequest")
	r.Header.Add("User-Agent", "omada_exporter")
	r.Header.Add("accept-encoding", "gzip, deflate, br")
	r.Header.Add("Connection", "keep-alive")
	r.Header.Add("Csrf-Token", crsfToken)
}
