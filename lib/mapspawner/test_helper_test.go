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
Walkable = true

[[tile]]
Name = "Wall"
Description = "壁タイル"
Walkable = false

[[tile]]
Name = "Empty"
Description = "空のタイル"
Walkable = false
`)
	return &rawMaster
}
