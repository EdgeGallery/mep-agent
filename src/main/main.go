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
	"github.com/astaxie/beego"
	log "github.com/sirupsen/logrus"
	"mep-agent/src/controllers"
	_ "mep-agent/src/router"
	"mep-agent/src/service"
	"mep-agent/src/util"
	"os"
)


func main() {
	// reading and cleaning the token from environment
	err := util.ReadTokenFromEnvironment()
	if err != nil {
		log.Error("Failed to read the token from environment variables")
		util.ClearMap()
		os.Exit(1)
	}

	// validate ACCESS_KEY and SECRET_KEY with nil, length and regex
	if util.AppConfig["ACCESS_KEY"] == nil || len(*util.AppConfig["ACCESS_KEY"]) == 0 {
		log.Info("The input param of ak is invalid")
		util.ClearMap()
		os.Exit(1)
	}

	if util.AppConfig["SECRET_KEY"] == nil || len(*util.AppConfig["SECRET_KEY"]) == 0 {
		log.Info("The input param of sk is invalid")
		util.ClearMap()
		os.Exit(1)
	}

	if err := util.ValidateAkSk(string(*util.AppConfig["ACCESS_KEY"]), util.AppConfig["SECRET_KEY"]); err != nil {
		log.Info("the input param of ak or sk do not pass the validation")
		util.ClearMap()
		os.Exit(1)
	}

	sk := util.AppConfig["SECRET_KEY"]
	// start main service
	go service.BeginService().Start("./conf/app_instance_info.yaml",
		string(*util.AppConfig["ACCESS_KEY"]), sk)

	log.Info("Starting server")
	beego.ErrorController(&controllers.ErrorController{})
	beego.Run()
}
