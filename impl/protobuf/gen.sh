#脚本所在目录
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
#echo "脚本所在目录: $SCRIPT_DIR"

protoc --go_out=. --go-grpc_out=. ./social/social.proto