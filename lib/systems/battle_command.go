package systems

import (
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// カード使用としてeffectに移したほうがいいかも
func BattleCommandSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.BattleCommand,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		cmd := gameComponents.BattleCommand.Get(entity).(*gc.BattleCommand)

		// wayから攻撃の属性を取り出す
		wayEntity := cmd.Way
		attack := gameComponents.Attack.Get(wayEntity).(*gc.Attack)
		if attack != nil {
			ownerEntity := cmd.Owner
			attrs := gameComponents.Attributes.Get(ownerEntity).(*gc.Attributes)
			damage := attack.Damage + attrs.Strength.Total
			effects.AddEffect(&ownerEntity, effects.Damage{Amount: damage}, effects.Single{Target: cmd.Target})
		}

		world.Manager.DeleteEntity(entity)
	}))
}
