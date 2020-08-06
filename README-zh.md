# Mep-Agent

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![Jenkins](https://img.shields.io/jenkins/build?jobUrl=http%3A%2F%2Fjenkins.edgegallery.org%2Fview%2FMEC-PLATFORM-BUILD%2Fjob%2Fmep-agent-docker-image-build-update-daily-master%2F)

## 介绍
Mep-Agent是为第三方应用程序提供代理服务的中间件。它可以帮助未实现ETSI接口的应用注册到MEP，并实现应用服务注册和发现。Mep-Agent将与应用程序容器同时启动，并读取conf/app_instance_info.yaml文件中的内容以自动注册服务。


## MEP-Agent代码目录

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

上面是MEP-Agent项目的目录树，其用法如下：
- conf: mep-agent配置文件 
- docker: dockerfile
- src: 源代码
- config: 配置文件
- main: main方法
- model: mep1服务注册模型和mepauth授权模型
- service: 请求令牌和注册服务
- test: 单元测试
- util: util工具文件

## 构建以及运行

Mep-Agent由Go语言开发，并以docker映像的形式提供服务。当启动时，它将读取配置文件，并将应用程序注册到MEP以实现服务注册和发现。


- ### 构建

    git clone mep-agent代码
    ```
    git clone https://gitee.com/edgegallery/mep-agent.git
    ```
  
    构建mep-agent镜像
    ```
    docker build -t mep-agent:latest -f docker/Dockerfile .
    ```
  
- ### 运行

    准备包含ACCESS_KEY和SECRET_KEY的证书文件和mepagent.properties，并执行
    ```
    docker run -tid --name mepagent \
               -e MEP_IP=<host IP> \
               -e MEP_APIGW_PORT=8443 \
               -e ENABLE_WAIT=true \
               -v mepagent.properties:/usr/mep/mepagent.properties \
               mep-agent:latest
    ```

有关mepagent.properties以及构建和安装过程的更多详细信息，请参阅 [HERE](https://gitee.com/edgegallery/docs/blob/master/MEP/EdgeGallery%E6%9C%AC%E5%9C%B0%E5%BC%80%E5%8F%91%E9%AA%8C%E8%AF%81%E6%9C%8D%E5%8A%A1%E8%AF%B4%E6%98%8E%E4%B9%A6.md).
  
## 注意

Mep-Agent用Go语言编写。为了使镜像最小化，它采用了静态编译然后打包的过程，而不依赖于基本的Go语言镜像，从而大大减小了镜像的大小。
