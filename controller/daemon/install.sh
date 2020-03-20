#!/bin/bash
echo "building binary"

pushd ../
./build
BUILDEXC=$?
popd

if [ ${BUILDEXC} -ne 0 ]
then
	echo "Failed to build binary"
	exit 1
fi

#BINARYSRC="../server"
#BINARYDES="/usr/local/sbin"

#echo "copying ${BINARYSRC} to ${BINARYDES}"
#cp ${BINARYSRC} ${BINARYDES}

SVC=audio-reactive-led-strip-server.service
echo "installing ${SVC}..."

systemctl is-enabled ${SVC}

#if [ $? -eq 0 ]
#then
#    echo "${SVC} already installed and enabled"
#    exit 0
#fi

SRC="${PWD}/${SVC}"
DEST="/lib/systemd/system/"

echo "copying ${SRC} to ${DEST}"
cp ${SRC} ${DEST}

if [ $? -ne 0 ]
then
    echo "mv failed"
    exit 1
fi

systemctl start ${SVC}
if [ $? -ne 0 ]
then
    echo "failed to start ${SVC}"
    exit 1
fi

systemctl enable ${SVC}
if [ $? -ne 0 ]
then
    echo "failed to enable ${SVC} at startup"
    exit 1
fi

echo "successfully installed ${SVC}"

systemctl daemon-reload
if [ $? -ne 0 ]
then
	echo "failed to reload daemon"
	exit 1
fi

exit 0
