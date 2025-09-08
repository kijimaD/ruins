package mapbuilder

import (
	"math"

	gc "github.com/kijimaD/ruins/lib/components"
)

// ForestBuilder は森風レイアウトを生成するビルダー
// 木々が点在し、自然な通路を持つ森を作成
type ForestBuilder struct{}

// BuildInitial は初期森マップをビルドする
func (f ForestBuilder) BuildInitial(buildData *BuilderMap) {
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)

	// 森の中に小さな空き地（部屋）をいくつか作成
	clearingCount := 3 + buildData.RandomSource.Intn(4)

	for i := 0; i < clearingCount; i++ {
		// 空き地のサイズとランダムな位置
		clearingWidth := 4 + buildData.RandomSource.Intn(6)
		clearingHeight := 4 + buildData.RandomSource.Intn(6)

		x := 3 + buildData.RandomSource.Intn(width-clearingWidth-6)
		y := 3 + buildData.RandomSource.Intn(height-clearingHeight-6)

		// 円形に近い空き地を作成
		centerX := x + clearingWidth/2
		centerY := y + clearingHeight/2
		radius := math.Min(float64(clearingWidth), float64(clearingHeight)) / 2.0

		clearing := Rect{
			X1: gc.Tile(centerX - int(radius)),
			Y1: gc.Tile(centerY - int(radius)),
			X2: gc.Tile(centerX + int(radius)),
			Y2: gc.Tile(centerY + int(radius)),
		}
		buildData.Rooms = append(buildData.Rooms, clearing)
	}
}

// ForestTerrain は森の基本地形を生成する
type ForestTerrain struct{}

// BuildMeta は森の基本地形をタイルに描画する
func (f ForestTerrain) BuildMeta(buildData *BuilderMap) {
	// まず全体を床で埋める（森の地面）
	for i := range buildData.Tiles {
		buildData.Tiles[i] = TileFloor
	}

	// 空き地を床として確保
	for _, clearing := range buildData.Rooms {
		f.createCircularClearing(buildData, clearing)
	}
}

// createCircularClearing は円形の空き地を作成する
func (f ForestTerrain) createCircularClearing(buildData *BuilderMap, clearing Rect) {
	centerX := float64(clearing.X1+clearing.X2) / 2.0
	centerY := float64(clearing.Y1+clearing.Y2) / 2.0
	radius := math.Min(float64(clearing.X2-clearing.X1), float64(clearing.Y2-clearing.Y1)) / 2.0

	for x := clearing.X1 - 1; x <= clearing.X2+1; x++ {
		for y := clearing.Y1 - 1; y <= clearing.Y2+1; y++ {
			dx := float64(x) - centerX
			dy := float64(y) - centerY
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance <= radius {
				idx := buildData.Level.XYTileIndex(x, y)
				buildData.Tiles[idx] = TileFloor
			}
		}
	}
}

// ForestTrees は森に木々を配置する
type ForestTrees struct{}

// BuildMeta は森に木を配置する
func (f ForestTrees) BuildMeta(buildData *BuilderMap) {
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)

	// 森全体に木を配置（60%の密度）
	for x := 1; x < width-1; x++ {
		for y := 1; y < height-1; y++ {
			idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))

			if buildData.Tiles[idx] == TileFloor {
				// 空き地の近くでは木の密度を下げる
				treeDensity := f.calculateTreeDensity(buildData, x, y)

				if buildData.RandomSource.Float64() < treeDensity {
					buildData.Tiles[idx] = TileWall // 木として壁タイルを使用

					// 大きな木の場合、周囲にも追加の木を配置
					if buildData.RandomSource.Float64() < 0.2 { // 20%の確率で大木
						f.placeLargeTree(buildData, x, y)
					}
				}
			}
		}
	}
}

// calculateTreeDensity は位置に基づいて木の密度を計算する
func (f ForestTrees) calculateTreeDensity(buildData *BuilderMap, x, y int) float64 {
	baseDensity := 0.6 // 基本密度60%

	// 空き地からの距離に基づいて密度を調整
	minDistanceToClearing := 1000.0

	for _, clearing := range buildData.Rooms {
		centerX := float64(clearing.X1+clearing.X2) / 2.0
		centerY := float64(clearing.Y1+clearing.Y2) / 2.0

		dx := float64(x) - centerX
		dy := float64(y) - centerY
		distance := math.Sqrt(dx*dx + dy*dy)

		if distance < minDistanceToClearing {
			minDistanceToClearing = distance
		}
	}

	// 空き地に近いほど木の密度を下げる
	if minDistanceToClearing < 5 {
		return baseDensity * 0.3 // 空き地の近くは30%
	} else if minDistanceToClearing < 10 {
		return baseDensity * 0.7 // 中距離は70%
	}
	return baseDensity // 遠くは基本密度
}

// placeLargeTree は大きな木を配置する
func (f ForestTrees) placeLargeTree(buildData *BuilderMap, centerX, centerY int) {
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)

	// 大木の周囲2x2または3x3エリアに追加の木を配置
	size := 1 + buildData.RandomSource.Intn(2) // 1または2

	for dx := -size; dx <= size; dx++ {
		for dy := -size; dy <= size; dy++ {
			x, y := centerX+dx, centerY+dy

			if x >= 0 && x < width && y >= 0 && y < height {
				idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))

				if buildData.Tiles[idx] == TileFloor && buildData.RandomSource.Float64() < 0.7 {
					buildData.Tiles[idx] = TileWall
				}
			}
		}
	}
}

// ForestPaths は森の中に自然な通路を作成する
type ForestPaths struct{}

// BuildMeta は空き地間に自然な通路を作成する
func (f ForestPaths) BuildMeta(buildData *BuilderMap) {
	if len(buildData.Rooms) < 2 {
		return
	}

	// 各空き地を他の空き地と繋ぐ
	for i := 0; i < len(buildData.Rooms); i++ {
		for j := i + 1; j < len(buildData.Rooms); j++ {
			// 距離が近い空き地のみ繋ぐ（50%の確率）
			if f.shouldCreatePath(buildData, buildData.Rooms[i], buildData.Rooms[j]) {
				f.createNaturalPath(buildData, buildData.Rooms[i], buildData.Rooms[j])
			}
		}
	}
}

// shouldCreatePath は通路を作成するかどうかを判定する
func (f ForestPaths) shouldCreatePath(buildData *BuilderMap, room1, room2 Rect) bool {
	// 空き地間の距離を計算
	center1X := float64(room1.X1+room1.X2) / 2.0
	center1Y := float64(room1.Y1+room1.Y2) / 2.0
	center2X := float64(room2.X1+room2.X2) / 2.0
	center2Y := float64(room2.Y1+room2.Y2) / 2.0

	dx := center1X - center2X
	dy := center1Y - center2Y
	distance := math.Sqrt(dx*dx + dy*dy)

	// 近い空き地ほど繋がりやすい
	if distance < 15 {
		return buildData.RandomSource.Float64() < 0.8
	} else if distance < 25 {
		return buildData.RandomSource.Float64() < 0.4
	}
	return buildData.RandomSource.Float64() < 0.1
}

// createNaturalPath は自然な曲線状の通路を作成する
func (f ForestPaths) createNaturalPath(buildData *BuilderMap, room1, room2 Rect) {
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)

	center1X := int(room1.X1+room1.X2) / 2
	center1Y := int(room1.Y1+room1.Y2) / 2
	center2X := int(room2.X1+room2.X2) / 2
	center2Y := int(room2.Y1+room2.Y2) / 2

	// ベジェ曲線風の自然な通路を作成
	steps := 50
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)

		// 制御点を追加してカーブを作る
		midX := (center1X + center2X) / 2
		midY := (center1Y + center2Y) / 2

		// ランダムな偏向を追加
		randomOffsetX := int(float64(buildData.RandomSource.Intn(11)-5) * (1.0 - math.Abs(t-0.5)*2))
		randomOffsetY := int(float64(buildData.RandomSource.Intn(11)-5) * (1.0 - math.Abs(t-0.5)*2))

		// 2次ベジェ曲線の近似
		x := int((1-t)*(1-t)*float64(center1X) + 2*(1-t)*t*float64(midX+randomOffsetX) + t*t*float64(center2X))
		y := int((1-t)*(1-t)*float64(center1Y) + 2*(1-t)*t*float64(midY+randomOffsetY) + t*t*float64(center2Y))

		// 通路を作成（少し幅を持たせる）
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				nx, ny := x+dx, y+dy
				if nx >= 0 && nx < width && ny >= 0 && ny < height {
					idx := buildData.Level.XYTileIndex(gc.Tile(nx), gc.Tile(ny))

					// 70%の確率で通路を作成（自然な感じに）
					if buildData.RandomSource.Float64() < 0.7 {
						buildData.Tiles[idx] = TileFloor
					}
				}
			}
		}
	}
}

// ForestWildlife は森に野生動物の痕跡（小さな空き地）を追加する
type ForestWildlife struct{}

// BuildMeta は森に小さな動物の痕跡を追加する
func (f ForestWildlife) BuildMeta(buildData *BuilderMap) {
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)

	// 小さな動物の通り道や巣穴を作成
	wildlifeSpotCount := 2 + buildData.RandomSource.Intn(4)

	for i := 0; i < wildlifeSpotCount; i++ {
		x := 2 + buildData.RandomSource.Intn(width-4)
		y := 2 + buildData.RandomSource.Intn(height-4)

		// 小さな円形の空き地を作成
		radius := 1 + buildData.RandomSource.Intn(2)

		for dx := -radius; dx <= radius; dx++ {
			for dy := -radius; dy <= radius; dy++ {
				nx, ny := x+dx, y+dy
				if nx >= 0 && nx < width && ny >= 0 && ny < height {
					distance := math.Sqrt(float64(dx*dx + dy*dy))
					if distance <= float64(radius) {
						idx := buildData.Level.XYTileIndex(gc.Tile(nx), gc.Tile(ny))
						buildData.Tiles[idx] = TileFloor
					}
				}
			}
		}
	}
}

// NewForestBuilder は森ビルダーを作成する
func NewForestBuilder(width gc.Tile, height gc.Tile, seed uint64) *BuilderChain {
	chain := NewBuilderChain(width, height, seed)
	chain.StartWith(ForestBuilder{})
	chain.With(ForestTerrain{})           // 基本地形を生成
	chain.With(ForestTrees{})             // 木を配置
	chain.With(ForestPaths{})             // 自然な通路を作成
	chain.With(ForestWildlife{})          // 野生動物の痕跡を追加
	chain.With(NewBoundaryWall(TileWall)) // 最外周を壁で囲む

	return chain
}
