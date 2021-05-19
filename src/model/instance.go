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

// Package model define the type information
package model

// AppInstanceInfo : Application instance info.
type AppInstanceInfo struct {
	ServiceInfoPosts                         []ServiceInfoPost                         `yaml:"serviceInfoPosts" json:"serviceInfoPosts"`
	SerAvailabilityNotificationSubscriptions []serAvailabilityNotificationSubscription `yaml:"serAvailabilityNotificationSubscriptions" json:"serAvailabilityNotificationSubscriptions"`
}

// ServiceInfoPost Service Information to be registered.
type ServiceInfoPost struct {
	SerName           string         `yaml:"serName" json:"serName"`
	SerCategory       categoryRef    `yaml:"serCategory" json:"serCategory"`
	Version           string         `yaml:"version" json:"version"`
	State             serviceState   `yaml:"state" json:"state"`
	TransportId       string         `yaml:"transportId" json:"transportId"`
	TransportInfo     transportInfo  `yaml:"transportInfo" json:"transportInfo"`
	Serializer        serializerType `yaml:"serializer" json:"serializer"`
	ScopeOfLocality   localityType   `yaml:"scopeOfLocality" json:"scopeOfLocality"`
	ConsumedLocalOnly bool           `yaml:"consumedLocalOnly" json:"consumedLocalOnly"`
	IsLocal           bool           `yaml:"isLocal" json:"isLocal"`
	LivenessInterval  int            `yaml:"livenessInterval" json:"livenessInterval,omitempty"`
	Links             _links         `json:"_links,omitempty"`
	SerInstanceId     string         `json:"serInstanceId,omitempty"`
}

// categoryRef Service category.
type categoryRef struct {
	Href    string `yaml:"href" json:"href"`
	Id      string `yaml:"id" json:"id"`
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
}

// serviceState Service state.
type serviceState string

// transportInfo Transport Information of the service to be registered.
type transportInfo struct {
	Id               string           `yaml:"id" json:"id"`
	Name             string           `yaml:"name" json:"name"`
	Description      string           `yaml:"description" json:"description"`
	TransportType    transportType    `yaml:"type" json:"type"`
	Protocol         string           `yaml:"protocol" json:"protocol"`
	Version          string           `yaml:"version" json:"version"`
	Endpoint         endPointInfo     `yaml:"endpoint" json:"endpoint"`
	Security         securityInfo     `yaml:"security" json:"security"`
	ImplSpecificInfo implSpecificInfo `yaml:"implSpecificInfo" json:"implSpecificInfo"`
}

type transportType string

// endPointInfo: Endpoint of the service to be registered.
type endPointInfo struct {
	Addresses []endPointInfoAddress `yaml:"addresses" json:"addresses"`
}

// endPointInfoAddress: Endpoint information.
type endPointInfoAddress struct {
	Host string `yaml:"host" json:"host"`
	Port uint32 `yaml:"port" json:"port"`
}

// securityInfo: Security information.
type securityInfo struct {
	OAuth2Info securityInfoOAuth2Info `yaml:"oAuth2Info" json:"oAuth2Info"`
}

// securityInfoOAuth2Info Security Auth2 information.
type securityInfoOAuth2Info struct {
	GrantTypes    []securityInfoOAuth2InfoGrantType `yaml:"grantTypes" json:"grantTypes"`
	TokenEndpoint string                            `yaml:"tokenEndpoint" json:"tokenEndpoint"`
}

// securityInfoOAuth2InfoGrantType Security Auth2 grant type.
type securityInfoOAuth2InfoGrantType string

// implSpecificInfo Impl specific info if any.
type implSpecificInfo struct {
}

// serializerType Serializer Type.
type serializerType string

// localityType : Scope of Locality Type.
type localityType string

// serAvailabilityNotificationSubscription : Service Availability Notification Subscription.
type serAvailabilityNotificationSubscription struct {
	SubscriptionType  string                                                   `yaml:"subscriptionType" json:"subscriptionType"`
	CallbackReference string                                                   `yaml:"callbackReference" json:"callbackReference"`
	Links             self                                                     `yaml:"links" json:"links"`
	FilteringCriteria serAvailabilityNotificationSubscriptionFilteringCriteria `yaml:"filteringCriteria" json:"filteringCriteria"`
}

// self link.
type self struct {
	Self linkType `yaml:"self" json:"self"`
}

// linkType Link type.
type linkType struct {
	Href string `yaml:"href" json:"href"`
}

// serAvailabilityNotificationSubscriptionFilteringCriteria : Service Availability Notification Subscription FilteringCriteria.
type serAvailabilityNotificationSubscriptionFilteringCriteria struct {
	SerInstanceIds []string       `yaml:"serInstanceIds" json:"serInstanceIds"`
	SerNames       []string       `yaml:"serNames" json:"serNames"`
	SerCategories  []categoryRef  `yaml:"serCategories" json:"serCategories"`
	States         []serviceState `yaml:"states" json:"states"`
	IsLocal        bool           `yaml:"isLocal" json:"isLocal"`
}

// _links : Liveness link.
type _links struct {
	Self livenessLinktype `yaml:"self" json:"self"`
}

// Liveness link type.
type livenessLinktype struct {
	Liveness string `json:"liveness"`
}
