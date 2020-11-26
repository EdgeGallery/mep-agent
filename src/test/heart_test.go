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
	log "github.com/sirupsen/logrus"
	"github.com/smartystreets/goconvey/convey"
	"mep-agent/src/config"
	"mep-agent/src/model"
	"mep-agent/src/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Tests register service to mep
func TestHeartBeatRequestToMep(t *testing.T)  {
	convey.Convey("HeartBeatRequestToMep", t, func() {
		registerResponse := string(`{"serName": "get", "livenessInterval": 30, "serInstanceId":"12345", "_links": { "self" : { "liveness": "/mec_service_mgmt/v1/applications/5abe4782-2c70-4e47-9a4e-0ee3a1a0fd1f/service/5abe4782-2c70-4e47-9a4e-0ee3a1a0fd1f/liveness" }}}`)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)

		}))
		patch1 := gomonkey.ApplyFunc(config.GetServerUrl, func() (config.ServerUrl, error) {
			return config.ServerUrl{MepHeartBeatUrl: ts.URL}, nil
		})
		patch2 := gomonkey.ApplyFunc(service.TlsConfig, func() (*tls.Config, error) {
			return nil, nil
		})

		defer ts.Close()
		defer patch1.Reset()
		defer patch2.Reset()

		var token = model.TokenModel{AccessToken: "akakak", TokenType: "Bear", ExpiresIn: 3600}

		serviceInfo := model.ServiceInfoPost{}
		errJsonUnMarshal := json.Unmarshal([]byte(registerResponse), &serviceInfo)
		if errJsonUnMarshal !=nil {
			log.Error("Failed to marshal service info to object", errJsonUnMarshal.Error())
		}
		service.HeartBeatRequestToMep(serviceInfo, &token)
	})
}
