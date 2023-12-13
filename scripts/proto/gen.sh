#脚本所在目录
#SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
#echo "脚本所在目录: $SCRIPT_DIR"

#生成 协议
cd ../../proto || exit
protoc --go_out=../../ message.proto
protoc --go_out=../../ common.proto
protoc --go_out=../../ ec.proto

#gate
protoc --go_out=../../ gate_message.proto
protoc --go_out=../../ gate_server.proto
protoc --go-grpc_out=../../ gate_server.proto

cd - || exit

#生成 pkg/ec/error.go 文件
genCommonProtoError(){
  #目标文件
  desFile="../pkg/ec/ec.go"

  #echo "gen Common Proto"
  echo "/*Code generated by gen.sh. DO NOT EDIT.*/" > ${desFile}
  echo "package ec" >> "${desFile}"
  echo "import liberror \"social/pkg/lib/error\"" >> ${desFile}
  echo "var (" >> ${desFile}

  for arg in "$@"
  do
    #echo "start gen: $arg"
    fileName=${arg}.proto
	  tmpFile="../pkg/ec/ec.tmp"

    grep -e "EC_" -e -v "//EC_" "${fileName}" > "${tmpFile}"

	  sed -i -r 's/[\s\t]*//g' ${tmpFile}
	  sed -i -r 's/\/\// \/\//g' ${tmpFile}
	  sed -i -r 's/;/ tagName tagDesc/g' ${tmpFile}

	  awk -F " " '{printf("%s = &liberror.Error {Code: %s,Name: \"%s\",Desc: \"%s\" } \r\n",$1,$3,$1,$6)}' ${tmpFile} > ${tmpFile}.xxx
	  cat ${tmpFile}.xxx > ${tmpFile}
	  unlink ${tmpFile}.xxx

	  sed -i -r 's/\/\///g' ${tmpFile}
	  cat ${tmpFile} >> ${desFile}

	  unlink ${tmpFile}
  done

  echo ")" >> ${desFile}

  go fmt ${desFile}
}

#生成CMD
genCMD(){
  packageName=${1}
  #源文件
  sourceFile=${2}
  #目标文件
  desFile=${3}

  #echo "gen CMD ... ${packageName} ${sourceFile} ${desFile}"

  echo "/*Code generated by gen.sh. DO NOT EDIT.*/" > "${desFile}"
  echo "package ${packageName}" >> "${desFile}"
  grep "#" "${sourceFile}" >> "${desFile}"

  sed -i "s/message /const /g" "${desFile}"
  sed -i "s/\/\//_CMD uint32 = /g" "${desFile}"
  sed -i "s/#/ \/\/ /g" "${desFile}"

  go fmt "${desFile}"
}

#生成CMD映射文件
genCMDMapAll(){
  #目标路径
  desPath=../pkg/proto
  #目标文件
  desFile=${desPath}/cmd_map.go
  #echo "gen CMD Map All"
  echo "/*Code generated by gen.sh. DO NOT EDIT.*/" > ${desFile}
  echo "package proto" >> ${desFile}
  echo "var CMDMap = map[uint32]string{" >> ${desFile}

  for arg in "$@"
  do
    #echo "start gen cmd map: $arg"
    fileName=${arg}_message.proto
    packageName=${arg}
    grep "#" "${fileName}" > "${desPath}/${packageName}/cmd.tmp"

    #sed -i -r 's/^[ \t]*message[ \t]+(.*)\/\/(.*)\#(.*?)/\2 : "\1",/g' ${desPath}/${packageName}/cmd.tmp
    sed -i  -r 's/^\s*message[ \t]+(.*)\/\/(.*)\#(.*?)/\2 : "\1",/g' "${desPath}/${packageName}/cmd.tmp"

    sort -t: -k1 "${desPath}/${packageName}/cmd.tmp" | awk -F ":" '($1==CMD){}($1!=CMD){CMD=$1;print$1,":",$2}' >> ${desFile}
    unlink "${desPath}/${packageName}/cmd.tmp"
  done

  echo "}" >> ${desFile}

  go fmt ${desFile}
}

cd ../../proto || exit
genCommonProtoError common
genCMD gate gate_message.proto ../pkg/proto/gate/cmd.go
cd - || exit

#合并CMD map
cd ../../proto || exit
genCMDMapAll gate
cd - || exit

echo "按下任意键退出脚本"
read -n 1 -s -r -p "按下任意键..."

exit 0
