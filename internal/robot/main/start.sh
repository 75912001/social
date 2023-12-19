#!/bin/bash

IP1=10.18.99.102
IP2=10.18.32.71
IP3=10.18.32.72
PORT=57522
USER=root

echo "start world service ..."

################################################
#IP2
ssh ${USER}@${IP2} -p ${PORT} << EOF

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

ssh ${USER}@${IP2} -p ${PORT} << EOF

cd /data/yoozoogame/admin/robot1
/bin/bash start.sh
cd -

ln -sf /data/yoozoogame/code/robot1/dawn-robot.exe.1 /data/yoozoogame/code/robot2/dawn-robot.exe.2

cd /data/yoozoogame/admin/robot2
/bin/bash start.sh
cd -

ln -sf /data/yoozoogame/code/robot1/dawn-robot.exe.1 /data/yoozoogame/code/robot3/dawn-robot.exe.3

cd /data/yoozoogame/admin/robot3
/bin/bash start.sh
cd -

ln -sf /data/yoozoogame/code/robot1/dawn-robot.exe.1 /data/yoozoogame/code/robot4/dawn-robot.exe.4

cd /data/yoozoogame/admin/robot4
/bin/bash start.sh
cd -

EOF

################################################
#see
echo "#############################${IP2}"
ssh  ${USER}@${IP2} -p ${PORT} "ps -ef | grep dawn"


echo "build world service done."
