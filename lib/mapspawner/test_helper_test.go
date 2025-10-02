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

[[Props]]
Name = "table"
Description = "テスト用テーブル"
SpriteKey = "prop_table"
BlockPass = true
BlockView = false

[[Props]]
Name = "chair"
Description = "テスト用椅子"
SpriteKey = "prop_chair"
BlockPass = true
BlockView = false

[[Props]]
Name = "bookshelf"
Description = "テスト用本棚"
SpriteKey = "prop_bookshelf"
BlockPass = true
BlockView = true

[[Props]]
Name = "barrel"
Description = "テスト用樽"
SpriteKey = "prop_barrel"
BlockPass = true
BlockView = false

[[Props]]
Name = "crate"
Description = "テスト用木箱"
SpriteKey = "prop_crate"
BlockPass = true
BlockView = false

[[Props]]
Name = "bed"
Description = "テスト用寝台"
SpriteKey = "prop_bed"
BlockPass = true
BlockView = false
`)
	return &rawMaster
}
