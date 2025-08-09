package effects

import (
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// BattleEncounter はフィールドで敵と接触した時の戦闘遷移エフェクト
type BattleEncounter struct {
	// TODO: 不要なので消す
	PlayerEntity ecs.Entity
	// フィールド上の敵シンボル。勝利時に削除するのに使う
	FieldEnemyEntity ecs.Entity
}

// Apply は戦闘遷移を実行する
func (e *BattleEncounter) Apply(world w.World, _ *Scope) error {
	// プレイヤーの動きを停止
	if world.Components.Velocity.Get(e.PlayerEntity) != nil {
		velocity := world.Components.Velocity.Get(e.PlayerEntity).(*gc.Velocity)
		velocity.ThrottleMode = gc.ThrottleModeNope
		velocity.Speed = 0
	}

	// フィールド敵シンボルの動きを停止
	if world.Components.Velocity.Get(e.FieldEnemyEntity) != nil {
		velocity := world.Components.Velocity.Get(e.FieldEnemyEntity).(*gc.Velocity)
		velocity.ThrottleMode = gc.ThrottleModeNope
		velocity.Speed = 0
	}

	// 戦闘開始イベントを設定
	gameResources := world.Resources.Game.(*resources.Game)
	gameResources.StateEvent = resources.StateEventBattleStart

	// 戦闘参加エンティティ情報を一時保存
	// BattleState作成時にこれらの情報を取得して使用する
	gameResources.BattleTempData = &resources.BattleTempData{
		PlayerEntity:     e.PlayerEntity,
		FieldEnemyEntity: e.FieldEnemyEntity,
	}

	return nil
}

// Validate は戦闘遷移前の妥当性を検証する
func (e *BattleEncounter) Validate(_ w.World, _ *Scope) error {
	// 実際の検証は省略（簡易実装）
	return nil
}

// String はエフェクトの文字列表現を返す
func (e *BattleEncounter) String() string {
	return "BattleEncounter"
}
