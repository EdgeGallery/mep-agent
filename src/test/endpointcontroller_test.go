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

package test

import (
	"encoding/json"
	"mep-agent/src/config"
	"mep-agent/src/controllers"
	"mep-agent/src/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey"

	"github.com/astaxie/beego/context"
	log "github.com/sirupsen/logrus"
)

func TestEndpointControllerGet(t *testing.T) {
	patch1 := gomonkey.ApplyFunc(config.GetServerURL, func() (config.ServerURL, error) {
		return config.ServerURL{MepServiceDiscoveryURL: "http://127.0.0.1:8088/mep/mec_service_mgmt/v1/services?ser_name="}, nil
	})
	defer patch1.Reset()

	patch2 := gomonkey.ApplyFunc(service.SendQueryRequest, func(service.RequestData, string) (string, error) {
		service := controllers.Service{
			TransportInfo: controllers.TransportInfo{
				Id:   "123",
				Name: "service1",
			},
		}
		str, _ := json.Marshal(service)
		return string(str), nil
	})
	defer patch2.Reset()
	c := getEndpointController()
	patch3 := gomonkey.ApplyFunc(c.ServeJSON, func(encoding ...bool) {

	})
	defer patch3.Reset()
	c.Get()

}

func getEndpointController() *controllers.EndpointController {
	c := &controllers.EndpointController{}
	c.Init(context.NewContext(), "", "", nil)
	req, err := http.NewRequest("POST", "http://127.0.0.1", strings.NewReader(""))
	if err != nil {
		log.Error("Prepare http request failed")
	}
	c.Ctx.Request = req
	c.Ctx.Request.Header.Set("X-Real-Ip", "127.0.0.1")
	c.Ctx.ResponseWriter = &context.Response{}
	c.Ctx.ResponseWriter.ResponseWriter = httptest.NewRecorder()
	c.Ctx.Output = context.NewOutput()
	c.Ctx.Input = context.NewInput()
	c.Ctx.Output.Reset(c.Ctx)
	c.Ctx.Input.Reset(c.Ctx)
	c.Ctx.Input.SetParam(":serName", "service1")
	return c
}
