package everbridgesdk

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/grokify/mogo/io/ioutilmore"
	"github.com/grokify/mogo/os/osutil"
)

type EndpointContacts struct {
	ClientCore ClientCore
}

func NewEndpointsContacts() EndpointContacts {
	ep := EndpointContacts{}
	return ep
}

func (ep *EndpointContacts) Get(organizationID int64, pageNumber int64) (*http.Response, error) {
	orgIDString := strconv.FormatInt(organizationID, 10)
	url := strings.Join([]string{ep.ClientCore.BaseURL, "contacts", orgIDString}, "/")
	url = url + "?sortBy=externalId&pageNumber=" + strconv.FormatInt(pageNumber, 10)

	req, err := ep.ClientCore.NewRequestForMethodAndURL(http.MethodGet, url)
	if err != nil {
		return &http.Response{}, err
	}

	return ep.ClientCore.NetHTTPClient.Do(req)
}

func (ep *EndpointContacts) GetStoreAll(organizationID int64, dir string) error {
	isDir, err := osutil.IsDir(dir)
	if err != nil {
		return err
	} else if isDir == false {
		str := fmt.Sprintf("500: Path Is Not Directory [%v]", dir)
		err = errors.New(str)
		return err
	}
	err = ioutilmore.RemoveAllChildren(dir)
	if err != nil {
		return err
	}
	contents, err := ep.getStoreOrgPage(organizationID, int64(1), dir)
	if err != nil {
		return err
	}
	epo := GetEprContactsForBody(contents)
	if epo.Page.TotalPageCount > 1 {
		for i := int64(2); i <= epo.Page.TotalPageCount; i++ {
			_, err := ep.getStoreOrgPage(organizationID, i, dir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (ep *EndpointContacts) getStoreOrgPage(organizationId int64, pageNumber int64, dir string) ([]byte, error) {
	res, err := ep.Get(organizationId, pageNumber)
	if err != nil {
		return []byte{}, err
	}
	filename := ep.GetFilenameForOrgIdAndPageNum(organizationId, pageNumber)
	filepath := path.Join(dir, filename)
	defer res.Body.Close()
	contents, err := io.ReadAll(res.Body)
	if err != nil {
		return contents, err
	}
	err = os.WriteFile(filepath, contents, 0600)
	return contents, err
}

func (ep *EndpointContacts) GetFilenameForOrgIdAndPageNum(organizationId int64, pageNumber int64) string {
	sOrgId := strconv.FormatInt(organizationId, 10)
	sPgNum := strconv.FormatInt(pageNumber, 10)
	parts := []string{"evb_contacts_org-id-", sOrgId, "_page-num-", sPgNum, ".json"}
	filename := strings.Join(parts, "")
	return filename
}
