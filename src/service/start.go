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

// Package service start service
package service

import (
	"mep-agent/src/model"
	"mep-agent/src/util"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// ser Empty service
type ser struct {
}

// BeginService Returns ser struct.
func BeginService() *ser {
	return &ser{}
}

// Start service entrance
func (ser *ser) Start(confPath string) {
	var wg = &sync.WaitGroup{}
	// read app_instance_info.yaml file and transform to AppInstanceInfo object
	conf, errGetConf := GetAppInstanceConf(confPath)
	if errGetConf != nil {
		log.Error("Parse app_instance_info.yaml failed.")

		return
	}
	_, errAppInst := util.GetAppInstanceID()
	if errAppInst != nil {
		log.Error("Get app instance id failed.")

		return
	}
	// signed ak and sk, then request the token
	var auth = model.Auth{SecretKey: util.AppConfig["SECRET_KEY"], AccessKey: string(*util.AppConfig["ACCESS_KEY"])}
	errGetMepToken := GetMepToken(auth)
	if errGetMepToken != nil {
		log.Error("Get token failed.")
		return
	}
	util.FirstToken = true

	// register service to mep with token
	// only ServiceInfo not nil
	if conf.ServiceInfoPosts != nil {
		responseBody, errRegisterToMep := RegisterToMep(conf, wg)
		if errRegisterToMep != nil {
			log.Error("Failed to register to mep: " + errRegisterToMep.Error())
			return
		}

		for _, serviceInfo := range responseBody {
			if serviceInfo.LivenessInterval != 0 && serviceInfo.Links.Self.Liveness != "" {
				wg.Add(1)
				heartBeatTicker(serviceInfo)
			} else {
				log.Info("Liveness is not configured or required")
			}
		}
	}

	wg.Wait()
}

func heartBeatTicker(serviceInfo model.ServiceInfoPost) {
	for range time.Tick(time.Duration(serviceInfo.LivenessInterval) * time.Second) {
		go HeartBeatRequestToMep(serviceInfo)
	}
}
