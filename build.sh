#!/bin/bash

export CGO_ENABLED=0
SCRIPT=$(readlink -f "$0")
SCRIPTPATH=$(dirname "$SCRIPT")

export GOARCH=arm && export GOOS=linux && go build -o $SCRIPTPATH/bin/linux_arm $SCRIPTPATH/src/*.go
export GOARCH=amd64 && export GOOS=linux && go build -o $SCRIPTPATH/bin/linux_amd64 $SCRIPTPATH/src/*.go
export GOARCH=386 && export GOOS=linux && go build -o $SCRIPTPATH/bin/linux_i386 $SCRIPTPATH/src/*.go
export GOARCH=mips64 && export GOOS=linux && go build -o $SCRIPTPATH/bin/linux_mips64 $SCRIPTPATH/src/*.go
export GOARCH=mips64le && export GOOS=linux && go build -o $SCRIPTPATH/bin/linux_mips64le $SCRIPTPATH/src/*.go

export GOARCH=arm && export GOOS=freebsd && go build -o $SCRIPTPATH/bin/freebsd_arm $SCRIPTPATH/src/*.go

export GOARCH=386 && export GOOS=freebsd && go build -o $SCRIPTPATH/bin/freebsd_i386 $SCRIPTPATH/src/*.go
export GOARCH=amd64 && export GOOS=freebsd && go build -o $SCRIPTPATH/bin/freebsd_amd64 $SCRIPTPATH/src/*.go

export GOARCH=amd64 && export GOOS=windows && go build -o $SCRIPTPATH/bin/windows_amd64.exe $SCRIPTPATH/src/*.go
export GOARCH=386 && export GOOS=windows && go build -o $SCRIPTPATH/bin/windows_i386.exe $SCRIPTPATH/src/*.go

export GOARCH=386 && export GOOS=darwin && go build -o $SCRIPTPATH/bin/darwin_i386 $SCRIPTPATH/src/*.go
export GOARCH=amd64 && export GOOS=darwin && go build -o $SCRIPTPATH/bin/darwin_amd64 $SCRIPTPATH/src/*.go
