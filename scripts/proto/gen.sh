#脚本所在目录
#SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
#echo "脚本所在目录: $SCRIPT_DIR"

#生成 gate 协议
cd ../../proto/gate
protoc --go_out=. --go-grpc_out=. *.proto
cd -