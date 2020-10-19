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
    "github.com/smartystreets/goconvey/convey"
    "mep-agent/src/util"
    "testing"
)

func TestValidateIp(t *testing.T) {
    convey.Convey("ValidateIp", t, func() {
        convey.Convey("ValidateIpSuccess", func() {
            convey.So(util.ValidateIp("127.0.0.1"), convey.ShouldBeNil)
        })

        convey.Convey("ValidateIpFail", func() {
            convey.So(util.ValidateIp("127.0.0"), convey.ShouldNotBeNil)
        })
    })
}

func TestValidateDomainName(t *testing.T) {
    convey.Convey("ValidateDomainName", t, func() {
        convey.Convey("ValidateDomainNameSuccess", func() {
            convey.So(util.ValidateDomainName("abcd.com"), convey.ShouldBeNil)
        })

        convey.Convey("ValidateDomainNameFail", func() {
            convey.So(util.ValidateDomainName("abcd$!!"), convey.ShouldNotBeNil)
            convey.So(util.ValidateDomainName("oikYVgrRbDZHZSaobOTo8ugCKsUSdVeMsg2d9b7Qr250q2HNBiET4WmecJ0MFavRA0cBzOWu8sObLha17auHoy6ULbAOgP50bDZapxOylTbr1kq8Z4m8uMztciGtq4e11GA0aEh0oLCR3kxFtV4EgOm4eZb7vmEQeMtBy4jaXl6miMJugoRqcfLo9ojDYk73lbCaP9ydUkO56fw8dUUYjeMvrzmIZPLdVjPm62R4AQFQ4CEs7vp6xafx9dRwPoym"), convey.ShouldNotBeNil)
        })
    })
}

func TestValidateAkSk(t *testing.T) {
    convey.Convey("ValidateAkSk", t, func() {
        convey.Convey("ValidateAkSkSuccess", func() {
            sk_sucess := []byte("DXPb4sqElKhcHe07Kw5uorayETwId1JOjjOIRomRs5wyszoCR5R7AtVa28KT3lSc")
            convey.So(util.ValidateAkSk("QVUJMSUMgS0VZLS0tLS0", &sk_sucess), convey.ShouldBeNil)
        })

        convey.Convey("ValidateAkSkFail", func() {
            sk_failed := []byte("DXPb4sqElKhcHe07Kw5uorayETwId1JOjjOIRomRs5wyszoCR5R7AtVa28KT3lScDXPb4sqElKhcHe07Kw5uorayETwId1JOjjOIRomRs5wyszoCR5R7AtVa28KT3lSc")
            convey.So(util.ValidateAkSk("QVUJMSUMgS0VZLS0tLS0", &sk_failed), convey.ShouldNotBeNil)
            convey.So(util.ValidateAkSk("DXPb4sqElKhcHe07Kw5uorayETwId1JOjjOIRomRs5wyszoCR5R7AtVa28KT3lSc", &sk_failed), convey.ShouldNotBeNil)
        })
    })
}
