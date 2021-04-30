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

// define the type information
package model

type AppInstanceInfo struct {
	ServiceInfoPosts                         []ServiceInfoPost                         `yaml:"serviceInfoPosts" json:"serviceInfoPosts"`
	SerAvailabilityNotificationSubscriptions []SerAvailabilityNotificationSubscription `yaml:"serAvailabilityNotificationSubscriptions" json:"serAvailabilityNotificationSubscriptions"`
}

// Service Information to be registered.
type ServiceInfoPost struct {
	SerName           string         `yaml:"serName" json:"serName"`
	SerCategory       CategoryRef    `yaml:"serCategory" json:"serCategory"`
	Version           string         `yaml:"version" json:"version"`
	State             ServiceState   `yaml:"state" json:"state"`
	TransportId       string         `yaml:"transportId" json:"transportId"`
	TransportInfo     TransportInfo  `yaml:"transportInfo" json:"transportInfo"`
	Serializer        SerializerType `yaml:"serializer" json:"serializer"`
	ScopeOfLocality   LocalityType   `yaml:"scopeOfLocality" json:"scopeOfLocality"`
	ConsumedLocalOnly bool           `yaml:"consumedLocalOnly" json:"consumedLocalOnly"`
	IsLocal           bool           `yaml:"isLocal" json:"isLocal"`
	LivenessInterval  int            `yaml:"livenessInterval" json:"livenessInterval,omitempty"`
	Links            _links       	 `json:"_links,omitempty"`
	SerInstanceId     string         `json:"serInstanceId,omitempty"`
}

type ServiceInfo struct {
	SerName           string         `json:"serName"`
	Version           string         `json:"version"`

}

type CategoryRef struct {
	Href    string `yaml:"href" json:"href"`
	Id      string `yaml:"id" json:"id"`
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
}

type ServiceState string

// Transport Information of the service to be registered.
type TransportInfo struct {
	Id               string           `yaml:"id" json:"id"`
	Name             string           `yaml:"name" json:"name"`
	Description      string           `yaml:"description" json:"description"`
	TransportType    TransportType    `yaml:"type" json:"type"`
	Protocol         string           `yaml:"protocol" json:"protocol"`
	Version          string           `yaml:"version" json:"version"`
	Endpoint         EndPointInfo     `yaml:"endpoint" json:"endpoint"`
	Security         SecurityInfo     `yaml:"security" json:"security"`
	ImplSpecificInfo ImplSpecificInfo `yaml:"implSpecificInfo" json:"implSpecificInfo"`
}

type TransportType string

// Endpoint of the service to be registered.
type EndPointInfo struct {
	Addresses   []EndPointInfoAddress `yaml:"addresses" json:"addresses"`
}

type EndPointInfoAddress struct {
	Host string `yaml:"host" json:"host"`
	Port uint32 `yaml:"port" json:"port"`
}

type SecurityInfo struct {
	OAuth2Info SecurityInfoOAuth2Info `yaml:"oAuth2Info" json:"oAuth2Info"`
}

type SecurityInfoOAuth2Info struct {
	GrantTypes    []SecurityInfoOAuth2InfoGrantType `yaml:"grantTypes" json:"grantTypes"`
	TokenEndpoint string                            `yaml:"tokenEndpoint" json:"tokenEndpoint"`
}

type SecurityInfoOAuth2InfoGrantType string

type ImplSpecificInfo struct {
}

type SerializerType string

type LocalityType string

type SerAvailabilityNotificationSubscription struct {
	SubscriptionType  string                                                   `yaml:"subscriptionType" json:"subscriptionType"`
	CallbackReference string                                                   `yaml:"callbackReference" json:"callbackReference"`
	Links             Self                                                     `yaml:"links" json:"links"`
	FilteringCriteria SerAvailabilityNotificationSubscriptionFilteringCriteria `yaml:"filteringCriteria" json:"filteringCriteria"`
}

type Self struct {
	Self LinkType `yaml:"self" json:"self"`
}

type LinkType struct {
	Href string `yaml:"href" json:"href"`
}

type SerAvailabilityNotificationSubscriptionFilteringCriteria struct {
	SerInstanceIds []string       `yaml:"serInstanceIds" json:"serInstanceIds"`
	SerNames       []string       `yaml:"serNames" json:"serNames"`
	SerCategories  []CategoryRef  `yaml:"serCategories" json:"serCategories"`
	States         []ServiceState `yaml:"states" json:"states"`
	IsLocal        bool           `yaml:"isLocal" json:"isLocal"`
}

type _links struct {
	Self LivenessLinktype `yaml:"self" json:"self"`
}

type LivenessLinktype struct {
	Liveness string `json:"liveness"`
}
