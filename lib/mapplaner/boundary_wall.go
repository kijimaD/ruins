package mapplaner

import "github.com/kijimaD/ruins/lib/resources"

// BoundaryWall は外枠をすべて指定タイルで覆うビルダー
// マップの最外枠のタイルを無条件で指定タイルタイプに変更する
type BoundaryWall struct {
	// WallTileType 外枠に使用するタイルタイプ
	WallTileType Tile
}

// NewBoundaryWall は新しいBoundaryWallビルダーを作成する
func NewBoundaryWall(wallTileType Tile) BoundaryWall {
	return BoundaryWall{
		WallTileType: wallTileType,
	}
}

// BuildMeta はメタデータをビルドする
func (b BoundaryWall) BuildMeta(buildData *PlannerMap) {
	// 全タイルをチェックして最外枠のタイルを壁で覆う
	for i := range buildData.Tiles {
		idx := resources.TileIdx(i)

		// 最外枠タイルの場合は無条件で壁にする
		if b.isBoundaryTile(buildData, idx) {
			buildData.Tiles[idx] = b.WallTileType
		}
	}
}

// isBoundaryTile はマップの最外枠のタイルかを判定する
func (b BoundaryWall) isBoundaryTile(buildData *PlannerMap, idx resources.TileIdx) bool {
	x, y := buildData.Level.XYTileCoord(idx)
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)

	return int(x) == 0 || int(x) == width-1 || int(y) == 0 || int(y) == height-1
}
