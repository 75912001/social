#!/bin/bash

cd ../
echo "当前目录是:$PWD"
gofmt -w -l -s .
cd - || exit
echo "执行完毕"
exit 0
