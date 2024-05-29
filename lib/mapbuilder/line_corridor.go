package mapbuilder

import (
	"math"
)

type LineCorridorBuilder struct{}

func (b LineCorridorBuilder) BuildMeta(buildData *BuilderMap) {
	b.BuildCorridors(buildData)
}

func (b LineCorridorBuilder) BuildCorridors(buildData *BuilderMap) {
	// 接続済みの部屋。通路を2重に計算しないようにする
	connected := map[int]bool{}
	// 廊下のスライス
	corridors := [][]int{}
	for i, room := range buildData.Rooms {
		roomDistances := map[int]float64{}
		centerX, centerY := room.Center()
		for j, otherRoom := range buildData.Rooms {
			isExist, _ := connected[j]
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
			ps2 := bresenhamPoints(point{x: centerX - 1, y: centerY - 1}, point{x: destCenterX - 1, y: destCenterY - 1})
			ps3 := bresenhamPoints(point{x: centerX + 1, y: centerY + 1}, point{x: destCenterX + 1, y: destCenterY + 1})
			points := []point{}
			for _, p := range ps1 {
				points = append(points, p)
			}
			for _, p := range ps2 {
				points = append(points, p)
			}
			for _, p := range ps3 {
				points = append(points, p)
			}
			corridor := []int{}
			for _, p := range points {
				idx := buildData.Level.XYTileIndex(p.x, p.y)
				if 0 < idx && idx < int(buildData.Level.TileWidth)*int(buildData.Level.TileHeight)-1 && buildData.Tiles[idx] == TileWall {
					buildData.Tiles[idx] = TileFloor
				}
				corridor = append(corridor, idx)
			}
			corridors = append(corridors, corridor)
		}
		connected[i] = true
	}
}

type point struct {
	x int
	y int
}

// https://gist.github.com/s1moe2/a85a5da7e2af25397de326d9714a6bbc#file-bresenham-go
func bresenhamPoints(p1, p2 point) []point {
	dx := int(math.Abs(float64(p2.x) - float64(p1.x)))
	sx := -1
	if p1.x < p2.x {
		sx = 1
	}

	dy := -int(math.Abs(float64(p2.y) - float64(p1.y)))
	sy := -1
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
