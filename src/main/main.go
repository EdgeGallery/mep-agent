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

// main package
package main

import (
	"mep-agent/src/config"
	controller "mep-agent/src/controllers"
	_ "mep-agent/src/router"
	"mep-agent/src/service"
	"mep-agent/src/util"
	"os"

	"github.com/astaxie/beego"
	log "github.com/sirupsen/logrus"
)

func main() {
	// reading and cleaning the token from environment
	err := util.ReadTokenFromEnvironment()
	if err != nil {
		log.Error("Failed to read the token from environment variables")
		util.ClearMap()
		os.Exit(1)
	}

	service.TLSConf, err = service.TLSConfig()
	if err != nil {
		log.Error("Failed to set TLS Configurations")
		util.ClearMap()
		os.Exit(1)
	}

	config.ServerURLConfig, err = config.GetServerURL()
	if err != nil {
		log.Error("Failed to get server url Configurations")
		util.ClearMap()
		os.Exit(1)
	}

	// start main service
	go service.BeginService().Start("./conf/app_instance_info.yaml")

	log.Info("Starting server")
	beego.ErrorController(&controller.ErrorController{})
	beego.Run()
}
