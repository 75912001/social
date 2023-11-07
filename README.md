# 社交系统 Social System


项目初始化

    go mod init social

安装grpc包

    go get google.golang.org/grpc

使用以下命令安装Go的协议编译器插件

    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
    相关link:https://grpc.io/docs/languages/go/quickstart/
    


## protobuf 工具

https://github.com/protocolbuffers/protobuf/releases

windows

    下载文件 protoc-25.0-win64.zip
    将 protoc-25.0-win64\bin\protoc.exe 配置到环境变量

生成代码

    protoc --go_out=. --go-grpc_out=. ./impl/protobuf/social/social.proto
    生成两个文件,一个用于普通的 Protocol Buffers,一个用于 gRPC
    相关路径定义在*.proto文件中
    如:option go_package = "proto/test";
    
