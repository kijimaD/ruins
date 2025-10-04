package actions

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/gamelog"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

const (
	// unknownEntityName は不明なエンティティ名
	unknownEntityName = "Unknown"
)

// UseItemActivity はActivityInterfaceの実装
type UseItemActivity struct{}

// Info はActivityInterfaceの実装
func (u *UseItemActivity) Info() ActivityInfo {
	return ActivityInfo{
		Name:            "アイテム使用",
		Description:     "アイテムを使う",
		Interruptible:   false,
		Resumable:       false,
		ActionPointCost: 100,
		TotalRequiredAP: 100,
	}
}

// String はActivityInterfaceの実装
func (u *UseItemActivity) String() string {
	return "UseItem"
}

// Validate はActivityInterfaceの実装
func (u *UseItemActivity) Validate(act *Activity, world w.World) error {
	if act.Target == nil {
		return ErrItemNotSet
	}

	item := *act.Target

	// アイテムエンティティにItemコンポーネントがあるかチェック
	if !item.HasComponent(world.Components.Item) {
		return ErrInvalidItem
	}

	// 何らかの効果があるかチェック
	hasEffect := item.HasComponent(world.Components.ProvidesHealing) ||
		item.HasComponent(world.Components.InflictsDamage)

	if !hasEffect {
		return ErrItemNoEffect
	}

	// アクターがPoolsコンポーネントを持っているかチェック
	if !act.Actor.HasComponent(world.Components.Pools) {
		return ErrActorNoPools
	}

	return nil
}

// Start はActivityInterfaceの実装
func (u *UseItemActivity) Start(act *Activity, _ w.World) error {
	act.Logger.Debug("アイテム使用開始", "actor", act.Actor, "item", *act.Target)
	return nil
}

// DoTurn はActivityInterfaceの実装
func (u *UseItemActivity) DoTurn(act *Activity, world w.World) error {
	if act.Target == nil {
		act.Cancel("アイテムが指定されていません")
		return ErrItemNotSet
	}

	item := *act.Target

	// 回復効果があるかチェック
	if healing := world.Components.ProvidesHealing.Get(item); healing != nil {
		healingComponent := healing.(*gc.ProvidesHealing)
		if err := u.applyHealing(act, world, healingComponent.Amount, item); err != nil {
			act.Cancel(fmt.Sprintf("回復処理エラー: %s", err.Error()))
			return err
		}
	}

	// ダメージ効果があるかチェック
	if damage := world.Components.InflictsDamage.Get(item); damage != nil {
		damageComponent := damage.(*gc.InflictsDamage)
		if err := u.applyDamage(act, world, damageComponent.Amount, item); err != nil {
			act.Cancel(fmt.Sprintf("ダメージ処理エラー: %s", err.Error()))
			return err
		}
	}

	// 消費可能アイテムの場合は削除
	if item.HasComponent(world.Components.Consumable) {
		world.Manager.DeleteEntity(item)
	}

	act.Complete()
	return nil
}

// Finish はActivityInterfaceの実装
func (u *UseItemActivity) Finish(act *Activity, _ w.World) error {
	act.Logger.Debug("アイテム使用完了", "actor", act.Actor)
	return nil
}

// Canceled はActivityInterfaceの実装
func (u *UseItemActivity) Canceled(act *Activity, _ w.World) error {
	act.Logger.Debug("アイテム使用キャンセル", "actor", act.Actor, "reason", act.CancelReason)
	return nil
}

// applyHealing は回復処理を適用する
func (u *UseItemActivity) applyHealing(act *Activity, world w.World, amounter gc.Amounter, item ecs.Entity) error {
	// Amounterから実際の回復量を取得
	var amount int
	switch amt := amounter.(type) {
	case gc.NumeralAmount:
		amount = amt.Calc()
	case gc.RatioAmount:
		// 最大HPに対する割合で回復
		pools := world.Components.Pools.Get(act.Actor).(*gc.Pools)
		amount = amt.Calc(pools.HP.Max)
	default:
		return fmt.Errorf("未対応のAmounterタイプ: %T", amounter)
	}
	pools := world.Components.Pools.Get(act.Actor).(*gc.Pools)

	beforeHP := pools.HP.Current
	pools.HP.Current += amount
	if pools.HP.Current > pools.HP.Max {
		pools.HP.Current = pools.HP.Max
	}
	actualHealing := pools.HP.Current - beforeHP

	// ログ出力
	u.logItemUse(act, world, item, actualHealing, true)

	return nil
}

// applyDamage はダメージ処理を適用する
func (u *UseItemActivity) applyDamage(act *Activity, world w.World, amount int, item ecs.Entity) error {
	pools := world.Components.Pools.Get(act.Actor).(*gc.Pools)

	pools.HP.Current -= amount
	if pools.HP.Current < 0 {
		pools.HP.Current = 0
	}

	// 死亡チェック
	if pools.HP.Current <= 0 {
		act.Actor.AddComponent(world.Components.Dead, &gc.Dead{})
		u.logDeath(world, act.Actor)
	}

	// ログ出力
	u.logItemUse(act, world, item, amount, false)

	return nil
}

// logItemUse はアイテム使用のログを出力する
func (u *UseItemActivity) logItemUse(act *Activity, world w.World, item ecs.Entity, amount int, isHealing bool) {
	// プレイヤーが関わる場合のみログ出力
	if !isPlayerActivity(act, world) {
		return
	}

	itemName := u.getItemName(item, world)
	actorName := u.getEntityName(act.Actor, world)

	logger := gamelog.New(gamelog.FieldLog)
	logger.Build(func(l *gamelog.Logger) {
		u.appendNameWithColor(l, act.Actor, actorName, world)
	}).Append(" は ").ItemName(itemName).Append(" を使った。")

	if isHealing {
		logger.Append(fmt.Sprintf(" HPが %d 回復した。", amount))
	} else {
		logger.Append(fmt.Sprintf(" %d のダメージを受けた。", amount))
	}

	logger.Log()
}

// logDeath は死亡ログを出力する
func (u *UseItemActivity) logDeath(world w.World, entity ecs.Entity) {
	name := u.getEntityName(entity, world)

	gamelog.New(gamelog.FieldLog).
		Build(func(l *gamelog.Logger) {
			u.appendNameWithColor(l, entity, name, world)
		}).
		Append(" は倒れた。").
		Log()
}

// getItemName はアイテムの名前を取得する
func (u *UseItemActivity) getItemName(item ecs.Entity, world w.World) string {
	name := world.Components.Name.Get(item)
	if name != nil {
		return name.(*gc.Name).Name
	}
	return "アイテム"
}

// getEntityName はエンティティの名前を取得する
func (u *UseItemActivity) getEntityName(entity ecs.Entity, world w.World) string {
	name := world.Components.Name.Get(entity)
	if name != nil {
		return name.(*gc.Name).Name
	}
	return unknownEntityName
}

// appendNameWithColor はエンティティの種類に応じて色付きで名前を追加する
func (u *UseItemActivity) appendNameWithColor(logger *gamelog.Logger, entity ecs.Entity, name string, world w.World) {
	if entity.HasComponent(world.Components.Player) {
		logger.PlayerName(name)
	} else if entity.HasComponent(world.Components.AIMoveFSM) {
		logger.NPCName(name)
	} else {
		logger.Append(name)
	}
}
