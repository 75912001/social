#!/bin/bash

cd ../
echo "当前目录是:$PWD"
gofmt -w -l -s .
cd - || exit

echo "按下任意键退出脚本"
read -n 1 -s -r -p "按下任意键..."

exit 0