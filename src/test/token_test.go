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
	"testing"
)

// Tests get mep token
func TestGetMepToken(t *testing.T) {

	convey.Convey("Start", t, func() {

		var token = model.TokenModel{AccessToken: "akakak", TokenType: "Bear", ExpiresIn: 3600}
		httpResponseBytes, errMarshal := json.Marshal(token)
		if errMarshal != nil {
			t.Error("Marshal http Response Error: " + errMarshal.Error())
		}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err2 := w.Write(httpResponseBytes)
			if err2 != nil {
				t.Error("Write Response Error")
			}
		}))

		defer ts.Close()
		api := ts.URL

		patch1 := gomonkey.ApplyFunc(config.GetServerURL, func() (config.ServerURL, error) {
			return config.ServerURL{MepAuthURL: api}, nil
		})
        config.ServerURLConfig, _ = config.GetServerURL()
		patch2 := gomonkey.ApplyFunc(service.TLSConfig, func() (*tls.Config, error) {
			return nil, nil
		})

		sk := []byte("sksksk")
		err := service.GetMepToken(model.Auth{AccessKey: "akakak", SecretKey: &sk})
		if err != nil {
			t.Error(err.Error())
		}

		defer patch1.Reset()
		defer patch2.Reset()
	})

}
