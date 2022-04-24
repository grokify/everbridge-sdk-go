package everbridgesdk

import (
	"encoding/json"
	"os"
	"path"
	"regexp"

	"github.com/grokify/mogo/os/osutil"
)

type EndpointContactsResponse struct {
	Message      string `json:"message"`
	FirstPageURI string `json:"firstPageUri"`
	NextPageURI  string `json:"nextPageUri"`
	LastPageURI  string `json:"lastPageUri"`
	Page         EndpointContactsResponseObjectPage
}

type EndpointContactsResponseObjectPage struct {
	Data           []EndpointContactsResponseObjectContactWrapper `json:"data"`
	PageSize       int64                                          `json:"pageSize"`
	Start          int64                                          `json:"start"`
	TotalCount     int64                                          `json:"totalCount"`
	TotalPageCount int64                                          `json:"totalPageCount"`
	CurrentPageNo  int64                                          `json:"currentPageNo"`
}

type EndpointContactsResponseObjectContactWrapper struct {
	LastModifiedTime  int64                                            `json:"lastModifiedTime"`
	OrganizationID    int64                                            `json:"organizationId"`
	CreatedDate       int64                                            `json:"createdDate"`
	Groups            []int64                                          `json:"groups"`
	CreatedName       string                                           `json:"createdName"`
	LastName          string                                           `json:"lastName"`
	Status            string                                           `json:"status"`
	Country           string                                           `json:"country"`
	RecordTypeID      int64                                            `json:"recordTypeId"`
	LastModifiedName  string                                           `json:"lastModifiedName"`
	AccountID         int64                                            `json:"accountId"`
	ExternalID        string                                           `json:"externalId"`
	ID                int64                                            `json:"id"`
	FirstName         string                                           `json:"firstName"`
	UploadProcessing  bool                                             `json:"uploadProcessing"`
	ResourceBundleID  int64                                            `json:"resourceBundleId"`
	CreatedID         int64                                            `json:"createdId"`
	LastModifiedID    int64                                            `json:"lastModifiedId"`
	LastModifiedDate  int64                                            `json:"lastModifiedDate"`
	ContactAttributes []EndpointContactsResponseObjectContactAttribute `json:"contactAttributes"`
	Paths             []EndpointContactsResponseObjectContactPath      `json:"paths"`
}

type EndpointContactsResponseObjectContactAttribute struct {
	Values    []string `json:"Values"`
	OrgAttrID int64    `json:"OrgAttrId"`
	Name      string   `json:"Name"`
}

type EndpointContactsResponseObjectContactPath struct {
	WaitTime    int64  `json:"WaitTime"`
	Status      string `json:"Status"`
	PathID      int64  `json:"PathId"`
	CountryCode string `json:"CountryCode"`
	Value       string `json:"Value"`
}

func GetEprContactsForBody(content []byte) (EndpointContactsResponse, error) {
	eprContacts := EndpointContactsResponse{}
	err := json.Unmarshal(content, &eprContacts)
	return eprContacts, err
}

func GetEpoContactsForBody(content []byte) ([]EndpointContactsResponseObjectContactWrapper, error) {
	eprContacts, err := GetEprContactsForBody(content)
	if err != nil {
		return []EndpointContactsResponseObjectContactWrapper{}, err
	}
	epoContacts := eprContacts.Page.Data
	return epoContacts, nil
}

func GetEpoContactsForPath(filepath string) ([]EndpointContactsResponseObjectContactWrapper, error) {
	bytContents, err := os.ReadFile(filepath)
	if err != nil {
		return []EndpointContactsResponseObjectContactWrapper{}, err
	}
	return GetEpoContactsForBody(bytContents)
}

func GetEpoContactsForDir(dir string) ([]EndpointContactsResponseObjectContactWrapper, error) {
	epoContacts := []EndpointContactsResponseObjectContactWrapper{}
	re1 := regexp.MustCompile(`^evb_contacts_org-id-[0-9]+_page-num-[0-9]+\.json$`)
	finfos, err := osutil.ReadDirMore(dir, re1, false, true, false)
	if err != nil {
		return epoContacts, nil
	}
	for _, fi := range finfos {
		filepath := path.Join(dir, fi.Name())
		epoContactsForPage, err := GetEpoContactsForPath(filepath)
		if err != nil {
			return epoContacts, err
		}
		epoContacts = append(epoContacts, epoContactsForPage...)
	}
	return epoContacts, nil
}
