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

[Tiles.SpriteRender]
SpriteSheetName = "tile"
SpriteKey = "floor"
Depth = 0

[[Tiles]]
Name = "Wall"
Description = "壁タイル"
Walkable = false
BlocksView = true

[Tiles.SpriteRender]
SpriteSheetName = "tile"
SpriteKey = "wall"
Depth = 2

[[Tiles]]
Name = "Empty"
Description = "空のタイル"
Walkable = false
BlocksView = true

[Tiles.SpriteRender]
SpriteSheetName = "tile"
SpriteKey = "empty"
Depth = 0

[[Tiles]]
Name = "Dirt"
Description = "土タイル"
Walkable = true

[Tiles.SpriteRender]
SpriteSheetName = "tile"
SpriteKey = "dirt"
Depth = 0

[[Props]]
Name = "table"
Description = "テスト用テーブル"
BlockPass = true
BlockView = false

[Props.SpriteRender]
SpriteSheetName = "field"
SpriteKey = "prop_table"
Depth = 1

[[Props]]
Name = "chair"
Description = "テスト用椅子"
BlockPass = true
BlockView = false

[Props.SpriteRender]
SpriteSheetName = "field"
SpriteKey = "prop_chair"
Depth = 1

[[Props]]
Name = "bookshelf"
Description = "テスト用本棚"
BlockPass = true
BlockView = true

[Props.SpriteRender]
SpriteSheetName = "field"
SpriteKey = "prop_bookshelf"
Depth = 1

[[Props]]
Name = "barrel"
Description = "テスト用樽"
BlockPass = true
BlockView = false

[Props.SpriteRender]
SpriteSheetName = "field"
SpriteKey = "prop_barrel"
Depth = 1

[[Props]]
Name = "crate"
Description = "テスト用木箱"
BlockPass = true
BlockView = false

[Props.SpriteRender]
SpriteSheetName = "field"
SpriteKey = "prop_crate"
Depth = 1

[[Props]]
Name = "bed"
Description = "テスト用寝台"
BlockPass = true
BlockView = false

[Props.SpriteRender]
SpriteSheetName = "field"
SpriteKey = "prop_bed"
Depth = 1

[[Props]]
Name = "lantern"
Description = "テスト用ランタン"
BlockPass = true
BlockView = false

[Props.SpriteRender]
SpriteSheetName = "field"
SpriteKey = "prop_lantern"
Depth = 1

[Props.LightSource]
Radius = 7
Enabled = true

[Props.LightSource.Color]
R = 255
G = 200
B = 150
A = 255

[[Props]]
Name = "warp_next"
Description = "次のフロアへ進むワープポータル"
BlockPass = false
BlockView = false

[Props.SpriteRender]
SpriteSheetName = "field"
SpriteKey = "warp_next"
Depth = 1

[Props.WarpNextInteraction]

[[Props]]
Name = "warp_escape"
Description = "脱出用ワープポータル"
BlockPass = false
BlockView = false

[Props.SpriteRender]
SpriteSheetName = "field"
SpriteKey = "warp_escape"
Depth = 1

[Props.WarpEscapeInteraction]

[[Members]]
Name = "老兵"
SpriteSheetName = "field"
SpriteKey = "old_soldier"
FactionType = "FactionNeutral"

[Members.Dialog]
MessageKey = "old_soldier_greeting"

[Members.Attributes]
Vitality = 10
Strength = 10
Sensation = 10
Dexterity = 10
Agility = 10
Defense = 10
`)
	return &rawMaster
}
