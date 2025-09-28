package mapspawner

import (
	"github.com/kijimaD/ruins/lib/raw"
)

// createMapspawnerTestRawMaster はテスト用のRawMasterを作成する
func createMapspawnerTestRawMaster() *raw.Master {
	rawMaster, _ := raw.Load(`
[[tile]]
Name = "Floor"
Description = "床タイル"
Type = "FLOOR"
Walkable = true
BlocksLOS = false
SpriteNumber = 1

[[tile]]
Name = "Wall"
Description = "壁タイル"
Type = "WALL"
Walkable = false
BlocksLOS = true
SpriteNumber = 2

[[tile]]
Name = "Empty"
Description = "空のタイル"
Type = "EMPTY"
Walkable = false
BlocksLOS = false
SpriteNumber = 0
`)
	return &rawMaster
}
