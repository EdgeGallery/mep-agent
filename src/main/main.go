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

func main() {

	// read mepagent.properties file to AppConfigProperties object
	configFilePath := filepath.FromSlash("/usr/mep/mepagent.properties")
	appConfig, err := readPropertiesFile(configFilePath)
	if err != nil {
		log.Error("Failed to read the config parameters from properties file")
		os.Exit(1)
	}

	// clear mepagent.properties file
	err = os.Truncate(configFilePath, 0)
	if err != nil {
		util.ClearMap(appConfig)
		os.Exit(1)
	}

	enableWait := false

	// wait before exit if only enabled by environment variable
	if waitStatus := os.Getenv("ENABLE_WAIT"); waitStatus != "" {
		if waitFlag, err := strconv.ParseBool(waitStatus); err == nil {
			enableWait = waitFlag
		}
	}

	// validate ACCESS_KEY and SECRET_KEY with nil, length and regex
	if appConfig["ACCESS_KEY"] == nil || len(*appConfig["ACCESS_KEY"]) == 0 {
		log.Info("The input param of ak is invalid")
		util.ClearMap(appConfig)
		os.Exit(1)
	}
	if appConfig["SECRET_KEY"] == nil || len(*appConfig["SECRET_KEY"]) == 0 {
		log.Info("The input param of sk is invalid")
		util.ClearMap(appConfig)
		os.Exit(1)
	}
	if err := util.ValidateAkSk(string(*appConfig["ACCESS_KEY"]), appConfig["SECRET_KEY"]); err != nil {
		log.Info("the input param of ak or sk do not pass the validation")
		util.ClearMap(appConfig)
		os.Exit(1)
	}

	sk := appConfig["SECRET_KEY"]
	var waitRoutineFinish sync.WaitGroup

	// start main service
	token := service.BeginService().Start("./conf/app_instance_info.yaml",
		string(*appConfig["ACCESS_KEY"]), sk, &waitRoutineFinish)
	waitRoutineFinish.Wait()
	if token != nil {
		// clear access_token
		accessToken := *(*[]byte)(unsafe.Pointer(&token.AccessToken))
		util.ClearByteArray(accessToken)
		*token = model.TokenModel{}
	}

	// service heart
	if enableWait {
		service.Heart()
	}
}

// read mepagent.properties file to AppConfigProperties object
func readPropertiesFile(filePath string) (AppConfigProperties, error) {
	var config AppConfigProperties

	// validate path length and open the file
	if len(filePath) == 0 {
		return config, nil
	}
	file, err := os.Open(filePath)
	if err != nil {
		log.Error("Failed to open the file.")
		return nil, err
	}
	defer file.Close()

	// scanner the file and transform to AppConfigProperties object
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		config = transformLineToAppConfigProperties(line, config)
	}
	if err := scanner.Err(); err != nil {
		log.Error("Failed to read the file.")
		util.ClearMap(config)
		return nil, err
	}
	return config, nil
}

func transformLineToAppConfigProperties(line []byte, config AppConfigProperties) AppConfigProperties {
	if bytes.Contains(line, []byte("=")) {
		keyVal := bytes.Split(line, []byte("="))
		key := bytes.TrimSpace(keyVal[0])
		val := bytes.TrimSpace(keyVal[1])
		config[string(key)] = &val
	}
	return config
}


