package effects

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/engine/entities"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/stretchr/testify/assert"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// テスト用のヘルパー関数：エンティティをPoolsコンポーネント付きで作成
func createTestPlayerEntity(world w.World, hp, sp int) ecs.Entity {
	componentList := entities.ComponentList{}
	componentList.Game = append(componentList.Game, gc.GameComponentList{
		Pools: &gc.Pools{
			HP: gc.Pool{Current: hp, Max: 100},
			SP: gc.Pool{Current: sp, Max: 50},
		},
		Name: &gc.Name{Name: "テストプレイヤー"},
	})

	entityList := entities.AddEntities(world, componentList)
	return entityList[0]
}

// テスト用のヘルパー関数：基本的なエンティティを作成
func createTestEntity(world w.World) ecs.Entity {
	entity := world.Manager.NewEntity()
	return entity
}

// テスト用のヘルパー関数：アイテムエンティティを作成
func createTestHealingItem(world w.World, healAmount int) ecs.Entity {
	componentList := entities.ComponentList{}
	componentList.Game = append(componentList.Game, gc.GameComponentList{
		Item:            &gc.Item{},
		ProvidesHealing: &gc.ProvidesHealing{Amount: gc.NumeralAmount{Numeral: healAmount}},
		Consumable:      &gc.Consumable{},
	})

	entityList := entities.AddEntities(world, componentList)
	return entityList[0]
}

// テスト用のヘルパー関数：基本アイテムエンティティを作成
func createTestBasicItem(world w.World, name string) ecs.Entity {
	componentList := entities.ComponentList{}
	componentList.Game = append(componentList.Game, gc.GameComponentList{
		Item: &gc.Item{},
		Name: &gc.Name{Name: name},
	})

	entityList := entities.AddEntities(world, componentList)
	return entityList[0]
}

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

		entity := createTestEntity(world)

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
		ctx := &Scope{Targets: []ecs.Entity{ecs.Entity(1)}}
		err = damage.Validate(world, ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ダメージは0以上")

		// ターゲットなしは無効
		validDamage := CombatDamage{Amount: 10, Source: DamageSourceWeapon}
		emptyCtx := &Scope{Targets: []ecs.Entity{}}
		err = validDamage.Validate(world, emptyCtx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ダメージ対象が指定されていません")
	})

	t.Run("Poolsコンポーネントの検証", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		// Poolsコンポーネントがないエンティティを作成
		entityWithoutPools := createTestEntity(world)

		// ダメージエフェクトでPoolsコンポーネントがないターゲットを検証
		damage := CombatDamage{Amount: 10, Source: DamageSourceWeapon}
		ctxWithInvalidTarget := &Scope{
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

		// 無効なアイテムIDの場合（存在しないエンティティID）
		useItemInvalid := UseItem{Item: ecs.Entity(9999)}
		ctx := &Scope{
			Targets: []ecs.Entity{ecs.Entity(1)},
		}
		err = useItemInvalid.Validate(world, ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "無効なアイテムエンティティです")
	})
}

func TestCombatEffects(t *testing.T) {
	t.Run("ダメージエフェクトの適用", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		// プレイヤーエンティティを作成
		player := createTestPlayerEntity(world, 100, 50)

		// ダメージエフェクトを適用
		damage := CombatDamage{Amount: 25, Source: DamageSourceWeapon}
		ctx := &Scope{Targets: []ecs.Entity{player}}

		err = damage.Validate(world, ctx)
		assert.NoError(t, err)

		err = damage.Apply(world, ctx)
		assert.NoError(t, err)

		// HPが減っていることを確認
		pools := world.Components.Pools.Get(player).(*gc.Pools)
		assert.Equal(t, 75, pools.HP.Current)
		assert.Equal(t, 100, pools.HP.Max)
	})

	t.Run("戦闘時回復エフェクトの適用", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		// ダメージを受けたプレイヤーを作成
		player := createTestPlayerEntity(world, 30, 50)

		// 戦闘時回復エフェクトを適用
		healing := CombatHealing{Amount: gc.NumeralAmount{Numeral: 40}}
		ctx := &Scope{Targets: []ecs.Entity{player}}

		err = healing.Validate(world, ctx)
		assert.NoError(t, err)

		err = healing.Apply(world, ctx)
		assert.NoError(t, err)

		// HPが回復していることを確認
		pools := world.Components.Pools.Get(player).(*gc.Pools)
		assert.Equal(t, 70, pools.HP.Current)
		assert.Equal(t, 100, pools.HP.Max)
	})

	t.Run("スタミナ消費エフェクトの適用", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		player := createTestPlayerEntity(world, 100, 50)

		// スタミナ消費エフェクトを適用
		consume := ConsumeStamina{Amount: gc.NumeralAmount{Numeral: 20}}
		ctx := &Scope{Targets: []ecs.Entity{player}}

		err = consume.Validate(world, ctx)
		assert.NoError(t, err)

		err = consume.Apply(world, ctx)
		assert.NoError(t, err)

		// SPが減っていることを確認
		pools := world.Components.Pools.Get(player).(*gc.Pools)
		assert.Equal(t, 30, pools.SP.Current)
		assert.Equal(t, 50, pools.SP.Max)
	})

	t.Run("スタミナ回復エフェクトの適用", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		player := createTestPlayerEntity(world, 100, 20)

		// スタミナ回復エフェクトを適用
		restore := RestoreStamina{Amount: gc.NumeralAmount{Numeral: 25}}
		ctx := &Scope{Targets: []ecs.Entity{player}}

		err = restore.Validate(world, ctx)
		assert.NoError(t, err)

		err = restore.Apply(world, ctx)
		assert.NoError(t, err)

		// SPが回復していることを確認
		pools := world.Components.Pools.Get(player).(*gc.Pools)
		assert.Equal(t, 45, pools.SP.Current)
		assert.Equal(t, 50, pools.SP.Max)
	})
}

func TestRecoveryEffects(t *testing.T) {
	t.Run("非戦闘時HP全回復", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		player := createTestPlayerEntity(world, 10, 50)

		fullRecovery := FullRecoveryHP{}
		ctx := &Scope{Targets: []ecs.Entity{player}}

		err = fullRecovery.Validate(world, ctx)
		assert.NoError(t, err)

		err = fullRecovery.Apply(world, ctx)
		assert.NoError(t, err)

		// HPが全回復していることを確認
		pools := world.Components.Pools.Get(player).(*gc.Pools)
		assert.Equal(t, 100, pools.HP.Current)
		assert.Equal(t, 100, pools.HP.Max)
	})

	t.Run("非戦闘時SP全回復", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		player := createTestPlayerEntity(world, 100, 5)

		fullRecovery := FullRecoverySP{}
		ctx := &Scope{Targets: []ecs.Entity{player}}

		err = fullRecovery.Validate(world, ctx)
		assert.NoError(t, err)

		err = fullRecovery.Apply(world, ctx)
		assert.NoError(t, err)

		// SPが全回復していることを確認
		pools := world.Components.Pools.Get(player).(*gc.Pools)
		assert.Equal(t, 50, pools.SP.Current)
		assert.Equal(t, 50, pools.SP.Max)
	})

	t.Run("非戦闘時HP部分回復", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		player := createTestPlayerEntity(world, 40, 50)

		recovery := RecoveryHP{Amount: gc.NumeralAmount{Numeral: 35}}
		ctx := &Scope{Targets: []ecs.Entity{player}}

		err = recovery.Validate(world, ctx)
		assert.NoError(t, err)

		err = recovery.Apply(world, ctx)
		assert.NoError(t, err)

		// HPが回復していることを確認
		pools := world.Components.Pools.Get(player).(*gc.Pools)
		assert.Equal(t, 75, pools.HP.Current)
		assert.Equal(t, 100, pools.HP.Max)
	})

	t.Run("割合回復のテスト", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		player := createTestPlayerEntity(world, 50, 25)

		// 50%回復
		recovery := RecoveryHP{Amount: gc.RatioAmount{Ratio: 0.5}}
		ctx := &Scope{Targets: []ecs.Entity{player}}

		err = recovery.Validate(world, ctx)
		assert.NoError(t, err)

		err = recovery.Apply(world, ctx)
		assert.NoError(t, err)

		// 50% (50HP)回復して100になることを確認
		pools := world.Components.Pools.Get(player).(*gc.Pools)
		assert.Equal(t, 100, pools.HP.Current)
	})
}

func TestItemEffects(t *testing.T) {
	t.Run("回復アイテムの使用", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		// 回復アイテムを作成
		healingItem := createTestHealingItem(world, 30)

		// プレイヤーを作成
		player := createTestPlayerEntity(world, 50, 50)

		useItem := UseItem{Item: healingItem}
		ctx := &Scope{Targets: []ecs.Entity{player}}

		err = useItem.Validate(world, ctx)
		assert.NoError(t, err)

		err = useItem.Apply(world, ctx)
		assert.NoError(t, err)

		// HPが回復していることを確認
		pools := world.Components.Pools.Get(player).(*gc.Pools)
		assert.Equal(t, 80, pools.HP.Current)

		// アイテムが消費されていることを確認
		assert.Nil(t, world.Components.ProvidesHealing.Get(healingItem))
	})

	t.Run("効果のないアイテムの検証", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		// 効果のないアイテムを作成（Itemコンポーネントはあるが効果なし）
		uselessItem := createTestBasicItem(world, "効果なし")
		player := createTestPlayerEntity(world, 50, 50)

		useItem := UseItem{Item: uselessItem}
		ctx := &Scope{Targets: []ecs.Entity{player}}

		err = useItem.Validate(world, ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "このアイテムには効果がありません")
	})

	t.Run("アイテム消費エフェクト", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		item := createTestBasicItem(world, "テストアイテム")

		consume := ConsumeItem{Item: item}
		ctx := &Scope{}

		err = consume.Validate(world, ctx)
		assert.NoError(t, err)

		err = consume.Apply(world, ctx)
		assert.NoError(t, err)

		// アイテムが削除されていることを確認
		assert.Nil(t, world.Components.Name.Get(item))
	})

	t.Run("アイテム生成エフェクト", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		// 生成前のアイテム数を数える（Name コンポーネントを持つエンティティ）
		initialCount := 0
		world.Manager.Join(world.Components.Name).Visit(ecs.Visit(func(entity ecs.Entity) {
			initialCount++
		}))

		create := CreateItem{ItemType: "回復薬", Quantity: 3}
		ctx := &Scope{}

		err = create.Validate(world, ctx)
		assert.NoError(t, err)

		err = create.Apply(world, ctx)
		assert.NoError(t, err)

		// アイテムが生成されたことを確認（3個増えている）
		finalCount := 0
		world.Manager.Join(world.Components.Name).Visit(ecs.Visit(func(entity ecs.Entity) {
			finalCount++
		}))
		assert.Equal(t, initialCount+3, finalCount)
	})
}

func TestMovementEffects(t *testing.T) {
	t.Run("次階層へのワープ", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		warp := MovementWarpNext{}
		ctx := &Scope{}

		err = warp.Validate(world, ctx)
		assert.NoError(t, err)

		err = warp.Apply(world, ctx)
		assert.NoError(t, err)

		// StateEventが設定されていることを確認（詳細は別パッケージで検証）
		assert.NotNil(t, world.Resources.Game)
	})

	t.Run("脱出ワープ", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		escape := MovementWarpEscape{}
		ctx := &Scope{}

		err = escape.Validate(world, ctx)
		assert.NoError(t, err)

		err = escape.Apply(world, ctx)
		assert.NoError(t, err)

		assert.NotNil(t, world.Resources.Game)
	})

	t.Run("特定階層へのワープ - 検証エラー", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		// 無効な階層（0以下）
		warpInvalid := MovementWarpToFloor{Floor: 0}
		ctx := &Scope{}

		err = warpInvalid.Validate(world, ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "階層は1以上である必要があります")
	})

	t.Run("特定階層へのワープ", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		// 階層3へのワープ
		warp := MovementWarpToFloor{Floor: 3}
		ctx := &Scope{}

		err = warp.Validate(world, ctx)
		assert.NoError(t, err)

		err = warp.Apply(world, ctx)
		assert.NoError(t, err)

		// 次階層イベントが設定されることを確認
		gameResources := world.Resources.Game.(*resources.Game)
		assert.Equal(t, resources.StateEventWarpNext, gameResources.StateEvent)
	})
}

func TestProcessor(t *testing.T) {
	t.Run("プロセッサーでのエフェクト実行", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		// プレイヤーエンティティを作成
		player := createTestPlayerEntity(world, 50, 50)

		processor := NewProcessor()

		// 複数のエフェクトをキューに追加
		healing := RecoveryHP{Amount: gc.NumeralAmount{Numeral: 20}}
		damage := CombatDamage{Amount: 10, Source: DamageSourceWeapon}

		processor.AddEffect(healing, nil, player)
		assert.Equal(t, 1, processor.QueueSize())

		processor.AddEffect(damage, nil, player)
		assert.Equal(t, 2, processor.QueueSize())

		// プロセッサーでエフェクトを実行
		err = processor.Execute(world)
		assert.NoError(t, err)
		assert.True(t, processor.IsEmpty())

		// エフェクトが順次適用されたことを確認（回復→ダメージ）
		pools := world.Components.Pools.Get(player).(*gc.Pools)
		// 50 + 20 - 10 = 60
		assert.Equal(t, 60, pools.HP.Current)
	})

	t.Run("検証失敗時のエラーハンドリング", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		processor := NewProcessor()

		// 無効なエフェクト（負のダメージ）を追加
		invalidDamage := CombatDamage{Amount: -10, Source: DamageSourceWeapon}
		processor.AddEffect(invalidDamage, nil, ecs.Entity(1))

		// 実行時に検証エラーが発生することを確認（Apply内でValidateが呼ばれるためエフェクト実行失敗になる）
		err = processor.Execute(world)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "エフェクト実行失敗")
	})

	t.Run("プロセッサーのクリア機能", func(t *testing.T) {
		processor := NewProcessor()

		// 複数のエフェクトを追加
		processor.AddEffect(RecoveryHP{Amount: gc.NumeralAmount{Numeral: 10}}, nil, ecs.Entity(1))
		processor.AddEffect(CombatDamage{Amount: 5, Source: DamageSourceWeapon}, nil, ecs.Entity(1))

		assert.Equal(t, 2, processor.QueueSize())
		assert.False(t, processor.IsEmpty())

		// クリア
		processor.Clear()
		assert.Equal(t, 0, processor.QueueSize())
		assert.True(t, processor.IsEmpty())
	})
}

func TestValidationErrors(t *testing.T) {
	t.Run("エフェクトパラメータの検証", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		player := createTestPlayerEntity(world, 100, 50)
		ctx := &Scope{Targets: []ecs.Entity{player}}

		// 回復量がnil
		healingInvalid := CombatHealing{Amount: nil}
		err = healingInvalid.Validate(world, ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "回復量が指定されていません")

		// スタミナ消費量がnil
		consumeInvalid := ConsumeStamina{Amount: nil}
		err = consumeInvalid.Validate(world, ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "スタミナ消費量が指定されていません")

		// スタミナ回復量がnil
		restoreInvalid := RestoreStamina{Amount: nil}
		err = restoreInvalid.Validate(world, ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "スタミナ回復量が指定されていません")
	})

	t.Run("ターゲットなしの検証", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		emptyCtx := &Scope{Targets: []ecs.Entity{}}

		// 各エフェクトでターゲットなしエラーをテスト
		effects := []Effect{
			CombatHealing{Amount: gc.NumeralAmount{Numeral: 10}},
			ConsumeStamina{Amount: gc.NumeralAmount{Numeral: 10}},
			RestoreStamina{Amount: gc.NumeralAmount{Numeral: 10}},
			FullRecoveryHP{},
			FullRecoverySP{},
			RecoveryHP{Amount: gc.NumeralAmount{Numeral: 10}},
			RecoverySP{Amount: gc.NumeralAmount{Numeral: 10}},
		}

		for _, effect := range effects {
			err = effect.Validate(world, emptyCtx)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "対象が指定されていません")
		}
	})
}
