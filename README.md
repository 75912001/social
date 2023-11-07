# 社交系统 Social System


## TODO
>> 

项目初始化

    go mod init social

安装包

    go get google.golang.org/grpc
    go get -v  github.com/pkg/errors

使用以下命令安装Go的协议编译器插件

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

    protoc --go_out=. --go-grpc_out=. ./impl/protobuf/social/social.proto
    生成两个文件,一个用于普通的 Protocol Buffers,一个用于 gRPC
    相关路径定义在*.proto文件中
    如:option go_package = "proto/test";


# 目录说明

##
```
├── api        --对外接口实现 
│   └── server
│       ├── methods.go    --RPC方法的入口
│       └── service.go   
├── bin         --二进制执行文件
│   └── server
|       └──log  --日志
├── cmd         --Main入口 
│   └── server
├── internal    --游戏服务业务逻辑 
│   └──  server
│        └── internal   --该目录下具体实现服务的各模块，内部可拆分子目录（如hero、skin等）
│            └── property    --静态数据存放在此目录下
│            └── skin
│            └── hero
├── pkg         --通用工具包
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
```


