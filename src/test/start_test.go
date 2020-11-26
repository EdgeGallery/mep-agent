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
	"github.com/agiledragon/gomonkey"
	"github.com/smartystreets/goconvey/convey"
	"mep-agent/src/model"
	"mep-agent/src/service"
	"sync"
	"testing"
)

// Tests start service
func TestStartSuccess(t *testing.T) {

	convey.Convey("Start", t, func() {
		var waitRoutineFinish sync.WaitGroup
		var dataStore []model.ServiceInfoPost
		data := model.ServiceInfoPost{LivenessInterval : 1}
		dataStore = append(dataStore,data)
		patch1 := gomonkey.ApplyFunc(service.GetAppInstanceConf, func(path string) (model.AppInstanceInfo, error) {
			return model.AppInstanceInfo{}, nil
		})
		patch2 := gomonkey.ApplyFunc(service.GetMepToken, func(auth model.Auth) (*model.TokenModel, error) {
			return &model.TokenModel{}, nil
		})
		patch3 := gomonkey.ApplyFunc(service.RegisterToMep, func(conf model.AppInstanceInfo,
			                                                     token *model.TokenModel,
			                                                     wg *sync.WaitGroup)([]model.ServiceInfoPost, error) {
			return dataStore, nil
		})
		skByte := []byte("secretKey")
		service.BeginService().Start("../../conf/app_instance_info.yaml", "accessKey", &skByte, &waitRoutineFinish)

		defer patch1.Reset()
		defer patch2.Reset()
		defer patch3.Reset()
	})


}
