package worldhelper

import (
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/kijimaD/ruins/lib/config"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/engine/entities"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/turns"
	ecs "github.com/x-hgg-x/goecs/v2"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
)

// 定数定義
const (
	// スプライト番号
	spriteNumberFloor      = 2  // 床
	spriteNumberPlayer     = 3  // プレイヤー
	spriteNumberWarpNext   = 4  // 進行ワープホール
	spriteNumberWarpEscape = 5  // 脱出ワープホール
	spriteNumberNPC        = 6  // NPC
	spriteNumberFieldItem  = 18 // フィールドアイテム

	// カメラスケール
	cameraInitialScale = 0.1 // カメラの初期スケール（ズームアウト）
	cameraNormalScale  = 1.0 // カメラの通常スケール

	// AI設定
	aiVisionDistance = 160.0 // AIの視界距離（ピクセル）

	// ステータス計算係数
	hpBaseValue        = 30   // HP計算の基本値
	hpVitalityMultiply = 8    // HP計算の体力係数
	spVitalityMultiply = 2    // SP計算の体力係数
	hpLevelGrowthRate  = 0.03 // HPのレベル成長率
	spLevelGrowthRate  = 0.02 // SPのレベル成長率
)

// エラー定義
var (
	ErrItemGeneration   = errors.New("アイテムの生成に失敗しました")
	ErrMemberGeneration = errors.New("メンバーの生成に失敗しました")
	ErrEnemyGeneration  = errors.New("敵の生成に失敗しました")
	ErrEffectGeneration = errors.New("エフェクトの生成に失敗しました")
)

// ========== 生成システム（旧spawner） ==========

// ================
// Field
// ================

// SpawnFloor はフィールド上に表示される床を生成する
func SpawnFloor(world w.World, x gc.Tile, y gc.Tile) (ecs.Entity, error) {
	componentList := entities.ComponentList{}
	componentList.Game = append(componentList.Game, gc.GameComponentList{
		GridElement: &gc.GridElement{X: x, Y: y},
		SpriteRender: &gc.SpriteRender{
			Name:         "field",
			SpriteNumber: spriteNumberFloor,
			Depth:        gc.DepthNumFloor,
		},
	})

	return entities.AddEntities(world, componentList)[0], nil
}

// SpawnFieldWallWithSprite は指定されたスプライト番号でフィールド上に表示される壁を生成する
func SpawnFieldWallWithSprite(world w.World, x gc.Tile, y gc.Tile, spriteNumber int) (ecs.Entity, error) {
	componentList := entities.ComponentList{}
	componentList.Game = append(componentList.Game, gc.GameComponentList{
		GridElement: &gc.GridElement{X: x, Y: y},
		SpriteRender: &gc.SpriteRender{
			Name:         "field",
			SpriteNumber: spriteNumber,
			Depth:        gc.DepthNumTaller,
		},
		BlockView: &gc.BlockView{},
		BlockPass: &gc.BlockPass{},
	})

	return entities.AddEntities(world, componentList)[0], nil
}

// SpawnFieldWarpNext はフィールド上に表示される進行ワープホールを生成する
func SpawnFieldWarpNext(world w.World, x gc.Tile, y gc.Tile) (ecs.Entity, error) {
	_, err := SpawnFloor(world, x, y) // 下敷き描画
	if err != nil {
		return ecs.Entity(0), fmt.Errorf("床の生成に失敗: %w", err)
	}

	componentList := entities.ComponentList{}
	componentList.Game = append(componentList.Game, gc.GameComponentList{
		GridElement: &gc.GridElement{X: x, Y: y},
		SpriteRender: &gc.SpriteRender{
			Name:         "field",
			SpriteNumber: spriteNumberWarpNext,
			Depth:        gc.DepthNumRug,
		},
		Warp: &gc.Warp{Mode: gc.WarpModeNext},
	})

	return entities.AddEntities(world, componentList)[0], nil
}

// SpawnFieldWarpEscape はフィールド上に表示される脱出ワープホールを生成する
func SpawnFieldWarpEscape(world w.World, x gc.Tile, y gc.Tile) (ecs.Entity, error) {
	_, err := SpawnFloor(world, x, y) // 下敷き描画
	if err != nil {
		return ecs.Entity(0), fmt.Errorf("床の生成に失敗: %w", err)
	}

	componentList := entities.ComponentList{}
	componentList.Game = append(componentList.Game, gc.GameComponentList{
		GridElement: &gc.GridElement{X: x, Y: y},
		SpriteRender: &gc.SpriteRender{
			Name:         "field",
			SpriteNumber: spriteNumberWarpEscape,
			Depth:        gc.DepthNumRug,
		},
		Warp: &gc.Warp{Mode: gc.WarpModeEscape},
	})

	return entities.AddEntities(world, componentList)[0], nil
}

// SpawnOperator はフィールド上に表示される操作対象キャラを生成する
// TODO: エンティティが重複しているときにエラーを返す。
// TODO: 置けるタイル以外が指定されるとエラーを返す。
// デバッグ用に任意の位置でスポーンさせたいことがあるためこの位置にある。スポーン可能なタイルかエンティティが重複してないかなどの判定はこの関数ではしていない。
func SpawnOperator(world w.World, tileX int, tileY int) error {
	{
		componentList := entities.ComponentList{}
		componentList.Game = append(componentList.Game, gc.GameComponentList{
			GridElement: &gc.GridElement{X: gc.Tile(tileX), Y: gc.Tile(tileY)},
			Player:      &gc.Player{},
			Operator:    &gc.Operator{},
			SpriteRender: &gc.SpriteRender{
				Name:         "field",
				SpriteNumber: spriteNumberPlayer,
				Depth:        gc.DepthNumOperator,
			},
			BlockPass: &gc.BlockPass{},
		})
		entities.AddEntities(world, componentList)
	}
	// カメラ
	{
		componentList := entities.ComponentList{}
		// config設定を確認
		cfg := config.Get()
		var scale, scaleTo float64
		if cfg.DisableAnimation {
			// アニメーション無効時は初期スケールを通常値に設定
			scale = cameraNormalScale
			scaleTo = cameraNormalScale
		} else {
			// アニメーション有効時はズームアウトアニメーション
			scale = cameraInitialScale
			scaleTo = cameraNormalScale
		}

		componentList.Game = append(componentList.Game, gc.GameComponentList{
			GridElement: &gc.GridElement{X: gc.Tile(tileX), Y: gc.Tile(tileY)},
			Camera:      &gc.Camera{Scale: scale, ScaleTo: scaleTo},
		})
		entities.AddEntities(world, componentList)
	}
	return nil
}

// SpawnNPC はフィールド上に表示されるNPCを生成する
// 接触すると戦闘開始する敵として動作する
func SpawnNPC(world w.World, tileX gc.Tile, tileY gc.Tile) error {
	{
		componentList := entities.ComponentList{}
		componentList.Game = append(componentList.Game, gc.GameComponentList{
			GridElement: &gc.GridElement{X: tileX, Y: tileY},
			SpriteRender: &gc.SpriteRender{
				Name:         "field",
				SpriteNumber: spriteNumberNPC,
				Depth:        gc.DepthNumTaller,
			},
			BlockPass: &gc.BlockPass{},
			AIMoveFSM: &gc.AIMoveFSM{},
			AIRoaming: &gc.AIRoaming{
				SubState:              gc.AIRoamingWaiting,
				StartSubStateTurn:     1,                // 初期ターン
				DurationSubStateTurns: 2 + rand.IntN(3), // 2-4ターン待機
			},
			AIVision: &gc.AIVision{
				ViewDistance: gc.Pixel(aiVisionDistance),
			},
			Attributes: &gc.Attributes{
				Vitality:  gc.Attribute{Base: 10, Modifier: 0, Total: 10},
				Strength:  gc.Attribute{Base: 10, Modifier: 0, Total: 10},
				Sensation: gc.Attribute{Base: 10, Modifier: 0, Total: 10},
				Dexterity: gc.Attribute{Base: 10, Modifier: 0, Total: 10},
				Agility:   gc.Attribute{Base: 10, Modifier: 0, Total: 10},
			},
			ActionPoints: &gc.ActionPoints{
				Current: 100, // あとで再計算される
			},
		})
		npcEntities := entities.AddEntities(world, componentList)

		// NPCのActionPointsを適切に初期化
		if len(npcEntities) > 0 {
			npcEntity := npcEntities[0]
			if world.Resources.TurnManager != nil {
				if turnManager, ok := world.Resources.TurnManager.(*turns.TurnManager); ok {
					if npcEntity.HasComponent(world.Components.ActionPoints) {
						actionPoints := world.Components.ActionPoints.Get(npcEntity).(*gc.ActionPoints)
						maxAP := turnManager.CalculateMaxActionPoints(world, npcEntity)
						actionPoints.Current = maxAP
					}
				}
			}
		}
	}
	return nil
}

// ================
// Item
// ================

// SpawnItem はアイテムを生成する
func SpawnItem(world w.World, name string, locationType gc.ItemLocationType) (ecs.Entity, error) {
	componentList := entities.ComponentList{}
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	gameComponent, err := rawMaster.GenerateItem(name, locationType)
	if err != nil {
		return ecs.Entity(0), fmt.Errorf("%w: %v", ErrItemGeneration, err)
	}
	componentList.Game = append(componentList.Game, gameComponent)
	entities := entities.AddEntities(world, componentList)

	return entities[len(entities)-1], nil
}

// SpawnPlayer はプレイヤーキャラクターを生成する
func SpawnPlayer(world w.World, name string) (ecs.Entity, error) {
	componentList := entities.ComponentList{}
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	memberComp, err := rawMaster.GeneratePlayer(name)
	if err != nil {
		return ecs.Entity(0), fmt.Errorf("%w: %v", ErrMemberGeneration, err)
	}
	componentList.Game = append(componentList.Game, memberComp)
	entities := entities.AddEntities(world, componentList)
	fullRecover(world, entities[len(entities)-1])

	return entities[len(entities)-1], nil
}

// SpawnEnemy は戦闘に参加する敵キャラを生成する
func SpawnEnemy(world w.World, name string) (ecs.Entity, error) {
	componentList := entities.ComponentList{}
	rawMaster := world.Resources.RawMaster.(*raw.Master)

	cl, err := rawMaster.GenerateEnemy(name)
	if err != nil {
		return ecs.Entity(0), fmt.Errorf("%w: %v", ErrEnemyGeneration, err)
	}
	componentList.Game = append(
		componentList.Game,
		cl,
	)
	entities := entities.AddEntities(world, componentList)
	fullRecover(world, entities[len(entities)-1])

	return entities[len(entities)-1], nil
}

// SpawnMaterial はmaterialを生成する
func SpawnMaterial(world w.World, name string, amount int, locationType gc.ItemLocationType) (ecs.Entity, error) {
	componentList := entities.ComponentList{}
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	gameComponent, err := rawMaster.GenerateMaterial(name, amount, locationType)
	if err != nil {
		return ecs.Entity(0), fmt.Errorf("%w: %v", ErrItemGeneration, err)
	}
	componentList.Game = append(componentList.Game, gameComponent)
	entities := entities.AddEntities(world, componentList)

	return entities[len(entities)-1], nil
}

// 完全回復させる
func fullRecover(world w.World, entity ecs.Entity) {
	// 新しく生成されたエンティティの最大HP/SPを設定
	_ = setMaxHPSP(world, entity) // エラーが発生した場合もリカバリーは続行する

	processor := effects.NewProcessor()

	// HP全回復
	hpEffect := effects.FullRecoveryHP{}
	processor.AddEffect(hpEffect, nil, entity)

	// SP全回復
	spEffect := effects.FullRecoverySP{}
	processor.AddEffect(spEffect, nil, entity)

	// エフェクト実行
	_ = processor.Execute(world) // エラーが発生した場合もリカバリーは続行する（ログ出力はProcessor内で行われる）

	// ActionPointsコンポーネントがある場合は最大APに設定
	if entity.HasComponent(world.Components.ActionPoints) {
		if world.Resources.TurnManager != nil {
			if turnManager, ok := world.Resources.TurnManager.(*turns.TurnManager); ok {
				actionPoints := world.Components.ActionPoints.Get(entity).(*gc.ActionPoints)
				maxAP := turnManager.CalculateMaxActionPoints(world, entity)
				actionPoints.Current = maxAP
			}
		}
	}
}

// 指定したエンティティの最大HP/SPを設定する
func setMaxHPSP(world w.World, entity ecs.Entity) error {

	if !entity.HasComponent(world.Components.Pools) || !entity.HasComponent(world.Components.Attributes) {
		return fmt.Errorf("entity %v does not have required components (Pools or Attributes)", entity)
	}

	pools := world.Components.Pools.Get(entity).(*gc.Pools)
	attrs := world.Components.Attributes.Get(entity).(*gc.Attributes)

	// Totalが設定されていない場合はBaseから初期化
	if attrs.Vitality.Total == 0 {
		attrs.Vitality.Total = attrs.Vitality.Base
	}
	if attrs.Strength.Total == 0 {
		attrs.Strength.Total = attrs.Strength.Base
	}
	if attrs.Sensation.Total == 0 {
		attrs.Sensation.Total = attrs.Sensation.Base
	}
	if attrs.Dexterity.Total == 0 {
		attrs.Dexterity.Total = attrs.Dexterity.Base
	}
	if attrs.Agility.Total == 0 {
		attrs.Agility.Total = attrs.Agility.Base
	}
	if attrs.Defense.Total == 0 {
		attrs.Defense.Total = attrs.Defense.Base
	}

	// 最大HP計算: base+(体力*multiplyV+力+感覚)*{1+(Lv-1)*growthRate}
	pools.HP.Max = int(hpBaseValue + float64(attrs.Vitality.Total*hpVitalityMultiply+attrs.Strength.Total+attrs.Sensation.Total)*(1+float64(pools.Level-1)*hpLevelGrowthRate))
	pools.HP.Current = pools.HP.Max

	// 最大SP計算: (体力*multiplyV+器用さ+素早さ)*{1+(Lv-1)*growthRate}
	pools.SP.Max = int(float64(attrs.Vitality.Total*spVitalityMultiply+attrs.Dexterity.Total+attrs.Agility.Total) * (1 + float64(pools.Level-1)*spLevelGrowthRate))
	pools.SP.Current = pools.SP.Max

	return nil
}

// SpawnAllMaterials は所持素材の個数を0で初期化する
func SpawnAllMaterials(world w.World) error {
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	for k := range rawMaster.MaterialIndex {
		componentList := entities.ComponentList{}
		gameComponent, err := rawMaster.GenerateMaterial(k, 0, gc.ItemLocationInBackpack)
		if err != nil {
			return fmt.Errorf("%w (material: %s): %v", ErrItemGeneration, k, err)
		}
		componentList.Game = append(componentList.Game, gameComponent)
		entities.AddEntities(world, componentList)
	}
	return nil
}

// SpawnAllRecipes はレシピ初期化
func SpawnAllRecipes(world w.World) error {
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	for k := range rawMaster.RecipeIndex {
		componentList := entities.ComponentList{}
		gameComponent, err := rawMaster.GenerateRecipe(k)
		if err != nil {
			return fmt.Errorf("%w (recipe: %s): %v", ErrItemGeneration, k, err)
		}
		componentList.Game = append(componentList.Game, gameComponent)
		entities.AddEntities(world, componentList)
	}
	return nil
}

// SpawnAllCards は敵が使う用。マスタとなるカードを初期化する
func SpawnAllCards(world w.World) error {
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	for k := range rawMaster.ItemIndex {
		componentList := entities.ComponentList{}
		gameComponent, err := rawMaster.GenerateItem(k, gc.ItemLocationNone)
		if err != nil {
			return fmt.Errorf("%w (card: %s): %v", ErrItemGeneration, k, err)
		}
		componentList.Game = append(componentList.Game, gameComponent)
		entities.AddEntities(world, componentList)
	}
	return nil
}

// SpawnFieldItem はフィールド上にアイテムを生成する
func SpawnFieldItem(world w.World, itemName string, x gc.Tile, y gc.Tile) (ecs.Entity, error) {
	_, err := SpawnFloor(world, x, y) // 下敷きの床を描画
	if err != nil {
		return ecs.Entity(0), fmt.Errorf("床の生成に失敗: %w", err)
	}

	// アイテムエンティティを生成
	item, err := SpawnItem(world, itemName, gc.ItemLocationOnField)
	if err != nil {
		return ecs.Entity(0), err
	}

	// フィールド表示用のコンポーネントを追加
	item.AddComponent(world.Components.GridElement, &gc.GridElement{X: x, Y: y})
	item.AddComponent(world.Components.SpriteRender, &gc.SpriteRender{
		Name:         "field", // フィールドスプライトシートを使用
		SpriteNumber: spriteNumberFieldItem,
		Depth:        gc.DepthNumRug,
	})

	return item, nil
}
