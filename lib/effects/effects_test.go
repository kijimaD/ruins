package effects

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/stretchr/testify/assert"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestEffectSystem(t *testing.T) {
	t.Run("プロセッサーの基本動作", func(t *testing.T) {
		processor := NewProcessor()
		assert.NotNil(t, processor)
		assert.True(t, processor.IsEmpty())
		assert.Equal(t, 0, processor.QueueSize())
	})
	
	t.Run("エフェクトの文字列表現", func(t *testing.T) {
		damage := CombatDamage{Amount: 50, Source: DamageSourceWeapon}
		assert.Equal(t, "Damage(50, 武器)", damage.String())
		
		healing := CombatHealing{Amount: gc.NumeralAmount{Numeral: 30}}
		assert.Contains(t, healing.String(), "Healing")
		
		recovery := FullRecoveryHP{}
		assert.Equal(t, "FullRecoveryHP", recovery.String())
		
		warp := MovementWarpNext{}
		assert.Equal(t, "MovementWarpNext", warp.String())
	})
	
	t.Run("ターゲットセレクタの文字列表現", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)
		
		entity := world.Manager.NewEntity()
		
		singleTarget := TargetSingle{Entity: entity}
		assert.Equal(t, "TargetSingle", singleTarget.String())
		
		partyTargets := TargetParty{}
		assert.Equal(t, "TargetParty", partyTargets.String())
		
		allEnemies := TargetAllEnemies{}
		assert.Equal(t, "TargetAllEnemies", allEnemies.String())
		
		noTarget := TargetNone{}
		assert.Equal(t, "TargetNone", noTarget.String())
	})
	
	t.Run("エフェクトの検証", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)
		
		// 負のダメージは無効
		damage := CombatDamage{Amount: -10, Source: DamageSourceWeapon}
		ctx := &Context{Targets: []ecs.Entity{ecs.Entity(1)}}
		err = damage.Validate(world, ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ダメージは0以上")
		
		// ターゲットなしは無効
		validDamage := CombatDamage{Amount: 10, Source: DamageSourceWeapon}
		emptyCtx := &Context{Targets: []ecs.Entity{}}
		err = validDamage.Validate(world, emptyCtx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ダメージ対象が指定されていません")
	})
	
	t.Run("Poolsコンポーネントの検証", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)
		
		// Poolsコンポーネントがないエンティティを作成
		entityWithoutPools := world.Manager.NewEntity()
		
		// ダメージエフェクトでPoolsコンポーネントがないターゲットを検証
		damage := CombatDamage{Amount: 10, Source: DamageSourceWeapon}
		ctxWithInvalidTarget := &Context{
			Targets: []ecs.Entity{entityWithoutPools},
		}
		
		err = damage.Validate(world, ctxWithInvalidTarget)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Poolsコンポーネントがありません")
		
		// 回復エフェクトでも同様にチェック
		healing := CombatHealing{Amount: gc.NumeralAmount{Numeral: 30}}
		err = healing.Validate(world, ctxWithInvalidTarget)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Poolsコンポーネントがありません")
		
		// 非戦闘時の回復エフェクトでもチェック
		fullRecovery := FullRecoveryHP{}
		err = fullRecovery.Validate(world, ctxWithInvalidTarget)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Poolsコンポーネントがありません")
	})
	
	t.Run("無効なアイテムの検証", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)
		
		// 無効なアイテムIDの場合
		useItemInvalid := UseItem{Item: ecs.Entity(0)}
		ctx := &Context{
			Targets: []ecs.Entity{ecs.Entity(1)},
		}
		err = useItemInvalid.Validate(world, ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "無効なアイテムエンティティです")
	})
}