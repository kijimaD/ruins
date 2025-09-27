package mapplanner

import (
	"math"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/resources"
)

// LineCorridorPlanner は直線廊下を生成するビルダー
type LineCorridorPlanner struct{}

// BuildMeta はメタデータをビルドする
func (b LineCorridorPlanner) BuildMeta(buildData *MetaPlan) {
	b.BuildCorridors(buildData)
}

// BuildCorridors は廊下をビルドする
func (b LineCorridorPlanner) BuildCorridors(buildData *MetaPlan) {
	// 接続済みの部屋。通路を2重に計算しないようにする
	connected := map[int]bool{}
	// 廊下のスライス
	for i, room := range buildData.Rooms {
		roomDistances := map[int]float64{}
		centerX, centerY := room.Center()
		for j, otherRoom := range buildData.Rooms {
			isExist := connected[j]
			if i != j && !isExist {
				oCenterX, oCenterY := otherRoom.Center()
				distance := math.Sqrt(math.Pow(float64(centerX-oCenterX), 2) + math.Pow(float64(centerY-oCenterY), 2))
				roomDistances[j] = float64(distance)
			}
		}

		if len(roomDistances) > 0 {
			var closestIdx int
			for k, v := range roomDistances {
				if roomDistances[closestIdx] < v {
					closestIdx = k
				}
			}
			destCenterX, destCenterY := buildData.Rooms[closestIdx].Center()

			points := []point{}
			corridorWidth := 3
			// 中心から上下左右にoffsetして複数のL字型廊下を生成
			for offsetX := -(corridorWidth / 2); offsetX <= corridorWidth/2; offsetX++ {
				for offsetY := -(corridorWidth / 2); offsetY <= corridorWidth/2; offsetY++ {
					startPoint := point{x: centerX + gc.Tile(offsetX), y: centerY + gc.Tile(offsetY)}
					endPoint := point{x: destCenterX + gc.Tile(offsetX), y: destCenterY + gc.Tile(offsetY)}
					corridorPoints := createLShapedCorridor(startPoint, endPoint)
					points = append(points, corridorPoints...)
				}
			}
			corridor := []resources.TileIdx{}
			for _, p := range points {
				idx := buildData.Level.XYTileIndex(p.x, p.y)
				if 0 < int(idx) && int(idx) < int(buildData.Level.TileWidth)*int(buildData.Level.TileHeight)-1 && buildData.Tiles[idx] == TileWall {
					buildData.Tiles[idx] = TileFloor
				}
				corridor = append(corridor, idx)
			}
			buildData.Corridors = append(buildData.Corridors, corridor)
		}
		connected[i] = true
	}
}

type point struct {
	x gc.Tile
	y gc.Tile
}

// createLShapedCorridor は横と縦のみのL字型廊下を生成する
func createLShapedCorridor(start, end point) []point {
	var points []point

	// 開始点から水平に移動
	current := start
	if start.x < end.x {
		// 右に移動
		for current.x <= end.x {
			points = append(points, current)
			if current.x == end.x {
				break
			}
			current.x++
		}
	} else {
		// 左に移動
		for current.x >= end.x {
			points = append(points, current)
			if current.x == end.x {
				break
			}
			current.x--
		}
	}

	// 水平移動後の位置から垂直に移動
	if start.y < end.y {
		// 下に移動
		for current.y < end.y {
			current.y++
			points = append(points, current)
		}
	} else if start.y > end.y {
		// 上に移動
		for current.y > end.y {
			current.y--
			points = append(points, current)
		}
	}

	return points
}
