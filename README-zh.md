# Mep-Agent

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![Jenkins](https://img.shields.io/jenkins/build?jobUrl=http%3A%2F%2Fjenkins.edgegallery.org%2Fview%2FMEC-PLATFORM-BUILD%2Fjob%2Fmep-agent-docker-image-build-update-daily-master%2F)

## Introduction
Mep-Agent is a middleware that provides proxy services for third-party apps. It can help apps, which do not implement the ETSI interface to register to MEP, and realize app service registration and discovery.
Mep-Agent will start at the same time as the application container, and read the content in the file conf/app_instance_info.yaml to automatically register the service.

## MEP-Agent code directory

```
.      
├─conf
├─docker
└─src
    ├─config
    ├─main
    ├─model
    ├─service
    ├─test
    └─util
```

Above is the directory tree of MEP-Agent project, their usage is as belows:
- conf: mep-agent config file 
- docker: dockerfile
- src: source code
- config: config files
- main: main method
- model: mep1 service register model and mepauth authorization model
- service: request for token and register service
- test: unit test
- util: util tool file

## Build & Run

Mep-Agent is developed by the Go language and provides services in the form of a docker image. When it starts, it will read the configuration file and register the App to the MEP to realize service registration and discovery.

- ### Build

    git clone from mep-agent master repo
    ```
    git clone https://gitee.com/edgegallery/mep-agent.git
    ```
  
    build the mep-agent image
    ```
    docker build -t mep-agent:latest -f docker/Dockerfile .
    ```
  
- ### Run

    Prepare the certificate files and mepagent.properties, which contains ACCESS_KEY and SECRET_KEY, and run with
    ```
    docker run -tid --name mepagent \
               -e MEP_IP=<host IP> \
               -e MEP_APIGW_PORT=8443 \
               -e ENABLE_WAIT=true \
               -v mepagent.properties:/usr/mep/mepagent.properties \
               mep-agent:latest
    ```

More details of the building and installation process please refer to [HERE](https://gitee.com/edgegallery/docs/blob/master/MEP/EdgeGallery%E6%9C%AC%E5%9C%B0%E5%BC%80%E5%8F%91%E9%AA%8C%E8%AF%81%E6%9C%8D%E5%8A%A1%E8%AF%B4%E6%98%8E%E4%B9%A6.md).
  
## Notice

Mep-Agent is written in Go language. In order to minimize the image, it adopts the process of statically compiling and then packaging, without relying on the basic Go language image, which greatly reduces the size of the image.
