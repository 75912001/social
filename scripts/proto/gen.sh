#脚本所在目录
#SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
#echo "脚本所在目录: $SCRIPT_DIR"

#生成 协议
cd ../../proto || exit
protoc --go_out=../../ struct.proto
protoc --go_out=../../ common.proto
protoc --go_out=../../ ec.proto

#gate
protoc --go_out=../../ gate.message.proto
protoc --go_out=../../ gate.server.proto
protoc --go-grpc_out=../../ gate.server.proto

#friend
protoc --go_out=../../ friend.message.proto
protoc --go_out=../../ friend.enum.proto
protoc --go_out=../../ friend.struct.proto
protoc --go_out=../../ friend.server.proto
protoc --go-grpc_out=../../ friend.server.proto

cd - || exit

#生成 pkg/ec/ec.go 文件
genCommonProtoError(){
  #目标文件
  desFile="../pkg/ec/ec.go"

  #echo "gen Common Proto"
  {
      echo '/*Package ec*/'
      echo "/*Code generated by gen.sh. DO NOT EDIT.*/"
      echo "package ec"
      echo "import liberror \"social/lib/error\""
      echo "var ("
  } > ${desFile}

  for arg in "$@"
  do
    #echo "start gen: $arg"
    fileName=${arg}.proto
	  tmpFile="../pkg/ec/ec.tmp"

    #grep "EC_" "${fileName}" | grep -v "//EC_" > ${tmpFile}
    grep "#tag_desc:" "${fileName}" > ${tmpFile}

	  sed -i -r 's/[\t]+//g' ${tmpFile}
	  sed -i -r 's/\/\// \/\//g' ${tmpFile}
	  sed -i -r 's/;/ tagName tagDesc/g' ${tmpFile}
	  awk -F " " '{printf("%s = liberror.CreateError(%s,\"%s\",\"%s\") \r\n",$1,$3,$1,$6)}' ${tmpFile} > ${tmpFile}.xxx
	  sed -i 's/#tag_desc://g' ${tmpFile}.xxx
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

  {
    echo "/*Package ${packageName}*/"
    echo "/*Code generated by gen.sh. DO NOT EDIT.*/"
    echo "package ${packageName}"
  } > "${desFile}"
  grep "#" "${sourceFile}" >> "${desFile}"

  sed -i "s/message /const /g" "${desFile}"
  sed -i "s/\/\//CMD uint32 = /g" "${desFile}"
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
  {
    echo '/*Package proto*/'
    echo "/*Code generated by gen.sh. DO NOT EDIT.*/"
    echo "package proto"
    echo "var CMDMap = map[uint32]string{"
  } > ${desFile}

  for arg in "$@"
  do
    echo "start gen cmd map: ${arg}"
    fileName=${arg}.message.proto
    packageName=${arg}

    grep "#" "${fileName}" > "${desPath}/${packageName}/cmd.tmp"

    #sed -i -r 's/^[ \t]*message[ \t]+(.*)\/\/(.*)\#(.*?)/\2 : "\1",/g' ${desPath}/${packageName}/cmd.tmp
    sed -i  -r 's/^\s*message[ \t]+(.*)\/\/(.*)\#(.*?)/\2 : "\1",/g' "${desPath}/${packageName}/cmd.tmp"
    argTag=("\"${arg}.")
    #sort -t: -k1 "${desPath}/${packageName}/cmd.tmp" | awk -F ":" '($1==CMD){}($1!=CMD){CMD=$1;print$1,":",$2}' >> ${desFile}
    sort -t: -k1 "${desPath}/${packageName}/cmd.tmp" | awk -F ":" -v arg="${argTag[0]}" '($1==CMD){}($1!=CMD){CMD=$1; $2 = arg substr($2, 3); print $1 " : " $2}' >> ${desFile}


    unlink "${desPath}/${packageName}/cmd.tmp"
  done

  echo "}" >> ${desFile}

  go fmt ${desFile}
}

cd ../../proto || exit
genCommonProtoError ec
genCMD gate gate.message.proto ../pkg/proto/gate/gate.cmd.go
genCMD friend friend.message.proto ../pkg/proto/friend/friend.cmd.go
cd - || exit

#合并CMD map
cd ../../proto || exit
genCMDMapAll gate friend
#genCMDMapAll friend
cd - || exit

echo "按下任意键退出脚本"
read -n 1 -s -r -p "按下任意键..."

exit 0
