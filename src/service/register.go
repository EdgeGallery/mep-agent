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

// Package service register to mep service
package service

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"mep-agent/src/config"
	"mep-agent/src/model"
	"mep-agent/src/util"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	retryTimes      int = 5
	retryPeriod     int = 30
	maxServiceCount int = 50
)

type registerData struct {
	token *model.TokenModel
	data  string
	url   string
}

// dataStore : Store service info posts.
var dataStore []model.ServiceInfoPost

// RegisterToMep Registers service to mep.
func RegisterToMep(conf model.AppInstanceInfo, wg *sync.WaitGroup) ([]model.ServiceInfoPost, error) {
	log.Info("Begin to register service to mep.")
	serviceInfos := conf.ServiceInfoPosts
	appInstanceID := util.AppInstanceID

	if len(serviceInfos) > maxServiceCount {
		log.Error("Failed to register all the services to mep, app instance id is " + appInstanceID)
		return nil, errors.New("registration of service failed, cannot contain more than " +
			strconv.Itoa(maxServiceCount) + " services in a single request")
	}

	if util.ValidateUUID(appInstanceID) != nil {
		return nil, errors.New("validate appInstanceId failed")
	}

	url := strings.Replace(config.ServerURLConfig.MepServerRegisterURL, "${appInstanceId}", appInstanceID, 1)

	for _, serviceInfo := range serviceInfos {
		data, errJSONMarshal := json.Marshal(serviceInfo)
		if errJSONMarshal != nil {
			log.Error("Failed to marshal service info to object")

			continue
		}
		var registerInfo = registerData{data: string(data), url: url, token: &util.MepToken}
		resBody, errPostRequest := postRegisterRequest(registerInfo)
		if errPostRequest != nil {
			log.Error("Failed to register to mep, app instance id is " + appInstanceID +
				", serviceName is " + serviceInfo.SerName)
			wg.Add(1)
			go retryRegister(registerInfo, appInstanceID, serviceInfo, wg)
		} else {
			log.Info("Register to mep success, app instance id is " + appInstanceID +
				", serviceName is " + serviceInfo.SerName)
			errPostRequest = storeRegisterData(resBody)
			if errPostRequest != nil {
				continue
			}
		}
	}
	log.Info("Services registered to mep count ", len(dataStore))

	return dataStore, nil
}

func retryRegister(registerData registerData, appInstanceID string, serviceInfo model.ServiceInfoPost,
		wg *sync.WaitGroup) {

	defer wg.Done()

	for i := 1; i < retryTimes; i++ {
		log.Warn("Failed to register to mep, register will retry 5 times, already register " + strconv.Itoa(i) +
			" times, the next register will begin after " + strconv.Itoa(retryPeriod*i) + " seconds.")
		time.Sleep(time.Duration(retryPeriod*i) * time.Second)

		resBody, errPostRequest := postRegisterRequest(registerData)
		if errPostRequest != nil {
			log.Error("Failed to register to mep, app instance id is " + appInstanceID +
				", serviceName is " + serviceInfo.SerName)
		} else {
			log.Info("Register to mep success, app instance id is " + appInstanceID +
				", serviceName is " + serviceInfo.SerName)
			errJSONUnMarshal := storeRegisterData(resBody)
			if errJSONUnMarshal != nil {
				log.Error("Failed to unmarshal object to service info " + errJSONUnMarshal.Error())
			}

			break
		}
	}
}

func storeRegisterData(resBody string) error {
	var data = model.ServiceInfoPost{}
	errJSONUnMarshal := json.Unmarshal([]byte(resBody), &data)

	if errJSONUnMarshal != nil {
		log.Error("Failed to unmarshal object to service info")

		return errJSONUnMarshal
	}
	dataStore = append(dataStore,data)

	return nil
}
