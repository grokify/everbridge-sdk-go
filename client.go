package everbridgesdk

import (
	"encoding/base64"
	"net/http"
	"strings"
)

type Client struct {
	ClientCore            ClientCore
	EndpointContacts      EndpointContacts
	EndpointOrganizations EndpointOrganizations
}

func NewClient(username string, password string, dataDir string) Client {
	cl := Client{}
	cl.ClientCore = NewClientCore(username, password, dataDir)
	cl.EndpointContacts = EndpointContacts{ClientCore: cl.ClientCore}
	cl.EndpointOrganizations = EndpointOrganizations{ClientCore: cl.ClientCore}
	return cl
}

type ClientCore struct {
	Protocol      string
	Hostname      string
	BaseURL       string
	BasicAuthz    string
	NetHTTPClient *http.Client
	DataDir       string
}

func NewClientCore(username string, password string, dataDir string) ClientCore {
	cc := ClientCore{}
	cc.Protocol = "https"
	cc.Hostname = "api.everbridge.net"
	cc.BaseURL = "https://api.everbridge.net/rest"
	cc.DataDir = dataDir
	cc.LoadClient()
	cc.LoadCredentials(username, password)
	return cc
}

func (cc *ClientCore) LoadClient() {
	client := &http.Client{}
	cc.NetHTTPClient = client
}

func (cc *ClientCore) LoadCredentials(username string, password string) {
	authorization := strings.Join([]string{username, ":", password}, "")
	cc.BasicAuthz = base64.StdEncoding.EncodeToString([]byte(authorization))
}

func (cc *ClientCore) NewRequestForMethodAndURL(method string, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err == nil {
		req.Header.Add("Authorization", strings.Join([]string{"Basic", cc.BasicAuthz}, " "))
	}
	return req, err
}
