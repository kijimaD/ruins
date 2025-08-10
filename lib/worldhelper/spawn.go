package worldhelper

import (
	"fmt"

	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/engine/entities"
	"github.com/kijimaD/ruins/lib/raw"
	ecs "github.com/x-hgg-x/goecs/v2"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
)

// ========== 生成システム（旧spawner） ==========

// ================
// Field
// ================

// SpawnFloor はフィールド上に表示される床を生成する
func SpawnFloor(world w.World, x gc.Row, y gc.Col) ecs.Entity {
	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := entities.ComponentList{}
	componentList.Game = append(componentList.Game, gc.GameComponentList{
		GridElement: &gc.GridElement{Row: x, Col: y},
		SpriteRender: &gc.SpriteRender{
			SpriteSheet:  &fieldSpriteSheet,
			SpriteNumber: 2,
			Depth:        gc.DepthNumFloor,
		},
	})

	return entities.AddEntities(world, componentList)[0]
}

// SpawnFieldWall はフィールド上に表示される壁を生成する
func SpawnFieldWall(world w.World, x gc.Row, y gc.Col) ecs.Entity {
	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := entities.ComponentList{}
	componentList.Game = append(componentList.Game, gc.GameComponentList{
		GridElement: &gc.GridElement{Row: x, Col: y},
		SpriteRender: &gc.SpriteRender{
			SpriteSheet:  &fieldSpriteSheet,
			SpriteNumber: 1,
			Depth:        gc.DepthNumTaller,
		},
		BlockView: &gc.BlockView{},
		BlockPass: &gc.BlockPass{},
	})

	return entities.AddEntities(world, componentList)[0]
}

// SpawnFieldWarpNext はフィールド上に表示される進行ワープホールを生成する
func SpawnFieldWarpNext(world w.World, x gc.Row, y gc.Col) ecs.Entity {
	SpawnFloor(world, x, y) // 下敷き描画

	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := entities.ComponentList{}
	componentList.Game = append(componentList.Game, gc.GameComponentList{
		GridElement: &gc.GridElement{Row: x, Col: y},
		SpriteRender: &gc.SpriteRender{
			SpriteSheet:  &fieldSpriteSheet,
			SpriteNumber: 4,
			Depth:        gc.DepthNumRug,
		},
		Warp: &gc.Warp{Mode: gc.WarpModeNext},
	})

	return entities.AddEntities(world, componentList)[0]
}

// SpawnFieldWarpEscape はフィールド上に表示される脱出ワープホールを生成する
func SpawnFieldWarpEscape(world w.World, x gc.Row, y gc.Col) ecs.Entity {
	SpawnFloor(world, x, y) // 下敷き描画

	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := entities.ComponentList{}
	componentList.Game = append(componentList.Game, gc.GameComponentList{
		GridElement: &gc.GridElement{Row: x, Col: y},
		SpriteRender: &gc.SpriteRender{
			SpriteSheet:  &fieldSpriteSheet,
			SpriteNumber: 5,
			Depth:        gc.DepthNumRug,
		},
		Warp: &gc.Warp{Mode: gc.WarpModeEscape},
	})

	return entities.AddEntities(world, componentList)[0]
}

// SpawnOperator はフィールド上に表示される操作対象キャラを生成する
// TODO: エンティティが重複しているときにエラーを返す。
// TODO: 置けるタイル以外が指定されるとエラーを返す。
// デバッグ用に任意の位置でスポーンさせたいことがあるためこの位置にある。スポーン可能なタイルかエンティティが重複してないかなどの判定はこの関数ではしていない。
func SpawnOperator(world w.World, x gc.Pixel, y gc.Pixel) {
	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	{
		componentList := entities.ComponentList{}
		componentList.Game = append(componentList.Game, gc.GameComponentList{
			Position: &gc.Position{X: x, Y: y},
			Velocity: &gc.Velocity{},
			Operator: &gc.Operator{},
			SpriteRender: &gc.SpriteRender{
				SpriteSheet:  &fieldSpriteSheet,
				SpriteNumber: 3,
				Depth:        gc.DepthNumOperator,
			},
			BlockPass: &gc.BlockPass{},
		})
		entities.AddEntities(world, componentList)
	}
	// カメラ
	{
		componentList := entities.ComponentList{}
		componentList.Game = append(componentList.Game, gc.GameComponentList{
			Position: &gc.Position{X: x, Y: y},
			Camera:   &gc.Camera{Scale: 0.1, ScaleTo: 1},
		})
		entities.AddEntities(world, componentList)
	}
}

// SpawnNPC はフィールド上に表示されるNPCを生成する
// 接触すると戦闘開始する敵として動作する
func SpawnNPC(world w.World, x gc.Pixel, y gc.Pixel) {
	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	{
		componentList := entities.ComponentList{}
		componentList.Game = append(componentList.Game, gc.GameComponentList{
			Position: &gc.Position{X: x, Y: y},
			Velocity: &gc.Velocity{},
			SpriteRender: &gc.SpriteRender{
				SpriteSheet:  &fieldSpriteSheet,
				SpriteNumber: 6,
				Depth:        gc.DepthNumTaller,
			},
			BlockPass: &gc.BlockPass{},
			AIMoveFSM: &gc.AIMoveFSM{},
			AIRoaming: &gc.AIRoaming{},
		})
		entities.AddEntities(world, componentList)
	}
}

// ================
// Item
// ================

// SpawnItem はアイテムを生成する
func SpawnItem(world w.World, name string, locationType gc.ItemLocationType) ecs.Entity {
	componentList := entities.ComponentList{}
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	gameComponent, err := rawMaster.GenerateItem(name, locationType)
	if err != nil {
		panic(err) // TODO: Handle error properly
	}
	componentList.Game = append(componentList.Game, gameComponent)
	entities := entities.AddEntities(world, componentList)

	return entities[len(entities)-1]
}

// SpawnMember はパーティに追加可能なキャラを生成する
func SpawnMember(world w.World, name string, inParty bool) ecs.Entity {
	componentList := entities.ComponentList{}
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	memberComp, err := rawMaster.GenerateMember(name, inParty)
	if err != nil {
		panic(fmt.Sprintf("GenerateMember failed: %v", err))
	}
	componentList.Game = append(componentList.Game, memberComp)
	entities := entities.AddEntities(world, componentList)
	fullRecover(world, entities[len(entities)-1])

	return entities[len(entities)-1]
}

// SpawnEnemy は戦闘に参加する敵キャラを生成する
func SpawnEnemy(world w.World, name string) ecs.Entity {
	componentList := entities.ComponentList{}
	rawMaster := world.Resources.RawMaster.(*raw.Master)

	cl, err := rawMaster.GenerateEnemy(name)
	if err != nil {
		panic(fmt.Sprintf("GenerateEnemy failed: %v", err))
	}
	componentList.Game = append(
		componentList.Game,
		cl,
	)
	entities := entities.AddEntities(world, componentList)
	fullRecover(world, entities[len(entities)-1])

	return entities[len(entities)-1]
}

// 完全回復させる
func fullRecover(world w.World, entity ecs.Entity) {
	// 新しく生成されたエンティティの最大HP/SPを設定
	setMaxHPSP(world, entity)

	processor := effects.NewProcessor()

	// HP全回復
	hpEffect := effects.FullRecoveryHP{}
	processor.AddEffect(hpEffect, nil, entity)

	// SP全回復
	spEffect := effects.FullRecoverySP{}
	processor.AddEffect(spEffect, nil, entity)

	// エフェクト実行
	_ = processor.Execute(world) // エラーが発生した場合もリカバリーは続行する（ログ出力はProcessor内で行われる）
}

// 指定したエンティティの最大HP/SPを設定する
func setMaxHPSP(world w.World, entity ecs.Entity) {

	if !entity.HasComponent(world.Components.Pools) || !entity.HasComponent(world.Components.Attributes) {
		return
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

	// 最大HP計算: 30+(体力*8+力+感覚)*{1+(Lv-1)*0.03}
	pools.HP.Max = int(30 + float64(attrs.Vitality.Total*8+attrs.Strength.Total+attrs.Sensation.Total)*(1+float64(pools.Level-1)*0.03))
	pools.HP.Current = pools.HP.Max

	// 最大SP計算: (体力*2+器用さ+素早さ)*{1+(Lv-1)*0.02}
	pools.SP.Max = int(float64(attrs.Vitality.Total*2+attrs.Dexterity.Total+attrs.Agility.Total) * (1 + float64(pools.Level-1)*0.02))
	pools.SP.Current = pools.SP.Max
}

// SpawnAllMaterials は所持素材の個数を0で初期化する
func SpawnAllMaterials(world w.World) {
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	for k := range rawMaster.MaterialIndex {
		componentList := entities.ComponentList{}
		gameComponent, err := rawMaster.GenerateMaterial(k, 0, gc.ItemLocationInBackpack)
		if err != nil {
			panic(err) // TODO: Handle error properly
		}
		componentList.Game = append(componentList.Game, gameComponent)
		entities.AddEntities(world, componentList)
	}
}

// SpawnAllRecipes はレシピ初期化
func SpawnAllRecipes(world w.World) {
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	for k := range rawMaster.RecipeIndex {
		componentList := entities.ComponentList{}
		gameComponent, err := rawMaster.GenerateRecipe(k)
		if err != nil {
			panic(err) // TODO: Handle error properly
		}
		componentList.Game = append(componentList.Game, gameComponent)
		entities.AddEntities(world, componentList)
	}
}

// SpawnAllCards は敵が使う用。マスタとなるカードを初期化する
func SpawnAllCards(world w.World) {
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	for k := range rawMaster.ItemIndex {
		componentList := entities.ComponentList{}
		gameComponent, err := rawMaster.GenerateItem(k, gc.ItemLocationNone)
		if err != nil {
			panic(err) // TODO: Handle error properly
		}
		componentList.Game = append(componentList.Game, gameComponent)
		entities.AddEntities(world, componentList)
	}
}
