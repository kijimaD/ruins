package mapspawner

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	mapplanner "github.com/kijimaD/ruins/lib/mapplanner"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// Spawn はMetaPlanからレベルを生成する
// タイル、NPC、Props、ワープポータル情報から効率的にエンティティを生成する
func Spawn(world w.World, metaPlan *mapplanner.MetaPlan) (resources.Level, error) {
	width := int(metaPlan.Level.TileWidth)
	height := int(metaPlan.Level.TileHeight)
	totalTiles := width * height

	level := resources.Level{
		TileWidth:  metaPlan.Level.TileWidth,
		TileHeight: metaPlan.Level.TileHeight,
		Entities:   make([]ecs.Entity, totalTiles),
	}

	// タイルからエンティティを直接生成
	for _i, tile := range metaPlan.Tiles {
		i := resources.TileIdx(_i)
		x, y := metaPlan.Level.XYTileCoord(i)
		tileX, tileY := gc.Tile(x), gc.Tile(y)

		var entity ecs.Entity
		var err error

		if tile.Walkable {
			// 床エンティティを生成
			entity, err = worldhelper.SpawnFloor(world, tileX, tileY)
		} else {
			// 隣接に床がある場合のみ壁エンティティを生成
			if metaPlan.AdjacentAnyFloor(i) {
				// 壁タイプを判定してスプライト番号を決定
				wallType := metaPlan.GetWallType(i)
				spriteNumber := getSpriteNumberForWallType(wallType)
				entity, err = worldhelper.SpawnWall(world, tileX, tileY, spriteNumber)
			}
		}

		if err != nil {
			return resources.Level{}, fmt.Errorf("タイルエンティティ生成エラー (%d, %d): %w", int(x), int(y), err)
		}

		if entity != ecs.Entity(0) {
			level.Entities[i] = entity
		}
	}

	// ワープポータルを生成
	for _, portal := range metaPlan.WarpPortals {
		tileX, tileY := gc.Tile(portal.X), gc.Tile(portal.Y)
		tileIdx := level.XYTileIndex(tileX, tileY)

		var entity ecs.Entity
		var err error
		switch portal.Type {
		case mapplanner.WarpPortalNext:
			entity, err = worldhelper.SpawnFieldWarpNext(world, tileX, tileY)
		case mapplanner.WarpPortalEscape:
			entity, err = worldhelper.SpawnFieldWarpEscape(world, tileX, tileY)
		}

		if err != nil {
			return resources.Level{}, fmt.Errorf("ワープポータル生成エラー (%d, %d): %w", portal.X, portal.Y, err)
		}

		level.Entities[tileIdx] = entity
	}

	// NPCを生成
	for _, npc := range metaPlan.NPCs {
		entity, err := worldhelper.SpawnEnemy(world, npc.X, npc.Y, npc.NPCType)
		if err != nil {
			return resources.Level{}, fmt.Errorf("NPC生成エラー (%d, %d): %w", npc.X, npc.Y, err)
		}

		tileIdx := level.XYTileIndex(gc.Tile(npc.X), gc.Tile(npc.Y))
		level.Entities[tileIdx] = entity
	}

	// アイテムを生成
	for _, item := range metaPlan.Items {
		tileX, tileY := gc.Tile(item.X), gc.Tile(item.Y)
		entity, err := worldhelper.SpawnFieldItem(world, item.ItemName, tileX, tileY)
		if err != nil {
			return resources.Level{}, fmt.Errorf("アイテム生成エラー (%d, %d): %w", item.X, item.Y, err)
		}

		tileIdx := level.XYTileIndex(tileX, tileY)
		level.Entities[tileIdx] = entity
	}

	// Propsを生成
	for _, prop := range metaPlan.Props {
		tileX, tileY := gc.Tile(prop.X), gc.Tile(prop.Y)
		entity, err := worldhelper.SpawnProp(world, prop.PropType, tileX, tileY)
		if err != nil {
			return resources.Level{}, fmt.Errorf("props生成エラー (%d, %d): %w", prop.X, prop.Y, err)
		}

		tileIdx := level.XYTileIndex(tileX, tileY)
		level.Entities[tileIdx] = entity
	}

	return level, nil
}
