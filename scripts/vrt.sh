#!/bin/bash
set -eux

###########################
# 各ステートのキャプチャを取る。
###########################
# Goテスト内で条件を変えてスクショしたいのだが、Ebitenの制限で1プロセス内で2回起動ができなかった。シェルスクリプトで複数プロセスで実行することにした。各ステートで簡単にスクショを取るだけだが、ないよりはマシだろう...

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
cmd FieldMenu
cmd HomeMenu
cmd Intro
cmd InventoryMenu
cmd MainMenu
cmd Dungeon
