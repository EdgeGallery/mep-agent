/*
 *  Copyright 2020-2021 Huawei Technologies Co., Ltd.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

// Package config server address url config
package config

import (
	"errors"
	"fmt"
	"mep-agent/src/util"
	"os"
	"strings"
)

// ServerURL : List of Urls.
type ServerURL struct {
	MepServerRegisterURL   string
	MepAuthURL             string
	MepHeartBeatURL        string
	MepServiceDiscoveryURL string
}

const (
	mepAuthApigwURL           string = "%s://%s:%s/mep/token"
	mepSerRegisterApigwURL    string = "%s://%s:%s/mep/mec_service_mgmt/v1/applications/${appInstanceId}/services"
	mepSerQueryByNameApigwURL string = "%s://%s:%s/mep/mec_service_mgmt/v1/services?ser_name="
	mepHeartBeatApigwURL      string = "%s://%s:%s"
)

// ServerURLConfig server Url Configuration.
var ServerURLConfig ServerURL

// GetServerURL returns server URL.
func GetServerURL() (ServerURL, error) {
	var serverURL ServerURL
	// validate the env params
	mepIPVal := os.Getenv("MEP_IP")
	if util.ValidateDNS(mepIPVal) != nil {
		return serverURL, errors.New("validate mep ip failed")
	}
	mepAPIGwPort := os.Getenv("MEP_APIGW_PORT")
	if util.ValidateByPattern(util.PortPattern, mepAPIGwPort) != nil {
		return serverURL, errors.New("validate mep api gw failed")
	}

	egProtocol := "https"
	if strings.EqualFold(os.Getenv("EG_PROTOCOL"), "http") {
		egProtocol = "http"
	}

	serverURL.MepServerRegisterURL = fmt.Sprintf(mepSerRegisterApigwURL, egProtocol, mepIPVal, mepAPIGwPort)

	serverURL.MepServerRegisterURL = fmt.Sprintf(mepAuthApigwURL, egProtocol, mepIPVal, mepAPIGwPort)

	serverURL.MepHeartBeatURL = fmt.Sprintf(mepHeartBeatApigwURL, egProtocol, mepIPVal, mepAPIGwPort)

	serverURL.MepServiceDiscoveryURL = fmt.Sprintf(mepSerQueryByNameApigwURL, egProtocol, mepIPVal, mepAPIGwPort)

	return serverURL, nil
}
