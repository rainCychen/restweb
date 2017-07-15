#!/bin/bash
MQ_PATH=$(cd "$(dirname "$0")"; pwd)	
CUR_PATH=`pwd`
DEPS_PATH='/home/Golang/'
OLD_GOPATH=$GOPATH
export GOPATH=$DEPS_PATH':'$MQ_PATH
echo $GOPATH
echo $MQ_PATH
echo $CUR_PATH
cd $MQ_PATH
go install $* ./src/cmd/web
set GOPATH = $OLD_GOPATH 
cd $CUR_PATH