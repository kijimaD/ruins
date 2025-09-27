// Package mapplanner の街用ビルダー
// 固定レイアウトの街マップを生成する
package mapplanner

import (
	gc "github.com/kijimaD/ruins/lib/components"
)

// TownEntityPlanner は街の固定マップ初期ビルダー
type TownEntityPlanner struct{}

// BuildInitial は市街地の固定マップ構造を初期化する
func (b TownEntityPlanner) BuildInitial(buildData *MetaPlan) {
	// 一般的な市街地レイアウト
	// - 中央に市庁舎・公園
	// - 周囲に住宅・商業・公共施設
	// - 機能的な街区配置

	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)
	centerX := width / 2
	centerY := height / 2

	// === 中央街区 ===

	// 市庁舎（中央公共施設）
	cityHall := gc.Rect{
		X1: gc.Tile(centerX - 6),
		Y1: gc.Tile(centerY - 6),
		X2: gc.Tile(centerX + 7),
		Y2: gc.Tile(centerY + 7),
	}
	buildData.Rooms = append(buildData.Rooms, cityHall)

	// === 北の文教区域 ===

	// 図書館（知識・学習施設）
	library := gc.Rect{
		X1: gc.Tile(centerX - 10),
		Y1: gc.Tile(centerY - 20),
		X2: gc.Tile(centerX + 3),
		Y2: gc.Tile(centerY - 10),
	}
	buildData.Rooms = append(buildData.Rooms, library)

	// 学校（教育施設）
	school := gc.Rect{
		X1: gc.Tile(centerX + 5),
		Y1: gc.Tile(centerY - 18),
		X2: gc.Tile(centerX + 13),
		Y2: gc.Tile(centerY - 9),
	}
	buildData.Rooms = append(buildData.Rooms, school)

	// === 東の居住区域 ===

	// 住宅1
	house1 := gc.Rect{
		X1: gc.Tile(centerX + 10),
		Y1: gc.Tile(centerY - 8),
		X2: gc.Tile(centerX + 20),
		Y2: gc.Tile(centerY + 2),
	}
	buildData.Rooms = append(buildData.Rooms, house1)

	// 住宅2
	house2 := gc.Rect{
		X1: gc.Tile(centerX + 9),
		Y1: gc.Tile(centerY + 4),
		X2: gc.Tile(centerX + 18),
		Y2: gc.Tile(centerY + 12),
	}
	buildData.Rooms = append(buildData.Rooms, house2)

	// === 南の公共区域 ===

	// 公民館（集会施設）
	communityHall := gc.Rect{
		X1: gc.Tile(centerX - 8),
		Y1: gc.Tile(centerY + 10),
		X2: gc.Tile(centerX + 9),
		Y2: gc.Tile(centerY + 22),
	}
	buildData.Rooms = append(buildData.Rooms, communityHall)

	// 事務所（管理施設）
	office := gc.Rect{
		X1: gc.Tile(centerX + 11),
		Y1: gc.Tile(centerY + 12),
		X2: gc.Tile(centerX + 19),
		Y2: gc.Tile(centerY + 20),
	}
	buildData.Rooms = append(buildData.Rooms, office)

	// === 西の商業区域 ===

	// 商店（商業施設）
	shop := gc.Rect{
		X1: gc.Tile(centerX - 20),
		Y1: gc.Tile(centerY - 6),
		X2: gc.Tile(centerX - 10),
		Y2: gc.Tile(centerY + 4),
	}
	buildData.Rooms = append(buildData.Rooms, shop)

	// 倉庫（物流施設）
	warehouse := gc.Rect{
		X1: gc.Tile(centerX - 18),
		Y1: gc.Tile(centerY + 6),
		X2: gc.Tile(centerX - 9),
		Y2: gc.Tile(centerY + 15),
	}
	buildData.Rooms = append(buildData.Rooms, warehouse)

	// === 郊外区域 ===

	// 小さな住宅（北西）
	cottage := gc.Rect{
		X1: gc.Tile(centerX - 16),
		Y1: gc.Tile(centerY - 15),
		X2: gc.Tile(centerX - 9),
		Y2: gc.Tile(centerY - 8),
	}
	buildData.Rooms = append(buildData.Rooms, cottage)

	// 公園（南東）
	park := gc.Rect{
		X1: gc.Tile(centerX + 12),
		Y1: gc.Tile(centerY + 10),
		X2: gc.Tile(centerX + 19),
		Y2: gc.Tile(centerY + 17),
	}
	buildData.Rooms = append(buildData.Rooms, park)
}

// TownMapDraw は街の固定マップを描画する
type TownMapDraw struct{}

// BuildMeta は街マップの描画を行う
func (b TownMapDraw) BuildMeta(buildData *MetaPlan) {
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)
	centerX := width / 2
	centerY := height / 2

	// 建物を描画
	b.drawRooms(buildData, width, height)

	// 道路網を描画
	b.drawRoadNetwork(buildData, width, height, centerX, centerY)

	// 斜めの空いている箇所を壁で埋める
	b.fillDiagonalGaps(buildData, width, height)
}

// drawRooms は部屋（建物）を描画する
func (b TownMapDraw) drawRooms(buildData *MetaPlan, width, height int) {
	for _, room := range buildData.Rooms {
		for x := room.X1; x < room.X2; x++ {
			for y := room.Y1; y < room.Y2; y++ {
				if x >= 0 && x < gc.Tile(width) && y >= 0 && y < gc.Tile(height) {
					idx := buildData.Level.XYTileIndex(x, y)
					buildData.Tiles[idx] = TileFloor
				}
			}
		}
	}
}

// drawRoadNetwork は街路網を描画する
func (b TownMapDraw) drawRoadNetwork(buildData *MetaPlan, width, height, centerX, centerY int) {
	// メイン通り（十字路）
	b.drawMainStreets(buildData, width, height, centerX, centerY)

	// 各地区の道路
	b.drawDistrictRoads(buildData, width, height, centerX, centerY)
}

// drawMainStreets はメインストリートを描画する
func (b TownMapDraw) drawMainStreets(buildData *MetaPlan, width, height, centerX, centerY int) {
	// メイン通り（南北）- 幅広の大通り
	for y := 0; y < height; y++ {
		for roadWidth := -2; roadWidth <= 2; roadWidth++ {
			x := centerX + roadWidth
			if x >= 0 && x < width && y >= 0 && y < height {
				idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
				buildData.Tiles[idx] = TileFloor
			}
		}
	}

	// メイン通り（東西）- 幅広の大通り
	for x := 0; x < width; x++ {
		for roadWidth := -2; roadWidth <= 2; roadWidth++ {
			y := centerY + roadWidth
			if x >= 0 && x < width && y >= 0 && y < height {
				idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
				buildData.Tiles[idx] = TileFloor
			}
		}
	}
}

// drawDistrictRoads は中央聖域から各区域への道を描画する
func (b TownMapDraw) drawDistrictRoads(buildData *MetaPlan, width, height, centerX, centerY int) {
	// 中央から北の学術区域への道
	b.drawScholarRoad(buildData, width, height, centerX, centerY)

	// 中央から東の工芸区域への道
	b.drawCraftRoad(buildData, width, height, centerX, centerY)

	// 中央から南の神殿区域への道
	b.drawTempleRoad(buildData, width, height, centerX, centerY)

	// 中央から西の交易区域への道
	b.drawTradeRoad(buildData, width, height, centerX, centerY)

	// 各エリア内の小道
	b.drawInnerPaths(buildData, width, height, centerX, centerY)
}

// drawScholarRoad は中央から北の学術区域への道を描画する
func (b TownMapDraw) drawScholarRoad(buildData *MetaPlan, width, _, centerX, centerY int) {
	// 中央から北へ向かう石畳の道（拡張された学術区域まで）
	for y := centerY - 6; y >= centerY-22 && y >= 0; y-- {
		for roadWidth := -1; roadWidth <= 1; roadWidth++ {
			x := centerX + roadWidth
			if x >= 0 && x < width {
				idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
				buildData.Tiles[idx] = TileFloor
			}
		}
	}
}

// drawCraftRoad は中央から東の工芸区域への道を描画する
func (b TownMapDraw) drawCraftRoad(buildData *MetaPlan, width, height, centerX, centerY int) {
	// 中央から東へ向かう石畳の道（拡張された工芸区域まで）
	for x := centerX + 7; x <= centerX+22 && x < width; x++ {
		for roadWidth := -1; roadWidth <= 1; roadWidth++ {
			y := centerY + roadWidth
			if y >= 0 && y < height {
				idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
				buildData.Tiles[idx] = TileFloor
			}
		}
	}
}

// drawTempleRoad は中央から南の神殿区域への大通りを描画する
func (b TownMapDraw) drawTempleRoad(buildData *MetaPlan, width, height, centerX, centerY int) {
	// 中央から南へ向かう幅広の参道（拡張された神殿区域まで）
	for y := centerY + 7; y <= centerY+24 && y < height; y++ {
		for roadWidth := -2; roadWidth <= 2; roadWidth++ {
			x := centerX + roadWidth
			if x >= 0 && x < width {
				idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
				buildData.Tiles[idx] = TileFloor
			}
		}
	}
}

// drawTradeRoad は中央から西の交易区域への道を描画する
func (b TownMapDraw) drawTradeRoad(buildData *MetaPlan, _, height, centerX, centerY int) {
	// 中央から西へ向かう石畳の道（拡張された交易区域まで）
	for x := centerX - 6; x >= centerX-22 && x >= 0; x-- {
		for roadWidth := -1; roadWidth <= 1; roadWidth++ {
			y := centerY + roadWidth
			if y >= 0 && y < height {
				idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
				buildData.Tiles[idx] = TileFloor
			}
		}
	}
}

// drawInnerPaths は各エリア内の小道を描画する
func (b TownMapDraw) drawInnerPaths(buildData *MetaPlan, width, height, centerX, centerY int) {
	// 北西の隠居者の庵への小道（拡張エリアに合わせて延長）
	for i := 0; i < 8; i++ {
		x := centerX - 6 - i
		y := centerY - 6 - i
		if x >= 0 && y >= 0 && x < width && y < height {
			idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
			buildData.Tiles[idx] = TileFloor
		}
	}

	// 南東の小さな祠への小道（拡張エリアに合わせて延長）
	for i := 0; i < 6; i++ {
		x := centerX + 7 + i
		y := centerY + 7 + i
		if x < width && y < height {
			idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
			buildData.Tiles[idx] = TileFloor
		}
	}

	// 各エリア間の横断路
	// 北東から南西への斜めの散策路（拡張）
	for i := -5; i <= 5; i++ {
		x := centerX + 12 + i
		y := centerY - 12 + i
		if x >= 0 && x < width && y >= 0 && y < height {
			idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
			buildData.Tiles[idx] = TileFloor
		}
	}

	// 学術区域の部屋間接続路
	for x := centerX - 10; x <= centerX+13; x++ {
		y := centerY - 15
		if x >= 0 && x < width && y >= 0 && y < height {
			idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
			buildData.Tiles[idx] = TileFloor
		}
	}

	// 東西エリアの部屋間接続路
	for y := centerY + 4; y <= centerY+12; y++ {
		x := centerX + 15
		if x >= 0 && x < width && y >= 0 && y < height {
			idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
			buildData.Tiles[idx] = TileFloor
		}
	}
}

// fillDiagonalGaps は斜めの空いている箇所を壁で埋める
func (b TownMapDraw) fillDiagonalGaps(buildData *MetaPlan, width, height int) {
	// マップ全体をスキャンして、斜めの空いている不自然な箇所を検出・修正
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			currentIdx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
			currentTile := buildData.Tiles[currentIdx]

			// 現在のタイルが壁の場合のみ処理
			if currentTile != TileWall {
				continue
			}

			// 斜めの隣接タイルが床で、直交する隣接タイルが壁の場合、
			// 斜めの移動が不自然になる箇所を特定
			if b.shouldFillDiagonalGap(buildData, x, y, width, height) {
				// 周囲の床タイルを壁に変更して繋がりを改善
				b.fillSurroundingGaps(buildData, x, y, width, height)
			}
		}
	}
}

// shouldFillDiagonalGap は斜めギャップを埋めるべきかを判定する
func (b TownMapDraw) shouldFillDiagonalGap(buildData *MetaPlan, x, y, _, _ int) bool {
	// 4つの直交方向の隣接タイル
	upIdx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y-1))
	downIdx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y+1))
	leftIdx := buildData.Level.XYTileIndex(gc.Tile(x-1), gc.Tile(y))
	rightIdx := buildData.Level.XYTileIndex(gc.Tile(x+1), gc.Tile(y))

	upTile := buildData.Tiles[upIdx]
	downTile := buildData.Tiles[downIdx]
	leftTile := buildData.Tiles[leftIdx]
	rightTile := buildData.Tiles[rightIdx]

	// 4つの斜め方向の隣接タイル
	upLeftIdx := buildData.Level.XYTileIndex(gc.Tile(x-1), gc.Tile(y-1))
	upRightIdx := buildData.Level.XYTileIndex(gc.Tile(x+1), gc.Tile(y-1))
	downLeftIdx := buildData.Level.XYTileIndex(gc.Tile(x-1), gc.Tile(y+1))
	downRightIdx := buildData.Level.XYTileIndex(gc.Tile(x+1), gc.Tile(y+1))

	upLeftTile := buildData.Tiles[upLeftIdx]
	upRightTile := buildData.Tiles[upRightIdx]
	downLeftTile := buildData.Tiles[downLeftIdx]
	downRightTile := buildData.Tiles[downRightIdx]

	// 斜めに床があるが、その両端の直交タイルが壁の場合、
	// 斜めの移動で行き詰まりが発生する可能性がある
	diagonalFloorCount := 0
	orthogonalWallCount := 0

	if upLeftTile.Walkable {
		diagonalFloorCount++
	}
	if upRightTile.Walkable {
		diagonalFloorCount++
	}
	if downLeftTile.Walkable {
		diagonalFloorCount++
	}
	if downRightTile.Walkable {
		diagonalFloorCount++
	}

	if upTile == TileWall {
		orthogonalWallCount++
	}
	if downTile == TileWall {
		orthogonalWallCount++
	}
	if leftTile == TileWall {
		orthogonalWallCount++
	}
	if rightTile == TileWall {
		orthogonalWallCount++
	}

	// 斜めに1つ以上の床があり、直交方向に3つ以上の壁がある場合、
	// 不自然な隙間の可能性が高い
	return diagonalFloorCount >= 1 && orthogonalWallCount >= 3
}

// fillSurroundingGaps は周囲の問題のあるギャップを埋める
func (b TownMapDraw) fillSurroundingGaps(buildData *MetaPlan, centerX, centerY, width, height int) {
	// 中心点から3x3の範囲を調査し、孤立した床タイルを壁に変更
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			x := centerX + dx
			y := centerY + dy

			if x < 0 || x >= width || y < 0 || y >= height {
				continue
			}

			idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
			tile := buildData.Tiles[idx]

			// 床タイルで、周囲に壁が多い場合は壁に変更
			if tile.Walkable && b.isSurroundedByWalls(buildData, x, y, width, height) {
				buildData.Tiles[idx] = TileWall
			}
		}
	}
}

// isSurroundedByWalls は指定位置が壁に囲まれているかを判定する
func (b TownMapDraw) isSurroundedByWalls(buildData *MetaPlan, x, y, width, height int) bool {
	wallCount := 0
	totalNeighbors := 0

	// 8方向の隣接タイルをチェック
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue // 中心点はスキップ
			}

			nx := x + dx
			ny := y + dy

			if nx < 0 || nx >= width || ny < 0 || ny >= height {
				wallCount++ // 境界外は壁として扱う
			} else {
				idx := buildData.Level.XYTileIndex(gc.Tile(nx), gc.Tile(ny))
				if buildData.Tiles[idx] == TileWall {
					wallCount++
				}
			}
			totalNeighbors++
		}
	}

	// 隣接タイルの75%以上が壁の場合、囲まれていると判定
	return float64(wallCount)/float64(totalNeighbors) >= 0.75
}
