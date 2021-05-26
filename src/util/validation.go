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

// Package util utility
package util

import (
	"errors"
	"net"
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	// PortPattern Regular expression for Port.
	PortPattern    string = `^([1-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$`
	// domainPattern Regular expression for Domain.
	domainPattern  string = `^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])(\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9]))*$`
	dnsPattern     string = `^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`
	maxHostNameLen        = 253
)

// ValidateDNS Validates Ip address.
func ValidateDNS(ip string) error {
	ipv := net.ParseIP(ip)
	if ipv != nil && (ipv.IsMulticast() || ipv.Equal(net.IPv4bcast)) {
		return errors.New("invalid dns ip")
	}
	if ipv != nil {
		return nil
	}

	return ValidateByPattern(dnsPattern, ip)
}

// ValidateDomainName Validates domain name.
func ValidateDomainName(name string) error {
	if len(name) == 0 || len(name) > maxHostNameLen {
		return errors.New("validate domain name failed")
	}

	return ValidateByPattern(domainPattern, name)
}

// ValidateByPattern Validates given string with pattern.
func ValidateByPattern(pattern string, param string) error {
	res, errMatch := regexp.MatchString(pattern, param)
	if errMatch != nil {
		return errMatch
	}
	if !res {
		return errors.New("validate failed")
	}

	return nil
}

// ValidateUUID Validates UUID.
func ValidateUUID(id string) error {
	if len(id) != 0 {
		validate := validator.New()
		res := validate.Var(id, "required,uuid")
		if res != nil {
			return errors.New("uuid validate failed")
		}
	} else {
		return errors.New("uuid validate failed")
	}

	return nil
}
