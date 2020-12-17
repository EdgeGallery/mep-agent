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

// util service
package service

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"mep-agent/src/model"
	"mep-agent/src/util"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// const
var cipherSuiteMap = map[string]uint16{
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256": tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384": tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
}

// get yaml and parse to AppInstanceInfo object
func GetAppInstanceConf(path string) (model.AppInstanceInfo, error) {
	yamlFile, err := ioutil.ReadFile(path)
	var info model.AppInstanceInfo
	if err != nil {
		return info, err
	}
	err = yaml.UnmarshalStrict(yamlFile, &info)
	if err != nil {
		return info, err
	}
	return info, nil
}

func getAPPConf(path string) (model.AppConfInfo, error) {
	yamlFile, err := ioutil.ReadFile(path)
	var info model.AppConfInfo
	if err != nil {
		return info, err
	}
	err = yaml.UnmarshalStrict(yamlFile, &info)
	if err != nil {
		return info, err
	}
	return info, nil
}

// register to mep
func PostRegisterRequest(registerData RegisterData) (string, error) {
	// construct http request
	req, errNewRequest := http.NewRequest("POST", registerData.url, strings.NewReader(registerData.data))
	if errNewRequest != nil {
		return "", errNewRequest
	}
	req.Header.Set("Authorization", registerData.token.TokenType+" "+registerData.token.AccessToken)

	// send http request
	response, errDo := DoRequest(req)
	if errDo != nil {
		return "", errDo
	}

	defer response.Body.Close()
	body, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		return "", err2
	}

	if response.StatusCode != http.StatusCreated {
		return "", errors.New("created failed, status is " + strconv.Itoa(response.StatusCode))
	}

	return string(body), nil
}

// get token from mep
func PostTokenRequest(param string, url string, auth model.Auth) (string, error) {
	log.Infof("PostTokenRequest param: %s, url: %s, auth: %s", param, url, auth)
	// construct http request
	req, errNewRequest := http.NewRequest("POST", url, strings.NewReader(param))
	if errNewRequest != nil {
		return "", errNewRequest
	}

	// request header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(util.DATE_HEADER, time.Now().Format(util.DATE_FORMAT))
	req.Header.Set("Host", req.Host)
	// calculate signature by safe algorithm
	sign := util.Sign{
		AccessKey: auth.AccessKey,
		SecretKey: auth.SecretKey,
	}
	authorization, errSign := sign.GetAuthorizationValueWithSign(req)
	if errSign != nil {
		return "", errSign
	}
	req.Header.Set("Authorization", authorization)

	// send http request
	response, errDo := DoRequest(req)
	if errDo != nil {
		return "", errDo
	}

	defer response.Body.Close()
	body, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		return "", err2
	}

	if response.StatusCode != http.StatusOK {
		log.Errorf("response: %s", response)
		return "", errors.New("request failed, status is " + strconv.Itoa(response.StatusCode))
	}
	log.Infof("response: %s", response)
	return string(body), nil
}

// do request
func DoRequest(req *http.Request) (*http.Response, error) {
	config, err := TlsConfig()
	if err != nil {
		log.Error("Unable to send request")
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: config,
	}
	client := &http.Client{Transport: tr}

	return client.Do(req)
}

// Constructs tls configuration
func TlsConfig() (*tls.Config, error) {
	appConf, errGetConf := getAPPConf("./conf/app_conf.yaml")
	if errGetConf != nil {
		log.Error("parse app_conf.yaml failed")
		return nil, errors.New("parse app_conf.yaml failed")
	}
	sslCiphers := appConf.SslCiphers
	if len(sslCiphers) == 0 {
		return nil, errors.New("TLS cipher configuration is not recommended or invalid")
	}
	cipherSuites := getCipherSuites(sslCiphers)
	if cipherSuites == nil {
		return nil, errors.New("TLS cipher configuration is not recommended or invalid")
	}
	domainName := os.Getenv("CA_CERT_DOMAIN_NAME")
	if util.ValidateDomainName(domainName) != nil {
		return nil, errors.New("Domain name validation failed")
	}
	return &tls.Config{
		ServerName:         domainName,
		MinVersion:         tls.VersionTLS12,
		CipherSuites:       cipherSuites,
		InsecureSkipVerify: true,
	}, nil
}

func getCipherSuites(sslCiphers string) []uint16 {
	cipherSuiteArr := make([]uint16, 0, 5)
	cipherSuiteNameList := strings.Split(sslCiphers, ",")
	for _, cipherName := range cipherSuiteNameList {
		cipherName = strings.TrimSpace(cipherName)
		if len(cipherName) == 0 {
			continue
		}
		mapValue, ok := cipherSuiteMap[cipherName]
		if !ok {
			log.Warn("Not recommended cipher suite.")
			return nil
		}
		cipherSuiteArr = append(cipherSuiteArr, mapValue)
	}
	if len(cipherSuiteArr) > 0 {
		return cipherSuiteArr
	}
	return nil
}

// Send Service heartbeat to MEP
func SendHeartBeatRequest(heartBeatData HeartBeatData) (string, error) {
	req, errNewRequest := http.NewRequest("PUT", heartBeatData.url, strings.NewReader(heartBeatData.data))
	if errNewRequest != nil {
		return "", errNewRequest
	}
	req.Header.Set("Authorization", heartBeatData.token.TokenType+" "+heartBeatData.token.AccessToken)

	response, errDo := DoRequest(req)
	if errDo != nil {
		return "", errDo
	}

	defer response.Body.Close()
	body, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		return "", err2
	}
	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusNoContent {
		return string(body), nil
	} else {
		return "", errors.New("heartbeat request failed, status is " + strconv.Itoa(response.StatusCode))
	}
}
