package api

type deviceResponse struct {
	Result []device `json:"result"`
}
type device struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Mac         string  `json:"mac"`
	Model       string  `json:"model"`
	Version     string  `json:"version"`
	Ip          string  `json:"ip"`
	CpuUtil     float64 `json:"cpuUtil"`
	MemUtil     float64 `json:"memUtil"`
	Uptime      float64 `json:"uptimeLong"`
	NeedUpgrade bool    `json:"needUpgrade"`
	TxRate      float64 `json:"txRate"`
	RxRate      float64 `json:"rxRate"`
	PoeRemain   float64 `json:"poeRemain"`
}

type clientResponse struct {
	Result data `json:"result"`
}
type data struct {
	Data []client `json:"data"`
}
type client struct {
	Name        string  `json:"name"`
	HostName    string  `json:"hostName"`
	Mac         string  `json:"mac"`
	Port        float64 `json:"port"`
	Ip          string  `json:"ip"`
	VlanId      float64 `json:"vid"`
	ApName      string  `json:"apName"`
	Wireless    bool    `json:"wireless"`
	SwitchMac   string  `json:"switchMac"`
	Vendor      string  `json:"vendor"`
	Activity    float64 `json:"activity"`
	SignalLevel float64 `json:"signalLevel"`
	WifiMode    float64 `json:"wifiMode"`
	Ssid        string  `json:"ssid"`
}

type portResponse struct {
	Result []Port `json:"result"`
}
type Port struct {
	Id          string     `json:"id"`
	SwitchId    string     `json:"switchId"`
	SwitchMac   string     `json:"switchMac"`
	Name        string     `json:"name"`
	PortStatus  portStatus `json:"portStatus"`
	Port        float64    `json:"port"`
	ProfileName string     `json:"profileName"`
}
type portStatus struct {
	Port       float64 `json:"id"`
	LinkStatus float64 `json:"linkStatus"`
	LinkSpeed  float64 `json:"linkSpeed"`
	PoePower   float64 `json:"poePower"`
	Poe        bool    `json:"poe"`
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
