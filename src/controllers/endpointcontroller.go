/*
 * Copyright 2021 Huawei Technologies Co., Ltd.
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

// Package controllers Endpoint controller
package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	log "github.com/sirupsen/logrus"
	"mep-agent/src/config"
	"mep-agent/src/service"
	"mep-agent/src/util"
	"net/http"
)

// EndpointController beego Endpoint Controller.
type EndpointController struct {
	beego.Controller
}

// Service Service information.
type Service struct {
	TransportInfo TransportInfo `yaml:"transportInfo" json:"transportInfo"`
}

// TransportInfo Transport information of the service.
type TransportInfo struct {
	Id          string       `yaml:"id" json:"id"`
	Name        string       `yaml:"name" json:"name"`
	Description string       `yaml:"description" json:"description"`
	Protocol    string       `yaml:"protocol" json:"protocol"`
	Version     string       `yaml:"version" json:"version"`
	Endpoint    endPointInfo `yaml:"endpoint" json:"endpoint"`
}

// endPointInfo End point of the service.
type endPointInfo struct {
	Uris        []string              `json:"uris" validate:"omitempty,dive,uri"`
	Addresses   []endPointInfoAddress `json:"addresses" validate:"omitempty,dive"`
	Alternative interface{}           `json:"alternative"`
}

// endPointInfoAddress Endpoint info address.
type endPointInfoAddress struct {
	Host string `json:"host" validate:"required"`
	Port uint32 `json:"port" validate:"required,gt=0,lte=65535"`
}

// Get handles endpoint request from app.
func (c *EndpointController) Get() {
	log.Info("Received get endpoint request from app")

	serName := c.Ctx.Input.Param(":serName")
	url := config.ServerURLConfig.MepServiceDiscoveryURL + serName
	requestData := service.RequestData{Data: "", URL: url, Token: &util.MepToken}
	resBody, errPostRequest := service.SendQueryRequest(requestData)
	if errPostRequest != nil {
		log.Error("Failed heart beat request to mep, URL is " + url)
	}
	var resBodyMap []Service

	log.Info("Response Body: ", resBody)

	err := json.Unmarshal([]byte(resBody), &resBodyMap)
	if err != nil {
		log.Error("Unmarshal failed")
		c.Data["json"] = "Service does not exist."
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
	} else {
		transportInfo := resBodyMap[0].TransportInfo
		log.Info("Endpoint: ", transportInfo.Endpoint)
		c.Data["json"] = transportInfo.Endpoint
		c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
	}

	c.ServeJSON()
}
