# 社交系统 Social System
提供社交服务,其中包括好友系统,博客系统



## TODO
>> robot 测试 
>>> 1.注册
>>> 2.注销
> >> 使用多客户端,不间断操作.并发送简单消息.
> 
> 整合log系统到social中
> 
> 在gate中创建管理器, 管理 friend system
> 

# 使用说明


## 项目初始化

    go mod init social

## 安装包

    go get google.golang.org/grpc
    go get -v github.com/pkg/errors
    go get go.etcd.io/etcd/client/v3
    go get github.com/google/uuid
    *go get github.com/agiledragon/gomonkey


## 使用以下命令安装Go的协议编译器插件

    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
    相关link:https://grpc.io/docs/languages/go/quickstart/

## protobuf 工具

https://github.com/protocolbuffers/protobuf/releases

windows

    下载文件 protoc-25.0-win64.zip
    将 protoc-25.0-win64\bin\protoc.exe 配置到环境变量

由*.proto文件生成代码
    
    运行 social/scripts/proto/gen.sh

## etcd windows部署(单例)

    下载文件
    https://github.com/etcd-io/etcd/releases/tag/v3.5.10
#
    启动服务
    解压运行 etcd.exe
#
    使用
    增
    ./etcdctl --endpoints=127.0.0.1:2379 put /test "Hello etcd"
    查
    ./etcdctl --endpoints=127.0.0.1:2379 get --prefix /
    删
    ./etcdctl --endpoints=1270.0.01:2379 del --prev-kv --prefix /

# 目录说明

##
```
├── api        --对外接口实现 
│   └── server
│       ├── methods.go    --RPC方法的入口
│       └── service.go   
├── bin         --二进制执行文件
│   └── server
|       └──log  --日志目录
|       └──bench.json.default  --服务配置(正式为bench.json)
├── cmd         --Main入口 
│   └── main  --程序入口
├── internal    --游戏服务业务逻辑 
│   └──  server
│        └── internal   --该目录下具体实现服务的各模块，内部可拆分子目录（如hero、skin等）
│            └── property    --静态数据存放在此目录下
│            └── skin
│            └── hero
├── pkg         --通用工具包
│   ├── bench   --配置文件方法定义
│   ├── common  --通用方法定义
│   ├── consts  --常量定义
│   ├── deps    --依赖封装
│   ├── ec      --错误码定义
│   └── proto   --各服务的proto文件生成的代码
├── proto       --PROTO文件
|   ├── server   
|        ├── server.proto    --service和rpc
│        └── message.proto  --服务的message
├── scripts     --脚本目录
│   ├── tpl
│   ├── proto  --生成协议
│   └── sql 
├── third_party --第三方依赖
│   └── protobuffer
│       ├── go-proto-validators
│       ├── google
│       │   ├── api
│       │   └── rpc
│       └── tevat
│           └── api
├── tools       --项目工具
└── vendor
└── test  --测试程序
```




