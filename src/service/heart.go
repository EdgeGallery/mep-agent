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

// heart service
package service

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"mep-agent/src/config"
	"mep-agent/src/model"
	"mep-agent/src/util"
)

type HeartBeatData struct {
	token *model.TokenModel
	data  string
	url   string
}

type ServiceLivenessUpdate struct {
	State string `json:"state"`
}

// Send service heartbeat to MEP
func HeartBeatRequestToMep(serviceInfo model.ServiceInfoPost) {

	heartBeatRequest := ServiceLivenessUpdate{State: "ACTIVE"}
	data, errJsonMarshal := json.Marshal(heartBeatRequest)
	if errJsonMarshal != nil {
		log.Error("Failed to marshal service info to object")
		return
	}

	url := config.ServerUrlConfig.MepHeartBeatUrl + serviceInfo.Links.Self.Liveness
	var heartBeatData = HeartBeatData{data: string(data), url: url, token: &util.MepToken}
	_, errPostRequest := SendHeartBeatRequest(heartBeatData)
	if errPostRequest != nil {
		log.Error("Failed heart beat request to mep, URL is " + url)
	}
}
