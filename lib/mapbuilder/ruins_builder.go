package mapbuilder

import (
	gc "github.com/kijimaD/ruins/lib/components"
)

// RuinsBuilder は廃墟風レイアウトを生成するビルダー
// 建物の残骸や瓦礫が散在する廃墟を作成
type RuinsBuilder struct{}

// BuildInitial は初期廃墟マップをビルドする
func (r RuinsBuilder) BuildInitial(buildData *BuilderMap) {
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)

	// まず大きな廃墟建物の基盤を作成（3-5個の大きな矩形）
	ruinCount := 3 + buildData.RandomSource.Intn(3)

	for i := 0; i < ruinCount; i++ {
		// ランダムな位置とサイズで廃墟の基盤を作成
		buildingWidth := 8 + buildData.RandomSource.Intn(10)
		buildingHeight := 6 + buildData.RandomSource.Intn(8)

		x := 2 + buildData.RandomSource.Intn(width-buildingWidth-4)
		y := 2 + buildData.RandomSource.Intn(height-buildingHeight-4)

		// 建物の外壁を作成
		room := Rect{
			X1: gc.Row(x),
			Y1: gc.Col(y),
			X2: gc.Row(x + buildingWidth - 1),
			Y2: gc.Col(y + buildingHeight - 1),
		}
		buildData.Rooms = append(buildData.Rooms, room)
	}
}

// RuinsDraw は廃墟構造を描画する
type RuinsDraw struct{}

// BuildMeta は廃墟構造をタイルに描画する
func (r RuinsDraw) BuildMeta(buildData *BuilderMap) {
	// まず全体を床で埋める（屋外エリア）
	for i := range buildData.Tiles {
		buildData.Tiles[i] = TileFloor
	}

	// 各廃墟建物を処理
	for _, room := range buildData.Rooms {
		r.drawRuinedBuilding(buildData, room)
	}
}

// drawRuinedBuilding は破損した建物を描画する
func (r RuinsDraw) drawRuinedBuilding(buildData *BuilderMap, building Rect) {
	// 建物の外壁を描画（一部欠損あり）
	for x := building.X1; x <= building.X2; x++ {
		// 上辺
		if y := building.Y1; buildData.RandomSource.Float64() > 0.3 { // 70%の確率で壁
			idx := buildData.Level.XYTileIndex(x, y)
			buildData.Tiles[idx] = TileWall
		}
		// 下辺
		if y := building.Y2; buildData.RandomSource.Float64() > 0.3 {
			idx := buildData.Level.XYTileIndex(x, y)
			buildData.Tiles[idx] = TileWall
		}
	}

	for y := building.Y1; y <= building.Y2; y++ {
		// 左辺
		if x := building.X1; buildData.RandomSource.Float64() > 0.3 {
			idx := buildData.Level.XYTileIndex(x, y)
			buildData.Tiles[idx] = TileWall
		}
		// 右辺
		if x := building.X2; buildData.RandomSource.Float64() > 0.3 {
			idx := buildData.Level.XYTileIndex(x, y)
			buildData.Tiles[idx] = TileWall
		}
	}

	// 建物内部に部屋の仕切りを作成（一部破損）
	r.addInteriorWalls(buildData, building)
}

// addInteriorWalls は建物内部に仕切り壁を追加する
func (r RuinsDraw) addInteriorWalls(buildData *BuilderMap, building Rect) {
	buildingWidth := int(building.X2 - building.X1 + 1)
	buildingHeight := int(building.Y2 - building.Y1 + 1)

	// 建物が十分大きい場合のみ内部の仕切りを作成
	if buildingWidth >= 10 && buildingHeight >= 8 {
		// 縦の仕切り
		midX := building.X1 + gc.Row(buildingWidth/2)
		for y := building.Y1 + 2; y <= building.Y2-2; y++ {
			if buildData.RandomSource.Float64() > 0.4 { // 60%の確率で壁
				idx := buildData.Level.XYTileIndex(midX, y)
				buildData.Tiles[idx] = TileWall
			}
		}

		// 横の仕切り
		midY := building.Y1 + gc.Col(buildingHeight/2)
		for x := building.X1 + 2; x <= building.X2-2; x++ {
			if buildData.RandomSource.Float64() > 0.4 { // 60%の確率で壁
				idx := buildData.Level.XYTileIndex(x, midY)
				buildData.Tiles[idx] = TileWall
			}
		}
	}
}

// RuinsDebris は瓦礫や破片を配置する
type RuinsDebris struct{}

// BuildMeta は廃墟に瓦礫を配置する
func (r RuinsDebris) BuildMeta(buildData *BuilderMap) {
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)

	// 屋外エリアに瓦礫を散乱させる
	for x := 1; x < width-1; x++ {
		for y := 1; y < height-1; y++ {
			idx := buildData.Level.XYTileIndex(gc.Row(x), gc.Col(y))

			if buildData.Tiles[idx] == TileFloor {
				// 建物から離れた場所ほど瓦礫が少ない
				debrisChance := r.calculateDebrisChance(buildData, x, y)

				if buildData.RandomSource.Float64() < debrisChance {
					buildData.Tiles[idx] = TileWall // 瓦礫として壁タイルを使用
				}
			}
		}
	}
}

// calculateDebrisChance は瓦礫の配置確率を計算する
func (r RuinsDebris) calculateDebrisChance(buildData *BuilderMap, x, y int) float64 {
	// 最寄りの建物までの距離を計算
	minDistance := 1000.0

	for _, room := range buildData.Rooms {
		// 建物の中心からの距離
		centerX := float64(room.X1+room.X2) / 2.0
		centerY := float64(room.Y1+room.Y2) / 2.0

		dx := float64(x) - centerX
		dy := float64(y) - centerY
		distance := dx*dx + dy*dy // 平方根は不要

		if distance < minDistance {
			minDistance = distance
		}
	}

	// 距離に基づいて瓦礫確率を計算（近いほど高確率）
	if minDistance < 25 { // 建物の近く
		return 0.15
	} else if minDistance < 100 { // 中距離
		return 0.08
	}
	// 遠距離
	return 0.03
}

// RuinsCorridors は廃墟間を繋ぐ通路を作成する
type RuinsCorridors struct{}

// BuildMeta は廃墟間に通路を作成する
func (r RuinsCorridors) BuildMeta(buildData *BuilderMap) {
	if len(buildData.Rooms) < 2 {
		return
	}

	// 各建物間に通路を作成
	for i := 0; i < len(buildData.Rooms); i++ {
		for j := i + 1; j < len(buildData.Rooms); j++ {
			room1 := buildData.Rooms[i]
			room2 := buildData.Rooms[j]

			// 30%の確率で通路を作成（すべての建物を繋ぐわけではない）
			if buildData.RandomSource.Float64() < 0.3 {
				r.createRuinedPath(buildData, room1, room2)
			}
		}
	}
}

// createRuinedPath は破損した通路を作成する
func (r RuinsCorridors) createRuinedPath(buildData *BuilderMap, room1, room2 Rect) {
	// 各建物の中心を計算
	center1X := (room1.X1 + room1.X2) / 2
	center1Y := (room1.Y1 + room1.Y2) / 2
	center2X := (room2.X1 + room2.X2) / 2
	center2Y := (room2.Y1 + room2.Y2) / 2

	// L字型の通路を作成（一部破損）
	currentX, currentY := center1X, center1Y

	// 水平方向に移動
	for currentX != center2X {
		if currentX < center2X {
			currentX++
		} else {
			currentX--
		}

		// 70%の確率で通路を作成（部分的に破損）
		if buildData.RandomSource.Float64() > 0.3 {
			idx := buildData.Level.XYTileIndex(currentX, currentY)
			if buildData.Tiles[idx] == TileWall {
				buildData.Tiles[idx] = TileFloor
			}
		}
	}

	// 垂直方向に移動
	for currentY != center2Y {
		if currentY < center2Y {
			currentY++
		} else {
			currentY--
		}

		// 70%の確率で通路を作成
		if buildData.RandomSource.Float64() > 0.3 {
			idx := buildData.Level.XYTileIndex(currentX, currentY)
			if buildData.Tiles[idx] == TileWall {
				buildData.Tiles[idx] = TileFloor
			}
		}
	}
}

// NewRuinsBuilder は廃墟ビルダーを作成する
func NewRuinsBuilder(width gc.Row, height gc.Col, seed uint64) *BuilderChain {
	chain := NewBuilderChain(width, height, seed)
	chain.StartWith(RuinsBuilder{})
	chain.With(NewFillAll(TileWall))      // 全体を壁で埋める
	chain.With(RuinsDraw{})               // 廃墟構造を描画
	chain.With(RuinsDebris{})             // 瓦礫を配置
	chain.With(RuinsCorridors{})          // 通路を作成
	chain.With(NewBoundaryWall(TileWall)) // 最外周を壁で囲む

	return chain
}
