#!/bin/bash
set -eux

##########################
# 各ステートのキャプチャを取る
##########################

cd `dirname $0`
cd ../

RUN=${RUN:-"go run ."}

cmd() {
    state=$1

    $RUN screenshot $state
}

cmd CraftMenu
cmd Intro
cmd DebugMenu
cmd DungeonSelect
cmd EquipMenu
cmd Field
cmd FieldMenu
cmd HomeMenu
cmd Intro
cmd InventoryMenu
cmd MainMenu
