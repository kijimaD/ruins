package mapbuilder

import (
	"fmt"
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
			fmt.Println(destCenterX, destCenterY)
			// TODO ================
			// 線を引く
			// 線をマス目ごとに分解してループ、タイルを切り替える
		}
		corridor := []int{}
		corridors = append(corridors, corridor)
		connected[i] = true
	}
}
