package systems

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/worldhelper/simple"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// カード使用としてeffectに移したほうがいいかも
// effects.ItemTrigger() 的な
func BattleCommandSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	if firstEntity := ecs.GetFirst(world.Manager.Join(
		gameComponents.BattleCommand,
	)); firstEntity != nil {
		entity := *firstEntity
		cmd := gameComponents.BattleCommand.Get(entity).(*gc.BattleCommand)

		// wayから攻撃の属性を取り出す
		wayEntity := cmd.Way
		attack := gameComponents.Attack.Get(wayEntity).(*gc.Attack)
		if attack != nil {
			{
				ownerName := simple.GetName(world, cmd.Owner)
				wayName := simple.GetName(world, cmd.Way)
				entry := fmt.Sprintf("%sは、%sで攻撃。", ownerName.Name, wayName.Name)
				gamelog.BattleLog.Append(entry)
			}

			ownerEntity := cmd.Owner
			attrs := gameComponents.Attributes.Get(ownerEntity).(*gc.Attributes)
			damage := attack.Damage + attrs.Strength.Total
			effects.AddEffect(&ownerEntity, effects.Damage{Amount: damage}, effects.Single{Target: cmd.Target})
			{
				targetName := simple.GetName(world, cmd.Target)
				entry := fmt.Sprintf("%sに%dのダメージ。", targetName.Name, damage)
				gamelog.BattleLog.Append(entry)
			}
		}

		world.Manager.DeleteEntity(entity)
	}
}
