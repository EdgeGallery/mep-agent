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
	mepAuthApigwURL           string = "https://${MEP_IP}:${MEP_APIGW_PORT}/mep/token"
	mepSerRegisterApigwURL    string = "https://${MEP_IP}:${MEP_APIGW_PORT}/mep/mec_service_mgmt/v1/applications/${appInstanceId}/services"
	mepSerQueryByNameApigwURL string = "https://${MEP_IP}:${MEP_APIGW_PORT}/mep/mec_service_mgmt/v1/services?ser_name="
	mepHeartBeatApigwURL      string = "https://${MEP_IP}:${MEP_APIGW_PORT}"
	mepIP                     string = "${MEP_IP}"
	mepApigwPort              string = "${MEP_APIGW_PORT}"
)

// ServerURLConfig server Url Configuration.
var ServerURLConfig ServerURL

// GetServerURL returns server URL.
func GetServerURL() (ServerURL, error) {
	var serverURL ServerURL
	// validate the env params
	mepIPVal := os.Getenv("MEP_IP")
	if util.ValidateDNS(mepIP) != nil {
		return serverURL, errors.New("validate mep ip failed")
	}
	mepAPIGwPort := os.Getenv("MEP_APIGW_PORT")
	if util.ValidateByPattern(util.PortPattern, mepAPIGwPort) != nil {
		return serverURL, errors.New("validate mep api gw failed")
	}

	serverURL.MepServerRegisterURL = strings.Replace(
		strings.Replace(mepSerRegisterApigwURL, mepIP, mepIPVal, 1),
		mepApigwPort, mepAPIGwPort, 1)

	serverURL.MepAuthURL = strings.Replace(
		strings.Replace(mepAuthApigwURL, mepIP, mepIPVal, 1),
		mepApigwPort, mepAPIGwPort, 1)

	serverURL.MepHeartBeatURL = strings.Replace(
		strings.Replace(mepHeartBeatApigwURL, mepIP, mepIPVal, 1),
		mepApigwPort, mepAPIGwPort, 1)

	serverURL.MepServiceDiscoveryURL = strings.Replace(
		strings.Replace(mepSerQueryByNameApigwURL, mepIP, mepIPVal, 1),
		mepApigwPort, mepAPIGwPort, 1)

	return serverURL, nil
}
