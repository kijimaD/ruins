package mapbuilder

import (
	"math"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/resources"
)

type LineCorridorBuilder struct{}

func (b LineCorridorBuilder) BuildMeta(buildData *BuilderMap) {
	b.BuildCorridors(buildData)
}

func (b LineCorridorBuilder) BuildCorridors(buildData *BuilderMap) {
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
			ps1 := bresenhamPoints(point{x: centerX, y: centerY}, point{x: destCenterX, y: destCenterY})
			ps2 := bresenhamPoints(point{x: centerX - 1, y: centerY}, point{x: destCenterX - 1, y: destCenterY})
			ps3 := bresenhamPoints(point{x: centerX, y: centerY - 1}, point{x: destCenterX, y: destCenterY - 1})
			points := []point{}
			points = append(points, ps1...)
			points = append(points, ps2...)
			points = append(points, ps3...)
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
	x gc.Row
	y gc.Col
}

// https://gist.github.com/s1moe2/a85a5da7e2af25397de326d9714a6bbc#file-bresenham-go
func bresenhamPoints(p1, p2 point) []point {
	dx := int(math.Abs(float64(p2.x) - float64(p1.x)))
	sx := gc.Row(-1)
	if p1.x < p2.x {
		sx = 1
	}

	dy := -int(math.Abs(float64(p2.y) - float64(p1.y)))
	sy := gc.Col(-1)
	if p1.y < p2.y {
		sy = 1
	}

	err := dx + dy

	points := []point{}

	for {
		points = append(points, point{p1.x, p1.y})

		if p1.x == p2.x && p1.y == p2.y {
			break
		}

		e2 := 2 * err

		if e2 >= dy {
			if p1.x == p2.x {
				break
			}
			err = err + dy
			p1.x = p1.x + sx
		}

		if e2 <= dx {
			if p1.y == p2.y {
				break
			}
			err = err + dx
			p1.y = p1.y + sy
		}
	}

	return points
}
