package mapspawner

import (
	"github.com/kijimaD/ruins/lib/raw"
)

// createMapspawnerTestRawMaster はテスト用のRawMasterを作成する
func createMapspawnerTestRawMaster() *raw.Master {
	rawMaster, _ := raw.Load(`
[[Tiles]]
Name = "Floor"
Description = "床タイル"
Walkable = true

[[Tiles]]
Name = "Wall"
Description = "壁タイル"
Walkable = false

[[Tiles]]
Name = "Empty"
Description = "空のタイル"
Walkable = false

[[Tiles]]
Name = "Dirt"
Description = "土タイル"
Walkable = true
`)
	return &rawMaster
}
