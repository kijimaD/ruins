package effects

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/game"
	w "github.com/kijimaD/ruins/lib/world"
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
		
		singleTarget := SingleTarget{Entity: entity}
		assert.Equal(t, "SingleTarget", singleTarget.String())
		
		partyTargets := PartyTargets{}
		assert.Equal(t, "PartyTargets", partyTargets.String())
		
		allEnemies := AllEnemies{}
		assert.Equal(t, "AllEnemies", allEnemies.String())
		
		noTarget := NoTarget{}
		assert.Equal(t, "NoTarget", noTarget.String())
	})
	
	t.Run("エフェクトの検証", func(t *testing.T) {
		// 負のダメージは無効
		damage := CombatDamage{Amount: -10, Source: DamageSourceWeapon}
		ctx := &Context{Targets: []ecs.Entity{ecs.Entity(1)}}
		err := damage.Validate(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ダメージは0以上")
		
		// ターゲットなしは無効
		validDamage := CombatDamage{Amount: 10, Source: DamageSourceWeapon}
		emptyCtx := &Context{Targets: []ecs.Entity{}}
		err = validDamage.Validate(emptyCtx)
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
			World:   world,
			Targets: []ecs.Entity{entityWithoutPools},
		}
		
		err = damage.Validate(ctxWithInvalidTarget)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Poolsコンポーネントがありません")
		
		// 回復エフェクトでも同様にチェック
		healing := CombatHealing{Amount: gc.NumeralAmount{Numeral: 30}}
		err = healing.Validate(ctxWithInvalidTarget)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Poolsコンポーネントがありません")
		
		// 非戦闘時の回復エフェクトでもチェック
		fullRecovery := FullRecoveryHP{}
		err = fullRecovery.Validate(ctxWithInvalidTarget)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Poolsコンポーネントがありません")
	})
	
	t.Run("Worldが設定されていない場合の検証", func(t *testing.T) {
		// Worldが設定されていない場合のエラー検証
		useItem := UseItem{Item: ecs.Entity(1)}
		ctxWithoutWorld := &Context{
			World: w.World{}, // 空のWorld（Manager == nil）
			Targets: []ecs.Entity{ecs.Entity(1)},
		}
		
		err := useItem.Validate(ctxWithoutWorld)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Worldが設定されていません")
		
		// 無効なアイテムIDの場合
		useItemInvalid := UseItem{Item: ecs.Entity(0)}
		err = useItemInvalid.Validate(ctxWithoutWorld)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "無効なアイテムエンティティです")
		
		// ダメージエフェクトでも同様
		damage := CombatDamage{Amount: 10, Source: DamageSourceWeapon}
		err = damage.Validate(ctxWithoutWorld)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Worldが設定されていません")
		
		// ワープエフェクトでも同様
		warp := MovementWarpNext{}
		err = warp.Validate(ctxWithoutWorld)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Worldが設定されていません")
	})
}