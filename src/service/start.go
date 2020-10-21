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

// start service
package service

import (
	"mep-agent/src/model"
	"mep-agent/src/util"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Ser struct {
}

// Returns ser struct
func BeginService() *Ser {
	return &Ser{}
}

// service entrance
func (ser *Ser) Start(confPath string, ak string, sk *[]byte, wg *sync.WaitGroup) *model.TokenModel {

	// read app_instance_info.yaml file and transform to AppInstanceInfo object
	conf, errGetConf := GetAppInstanceConf(confPath)
	if errGetConf != nil {
		// clear sk
		util.ClearByteArray(*sk)
		log.Error("parse app_instance_info.yaml failed.")
		return nil
	}

	// signed ak and sk, then request the token
	var auth = model.Auth{AccessKey: ak, SecretKey: sk}
	token, errGetMepToken := GetMepToken(auth)
	if errGetMepToken != nil {
		log.Error("get token failed.")
		return nil
	}

	// register service to mep with token
	errRegisterToMep := RegisterToMep(conf, token, wg)
	if errRegisterToMep != nil {
		log.Error("failed to register to mep: " + errRegisterToMep.Error())
		return token
	}
	return token
}
