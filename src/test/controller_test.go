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

//Test package
package test

import (
	"github.com/agiledragon/gomonkey"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	controller "mep-agent/src/controllers"
	"mep-agent/src/model"
	svc "mep-agent/src/service"
	"mep-agent/src/util"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const panicFormatString = "Panic: %v"

func TestGetForTokenNil(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Errorf(panicFormatString, r)
		}
	}()

	getBeegoController := beego.Controller{Ctx: &context.Context{ResponseWriter: &context.Response{ResponseWriter: httptest.NewRecorder()}},
		Data: make(map[interface{}]interface{})}
	getController := &controller.TokenController{Controller: getBeegoController}
	getController.Get()
	response := getController.Ctx.ResponseWriter.ResponseWriter.(*httptest.ResponseRecorder)

	assert.Equal(t, 412, response.Code, "get failed")
}


func TestGet(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Errorf(panicFormatString, r)
		}
	}()
	util.FirstToken = true
	AppConfigProperties := make(map[string]*[]byte)
	sk := []byte("sksksk")
	ak := []byte("akakak")

	AppConfigProperties["SECRET_KEY"] = &sk
	AppConfigProperties["ACCESS_KEY"] = &ak
	util.AppConfig = AppConfigProperties

	util.MepToken = model.TokenModel{AccessToken: "akakak", TokenType: "Bearer", ExpiresIn: 3600}
	patch1 := gomonkey.ApplyFunc(svc.GetMepToken, func(model.Auth) (error) {
		return nil
	})
	defer patch1.Reset()



	//getBeegoController := beego.Controller{Ctx: &context.Context{ResponseWriter: &context.Response{ResponseWriter: httptest.NewRecorder()}},
	//	Data: make(map[interface{}]interface{})}
	getController := getController()//&controller.TokenController{Controller: getBeegoController}
	//getController.Data["json"] = &util.MepToken
	//getController.Ctx.Output.JSON(getController.Data["json"], true, true)
	response := getController.Ctx.ResponseWriter.ResponseWriter.(*httptest.ResponseRecorder)
	patch2 := gomonkey.ApplyFunc(getController.ServeJSON, func(encoding ...bool) {

	})
	defer patch2.Reset()
	//getController.ServeJSON()
	getController.Get()


	assert.Equal(t, 200, response.Code, "get failed")
}

func getController() *controller.TokenController {
	c := &controller.TokenController{}
	c.Init(context.NewContext(), "", "", nil)
	req, err := http.NewRequest("POST", "http://127.0.0.1", strings.NewReader(""))
	if err != nil {
		log.Error("prepare http request failed")
	}
	c.Ctx.Request = req
	c.Ctx.Request.Header.Set("X-Real-Ip", "127.0.0.1")
	c.Ctx.ResponseWriter = &context.Response{}
	c.Ctx.ResponseWriter.ResponseWriter = httptest.NewRecorder()
	c.Ctx.Output = context.NewOutput()
	c.Ctx.Input = context.NewInput()
	c.Ctx.Output.Reset(c.Ctx)
	c.Ctx.Input.Reset(c.Ctx)
	return c
}