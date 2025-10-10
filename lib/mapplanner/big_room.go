package mapplanner

import (
	gc "github.com/kijimaD/ruins/lib/components"
)

// BigRoomPlanner は大部屋を生成するビルダー
// マップ全体の大部分を1つの部屋にする
type BigRoomPlanner struct{}

// PlanInitial は初期マップをプランする
func (b BigRoomPlanner) PlanInitial(planData *MetaPlan) error {
	// マップの境界を考慮して大きな部屋を1つ作成
	// 外周に1タイル分の壁を残す
	margin := 2
	room := gc.Rect{
		X1: gc.Tile(margin),
		Y1: gc.Tile(margin),
		X2: gc.Tile(int(planData.Level.TileWidth) - margin - 1),
		Y2: gc.Tile(int(planData.Level.TileHeight) - margin - 1),
	}

	// 部屋をリストに追加
	planData.Rooms = append(planData.Rooms, room)
	return nil
}

// BigRoomDraw は大部屋を描画し、ランダムにバリエーションを適用するビルダー
type BigRoomDraw struct{}

// PlanMeta は大部屋をタイルに描画し、ランダムにバリエーションを適用する
func (b BigRoomDraw) PlanMeta(planData *MetaPlan) {
	// まず基本の大部屋を描画
	b.drawBasicBigRoom(planData)

	// ランダムにバリエーションを選択して適用
	variantType := planData.RandomSource.Intn(5)

	switch variantType {
	case 0:
		// 通常の大部屋（何も追加しない）
	case 1:
		// 柱を追加
		b.applyPillars(planData)
	case 2:
		// 障害物を追加
		b.applyObstacles(planData)
	case 3:
		// 迷路パターンを追加
		b.applyMazePattern(planData)
	case 4:
		// 中央台座を追加
		b.applyCenterPlatform(planData)
	}
}

// drawBasicBigRoom は基本の大部屋を描画する
func (b BigRoomDraw) drawBasicBigRoom(planData *MetaPlan) {
	for _, room := range planData.Rooms {
		// 部屋の内部を床タイルで埋める
		for x := room.X1; x <= room.X2; x++ {
			for y := room.Y1; y <= room.Y2; y++ {
				idx := planData.Level.XYTileIndex(x, y)
				planData.Tiles[idx] = planData.GetTile(floorTileType)
			}
		}

		// 部屋の境界を壁で囲む
		// 上辺と下辺
		for x := room.X1 - 1; x <= room.X2+1; x++ {
			// 上辺
			if y := room.Y1 - 1; y >= 0 {
				idx := planData.Level.XYTileIndex(x, y)
				if planData.Tiles[idx].Name != floorTileType {
					planData.Tiles[idx] = planData.GetTile("Wall")
				}
			}
			// 下辺
			if y := room.Y2 + 1; int(y) < int(planData.Level.TileHeight) {
				idx := planData.Level.XYTileIndex(x, y)
				if planData.Tiles[idx].Name != floorTileType {
					planData.Tiles[idx] = planData.GetTile("Wall")
				}
			}
		}

		// 左辺と右辺
		for y := room.Y1; y <= room.Y2; y++ {
			// 左辺
			if x := room.X1 - 1; x >= 0 {
				idx := planData.Level.XYTileIndex(x, y)
				if planData.Tiles[idx].Name != floorTileType {
					planData.Tiles[idx] = planData.GetTile("Wall")
				}
			}
			// 右辺
			if x := room.X2 + 1; int(x) < int(planData.Level.TileWidth) {
				idx := planData.Level.XYTileIndex(x, y)
				if planData.Tiles[idx].Name != floorTileType {
					planData.Tiles[idx] = planData.GetTile("Wall")
				}
			}
		}
	}
}

// applyPillars は部屋に柱を追加する
func (b BigRoomDraw) applyPillars(planData *MetaPlan) {
	// 柱の間隔をランダムに決定（3-6の範囲）
	spacing := 3 + planData.RandomSource.Intn(4)

	for _, room := range planData.Rooms {
		// 柱の開始位置を計算（部屋の中心から対称に配置）
		startX := int(room.X1) + spacing
		startY := int(room.Y1) + spacing

		// 規則的に柱を配置
		for x := startX; x < int(room.X2); x += spacing + 1 {
			for y := startY; y < int(room.Y2); y += spacing + 1 {
				idx := planData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
				planData.Tiles[idx] = planData.GetTile("Wall")
			}
		}
	}
}

// applyObstacles は部屋にランダムな障害物を追加する
func (b BigRoomDraw) applyObstacles(planData *MetaPlan) {
	for _, room := range planData.Rooms {
		// 障害物の数を部屋のサイズに基づいて決定
		roomWidth := int(room.X2 - room.X1)
		roomHeight := int(room.Y2 - room.Y1)
		obstacleCount := (roomWidth * roomHeight) / 30 // 面積の1/30程度

		for i := 0; i < obstacleCount; i++ {
			// 部屋内のランダムな位置に障害物を配置
			x := int(room.X1) + 1 + planData.RandomSource.Intn(roomWidth-2)
			y := int(room.Y1) + 1 + planData.RandomSource.Intn(roomHeight-2)

			idx := planData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
			planData.Tiles[idx] = planData.GetTile("Wall")
		}
	}
}

// applyMazePattern は部屋に迷路パターンを追加する
func (b BigRoomDraw) applyMazePattern(planData *MetaPlan) {
	for _, room := range planData.Rooms {
		// 格子状に壁を配置し、ランダムに開口部を作る
		for x := int(room.X1) + 2; x < int(room.X2)-1; x += 3 {
			for y := int(room.Y1) + 1; y < int(room.Y2); y++ {
				// 縦の壁を配置（ランダムに開口部を作る）
				if planData.RandomSource.Float64() > 0.3 {
					idx := planData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
					planData.Tiles[idx] = planData.GetTile("Wall")
				}
			}
		}

		for y := int(room.Y1) + 2; y < int(room.Y2)-1; y += 3 {
			for x := int(room.X1) + 1; x < int(room.X2); x++ {
				// 横の壁を配置（ランダムに開口部を作る）
				if planData.RandomSource.Float64() > 0.3 {
					idx := planData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
					planData.Tiles[idx] = planData.GetTile("Wall")
				}
			}
		}
	}
}

// applyCenterPlatform は部屋に中央台座を追加する
func (b BigRoomDraw) applyCenterPlatform(planData *MetaPlan) {
	for _, room := range planData.Rooms {
		centerX := int(room.X1+room.X2) / 2
		centerY := int(room.Y1+room.Y2) / 2

		// 台座のサイズを部屋のサイズに基づいて決定
		platformSize := 2 + planData.RandomSource.Intn(3) // 2-4タイルの台座

		// 円形の台座を作成
		for dx := -platformSize; dx <= platformSize; dx++ {
			for dy := -platformSize; dy <= platformSize; dy++ {
				distance := dx*dx + dy*dy
				if distance <= platformSize*platformSize {
					x := centerX + dx
					y := centerY + dy
					if x >= int(room.X1) && x <= int(room.X2) &&
						y >= int(room.Y1) && y <= int(room.Y2) {
						idx := planData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
						// 外周は壁、内部は床のまま
						if distance >= (platformSize-1)*(platformSize-1) {
							planData.Tiles[idx] = planData.GetTile("Wall")
						}
					}
				}
			}
		}
	}
}
