package everbridgesdk

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type EndpointOrganizations struct {
	ClientCore ClientCore
}

func (epo *EndpointOrganizations) Get() (*http.Response, error) {
	url := strings.Join([]string{epo.ClientCore.BaseURL, "organizations"}, "/")

	req, err := epo.ClientCore.NewRequestForMethodAndURL("GET", url)
	if err != nil {
		return &http.Response{}, err
	}

	return epo.ClientCore.NetHTTPClient.Do(req)
}

func (epo *EndpointOrganizations) GetOrganizationIds() ([]int64, error) {
	res, err := epo.Get()
	if err != nil {
		return []int64{}, err
	}
	defer res.Body.Close()
	contents, err := io.ReadAll(res.Body)
	if err != nil {
		return []int64{}, err
	}

	root := map[string]interface{}{}
	err = json.Unmarshal(contents, &root)
	if err != nil {
		return []int64{}, err
	}

	data := root["page"].(map[string]interface{})["data"].([]interface{})

	organizationIds := []int64{}
	for _, org := range data {
		id := int64(org.(map[string]interface{})["organizationId"].(float64))
		organizationIds = append(organizationIds, id)
	}
	return organizationIds, nil
}
