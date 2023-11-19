#脚本所在目录
#SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
#echo "脚本所在目录: $SCRIPT_DIR"

#生成 协议
cd ../../proto
protoc --go_out=../../ message.proto
protoc --go_out=../../ common.proto

#gate
protoc --go_out=../../ gate_message.proto
protoc --go_out=../../ gate_server.proto
protoc --go-grpc_out=../../ gate_server.proto

cd -



