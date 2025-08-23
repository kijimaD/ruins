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
	"github.com/kijimaD/ruins/lib/config"
	"github.com/kijimaD/ruins/lib/effects"
	gs "github.com/kijimaD/ruins/lib/systems"

	"github.com/hajimehoshi/ebiten/v2"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/kijimaD/ruins/lib/worldhelper"
)

// エラーを返さないと実行終了しないため
var errRegularTermination = errors.New("テスト環境における、想定どおりの終了")

// TestGame はビジュアルリグレッションテスト用のゲーム構造体
type TestGame struct {
	game.MainGame
	gameCount  int
	outputPath string
}

// Update はゲームの更新処理を行う
func (g *TestGame) Update() error {
	// テストの前に実行される
	g.StateMachine.Update(g.World)

	// 10フレームだけ実行する。更新→描画の順なので、1度は更新しないと描画されない
	if g.gameCount < 10 {
		g.gameCount++
		return nil
	}

	// エラーを返さないと、実行終了しない
	return errRegularTermination
}

const outputDirName = "vrtimages"
const dirPerm = 0o755

// Draw はゲームの描画処理を行う
func (g *TestGame) Draw(screen *ebiten.Image) {
	g.StateMachine.Draw(g.World, screen)

	// テストでは保存しない
	if flag.Lookup("test.v") != nil {
		return
	}

	if err := os.Mkdir(outputDirName, dirPerm); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	file, err := os.Create(path.Join(outputDirName, fmt.Sprintf("%s.png", g.outputPath)))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}()

	err = png.Encode(file, screen)
	if err != nil {
		log.Fatal(err)
	}
}

// RunTestGame はテストゲームを実行してスクリーンショットを保存する
func RunTestGame(state es.State, outputPath string) {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}
	// VRT用にアニメーションを無効化
	cfg.DisableAnimation = true

	world, err := game.InitWorld(960, 720)
	if err != nil {
		panic(fmt.Sprintf("InitWorld failed: %v", err))
	}

	worldhelper.SpawnItem(world, "木刀", gc.ItemLocationInBackpack)
	worldhelper.SpawnItem(world, "ハンドガン", gc.ItemLocationInBackpack)
	worldhelper.SpawnItem(world, "レイガン", gc.ItemLocationInBackpack)
	armor := worldhelper.SpawnItem(world, "西洋鎧", gc.ItemLocationInBackpack)
	worldhelper.SpawnItem(world, "作業用ヘルメット", gc.ItemLocationInBackpack)
	worldhelper.SpawnItem(world, "革のブーツ", gc.ItemLocationInBackpack)
	worldhelper.SpawnItem(world, "ルビー原石", gc.ItemLocationInBackpack)
	worldhelper.SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)
	worldhelper.SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)
	worldhelper.SpawnItem(world, "回復スプレー", gc.ItemLocationInBackpack)
	worldhelper.SpawnItem(world, "回復スプレー", gc.ItemLocationInBackpack)
	worldhelper.SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
	worldhelper.SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
	worldhelper.SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
	worldhelper.SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
	ishihara := worldhelper.SpawnMember(world, "イシハラ", true)
	worldhelper.SpawnMember(world, "シラセ", true)
	worldhelper.SpawnMember(world, "タチバナ", true)
	worldhelper.SpawnMember(world, "ハンス", false)
	worldhelper.SpawnMember(world, "カイン", false)
	worldhelper.SpawnMember(world, "メイ", false)
	worldhelper.SpawnAllMaterials(world)
	worldhelper.PlusAmount("鉄", 40, world)
	worldhelper.PlusAmount("鉄くず", 4, world)
	worldhelper.PlusAmount("緑ハーブ", 2, world)
	worldhelper.PlusAmount("フェライトコア", 30, world)
	worldhelper.SpawnAllRecipes(world)
	worldhelper.Equip(world, armor, ishihara, gc.EquipmentSlotNumber(0))

	// 装備変更後にステータスを更新
	if changed := gs.EquipmentChangedSystem(world); !changed {
		// 装備変更が期待されていた場合はログ出力
		log.Println("Equipment change was not detected")
	}

	// 完全回復
	processor := effects.NewProcessor()
	healingEffect := effects.Healing{Amount: gc.RatioAmount{Ratio: float64(1.0)}}
	staminaEffect := effects.RestoreStamina{Amount: gc.RatioAmount{Ratio: float64(1.0)}}

	partySelector := effects.TargetParty{}
	if err := processor.AddTargetedEffect(healingEffect, nil, partySelector, world); err != nil {
		log.Printf("回復エフェクト追加エラー: %v", err)
	}
	if err := processor.AddTargetedEffect(staminaEffect, nil, partySelector, world); err != nil {
		log.Printf("スタミナ回復エフェクト追加エラー: %v", err)
	}
	if err := processor.Execute(world); err != nil {
		log.Printf("回復エフェクト実行エラー: %v", err)
	}

	g := &TestGame{
		MainGame: game.MainGame{
			World:        world,
			StateMachine: es.Init(state, world),
		},
		gameCount:  0,
		outputPath: outputPath,
	}

	if err := ebiten.RunGame(g); err != nil && err != errRegularTermination {
		log.Fatal(err)
	}
}
