/*
 * Copyright 2020 Huawei Technologies Co., Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Clear util
package util

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"mep-agent/src/model"
	"os"
	"time"
)
//App configuration properties
type AppConfigProperties map[string]*[]byte

//Application configurations
var AppConfig = AppConfigProperties{}

// Token
var MepToken = model.TokenModel{}

//Mark initial token has been received
var FirstToken      = false

//App instance ID from configuration
var AppInstanceId    string

//Refresh token timer
var RefresTimer      *time.Timer

//Timer buffer 5 sec
const RefreshTimeBuffer = 5

// Clears byte array
func ClearByteArray(data []byte) {
	if data == nil {
		return
	}
	for i := 0; i < len(data); i++ {
		data[i] = 0
	}
}

// clear [string, *[]byte] map, called only in error case
func ClearMap(){
	for _, element := range AppConfig {
		ClearByteArray(*element)
	}
}

//read and clearing the variable from the environment
func ReadTokenFromEnvironment() error {
	if len(os.Getenv("AK")) == 0 || len(os.Getenv("SK")) == 0 {
		err := errors.New("AK and SK keys should be set in env variable")
		log.Error("Keys should not be empty")
		return err
	}
	AK := []byte(os.Getenv("AK"))
	AppConfig["ACCESS_KEY"] = &AK
	SK := []byte(os.Getenv("SK"))
	AppConfig["SECRET_KEY"] = &SK

	//clean the environment
	os.Unsetenv("AK")
	os.Unsetenv("SK")
	return nil
}

//Read application instanceId
func GetAppInstanceId() (string, error) {
	if len(os.Getenv("APPINSTID")) == 0  {
		err := errors.New("APPINSTID should be set in env variable")
		log.Error("AppInstanceId must be set")
		return "", err
	}
	AppInstanceId = os.Getenv("APPINSTID")
	return AppInstanceId, nil
}