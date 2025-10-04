package worldhelper

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"sort"

	"github.com/kijimaD/ruins/lib/config"
	"github.com/kijimaD/ruins/lib/engine/entities"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/turns"
	ecs "github.com/x-hgg-x/goecs/v2"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
)

// 定数定義
const (
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

// ================
// Field
// ================

// SpawnFloor は指定されたスプライトキーでフィールド上に表示される床を生成する
func SpawnFloor(world w.World, x gc.Tile, y gc.Tile, sheetName, spriteKey string) (ecs.Entity, error) {
	componentList := entities.ComponentList[gc.EntitySpec]{}
	componentList.Entities = append(componentList.Entities, gc.EntitySpec{
		GridElement: &gc.GridElement{X: x, Y: y},
		SpriteRender: &gc.SpriteRender{
			SpriteSheetName: sheetName,
			SpriteKey:       spriteKey,
			Depth:           gc.DepthNumFloor,
		},
	})

	return entities.AddEntities(world, componentList)[0], nil
}

// SpawnWall は指定されたスプライトキーで壁を生成する
func SpawnWall(world w.World, x gc.Tile, y gc.Tile, sheetName, spriteKey string) (ecs.Entity, error) {
	componentList := entities.ComponentList[gc.EntitySpec]{}
	componentList.Entities = append(componentList.Entities, gc.EntitySpec{
		GridElement: &gc.GridElement{X: x, Y: y},
		SpriteRender: &gc.SpriteRender{
			SpriteSheetName: sheetName,
			SpriteKey:       spriteKey,
			Depth:           gc.DepthNumTaller,
		},
		BlockView: &gc.BlockView{},
		BlockPass: &gc.BlockPass{},
	})

	return entities.AddEntities(world, componentList)[0], nil
}

// SpawnFieldWarpNext はフィールド上に表示される進行ワープホールを生成する
func SpawnFieldWarpNext(world w.World, x gc.Tile, y gc.Tile) (ecs.Entity, error) {
	_, err := SpawnFloor(world, x, y, "field", "floor") // 下敷き描画
	if err != nil {
		return ecs.Entity(0), fmt.Errorf("床の生成に失敗: %w", err)
	}

	componentList := entities.ComponentList[gc.EntitySpec]{}
	componentList.Entities = append(componentList.Entities, gc.EntitySpec{
		GridElement: &gc.GridElement{X: x, Y: y},
		SpriteRender: &gc.SpriteRender{
			SpriteSheetName: "field",
			SpriteKey:       "warp_next",
			Depth:           gc.DepthNumRug,
		},
		Warp: &gc.Warp{Mode: gc.WarpModeNext},
	})

	return entities.AddEntities(world, componentList)[0], nil
}

// SpawnFieldWarpEscape はフィールド上に表示される脱出ワープホールを生成する
func SpawnFieldWarpEscape(world w.World, x gc.Tile, y gc.Tile) (ecs.Entity, error) {
	_, err := SpawnFloor(world, x, y, "field", "floor") // 下敷き描画
	if err != nil {
		return ecs.Entity(0), fmt.Errorf("床の生成に失敗: %w", err)
	}

	componentList := entities.ComponentList[gc.EntitySpec]{}
	componentList.Entities = append(componentList.Entities, gc.EntitySpec{
		GridElement: &gc.GridElement{X: x, Y: y},
		SpriteRender: &gc.SpriteRender{
			SpriteSheetName: "field",
			SpriteKey:       "warp_escape",
			Depth:           gc.DepthNumRug,
		},
		Warp: &gc.Warp{Mode: gc.WarpModeEscape},
	})

	return entities.AddEntities(world, componentList)[0], nil
}

// ================
// Characters
// ================

// SpawnPlayer はプレイヤーキャラクターを生成する
func SpawnPlayer(world w.World, tileX int, tileY int, name string) (ecs.Entity, error) {
	componentList := entities.ComponentList[gc.EntitySpec]{}
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	gcl, err := rawMaster.GeneratePlayer(name)
	if err != nil {
		return ecs.Entity(0), fmt.Errorf("%w: %v", ErrMemberGeneration, err)
	}
	gcl.GridElement = &gc.GridElement{X: gc.Tile(tileX), Y: gc.Tile(tileY)}
	// カメラ
	{
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
		gcl.Camera = &gc.Camera{Scale: scale, ScaleTo: scaleTo}
	}
	componentList.Entities = append(componentList.Entities, gcl)
	entities := entities.AddEntities(world, componentList)
	fullRecover(world, entities[len(entities)-1])

	return entities[len(entities)-1], nil
}

// SpawnEnemy はフィールド上に敵キャラクターを生成する
func SpawnEnemy(world w.World, tileX int, tileY int, name string) (ecs.Entity, error) {
	componentList := entities.ComponentList[gc.EntitySpec]{}
	rawMaster := world.Resources.RawMaster.(*raw.Master)

	// raw.Masterから敵データを取得
	cl, err := rawMaster.GenerateEnemy(name)
	if err != nil {
		return ecs.Entity(0), fmt.Errorf("%w: %v", ErrEnemyGeneration, err)
	}

	// フィールド用のコンポーネントを設定
	cl.GridElement = &gc.GridElement{X: gc.Tile(tileX), Y: gc.Tile(tileY)}
	cl.BlockPass = &gc.BlockPass{}
	cl.AIMoveFSM = &gc.AIMoveFSM{}
	cl.AIRoaming = &gc.AIRoaming{
		SubState:              gc.AIRoamingWaiting,
		StartSubStateTurn:     1,                // 初期ターン
		DurationSubStateTurns: 2 + rand.IntN(3), // 2-4ターン待機
	}
	cl.AIVision = &gc.AIVision{
		ViewDistance: gc.Pixel(aiVisionDistance),
	}

	componentList.Entities = append(componentList.Entities, cl)
	entities := entities.AddEntities(world, componentList)

	// 全回復
	npcEntity := entities[len(entities)-1]
	fullRecover(world, npcEntity)

	// ActionPointsを初期化
	if world.Resources.TurnManager != nil {
		if turnManager, ok := world.Resources.TurnManager.(*turns.TurnManager); ok {
			if npcEntity.HasComponent(world.Components.TurnBased) {
				actionPoints := world.Components.TurnBased.Get(npcEntity).(*gc.TurnBased)
				maxAP := turnManager.CalculateMaxActionPoints(world, npcEntity)
				actionPoints.AP.Current = maxAP
				actionPoints.AP.Max = maxAP
			}
		}
	}

	return npcEntity, nil
}

// ================
// Items
// ================

// SpawnItem はアイテムを生成する
func SpawnItem(world w.World, name string, locationType gc.ItemLocationType) (ecs.Entity, error) {
	componentList := entities.ComponentList[gc.EntitySpec]{}
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	gameComponent, err := rawMaster.GenerateItem(name, locationType)
	if err != nil {
		return ecs.Entity(0), fmt.Errorf("%w: %v", ErrItemGeneration, err)
	}
	componentList.Entities = append(componentList.Entities, gameComponent)
	entities := entities.AddEntities(world, componentList)

	return entities[len(entities)-1], nil
}

// SpawnMaterial はmaterialを生成する
func SpawnMaterial(world w.World, name string, amount int, locationType gc.ItemLocationType) (ecs.Entity, error) {
	componentList := entities.ComponentList[gc.EntitySpec]{}
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	gameComponent, err := rawMaster.GenerateMaterial(name, amount, locationType)
	if err != nil {
		return ecs.Entity(0), fmt.Errorf("%w: %v", ErrItemGeneration, err)
	}
	componentList.Entities = append(componentList.Entities, gameComponent)
	entities := entities.AddEntities(world, componentList)

	return entities[len(entities)-1], nil
}

// 完全回復させる
func fullRecover(world w.World, entity ecs.Entity) {
	// 新しく生成されたエンティティの最大HP/SPを設定
	_ = setMaxHPSP(world, entity) // エラーが発生した場合もリカバリーは続行する

	// Poolsコンポーネントを取得
	poolsComponent := world.Components.Pools.Get(entity)
	if poolsComponent == nil {
		return // Poolsがない場合は何もしない
	}

	pools := poolsComponent.(*gc.Pools)

	// HP全回復
	pools.HP.Current = pools.HP.Max

	// SP全回復
	pools.SP.Current = pools.SP.Max

	// ActionPointsコンポーネントがある場合は最大APに設定
	if entity.HasComponent(world.Components.TurnBased) {
		if world.Resources.TurnManager != nil {
			if turnManager, ok := world.Resources.TurnManager.(*turns.TurnManager); ok {
				actionPoints := world.Components.TurnBased.Get(entity).(*gc.TurnBased)
				maxAP := turnManager.CalculateMaxActionPoints(world, entity)
				actionPoints.AP.Current = maxAP
				actionPoints.AP.Max = maxAP
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

	// 最大HP計算: base+(体力*multiplyV+力+感覚)
	pools.HP.Max = int(hpBaseValue) + attrs.Vitality.Total*hpVitalityMultiply + attrs.Strength.Total + attrs.Sensation.Total
	pools.HP.Current = pools.HP.Max

	// 最大SP計算: (体力*multiplyV+器用さ+素早さ)
	pools.SP.Max = attrs.Vitality.Total*spVitalityMultiply + attrs.Dexterity.Total + attrs.Agility.Total
	pools.SP.Current = pools.SP.Max

	return nil
}

// SpawnAllMaterials は所持素材の個数を0で初期化する
func SpawnAllMaterials(world w.World) error {
	rawMaster := world.Resources.RawMaster.(*raw.Master)

	// マップのキーをソートして決定的な順序にする
	keys := make([]string, 0, len(rawMaster.MaterialIndex))
	for k := range rawMaster.MaterialIndex {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// ソート済みの順序でマテリアルを生成
	for _, k := range keys {
		componentList := entities.ComponentList[gc.EntitySpec]{}
		gameComponent, err := rawMaster.GenerateMaterial(k, 0, gc.ItemLocationInBackpack)
		if err != nil {
			return fmt.Errorf("%w (material: %s): %v", ErrItemGeneration, k, err)
		}
		componentList.Entities = append(componentList.Entities, gameComponent)
		entities.AddEntities(world, componentList)
	}
	return nil
}

// SpawnAllRecipes はレシピ初期化
func SpawnAllRecipes(world w.World) error {
	rawMaster := world.Resources.RawMaster.(*raw.Master)

	// マップのキーをソートして決定的な順序にする
	keys := make([]string, 0, len(rawMaster.RecipeIndex))
	for k := range rawMaster.RecipeIndex {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// ソート済みの順序でレシピを生成
	for _, k := range keys {
		componentList := entities.ComponentList[gc.EntitySpec]{}
		gameComponent, err := rawMaster.GenerateRecipe(k)
		if err != nil {
			return fmt.Errorf("%w (recipe: %s): %v", ErrItemGeneration, k, err)
		}
		componentList.Entities = append(componentList.Entities, gameComponent)
		entities.AddEntities(world, componentList)
	}
	return nil
}

// SpawnAllCards は敵が使う用。マスタとなるカードを初期化する
func SpawnAllCards(world w.World) error {
	rawMaster := world.Resources.RawMaster.(*raw.Master)

	// マップのキーをソートして決定的な順序にする
	keys := make([]string, 0, len(rawMaster.ItemIndex))
	for k := range rawMaster.ItemIndex {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// ソート済みの順序でカードを生成
	for _, k := range keys {
		componentList := entities.ComponentList[gc.EntitySpec]{}
		gameComponent, err := rawMaster.GenerateItem(k, gc.ItemLocationNone)
		if err != nil {
			return fmt.Errorf("%w (card: %s): %v", ErrItemGeneration, k, err)
		}
		componentList.Entities = append(componentList.Entities, gameComponent)
		entities.AddEntities(world, componentList)
	}
	return nil
}

// SpawnFieldItem はフィールド上にアイテムを生成する
func SpawnFieldItem(world w.World, itemName string, x gc.Tile, y gc.Tile) (ecs.Entity, error) {
	_, err := SpawnFloor(world, x, y, "field", "floor") // 下敷きの床を描画
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

	return item, nil
}

// MovePlayerToPosition は既存のプレイヤーエンティティを指定位置に移動させる
func MovePlayerToPosition(world w.World, tileX int, tileY int) error {
	// 既存のプレイヤーエンティティを検索
	var playerEntity ecs.Entity
	var found bool

	world.Manager.Join(
		world.Components.Player,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		if !found {
			playerEntity = entity
			found = true
		}
	}))

	if !found {
		return errors.New("プレイヤーエンティティが見つかりません")
	}

	// プレイヤーの位置を更新
	gridElement := world.Components.GridElement.Get(playerEntity).(*gc.GridElement)
	gridElement.X = gc.Tile(tileX)
	gridElement.Y = gc.Tile(tileY)

	return nil
}
