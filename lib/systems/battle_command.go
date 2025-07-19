package systems

import (
	"fmt"
	"log"
	"math/rand/v2"
	"sort"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/gamelog"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// BattleCommandSystem は戦闘中のコマンド処理を行う
// 1回1回実行ごとにコマンドを取り出して結果を得られるようになっている
// クリックごとにコマンドの結果を見られるようにするため
func BattleCommandSystem(world w.World) {

	// 持ち主が死んでいるBattleCommandを削除する
	world.Manager.Join(
		world.Components.BattleCommand,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		cmd := world.Components.BattleCommand.Get(entity).(*gc.BattleCommand)
		ownerPools := world.Components.Pools.Get(cmd.Owner).(*gc.Pools)
		if ownerPools.HP.Current == 0 {
			world.Manager.DeleteEntity(entity)
		}
	}))

	// ownerの素早さが一番高いものでソートする
	bcEntities := []ecs.Entity{}
	world.Manager.Join(
		world.Components.BattleCommand,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		bcEntities = append(bcEntities, entity)
	}))
	if len(bcEntities) == 0 {
		return
	}
	sort.Slice(bcEntities, func(i, j int) bool {
		ibc := world.Components.BattleCommand.Get(bcEntities[i]).(*gc.BattleCommand)
		jbc := world.Components.BattleCommand.Get(bcEntities[j]).(*gc.BattleCommand)

		iOwnerAttributes := world.Components.Attributes.Get(ibc.Owner).(*gc.Attributes)
		jOwnerAttributes := world.Components.Attributes.Get(jbc.Owner).(*gc.Attributes)

		// ランダムな小数を付加して等しくならないようにする
		isum := float64(iOwnerAttributes.Agility.Total) + rand.Float64()
		jsum := float64(jOwnerAttributes.Agility.Total) + rand.Float64()

		return isum < jsum
	})

	// 最も素早さが高いコマンドを実行する
	entity := bcEntities[0]
	cmd := world.Components.BattleCommand.Get(entity).(*gc.BattleCommand)
	{
		targetPools := world.Components.Pools.Get(cmd.Target).(*gc.Pools)
		// ターゲットが死んでいる場合は同じ派閥の別の生存エンティティに変更する
		if targetPools.HP.Current == 0 {
			p, err := worldhelper.NewByEntity(world, cmd.Target)
			if err != nil {
				log.Fatal(err)
			}

			var newTarget ecs.Entity
			if p.LivesLen() == 1 {
				newTarget = *p.Value()
			} else {
				newTarget, err = p.GetPrev()
				if err != nil {
					var err2 error
					newTarget, err2 = p.GetNext()
					if err2 != nil {
						log.Fatal(err)
					}
				}
			}
			cmd.Target = newTarget
		}
	}

	// wayから攻撃の属性を取り出す
	attack := world.Components.Attack.Get(cmd.Way).(*gc.Attack)
	card := world.Components.Card.Get(cmd.Way).(*gc.Card)
	if attack != nil {
		{
			ownerName := world.Components.Name.Get(cmd.Owner).(*gc.Name)
			wayName := world.Components.Name.Get(cmd.Way).(*gc.Name)
			entry := fmt.Sprintf("%sは、%sで攻撃。", ownerName.Name, wayName.Name)
			gamelog.BattleLog.Append(entry)
		}

		ownerEntity := cmd.Owner
		attrs := world.Components.Attributes.Get(ownerEntity).(*gc.Attributes)
		damage := attack.Damage + attrs.Strength.Total

		// 新しいエフェクトシステムを使用
		processor := effects.NewProcessor()
		damageEffect := effects.Damage{Amount: damage, Source: effects.DamageSourceWeapon}
		staminaEffect := effects.ConsumeStamina{Amount: gc.NumeralAmount{Numeral: card.Cost}}

		processor.AddEffectWithLogger(damageEffect, &ownerEntity, &gamelog.BattleLog, cmd.Target)
		processor.AddEffect(staminaEffect, &ownerEntity, cmd.Owner)
		if err := processor.Execute(world); err != nil {
			log.Printf("エフェクト実行エラー: %v", err)
		}
	}

	world.Manager.DeleteEntity(entity)
}
