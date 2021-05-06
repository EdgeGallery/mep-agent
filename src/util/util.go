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
var FirstToken = false

//App instance ID from configuration
var AppInstanceId string

//Refresh token timer
var RefreshTimer *time.Timer

//Timer buffer 5 sec
const RefreshTimeBuffer = 5

const (
	AK        string = "AK"
	SK        string = "SK"
	APPINSTID string = "APPINSTID"
)

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
func ClearMap() {
	for _, element := range AppConfig {
		ClearByteArray(*element)
	}
}

//read and clearing the variable from the environment
func ReadTokenFromEnvironment() error {
	//clean the environment
	defer os.Unsetenv(AK)
	defer os.Unsetenv(SK)

	ak := os.Getenv(AK)
	sk := os.Getenv(SK)

	if len(ak) == 0 || len(sk) == 0 {
		err := errors.New("ak and sk keys should be set in env variable")
		log.Error("Keys should not be empty")
		return err
	}
	akByte := []byte(ak)
	AppConfig["ACCESS_KEY"] = &akByte
	skByte := []byte(sk)
	AppConfig["SECRET_KEY"] = &skByte
	log.Infof("Ak: %s", akByte)
	return nil
}

//Read application instanceId
func GetAppInstanceId() (string, error) {
	defer os.Unsetenv(APPINSTID)
	instId := os.Getenv(APPINSTID)
	if len(instId) == 0 {
		err := errors.New("appInstanceId should be set in env variable")
		log.Error("appInstanceId must be set")
		return "", err
	}

	return instId, nil
}
