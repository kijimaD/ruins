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
SpriteSheetName = "field"
SpriteKey = "prop_table"
BlockPass = true
BlockView = false

[[Props]]
Name = "chair"
Description = "テスト用椅子"
SpriteSheetName = "field"
SpriteKey = "prop_chair"
BlockPass = true
BlockView = false

[[Props]]
Name = "bookshelf"
Description = "テスト用本棚"
SpriteSheetName = "field"
SpriteKey = "prop_bookshelf"
BlockPass = true
BlockView = true

[[Props]]
Name = "barrel"
Description = "テスト用樽"
SpriteSheetName = "field"
SpriteKey = "prop_barrel"
BlockPass = true
BlockView = false

[[Props]]
Name = "crate"
Description = "テスト用木箱"
SpriteSheetName = "field"
SpriteKey = "prop_crate"
BlockPass = true
BlockView = false

[[Props]]
Name = "bed"
Description = "テスト用寝台"
SpriteSheetName = "field"
SpriteKey = "prop_bed"
BlockPass = true
BlockView = false

[[Props]]
Name = "lantern"
Description = "テスト用ランタン"
SpriteSheetName = "field"
SpriteKey = "prop_lantern"
BlockPass = true
BlockView = false

[Props.LightSource]
Radius = 7
Enabled = true

[Props.LightSource.Color]
R = 255
G = 200
B = 150
A = 255
`)
	return &rawMaster
}
