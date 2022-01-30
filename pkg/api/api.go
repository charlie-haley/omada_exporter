package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var httpclient http.Client
var token string

func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Error("Failed to init cookiejar")
	}
	token = ""
	httpclient = http.Client{Timeout: time.Duration(5) * time.Second, Jar: jar}

	insecure := false
	if os.Getenv("OMADA_INSECURE") == "true" {
		insecure = true
	}

	if insecure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
}

func setHeaders(r *http.Request) {
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json; charset=UTF-8")
	r.Header.Add("X-Requested-With", "XMLHttpRequest")
	r.Header.Add("User-Agent", "omada_exporter")
	r.Header.Add("accept-encoding", "gzip, deflate, br")
	r.Header.Add("Connection", "keep-alive")
}

func isLoggedIn() (bool, error) {
	loginstatus := loginStatus{}

	url := fmt.Sprintf("%s/api/v2/loginStatus", os.Getenv("OMADA_HOST"))
	log.Error(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	q := req.URL.Query()
	q.Add("token", token)
	req.URL.RawQuery = q.Encode()

	setHeaders(req)

	res, err := httpclient.Do(req)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	log.Error(string(body))
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(body, &loginstatus)
	if loginstatus.ErrorCode != 0 {
		return false, errors.New(fmt.Sprintf("Invalid error code returned from API. Response Body: %s", string(body)))
	}

	return loginstatus.Result.Login, err
}

func Login() (string, error) {
	logindata := loginResponse{}

	url := fmt.Sprintf("%s/api/v2/login", os.Getenv("OMADA_HOST"))
	jsonStr := []byte(fmt.Sprintf(`{"username":"%s","password":"%s"}`, os.Getenv("OMADA_USER"), os.Getenv("OMADA_PASS")))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}

	setHeaders(req)
	res, err := httpclient.Do(req)
	if err != nil {
		log.Error(err)
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	err = json.Unmarshal(body, &logindata)
	// if err != nil {

	// }

	return logindata.Result.Token, err
}

func GetDevices() ([]device, error) {
	loggedIn, err := isLoggedIn()
	if !loggedIn {
		log.Info(fmt.Sprintf("Not logged in, logging in with user: %s...", os.Getenv("OMADA_USER")))
		token, err = Login()
	}

	devicedata := deviceResponse{}
	if err != nil || token == "" {
		log.Error(fmt.Sprintf("Failed to login: %s", err))
		return devicedata.Result, err
	}

	url := fmt.Sprintf("%s/api/v2/sites/%s/devices", os.Getenv("OMADA_HOST"), os.Getenv("OMADA_SITE"))
	req, err := http.NewRequest("GET", url, nil)

	q := req.URL.Query()
	q.Add("token", token)
	req.URL.RawQuery = q.Encode()

	setHeaders(req)
	resp, err := httpclient.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &devicedata)

	return devicedata.Result, err
}

func GetClients() ([]client, error) {
	token, err := Login()
	clientdata := clientResponse{}
	if err != nil || token == "" {
		log.Error(fmt.Sprintf("Failed to login: %s", err))
		return clientdata.Result.Data, err
	}

	url := fmt.Sprintf("%s/api/v2/sites/%s/clients", os.Getenv("OMADA_HOST"), os.Getenv("OMADA_SITE"))
	req, err := http.NewRequest("GET", url, nil)
	q := req.URL.Query()
	q.Add("token", token)
	q.Add("currentPage", "1")
	q.Add("currentPageSize", "10000")
	q.Add("filters.active", "true")
	req.URL.RawQuery = q.Encode()

	setHeaders(req)
	resp, err := httpclient.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &clientdata)

	return clientdata.Result.Data, err
}

func GetPorts(switchMac string) ([]Port, error) {
	token, err := Login()
	portdata := portResponse{}
	if err != nil || token == "" {
		log.Error(fmt.Sprintf("Failed to login: %s", err))
		return portdata.Result, err
	}

	url := fmt.Sprintf("%s/api/v2/sites/%s/switches/%s/ports", os.Getenv("OMADA_HOST"), os.Getenv("OMADA_SITE"), switchMac)
	req, err := http.NewRequest("GET", url, nil)
	q := req.URL.Query()
	q.Add("token", token)
	req.URL.RawQuery = q.Encode()

	setHeaders(req)
	resp, err := httpclient.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &portdata)

	return portdata.Result, err
}
