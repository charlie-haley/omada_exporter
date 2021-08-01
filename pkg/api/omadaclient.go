package omadaclient

import (
    "io/ioutil"
    "net/http"
    "net/http/cookiejar"
    "time"
    "encoding/json"
    "crypto/tls"
    "fmt"
    "bytes"
    "os"

    log "github.com/sirupsen/logrus"
    "omada_exporter/pkg/api/structs"
)

var client http.Client
var token string

func init(){
    jar, err := cookiejar.New(nil)
    if err != nil {
        log.Error("Failed to init cookiejar")
    }
    token = ""
    client = http.Client{Timeout: time.Duration(5) * time.Second,  Jar: jar}

    insecure = false
    if os.Getenv("OMADA_INSECURE") == "true"{
        insecure = true
    }
    
    if insecure == true{
        http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
    }
}

func setHeaders(r *http.Request) {
    r.Header.Add("Accept", "application/json")
    r.Header.Add("Content-Type", "application/json; charset=UTF-8")
    r.Header.Add("X-Requested-With", "XMLHttpRequest")
    r.Header.Add("User-Agent", "omada_exporter")
    r.Header.Add("accept-encoding", "gzip, deflate")
    r.Header.Add("Connection", "close")
}

func isLoggedIn() (bool, error) {
    loginstatus := structs.LoginStatus{}

    url := fmt.Sprintf("%s/api/v2/loginStatus", os.Getenv("OMADA_HOST"))
    req, err := http.NewRequest("GET", url, nil)
    q := req.URL.Query()
    q.Add("token", token)
    req.URL.RawQuery = q.Encode()

    setHeaders(req)

    res, err := client.Do(req)

    defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)

    err = json.Unmarshal(body, &loginstatus)

    return loginstatus.Result.Login, err
}

func Login() (string, error) {
    logindata := structs.LoginResponse{}

    url := fmt.Sprintf("%s/api/v2/login", os.Getenv("OMADA_HOST"))
    var jsonStr = []byte( fmt.Sprintf(`{"username":"%s","password":"%s"}`, os.Getenv("OMADA_USER"), os.Getenv("OMADA_PASS")))
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

    setHeaders(req)
    res, err := client.Do(req)
    if err != nil {
        return "", err
    }

    defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)

    err = json.Unmarshal(body, &logindata)

    return logindata.Result.Token, err
}

func GetDevices() ([]structs.Device, error) {
    loggedIn, err := isLoggedIn()
    if loggedIn == false {
        log.Info(fmt.Sprintf("Not logged in, logging in with user: %s...", os.Getenv("OMADA_USER")))
        token, err = Login()
    }

    devicedata := structs.DeviceResponse{}
    if err != nil || token == ""  {
        log.Error(fmt.Sprintf("Failed to login: %s", err))
        return devicedata.Result, err
    }

    url := fmt.Sprintf("%s/api/v2/sites/%s/devices", os.Getenv("OMADA_HOST"), os.Getenv("OMADA_SITE"))
    req, err := http.NewRequest("GET", url, nil)
    q := req.URL.Query()
    q.Add("token", token)
    req.URL.RawQuery = q.Encode()

    setHeaders(req)
    resp, err := client.Do(req)

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    err = json.Unmarshal(body, &devicedata)

    return devicedata.Result, err
}

func GetClients() ([]structs.Client, error) {
    token, err := Login()
    clientdata := structs.ClientResponse{}
    if err != nil || token == ""  {
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
    resp, err := client.Do(req)

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    err = json.Unmarshal(body, &clientdata)

    return clientdata.Result.Data, err
}

func GetPorts(switchMac string) ([]structs.Port, error) {
    token, err := Login()
    portdata := structs.PortResponse{}
    if err != nil || token == ""  {
        log.Error(fmt.Sprintf("Failed to login: %s", err))
        return portdata.Result, err
    }

    url := fmt.Sprintf("%s/api/v2/sites/%s/switches/%s/ports", os.Getenv("OMADA_HOST"), os.Getenv("OMADA_SITE"), switchMac)
    req, err := http.NewRequest("GET", url, nil)
    q := req.URL.Query()
    q.Add("token", token)
    req.URL.RawQuery = q.Encode()

    setHeaders(req)
    resp, err := client.Do(req)

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    err = json.Unmarshal(body, &portdata)

    return portdata.Result, err
}