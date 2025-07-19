package effects

import (
	"fmt"

	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
)

// MovementWarpNext は次の階層に移動するエフェクト
type MovementWarpNext struct{}

func (m MovementWarpNext) Apply(world w.World, ctx *Context) error {
	gameResources := world.Resources.Game.(*resources.Game)
	gameResources.StateEvent = resources.StateEventWarpNext
	return nil
}

func (m MovementWarpNext) Validate(world w.World, ctx *Context) error {
	return nil
}

func (w MovementWarpNext) String() string {
	return "MovementWarpNext"
}

// MovementWarpEscape はダンジョンから脱出するエフェクト
type MovementWarpEscape struct{}

func (m MovementWarpEscape) Apply(world w.World, ctx *Context) error {
	gameResources := world.Resources.Game.(*resources.Game)
	gameResources.StateEvent = resources.StateEventWarpEscape
	return nil
}

func (m MovementWarpEscape) Validate(world w.World, ctx *Context) error {
	return nil
}

func (w MovementWarpEscape) String() string {
	return "MovementWarpEscape"
}

// MovementWarpToFloor は特定の階層にワープするエフェクト（将来拡張用）
type MovementWarpToFloor struct {
	Floor int // ワープ先の階層
}

func (m MovementWarpToFloor) Apply(world w.World, ctx *Context) error {
	gameResources := world.Resources.Game.(*resources.Game)
	
	// 現在の階層と比較して適切なワープを決定
	currentDepth := gameResources.Depth
	
	if m.Floor > currentDepth {
		// より深い階層へのワープ（次へ進む）
		gameResources.StateEvent = resources.StateEventWarpNext
	} else if m.Floor < currentDepth {
		// より浅い階層へのワープ（脱出として扱う）
		gameResources.StateEvent = resources.StateEventWarpEscape
	} else {
		// 同じ階層の場合は何もしない
		return nil
	}
	
	return nil
}

func (m MovementWarpToFloor) Validate(world w.World, ctx *Context) error {
	if m.Floor < 1 {
		return fmt.Errorf("階層は1以上である必要があります: %d", m.Floor)
	}
	return nil
}

func (w MovementWarpToFloor) String() string {
	return fmt.Sprintf("MovementWarpToFloor(%d)", w.Floor)
}
