package mapplaner

import (
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/resources"
)

// CavePlanner は洞窟風レイアウトを生成するビルダー
// セルラーオートマトンを使用して有機的な洞窟形状を作成
// TODO: 全然わかってないので理解する
type CavePlanner struct{}

// BuildInitial は初期洞窟マップをビルドする
func (c CavePlanner) BuildInitial(buildData *PlannerMap) {
	// 初期状態をランダムに設定（30%の確率で壁、より広い空間を確保）
	for i := range buildData.Tiles {
		if buildData.RandomSource.Float64() < 0.30 {
			buildData.Tiles[i] = TileWall
		} else {
			buildData.Tiles[i] = TileFloor
		}
	}
}

// CaveCellularAutomata はセルラーオートマトンによる洞窟生成
type CaveCellularAutomata struct {
	// 反復回数
	Iterations int
}

// BuildMeta はセルラーオートマトンで洞窟を生成する
func (c CaveCellularAutomata) BuildMeta(buildData *PlannerMap) {
	iterations := c.Iterations
	if iterations <= 0 {
		iterations = 5 // デフォルト反復回数
	}

	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)

	// セルラーオートマトンを指定回数実行
	for iter := 0; iter < iterations; iter++ {
		newTiles := make([]Tile, len(buildData.Tiles))

		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))

				// 境界は常に壁
				if x == 0 || x == width-1 || y == 0 || y == height-1 {
					newTiles[idx] = TileWall
					continue
				}

				// 周囲の壁の数を数える
				wallCount := c.countWallsInRadius(buildData, x, y, 1)

				// ルール：周囲に6つ以上の壁があれば壁、そうでなければ床（より通路を確保）
				if wallCount >= 6 {
					newTiles[idx] = TileWall
				} else {
					newTiles[idx] = TileFloor
				}
			}
		}

		buildData.Tiles = newTiles
	}

	// 生成された洞窟から部屋領域を抽出
	c.extractCaveRooms(buildData)
}

// countWallsInRadius は指定半径内の壁タイル数を数える
func (c CaveCellularAutomata) countWallsInRadius(buildData *PlannerMap, centerX, centerY, radius int) int {
	wallCount := 0
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)

	for x := centerX - radius; x <= centerX+radius; x++ {
		for y := centerY - radius; y <= centerY+radius; y++ {
			// 境界外は壁とみなす
			if x < 0 || x >= width || y < 0 || y >= height {
				wallCount++
				continue
			}

			idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
			if buildData.Tiles[idx] == TileWall {
				wallCount++
			}
		}
	}

	return wallCount
}

// extractCaveRooms は洞窟から部屋領域を抽出する
func (c CaveCellularAutomata) extractCaveRooms(buildData *PlannerMap) {
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)
	visited := make([]bool, len(buildData.Tiles))

	// 連結している床領域を見つけて部屋として登録
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))

			if buildData.Tiles[idx] == TileFloor && !visited[idx] {
				// 洪水塗りつぶしで連結領域を見つける
				floorTiles := c.floodFill(buildData, x, y, visited)

				// 十分な大きさの領域のみ部屋として認識
				if len(floorTiles) >= 10 {
					// 境界ボックスを計算
					minX, maxX := x, x
					minY, maxY := y, y

					for _, tilePos := range floorTiles {
						tileX, tileY := buildData.Level.XYTileCoord(resources.TileIdx(tilePos))
						if int(tileX) < minX {
							minX = int(tileX)
						}
						if int(tileX) > maxX {
							maxX = int(tileX)
						}
						if int(tileY) < minY {
							minY = int(tileY)
						}
						if int(tileY) > maxY {
							maxY = int(tileY)
						}
					}

					// 部屋として登録
					room := gc.Rect{
						X1: gc.Tile(minX),
						Y1: gc.Tile(minY),
						X2: gc.Tile(maxX),
						Y2: gc.Tile(maxY),
					}
					buildData.Rooms = append(buildData.Rooms, room)
				}
			}
		}
	}
}

// floodFill は洪水塗りつぶしで連結する床タイルを見つける
func (c CaveCellularAutomata) floodFill(buildData *PlannerMap, startX, startY int, visited []bool) []int {
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)
	var result []int
	var queue [][2]int

	startIdx := buildData.Level.XYTileIndex(gc.Tile(startX), gc.Tile(startY))
	queue = append(queue, [2]int{startX, startY})
	visited[startIdx] = true

	// 4方向の隣接タイル
	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		x, y := current[0], current[1]
		idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
		result = append(result, int(idx))

		// 隣接タイルをチェック
		for _, dir := range directions {
			nx, ny := x+dir[0], y+dir[1]

			if nx >= 0 && nx < width && ny >= 0 && ny < height {
				nIdx := buildData.Level.XYTileIndex(gc.Tile(nx), gc.Tile(ny))

				if !visited[nIdx] && buildData.Tiles[nIdx] == TileFloor {
					visited[nIdx] = true
					queue = append(queue, [2]int{nx, ny})
				}
			}
		}
	}

	return result
}

// CavePathWidener は通路を広げる
type CavePathWidener struct{}

// BuildMeta は狭い通路を広げる
func (c CavePathWidener) BuildMeta(buildData *PlannerMap) {
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)

	// 床タイルの周囲1マスを床にして通路を広げる
	newTiles := make([]Tile, len(buildData.Tiles))
	copy(newTiles, buildData.Tiles)

	for x := 1; x < width-1; x++ {
		for y := 1; y < height-1; y++ {
			idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))

			// 現在が壁で、隣接に床がある場合
			if buildData.Tiles[idx] == TileWall {
				adjacentFloorCount := c.countAdjacentFloors(buildData, x, y)

				// 隣接する床が2個以上ある場合、30%の確率で床にする
				if adjacentFloorCount >= 2 && buildData.RandomSource.Float64() < 0.3 {
					newTiles[idx] = TileFloor
				}
			}
		}
	}

	buildData.Tiles = newTiles
}

// countAdjacentFloors は隣接する床タイルの数を数える
func (c CavePathWidener) countAdjacentFloors(buildData *PlannerMap, centerX, centerY int) int {
	count := 0
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)

	// 4方向の隣接タイルをチェック
	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	for _, dir := range directions {
		x, y := centerX+dir[0], centerY+dir[1]

		if x >= 0 && x < width && y >= 0 && y < height {
			idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
			if buildData.Tiles[idx] == TileFloor {
				count++
			}
		}
	}

	return count
}

// CaveStalactites は鍾乳石や岩の障害物を配置する
type CaveStalactites struct{}

// BuildMeta は洞窟内に鍾乳石を配置する
func (c CaveStalactites) BuildMeta(buildData *PlannerMap) {
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)

	// 床タイルの一部を鍾乳石（壁）に変換
	for x := 2; x < width-2; x++ {
		for y := 2; y < height-2; y++ {
			idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))

			if buildData.Tiles[idx] == TileFloor {
				// 2%の確率で鍾乳石を配置（確率を下げてより通行可能に）
				if buildData.RandomSource.Float64() < 0.02 {
					buildData.Tiles[idx] = TileWall
				}
			}
		}
	}
}

// CaveConnector は隔離された洞窟領域を接続する
type CaveConnector struct{}

// BuildMeta は隔離された洞窟領域を接続する
func (c CaveConnector) BuildMeta(buildData *PlannerMap) {
	if len(buildData.Rooms) < 2 {
		return
	}

	// 各部屋を最低1つの他の部屋と接続する
	for i := 0; i < len(buildData.Rooms)-1; i++ {
		room1 := buildData.Rooms[i]
		room2 := buildData.Rooms[i+1]

		c.createCaveTunnel(buildData, room1, room2)
	}
}

// createCaveTunnel は2つの洞窟領域間にトンネルを作成する
func (c CaveConnector) createCaveTunnel(buildData *PlannerMap, room1, room2 gc.Rect) {
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)

	// 各部屋の中心を計算
	center1X := int(room1.X1+room1.X2) / 2
	center1Y := int(room1.Y1+room1.Y2) / 2
	center2X := int(room2.X1+room2.X2) / 2
	center2Y := int(room2.Y1+room2.Y2) / 2

	// L字型のトンネルを作成（太さ2タイル）
	currentX, currentY := center1X, center1Y

	// 水平方向に移動
	for currentX != center2X {
		if currentX < center2X {
			currentX++
		} else {
			currentX--
		}

		// トンネルを太くする（上下1タイルずつ）
		for dy := -1; dy <= 1; dy++ {
			y := currentY + dy
			if y >= 1 && y < height-1 && currentX >= 1 && currentX < width-1 {
				idx := buildData.Level.XYTileIndex(gc.Tile(currentX), gc.Tile(y))
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

		// トンネルを太くする（左右1タイルずつ）
		for dx := -1; dx <= 1; dx++ {
			x := currentX + dx
			if x >= 1 && x < width-1 && currentY >= 1 && currentY < height-1 {
				idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(currentY))
				buildData.Tiles[idx] = TileFloor
			}
		}
	}
}

// NewCavePlanner は洞窟ビルダーを作成する
func NewCavePlanner(width gc.Tile, height gc.Tile, seed uint64) *PlannerChain {
	chain := NewPlannerChain(width, height, seed)
	chain.StartWith(CavePlanner{})
	chain.With(CaveCellularAutomata{Iterations: 3}) // セルラーオートマトン
	chain.With(CavePathWidener{})                   // 通路を広げる
	chain.With(CaveConnector{})                     // 隔離領域を接続
	chain.With(CaveStalactites{})                   // 鍾乳石配置
	chain.With(NewBoundaryWall(TileWall))           // 最外周を壁で囲む

	return chain
}
