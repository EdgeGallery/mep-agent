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

// Package controllers TokenController controller package
package controllers

import (
	"github.com/astaxie/beego"
	log "github.com/sirupsen/logrus"
	"mep-agent/src/util"
	"net/http"
)

// TokenController handles token request.
type TokenController struct {
	beego.Controller
}

// Get /mep-agent/v1/token function.
func (c *TokenController) Get() {
	log.Info("Received get token request from app")
	if !util.FirstToken {
		log.Error("First token not yet received.")
		c.Ctx.ResponseWriter.WriteHeader(http.StatusPreconditionFailed)

		return
	}
	// Get the last token
	c.Data["json"] = &util.MepToken
	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
	c.ServeJSON()
}
