#!/bin/bash
set -eux

##################################
# クロスコンパイルするスクリプト
##################################

SOURCE=app
BUILD_STAGE_TARGET=base
BUILDER_IMAGE_NAME=base

cd `dirname $0`

# cmd <command> <GOOS> <GOARCH> <CGO>
cmd() {
    command=$1
    goos=$2
    goarch=$3
    cgo=$4

    docker run \
           --rm \
           -v $PWD:/work \
           -w /work \
           --env GOOS=$goos \
           --env GOARCH=$goarch \
           --env CGO_ENABLED=$cgo \
           $BUILDER_IMAGE_NAME \
           $command
}

start() {
    docker build . --target $BUILD_STAGE_TARGET -t $BUILDER_IMAGE_NAME

    # cmd "go build -buildvcs=false -o ${SOURCE}_darwin_amd64 ." darwin amd64 1
    # cmd "go build -o ${SOURCE}_darwin_arm64 ." darwin arm64 1

    cmd "go build -buildvcs=false -o ${SOURCE}_linux_amd64 ." linux amd64 1
    # cmd "go build -buildvcs=false -o ${SOURCE}_linux_arm64 ." linux arm64 1

    cmd "go build -buildvcs=false -o ${SOURCE}_windows_amd64.exe ." windows amd64 0
    # cmd "go build -buildvcs=false -o ${SOURCE}_windows_arm64.exe ." windows arm64 0

    cmd "go build -buildvcs=false -o ${SOURCE}_js_wasm.wasm ." js wasm 0
}

start
