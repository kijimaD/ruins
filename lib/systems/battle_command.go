package systems

import (
	"fmt"
	"math/rand/v2"
	"sort"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/gamelog"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 1回1回実行して結果を得られるようになっている
// クリックごとにコマンドの結果を見るということができる
// TODO: Targetがすでに死んでいたときを考慮していない。死んでいた場合は次の選択肢の敵にターゲットを変えるのが自然だろう
func BattleCommandSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	// ownerの素早さが一番高いものでソートする
	bcEntities := []ecs.Entity{}
	world.Manager.Join(
		gameComponents.BattleCommand,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		bcEntities = append(bcEntities, entity)
	}))
	if len(bcEntities) == 0 {
		return
	}
	sort.Slice(bcEntities, func(i, j int) bool {
		ibc := gameComponents.BattleCommand.Get(bcEntities[i]).(*gc.BattleCommand)
		jbc := gameComponents.BattleCommand.Get(bcEntities[j]).(*gc.BattleCommand)

		iOwnerAttributes := gameComponents.Attributes.Get(ibc.Owner).(*gc.Attributes)
		jOwnerAttributes := gameComponents.Attributes.Get(jbc.Owner).(*gc.Attributes)

		// ランダムな小数を付加して等しくならないようにする
		isum := float64(iOwnerAttributes.Agility.Total) + rand.Float64()
		jsum := float64(jOwnerAttributes.Agility.Total) + rand.Float64()

		return isum < jsum
	})

	entity := bcEntities[0]
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
	}

	world.Manager.DeleteEntity(entity)
}
