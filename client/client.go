// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// TODO(cAdvisor): Package comment.
package cadvisor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/cadvisor/info"
)

// Client represents the base URL for a cAdvisor client.
type Client struct {
	baseUrl string
}

// NewClient returns a new client with the specified base URL.
// TODO(cAdvisor): Currently the error result is always nil.
func NewClient(URL string) (*Client, error) {
	return &Client{
		baseUrl: strings.Join([]string{
			URL,
			"api/v1.0",
		}, "/"),
	}, nil
}

func (self *Client) machineInfoUrl() string {
	return strings.Join([]string{self.baseUrl, "machine"}, "/")
}

// MachineInfo returns the JSON machine information for this client.
// A non-nil error result indicates a problem with obtaining
// the JSON machine information data.
func (self *Client) MachineInfo() (minfo *info.MachineInfo, err error) {
	u := self.machineInfoUrl()
	ret := new(info.MachineInfo)
	if err = self.httpGetJsonData(ret, nil, u, "machine info"); err != nil {
		return
	}
	minfo = ret
	return
}

func (self *Client) containerInfoUrl(name string) string {
	if name[0] == '/' {
		name = name[1:]
	}
	return strings.Join([]string{self.baseUrl, "containers", name}, "/")
}

func (self *Client) httpGetJsonData(data, postData interface{}, url, infoName string) error {
	var resp *http.Response
	var err error

	if postData != nil {
		data, err := json.Marshal(postData)
		if err != nil {
			return fmt.Errorf("unable to marshal data: %v", err)
		}
		resp, err = http.Post(url, "application/json", bytes.NewBuffer(data))
	} else {
		resp, err = http.Get(url)
	}
	if err != nil {
		err = fmt.Errorf("unable to get %v: %v", infoName, err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("unable to read all %v: %v", infoName, err)
		return err
	}
	if err = json.Unmarshal(body, data); err != nil {
		err = fmt.Errorf("unable to unmarshal %v (%v): %v", infoName, string(body), err)
		return err
	}
	return nil
}

// ContainerInfo returns the JSON container information for the specified
// container and request.
func (self *Client) ContainerInfo(name string,
	query *info.ContainerInfoRequest) (cinfo *info.ContainerInfo, err error) {
	u := self.containerInfoUrl(name)
	ret := new(info.ContainerInfo)
	if err = self.httpGetJsonData(ret, query, u, fmt.Sprintf("container info for %v", name)); err != nil {
		return
	}
	cinfo = ret
	return
}
