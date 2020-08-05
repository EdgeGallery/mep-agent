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

package test

import (
	"crypto/tls"
	"encoding/json"
	"github.com/agiledragon/gomonkey"
	"github.com/smartystreets/goconvey/convey"
	"mep-agent/src/config"
	"mep-agent/src/model"
	"mep-agent/src/service"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// Tests register service to mep
func TestRegisterToMep(t *testing.T)  {
	convey.Convey("RegisterToMep", t, func() {
		respString := "http response"
		var waitRoutineFinish sync.WaitGroup
		httpResponseBytes, errMarshal := json.Marshal(respString)
		if errMarshal != nil {
			t.Error("Marshal http Response Error: " + errMarshal.Error())
		}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			_, err2 := w.Write(httpResponseBytes)
			if err2 != nil {
				t.Error("Write Response Error")
			}
		}))

		patch1 := gomonkey.ApplyFunc(config.GetServerUrl, func() (config.ServerUrl, error) {
			return config.ServerUrl{MepServerRegisterUrl: ts.URL}, nil
		})
		patch2 := gomonkey.ApplyFunc(service.TlsConfig, func() (*tls.Config, error) {
			return nil, nil
		})

		defer ts.Close()
		defer patch1.Reset()
		defer patch2.Reset()

		conf, errGetConf := service.GetAppInstanceConf("../../conf/app_instance_info.yaml")
		if errGetConf != nil {
			t.Error(errGetConf.Error())
		}

		var token = model.TokenModel{AccessToken: "akakak", TokenType: "Bear", ExpiresIn: 3600}
		errRegister := service.RegisterToMep(conf, &token, &waitRoutineFinish)
		if errRegister != nil {
			t.Error(errRegister.Error())
		}
	})
}
