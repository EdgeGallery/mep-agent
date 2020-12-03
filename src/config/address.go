/*
 *  Copyright 2020 Huawei Technologies Co., Ltd.
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

// server address url config
package config

import (
	"errors"
	"mep-agent/src/util"
	"os"
	"strings"
)

type ServerUrl struct {
	MepServerRegisterUrl string
	MepAuthUrl           string
	MepHeartBeatUrl      string
}

const (
	MEP_AUTH_APIGW_URL         string = "https://${MEP_IP}:${MEP_APIGW_PORT}/mepauth/mepauth/v1/token"
	MEP_SER_REGISTER_APIGW_URL string = "https://${MEP_IP}:${MEP_APIGW_PORT}/mepserver/mec_service_mgmt/v1/applications/${appInstanceId}/services"
	MEP_HEART_BEAT_APIGW_URL   string = "https://${MEP_IP}:${MEP_APIGW_PORT}"
	MEP_IP                     string = "${MEP_IP}"
	MEP_APIGW_PORT             string = "${MEP_APIGW_PORT}"
)

// Returns server URL
func GetServerUrl() (ServerUrl, error) {

	var serverUrl ServerUrl
	// validate the env params
	mepIp := os.Getenv("MEP_IP")
	if util.ValidateIp(mepIp) != nil {
		return serverUrl, errors.New("validate MEP_IP failed")
	}
	mepApiGwPort := os.Getenv("MEP_APIGW_PORT")
	if util.ValidateByPattern(util.PORT_PATTERN, mepApiGwPort) != nil {
		return serverUrl, errors.New("validate MEP_APIGW_PORT failed")
	}

	serverUrl.MepServerRegisterUrl = strings.Replace(
		strings.Replace(MEP_SER_REGISTER_APIGW_URL, MEP_IP, mepIp, 1),
		MEP_APIGW_PORT, mepApiGwPort, 1)

	serverUrl.MepAuthUrl = strings.Replace(
		strings.Replace(MEP_AUTH_APIGW_URL, MEP_IP, mepIp, 1),
		MEP_APIGW_PORT, mepApiGwPort, 1)

	serverUrl.MepHeartBeatUrl = strings.Replace(
		strings.Replace(MEP_HEART_BEAT_APIGW_URL, MEP_IP, mepIp, 1),
		MEP_APIGW_PORT, mepApiGwPort, 1)

	return serverUrl, nil
}
