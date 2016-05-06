#!/usr/bin/env bash

set -e

PJ=$(cd `dirname $0`; pwd)

SRV=$PJ/server/src/chat_server
BIN=$PJ/server/bin

OLD_GOPATH=$GOPATH
export GOPATH=$PJ/server

echo "compiling server..."
cd $SRV

echo $GOPATH
echo $SRV
go install chat_server

export GOPATH=$OLD_GOPATH

echo "Successed in `date "+%Y-%m-%d %H:%M:%S"`"
