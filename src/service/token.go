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

// get token service
package service

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"mep-agent/src/config"
	"mep-agent/src/model"
	"mep-agent/src/util"
	"time"
	"unsafe"
)

// Request token from mep_auth
func GetMepToken(auth model.Auth) (error) {
	// get request url
	server, errGetServer := config.GetServerUrl()
	if errGetServer != nil {
		return errGetServer
	}

	// construct http request and send
	resp, errPostRequest := PostTokenRequest("", server.MepAuthUrl, auth)
	if errPostRequest != nil {
		return errPostRequest
	}

	// unmarshal resp to object
	errJson := json.Unmarshal([]byte(resp), &util.MepToken)

	// clear resp
	util.ClearByteArray(*(*[]byte)(unsafe.Pointer(&resp)))
	if errJson != nil {
		return errJson
	}

	log.Info("get token success.")

	//start timer to refresh token
	go startRefreshTimer()
	return nil
}

// This function will be only call by GetMepToken
func startRefreshTimer() {
	defer func() {
		if err1 := recover(); err1 != nil {
			log.Error("panic handled:", err1)
		}
	}()
	if util.RefresTimer != nil {
		ok := util.RefresTimer.Stop()
		if ok {
			log.Info("timer stopped")
		} else {
			log.Info("timer not yet started")
		}
	}
	//start timer with latest token expiry value - buffertime
	util.RefresTimer = time.NewTimer(time.Duration(util.MepToken.ExpiresIn - util.RefreshTimeBuffer) * time.Second)
	log.Info("Refresh timer started")
	go func() {
		_, ok := <-util.RefresTimer.C
		if !ok {
			log.Error("Timer C channel closed")
		}
		var auth = model.Auth{SecretKey: util.AppConfig["SECRET_KEY"], AccessKey: string(*util.AppConfig["ACCESS_KEY"])}
		errGetMepToken := GetMepToken(auth)
		if errGetMepToken != nil {
			log.Error("Get token failed.")
			return
		}
	}()
}

