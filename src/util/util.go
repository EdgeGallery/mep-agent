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

// Package util utility package
package util

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"mep-agent/src/model"
	"os"
	"time"
)

// appConfigProperties App configuration properties.
type appConfigProperties map[string]*[]byte

// AppConfig Application configurations.
var AppConfig = appConfigProperties{}

// MepToken Token.
var MepToken = model.TokenModel{}

// FirstToken Mark initial token has been received.
var FirstToken = false

// AppInstanceID App instance ID from configuration.
var AppInstanceID string

// RefreshTimer Refresh token timer.
var RefreshTimer *time.Timer

// RefreshTimeBuffer Timer buffer 5 sec.
const RefreshTimeBuffer = 5

const (
	ak        string = "ak"
	sk        string = "sk"
	appInstID string = "APPINSTID"
)

// ClearByteArray Clears byte array.
func ClearByteArray(data []byte) {
	if data == nil {
		return
	}
	for i := 0; i < len(data); i++ {
		data[i] = 0
	}
}

// ClearMap clear [string, *[]byte] map, called only in error case.
func ClearMap() {
	for _, element := range AppConfig {
		ClearByteArray(*element)
	}
}

// ReadTokenFromEnvironment read and clearing the variable from the environment.
func ReadTokenFromEnvironment() error {
	// clean the environment
	defer os.Unsetenv(ak)
	defer os.Unsetenv(sk)

	akVal := os.Getenv(ak)
	skVal := os.Getenv(sk)

	if len(akVal) == 0 || len(skVal) == 0 {
		err := errors.New("ak and sk keys should be set in env variable")
		log.Error("Keys should not be empty")

		return err
	}
	akByte := []byte(akVal)
	AppConfig["ACCESS_KEY"] = &akByte
	skByte := []byte(skVal)
	AppConfig["SECRET_KEY"] = &skByte
	log.Infof("Ak: %s", akByte)

	return nil
}

// GetAppInstanceID Read application instanceId.
func GetAppInstanceID() (string, error) {
	defer os.Unsetenv(appInstID)
	instID := os.Getenv(appInstID)
	if len(instID) == 0 {
		err := errors.New("app instance id should be set in env variable")
		log.Error("App instance id must be set")

		return "", err
	}
	AppInstanceID = instID
	return instID, nil
}
