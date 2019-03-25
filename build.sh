#!/bin/bash

export SCRIPTPATH=$(dirname $(realpath $0))
export GOPATH=${SCRIPTPATH}/deps
mkdir --parent ${GOPATH}
mkdir --parent ${SCRIPTPATH}/bin

export "CGO_ENABLED=0"

go get github.com/vishvananda/netlink

export GOARCH=amd64 && export GOOS=linux && go build -o $SCRIPTPATH/bin/linux_amd64 $SCRIPTPATH/src/*.go
