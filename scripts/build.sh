#!/bin/bash
set -eux

##################################
# クロスコンパイルするスクリプト
##################################

# 定数的
SOURCE=app
BUILD_STAGE_TARGET=base
BUILDER_IMAGE_NAME=base
APP_NAME=ruins

# 変数的
APP_VERSION=v0.0.0

cd `dirname $0`
cd ../

# ================

function is_git_repo {
    echo `git rev-parse --is-inside-work-tree`
}

if [ $(is_git_repo) = "true" ]; then
    APP_VERSION=`git describe --tag --abbrev=0`
fi

# ================

# cmd <output> <GOOS> <GOARCH> <CGO>
cmd() {
    output=$1
    goos=$2
    goarch=$3
    cgo=$4

    docker run \
           --rm \
           -w /work \
           -v $PWD:/work \
           -v $HOME/go/pkg/mod:/go/pkg/mod \
           -v $HOME/.cache/go-build:/root/.cache/go-build \
           --env GOOS=$goos \
           --env GOARCH=$goarch \
           --env CGO_ENABLED=$cgo \
           $BUILDER_IMAGE_NAME \
           go build -o $output -buildvcs=false -ldflags "-X github.com/kijimaD/ruins/lib/consts.AppVersion=$APP_VERSION" .
}

start() {
    docker build . --target $BUILD_STAGE_TARGET -t $BUILDER_IMAGE_NAME

    cmd "bin/${APP_NAME}_linux_amd64" linux amd64 1
    cmd "bin/${APP_NAME}_windows_amd64" windows amd64 0

    # no such instruction になる...
    # cmd "${APP_NAME}_linux_arm64" linux arm64 1
    # cmd "${APP_NAME}_windows_arm64" windows arm64 0

    cmd "wasm/game.wasm" js wasm 0
}

start
