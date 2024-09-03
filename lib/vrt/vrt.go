package vrt

import (
	"errors"
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"
	"path"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	gs "github.com/kijimaD/ruins/lib/systems"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/engine/states"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/worldhelper/equips"
	"github.com/kijimaD/ruins/lib/worldhelper/material"
	"github.com/kijimaD/ruins/lib/worldhelper/spawner"
)

var regularTermination = errors.New("テスト環境における、想定どおりの終了")

type TestGame struct {
	game.MainGame
	gameCount  int
	outputPath string
}

func (g *TestGame) Update() error {
	// テストの前に実行される
	g.StateMachine.Update(g.World)

	// 10フレームだけ実行する。更新→描画の順なので、1度は更新しないと描画されない
	if g.gameCount < 10 {
		g.gameCount += 1
		return nil
	}

	// エラーを返さないと、実行終了しない
	return regularTermination
}

const outputDirName = "vrtimages"
const dirPerm = 0o755

func (g *TestGame) Draw(screen *ebiten.Image) {
	g.StateMachine.Draw(g.World, screen)

	// テストでは保存しない
	if flag.Lookup("test.v") != nil {
		return
	}

	_ = os.Mkdir(outputDirName, dirPerm)
	file, err := os.Create(path.Join(outputDirName, fmt.Sprintf("%s.png", g.outputPath)))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = png.Encode(file, screen)
	if err != nil {
		log.Fatal(err)
	}
}

func RunTestGame(state states.State, outputPath string) {
	world := game.InitWorld(960, 720)

	spawner.SpawnItem(world, "木刀", raw.SpawnInBackpack)
	spawner.SpawnItem(world, "ハンドガン", raw.SpawnInBackpack)
	spawner.SpawnItem(world, "レイガン", raw.SpawnInBackpack)
	armor := spawner.SpawnItem(world, "西洋鎧", raw.SpawnInBackpack)
	spawner.SpawnItem(world, "作業用ヘルメット", raw.SpawnInBackpack)
	spawner.SpawnItem(world, "革のブーツ", raw.SpawnInBackpack)
	spawner.SpawnItem(world, "ルビー原石", raw.SpawnInBackpack)
	spawner.SpawnItem(world, "回復薬", raw.SpawnInBackpack)
	spawner.SpawnItem(world, "回復薬", raw.SpawnInBackpack)
	spawner.SpawnItem(world, "回復スプレー", raw.SpawnInBackpack)
	spawner.SpawnItem(world, "回復スプレー", raw.SpawnInBackpack)
	spawner.SpawnItem(world, "手榴弾", raw.SpawnInBackpack)
	spawner.SpawnItem(world, "手榴弾", raw.SpawnInBackpack)
	spawner.SpawnItem(world, "手榴弾", raw.SpawnInBackpack)
	spawner.SpawnItem(world, "手榴弾", raw.SpawnInBackpack)
	ishihara := spawner.SpawnMember(world, "イシハラ", true)
	spawner.SpawnMember(world, "シラセ", true)
	spawner.SpawnMember(world, "ヒラヤマ", true)
	spawner.SpawnAllMaterials(world)
	material.PlusAmount("鉄", 40, world)
	material.PlusAmount("鉄くず", 4, world)
	material.PlusAmount("緑ハーブ", 2, world)
	material.PlusAmount("フェライトコア", 30, world)
	spawner.SpawnAllRecipes(world)
	equips.Equip(world, armor, ishihara, gc.EquipmentSlotZero)

	_ = gs.EquipmentChangedSystem(world)

	// 完全回復
	effects.AddEffect(nil, effects.Healing{Amount: gc.RatioAmount{Ratio: float64(1.0)}}, effects.Party{})
	effects.AddEffect(nil, effects.RecoveryStamina{Amount: gc.RatioAmount{Ratio: float64(1.0)}}, effects.Party{})

	effects.RunEffectQueue(world)

	g := &TestGame{
		MainGame: game.MainGame{
			World:        world,
			StateMachine: es.Init(state, world),
		},
		gameCount:  0,
		outputPath: outputPath,
	}

	if err := ebiten.RunGame(g); err != nil && err != regularTermination {
		log.Fatal(err)
	}
}
