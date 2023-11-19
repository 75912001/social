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

#生成 common_proto error 文件
genCommonProtoError(){
  #目标文件
  desFile="../pkg/error_code/error.go"

  echo "gen Common Proto"
  echo "/*Code generated by gen.sh. DO NOT EDIT.*/" > ${desFile}
  echo "package error_code" >> ${desFile}
  echo "import xrerror \"social/pkg/lib/error\"" >> ${desFile}
  echo "var (" >> ${desFile}

  for arg in $*
  do
    echo "start gen: $arg"
    fileName=${arg}.proto
	  tmpFile="../pkg/error_code/error.tmp"

    cat ${fileName} | grep "EC_" | grep -v "//EC_" > ${tmpFile}
	  sed -i -r 's/[\s\t]*//g' ${tmpFile}
	  sed -i -r 's/\/\// \/\//g' ${tmpFile}
	  sed -i -r 's/;/ tagName tagDesc/g' ${tmpFile}

	  awk -F " " '{printf("%s = &xrerror.Error {Code: %s,Name: \"%s\",Desc: \"%s\" } \r\n",$1,$3,$1,$6)}' ${tmpFile} > ${tmpFile}.xxx
	  cat ${tmpFile}.xxx > ${tmpFile}
	  unlink ${tmpFile}.xxx

	  sed -i -r 's/\/\///g' ${tmpFile}
	  cat ${tmpFile} >> ${desFile}

	  unlink ${tmpFile}
  done

  echo ")" >> ${desFile}

  go fmt ${desFile}
}

cd ../../proto
genCommonProtoError common
cd -
