package spawner

import (
	"github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/engine/loader"
	"github.com/kijimaD/ruins/lib/raw"
	ecs "github.com/x-hgg-x/goecs/v2"

	gc "github.com/kijimaD/ruins/lib/components"
	ec "github.com/kijimaD/ruins/lib/engine/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
)

// ================
// Field
// ================

// フィールド上に表示される床を生成する
func SpawnFloor(world w.World, x gc.Row, y gc.Col) ecs.Entity {
	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := loader.EntityComponentList{}
	componentList.Game = append(componentList.Game, components.GameComponentList{
		GridElement: &gc.GridElement{Row: x, Col: y},
		SpriteRender: &ec.SpriteRender{
			SpriteSheet:  &fieldSpriteSheet,
			SpriteNumber: 2,
			Depth:        ec.DepthNumFloor,
		},
	})

	return loader.AddEntities(world, componentList)[0]
}

// フィールド上に表示される壁を生成する
func SpawnFieldWall(world w.World, x gc.Row, y gc.Col) ecs.Entity {
	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := loader.EntityComponentList{}
	componentList.Game = append(componentList.Game, components.GameComponentList{
		GridElement: &gc.GridElement{Row: x, Col: y},
		SpriteRender: &ec.SpriteRender{
			SpriteSheet:  &fieldSpriteSheet,
			SpriteNumber: 1,
			Depth:        ec.DepthNumTaller,
		},
		BlockView: &gc.BlockView{},
		BlockPass: &gc.BlockPass{},
	})

	return loader.AddEntities(world, componentList)[0]
}

// フィールド上に表示される進行ワープホールを生成する
func SpawnFieldWarpNext(world w.World, x gc.Row, y gc.Col) ecs.Entity {
	SpawnFloor(world, x, y) // 下敷き描画

	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := loader.EntityComponentList{}
	componentList.Game = append(componentList.Game, components.GameComponentList{
		GridElement: &gc.GridElement{Row: x, Col: y},
		SpriteRender: &ec.SpriteRender{
			SpriteSheet:  &fieldSpriteSheet,
			SpriteNumber: 4,
			Depth:        ec.DepthNumRug,
		},
		Warp: &gc.Warp{Mode: gc.WarpModeNext},
	})

	return loader.AddEntities(world, componentList)[0]
}

// フィールド上に表示される脱出ワープホールを生成する
func SpawnFieldWarpEscape(world w.World, x gc.Row, y gc.Col) ecs.Entity {
	SpawnFloor(world, x, y) // 下敷き描画

	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := loader.EntityComponentList{}
	componentList.Game = append(componentList.Game, components.GameComponentList{
		GridElement: &gc.GridElement{Row: x, Col: y},
		SpriteRender: &ec.SpriteRender{
			SpriteSheet:  &fieldSpriteSheet,
			SpriteNumber: 5,
			Depth:        ec.DepthNumRug,
		},
		Warp: &gc.Warp{Mode: gc.WarpModeEscape},
	})

	return loader.AddEntities(world, componentList)[0]
}

// フィールド上に表示されるプレイヤーを生成する。
// TODO: エンティティが重複しているときにエラーを返す。
// TODO: 置けるタイル以外が指定されるとエラーを返す。
// デバッグ用に任意の位置でスポーンさせたいことがあるためこの位置にある。スポーン可能なタイルかエンティティが重複してないかなどの判定はこの関数ではしていない。
func SpawnPlayer(world w.World, x gc.Pixel, y gc.Pixel) {
	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	{
		componentList := loader.EntityComponentList{}
		componentList.Game = append(componentList.Game, components.GameComponentList{
			Position: &gc.Position{X: x, Y: y},
			Player:   &gc.Player{},
			SpriteRender: &ec.SpriteRender{
				SpriteSheet:  &fieldSpriteSheet,
				SpriteNumber: 3,
				Depth:        ec.DepthNumPlayer,
			},
		})
		loader.AddEntities(world, componentList)
	}
	// カメラ
	{
		componentList := loader.EntityComponentList{}
		componentList.Game = append(componentList.Game, components.GameComponentList{
			Position: &gc.Position{X: x, Y: y},
			Camera:   &gc.Camera{Scale: 0.1, ScaleTo: 1},
		})
		loader.AddEntities(world, componentList)
	}
}

// フィールド上に表示されるNPCを生成する
// TODO: 接触すると戦闘開始するようにする
func SpawnNPC(world w.World, x gc.Pixel, y gc.Pixel) {
	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	{
		componentList := loader.EntityComponentList{}
		componentList.Game = append(componentList.Game, components.GameComponentList{
			Position: &gc.Position{X: x, Y: y},
			SpriteRender: &ec.SpriteRender{
				SpriteSheet:  &fieldSpriteSheet,
				SpriteNumber: 6,
				Depth:        ec.DepthNumTaller,
			},
			BlockPass: &gc.BlockPass{},
		})
		loader.AddEntities(world, componentList)
	}
}

// ================
// Item
// ================

// アイテムを生成する
func SpawnItem(world w.World, name string, spawnType raw.SpawnType) ecs.Entity {
	componentList := loader.EntityComponentList{}
	rawMaster := world.Resources.RawMaster.(raw.RawMaster)
	componentList.Game = append(componentList.Game, rawMaster.GenerateItem(name, spawnType))
	entities := loader.AddEntities(world, componentList)

	return entities[len(entities)-1]
}

// パーティに追加可能なキャラを生成する
func SpawnMember(world w.World, name string, inParty bool) ecs.Entity {
	componentList := loader.EntityComponentList{}
	rawMaster := world.Resources.RawMaster.(raw.RawMaster)
	componentList.Game = append(componentList.Game, rawMaster.GenerateMember(name, inParty))
	entities := loader.AddEntities(world, componentList)

	return entities[len(entities)-1]
}

// 所持素材の個数を0で初期化する
func SpawnAllMaterials(world w.World) {
	rawMaster := world.Resources.RawMaster.(raw.RawMaster)
	for k, _ := range rawMaster.MaterialIndex {
		componentList := loader.EntityComponentList{}
		componentList.Game = append(componentList.Game, rawMaster.GenerateMaterial(k, 0, raw.SpawnInBackpack))
		loader.AddEntities(world, componentList)
	}
}

// レシピ初期化
func SpawnAllRecipes(world w.World) {
	rawMaster := world.Resources.RawMaster.(raw.RawMaster)
	for k, _ := range rawMaster.RecipeIndex {
		componentList := loader.EntityComponentList{}
		componentList.Game = append(componentList.Game, rawMaster.GenerateRecipe(k))
		loader.AddEntities(world, componentList)
	}
}