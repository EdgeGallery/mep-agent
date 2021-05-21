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

// Test package
package test

import (
	"mep-agent/src/config"
	"mep-agent/src/util"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/agiledragon/gomonkey"
)

func TestGetServerUrl(t *testing.T) {
	patch1 := gomonkey.ApplyFunc(os.Getenv, func(string) string {
		return ""
	})
	defer patch1.Reset()

	patch2 := gomonkey.ApplyFunc(util.ValidateDNS, func(string) error {
		return nil
	})
	defer patch2.Reset()

	patch3 := gomonkey.ApplyFunc(util.ValidateByPattern, func(string, string) error {
		return nil
	})
	defer patch3.Reset()

	_, err := config.GetServerURL()
	assert.Equal(t, nil, err)
}
