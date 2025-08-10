package effects

import (
	"fmt"

	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
)

// MovementWarpNext は次の階層に移動するエフェクト
type MovementWarpNext struct{}

// Apply は次の階層へのワープエフェクトを適用する
func (m MovementWarpNext) Apply(world w.World, scope *Scope) error {
	if err := m.Validate(world, scope); err != nil {
		return err
	}

	gameResources := world.Resources.Game.(*resources.Game)
	gameResources.SetStateEvent(resources.StateEventWarpNext)
	return nil
}

// Validate は次の階層へのワープエフェクトの妥当性を検証する
func (m MovementWarpNext) Validate(_ w.World, _ *Scope) error {
	return nil
}

func (m MovementWarpNext) String() string {
	return "MovementWarpNext"
}

// MovementWarpEscape はダンジョンから脱出するエフェクト
type MovementWarpEscape struct{}

// Apply はダンジョンからの脱出ワープエフェクトを適用する
func (m MovementWarpEscape) Apply(world w.World, scope *Scope) error {
	if err := m.Validate(world, scope); err != nil {
		return err
	}

	gameResources := world.Resources.Game.(*resources.Game)
	gameResources.SetStateEvent(resources.StateEventWarpEscape)
	return nil
}

// Validate はダンジョンからの脱出ワープエフェクトの妥当性を検証する
func (m MovementWarpEscape) Validate(_ w.World, _ *Scope) error {
	return nil
}

func (m MovementWarpEscape) String() string {
	return "MovementWarpEscape"
}

// MovementWarpToFloor は特定の階層にワープするエフェクト（将来拡張用）
type MovementWarpToFloor struct {
	Floor int // ワープ先の階層
}

// Apply は特定の階層へのワープエフェクトを適用する
func (m MovementWarpToFloor) Apply(world w.World, scope *Scope) error {
	if err := m.Validate(world, scope); err != nil {
		return err
	}

	gameResources := world.Resources.Game.(*resources.Game)
	gameResources.SetStateEvent(resources.StateEventWarpNext)
	return nil
}

// Validate は特定の階層へのワープエフェクトの妥当性を検証する
func (m MovementWarpToFloor) Validate(_ w.World, _ *Scope) error {
	if m.Floor < 1 {
		return fmt.Errorf("階層は1以上である必要があります: %d", m.Floor)
	}
	return nil
}

func (m MovementWarpToFloor) String() string {
	return fmt.Sprintf("MovementWarpToFloor(%d)", m.Floor)
}
