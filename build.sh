#!/bin/bash

export SCRIPTPATH=$(dirname $(realpath $0))
mkdir --parent ${SCRIPTPATH}/bin
export "CGO_ENABLED=0"

go build -o $SCRIPTPATH/bin/linux_amd64 $SCRIPTPATH/src/*.go
