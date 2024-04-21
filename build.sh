#!/bin/bash
set -eux

##################################
# クロスコンパイルするスクリプト
##################################

SOURCE=app
BUILD_STAGE_TARGET=base
BUILDER_IMAGE_NAME=base
APP_NAME=ruins

APP_VERSION=v0.0.0

cd `dirname $0`

# ================

function is_git_repo {
    echo `git rev-parse --is-inside-work-tree`
}

if [ $(is_git_repo) = "true" ]; then
    APP_VERSION=`git describe --tag --abbrev=0`
else
    APP_VERSION=`cat ../.versions`
fi

# ================

# cmd <command> <GOOS> <GOARCH> <CGO>
cmd() {
    output=$1
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
           go build -o $output -buildvcs=false -ldflags "-X github.com/kijimaD/ruins/lib/utils/consts.AppVersion=$APP_VERSION" .
}

start() {
    docker build . --target $BUILD_STAGE_TARGET -t $BUILDER_IMAGE_NAME

    cmd "${APP_NAME}_linux_amd64" linux amd64 1
    # cmd "${APP_NAME}_linux_arm64" linux arm64 1

    cmd "${APP_NAME}_windows_amd64" windows amd64 0
    # cmd "${APP_NAME}_windows_arm64" windows arm64 0

    cmd "${APP_NAME}_js_wasm" js wasm 0
}

start
