#!/bin/bash

IP1=10.18.99.102
IP2=10.18.32.71
IP3=10.18.32.72
PORT=57522
USER=root

echo "build world service ..."

CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./dawn-robot-linux.exe ./main.go

################################################
#IP2
ssh ${USER}@${IP2} -p ${PORT} << EOF
mkdir -p /data/yoozoogame/code/robot1/log
mkdir -p /data/yoozoogame/code/robot2/log
mkdir -p /data/yoozoogame/code/robot3/log
mkdir -p /data/yoozoogame/code/robot4/log
mkdir -p /data/yoozoogame/log
mkdir -p /data/yoozoogame/admin/robot1
mkdir -p /data/yoozoogame/admin/robot2
mkdir -p /data/yoozoogame/admin/robot3
mkdir -p /data/yoozoogame/admin/robot4

cd /data/yoozoogame/admin/robot1
/bin/bash stop.sh
cd -

cd /data/yoozoogame/admin/robot2
/bin/bash stop.sh
cd -

cd /data/yoozoogame/admin/robot3
/bin/bash stop.sh
cd -

cd /data/yoozoogame/admin/robot4
/bin/bash stop.sh
cd -

EOF

scp -P ${PORT} ./dawn-robot-linux.exe ${USER}@${IP2}:/data/yoozoogame/code/robot1/dawn-robot.exe.1


################################################
#see
echo "#############################${IP2}"
ssh  ${USER}@${IP2} -p ${PORT} "ps -ef | grep dawn"


echo "build world service done."
