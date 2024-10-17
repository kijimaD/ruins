package systems

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/gamelog"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 1回1回実行して結果を得られるようになっている
// クリックごとにコマンドの結果を見るということができる
func BattleCommandSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	if firstEntity := ecs.GetFirst(world.Manager.Join(
		gameComponents.BattleCommand,
	)); firstEntity != nil {
		entity := *firstEntity
		cmd := gameComponents.BattleCommand.Get(entity).(*gc.BattleCommand)

		// wayから攻撃の属性を取り出す
		attack := gameComponents.Attack.Get(cmd.Way).(*gc.Attack)
		card := gameComponents.Card.Get(cmd.Way).(*gc.Card)
		if attack != nil {
			{
				ownerName := gameComponents.Name.Get(cmd.Owner).(*gc.Name)
				wayName := gameComponents.Name.Get(cmd.Way).(*gc.Name)
				entry := fmt.Sprintf("%sは、%sで攻撃。", ownerName.Name, wayName.Name)
				gamelog.BattleLog.Append(entry)
			}

			ownerEntity := cmd.Owner
			attrs := gameComponents.Attributes.Get(ownerEntity).(*gc.Attributes)
			damage := attack.Damage + attrs.Strength.Total
			effects.AddEffect(&ownerEntity, effects.Damage{Amount: damage}, effects.Single{Target: cmd.Target})
			effects.AddEffect(&ownerEntity, effects.ConsumptionStamina{Amount: gc.NumeralAmount{Numeral: card.Cost}}, effects.Single{Target: cmd.Owner})
			{
				targetName := gameComponents.Name.Get(cmd.Target).(*gc.Name)
				entry := fmt.Sprintf("%sに%dのダメージ。", targetName.Name, damage)
				gamelog.BattleLog.Append(entry)
			}
		}

		world.Manager.DeleteEntity(entity)
	}
}
