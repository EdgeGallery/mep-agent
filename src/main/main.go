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

package main

import (
	"bufio"
	"bytes"
	"mep-agent/src/model"
	"mep-agent/src/service"
	"mep-agent/src/util"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"unsafe"

	log "github.com/sirupsen/logrus"
)

type AppConfigProperties map[string]*[]byte

func getMepAgentAppConfig(line []byte, config AppConfigProperties) AppConfigProperties {
	if bytes.Contains(line, []byte("=")) {
		keyVal := bytes.Split(line, []byte("="))
		key := bytes.TrimSpace(keyVal[0])
		val := bytes.TrimSpace(keyVal[1])
		config[string(key)] = &val
	}
	return config
}

func readPropertiesFile(filename string) (AppConfigProperties, error) {
	config := AppConfigProperties{}

	if len(filename) == 0 {
		return config, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Error("Failed to open the file.")
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		config = getMepAgentAppConfig(line, config)
	}
	if err := scanner.Err(); err != nil {
		log.Error("Failed to read the file.")
		clearAppConfigOnExit(config)
		return nil, err
	}
	return config, nil
}
func main() {
	configFilePath := filepath.FromSlash("/usr/mep/mepagent.properties")
	appConfig, err := readPropertiesFile(configFilePath)
	if err != nil {
		log.Error("Failed to read the config parameters from properties file")
		os.Exit(1)
	}
	err = os.Truncate(configFilePath, 0)
	if err != nil {
		clearAppConfigOnExit(appConfig)
		os.Exit(1)
	}
	var waitRoutineFinish sync.WaitGroup
	enableWait := false

	// wait before exit if only enabled by environment variable
	if waitStatus := os.Getenv("ENABLE_WAIT"); waitStatus != "" {
		if waitFlag, err := strconv.ParseBool(waitStatus); err == nil {
			enableWait = waitFlag
		}
	}

	if appConfig["ACCESS_KEY"] == nil || len(*appConfig["ACCESS_KEY"]) == 0 {
		log.Info("The input param of ak is invalid")
		clearAppConfigOnExit(appConfig)
		os.Exit(1)
	}

	if appConfig["SECRET_KEY"] == nil || len(*appConfig["SECRET_KEY"]) == 0 {
		log.Info("The input param of sk is invalid")
		clearAppConfigOnExit(appConfig)
		os.Exit(1)
	}

	if err := util.ValidateAkSk(string(*appConfig["ACCESS_KEY"]), appConfig["SECRET_KEY"]); err != nil {
		log.Info("the input param of ak or sk do not pass the validation")
		clearAppConfigOnExit(appConfig)
		os.Exit(1)
	}
	sk := appConfig["SECRET_KEY"]
	token := service.BeginService().Start("./conf/app_instance_info.yaml",
		string(*appConfig["ACCESS_KEY"]), sk, &waitRoutineFinish)
	waitRoutineFinish.Wait()
	if token != nil {
		// clear skValue
		skByteVal := *(*[]byte)(unsafe.Pointer(&token.AccessToken))
		util.ClearByteArray(skByteVal)
		*token = model.TokenModel{}
	}
	if enableWait {
		service.Heart()
	}
}

func clearAppConfigOnExit(appConfig AppConfigProperties){
	for _, element := range appConfig {
		util.ClearByteArray(*element)
	}
}
