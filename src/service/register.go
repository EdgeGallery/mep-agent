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

// register to mep service
package service

import (
	"encoding/json"
	"errors"
	"mep-agent/src/config"
	"mep-agent/src/model"
	"mep-agent/src/util"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	RETRY_TIMES       int = 5
	RETRY_PERIOD      int = 30
	MAX_SERVICE_COUNT int = 50
)

type RegisterData struct {
	token *model.TokenModel
	data  string
	url   string
}

var dataStore []model.ServiceInfoPost

// Registers service to mep
func RegisterToMep(conf model.AppInstanceInfo, wg *sync.WaitGroup) ([]model.ServiceInfoPost, error) {
	log.Info("begin to register service to mep")
	serviceInfos := conf.ServiceInfoPosts
	appInstanceId := util.AppInstanceId

	if len(serviceInfos) > MAX_SERVICE_COUNT {
		log.Error("Failed to register all the services to mep, appInstanceId is " + appInstanceId)
		return nil, errors.New("Registration of service failed, cannot contain more than " +
			strconv.Itoa(MAX_SERVICE_COUNT) + " services in a single request")
	}
	server, errGetServer := config.GetServerUrl()
	if errGetServer != nil {
		return nil, errGetServer
	}

	if util.ValidateUUID(appInstanceId) != nil {
		return nil, errors.New("validate appInstanceId failed")
	}

	url := strings.Replace(server.MepServerRegisterUrl, "${appInstanceId}", appInstanceId, 1)

	for _, serviceInfo := range serviceInfos {
		data, errJsonMarshal := json.Marshal(serviceInfo)
		if errJsonMarshal != nil {
			log.Error("Failed to marshal service info to object")
			continue
		}
		var registerData = RegisterData{data: string(data), url: url, token: &util.MepToken}
		resBody, errPostRequest := PostRegisterRequest(registerData)
		if errPostRequest != nil {
			log.Error("failed to register to mep, appInstanceId is " + appInstanceId +
				", serviceName is " + serviceInfo.SerName)
			wg.Add(1)
			go retryRegister(registerData, appInstanceId, serviceInfo, wg)
		} else {
			log.Info("register to mep success, appInstanceId is " + appInstanceId +
				", serviceName is " + serviceInfo.SerName)
			_, errPostRequest = storeRegisterData(resBody)
			if errPostRequest != nil {
				continue
			}

		}

	}
	log.Info("services registered to mep count ", len(dataStore))
	return dataStore, nil
}

func retryRegister(registerData RegisterData, appInstanceId string, serviceInfo model.ServiceInfoPost,
		wg *sync.WaitGroup) {

	defer wg.Done()

	for i := 1; i < RETRY_TIMES; i++ {
		log.Warn("Failed to register to mep, register will retry 5 times, already register " + strconv.Itoa(i) +
			" times, the next register will begin after " + strconv.Itoa(RETRY_PERIOD*i) + " seconds.")
		time.Sleep(time.Duration(RETRY_PERIOD*i) * time.Second)

		resBody, errPostRequest := PostRegisterRequest(registerData)
		if errPostRequest != nil {
			log.Error("Failed to register to mep, appInstanceId is " + appInstanceId +
				", serviceName is " + serviceInfo.SerName)
		} else {
			log.Info("Register to mep success, appInstanceId is " + appInstanceId +
				", serviceName is " + serviceInfo.SerName)
			_, errJsonUnMarshal := storeRegisterData(resBody)
			if errJsonUnMarshal != nil {
				log.Error("Failed to unmarshal object to service info " + errJsonUnMarshal.Error())
			}
			break
		}
	}
}

func storeRegisterData(resBody string) ([]model.ServiceInfoPost, error) {
	data := model.ServiceInfoPost{}
	errJsonUnMarshal := json.Unmarshal([]byte(resBody), &data)

	if errJsonUnMarshal != nil {
		log.Error("Failed to unmarshal object to service info")
		return nil, errJsonUnMarshal
	}
	dataStore = append(dataStore,data)
	return dataStore, nil
}
