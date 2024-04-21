#!/bin/bash
set -eux

##########################
# 各ステートのキャプチャを取る
##########################

cd `dirname $0`
cd ../

cmd() {
    state=$1

    go run . screenshot $state
}

cmd Intro
