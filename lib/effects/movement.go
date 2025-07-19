package effects

import (
	"errors"
	"fmt"

	"github.com/kijimaD/ruins/lib/resources"
)

// MovementWarpNext は次の階層に移動するエフェクト
type MovementWarpNext struct{}

func (w MovementWarpNext) Apply(ctx *Context) error {
	gameResources := ctx.World.Resources.Game.(*resources.Game)
	gameResources.StateEvent = resources.StateEventWarpNext
	return nil
}

func (w MovementWarpNext) Validate(ctx *Context) error {
	if ctx.World.Manager == nil {
		return errors.New("Worldが設定されていません")
	}
	return nil
}

func (w MovementWarpNext) String() string {
	return "MovementWarpNext"
}

// MovementWarpEscape はダンジョンから脱出するエフェクト
type MovementWarpEscape struct{}

func (w MovementWarpEscape) Apply(ctx *Context) error {
	gameResources := ctx.World.Resources.Game.(*resources.Game)
	gameResources.StateEvent = resources.StateEventWarpEscape
	return nil
}

func (w MovementWarpEscape) Validate(ctx *Context) error {
	if ctx.World.Manager == nil {
		return errors.New("Worldが設定されていません")
	}
	return nil
}

func (w MovementWarpEscape) String() string {
	return "MovementWarpEscape"
}

// MovementWarpToFloor は特定の階層にワープするエフェクト（将来拡張用）
type MovementWarpToFloor struct {
	Floor int // ワープ先の階層
}

func (w MovementWarpToFloor) Apply(ctx *Context) error {
	// TODO: 特定階層へのワープ機能を実装
	// 現在は次階層へのワープと同じ動作
	gameResources := ctx.World.Resources.Game.(*resources.Game)
	gameResources.StateEvent = resources.StateEventWarpNext
	return nil
}

func (w MovementWarpToFloor) Validate(ctx *Context) error {
	if w.Floor < 1 {
		return fmt.Errorf("階層は1以上である必要があります: %d", w.Floor)
	}
	if ctx.World.Manager == nil {
		return errors.New("Worldが設定されていません")
	}
	return nil
}

func (w MovementWarpToFloor) String() string {
	return fmt.Sprintf("MovementWarpToFloor(%d)", w.Floor)
}
