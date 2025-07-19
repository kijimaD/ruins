package effects

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/engine/entities"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/stretchr/testify/assert"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// MockLogger はテスト用のゲームログ出力先（GameLogAppenderインターフェースを実装）
type MockLogger struct {
	Entries []string
}

func (m *MockLogger) Append(entry string) {
	m.Entries = append(m.Entries, entry)
}

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

		healing := Healing{Amount: gc.NumeralAmount{Numeral: 30}}
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
		scope := &Scope{Targets: []ecs.Entity{ecs.Entity(1)}}
		err = damage.Validate(world, scope)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ダメージは0以上")

		// ターゲットなしは無効
		validDamage := CombatDamage{Amount: 10, Source: DamageSourceWeapon}
		emptyScope := &Scope{Targets: []ecs.Entity{}}
		err = validDamage.Validate(world, emptyScope)
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
		healing := Healing{Amount: gc.NumeralAmount{Numeral: 30}}
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
		scope := &Scope{
			Targets: []ecs.Entity{ecs.Entity(1)},
		}
		err = useItemInvalid.Validate(world, scope)
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
		scope := &Scope{Targets: []ecs.Entity{player}}

		err = damage.Validate(world, scope)
		assert.NoError(t, err)

		err = damage.Apply(world, scope)
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
		healing := Healing{Amount: gc.NumeralAmount{Numeral: 40}}
		scope := &Scope{Targets: []ecs.Entity{player}}

		err = healing.Validate(world, scope)
		assert.NoError(t, err)

		err = healing.Apply(world, scope)
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
		scope := &Scope{Targets: []ecs.Entity{player}}

		err = consume.Validate(world, scope)
		assert.NoError(t, err)

		err = consume.Apply(world, scope)
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
		scope := &Scope{Targets: []ecs.Entity{player}}

		err = restore.Validate(world, scope)
		assert.NoError(t, err)

		err = restore.Apply(world, scope)
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
		scope := &Scope{Targets: []ecs.Entity{player}}

		err = fullRecovery.Validate(world, scope)
		assert.NoError(t, err)

		err = fullRecovery.Apply(world, scope)
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
		scope := &Scope{Targets: []ecs.Entity{player}}

		err = fullRecovery.Validate(world, scope)
		assert.NoError(t, err)

		err = fullRecovery.Apply(world, scope)
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

		recovery := Healing{Amount: gc.NumeralAmount{Numeral: 35}}
		scope := &Scope{Targets: []ecs.Entity{player}}

		err = recovery.Validate(world, scope)
		assert.NoError(t, err)

		err = recovery.Apply(world, scope)
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
		recovery := Healing{Amount: gc.RatioAmount{Ratio: 0.5}}
		scope := &Scope{Targets: []ecs.Entity{player}}

		err = recovery.Validate(world, scope)
		assert.NoError(t, err)

		err = recovery.Apply(world, scope)
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
		scope := &Scope{Targets: []ecs.Entity{player}}

		err = useItem.Validate(world, scope)
		assert.NoError(t, err)

		err = useItem.Apply(world, scope)
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
		scope := &Scope{Targets: []ecs.Entity{player}}

		err = useItem.Validate(world, scope)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "このアイテムには効果がありません")
	})

	t.Run("アイテム消費エフェクト", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		item := createTestBasicItem(world, "テストアイテム")

		consume := ConsumeItem{Item: item}
		scope := &Scope{}

		err = consume.Validate(world, scope)
		assert.NoError(t, err)

		err = consume.Apply(world, scope)
		assert.NoError(t, err)

		// アイテムが削除されていることを確認
		assert.Nil(t, world.Components.Name.Get(item))
	})

	t.Run("アイテム生成エフェクト", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		// 生成前のアイテム数を数える（Name コンポーネントを持つエンティティ）
		initialCount := 0
		world.Manager.Join(world.Components.Name).Visit(ecs.Visit(func(_ ecs.Entity) {
			initialCount++
		}))

		create := CreateItem{ItemType: "回復薬", Quantity: 3}
		scope := &Scope{}

		err = create.Validate(world, scope)
		assert.NoError(t, err)

		err = create.Apply(world, scope)
		assert.NoError(t, err)

		// アイテムが生成されたことを確認（3個増えている）
		finalCount := 0
		world.Manager.Join(world.Components.Name).Visit(ecs.Visit(func(_ ecs.Entity) {
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
		scope := &Scope{}

		err = warp.Validate(world, scope)
		assert.NoError(t, err)

		err = warp.Apply(world, scope)
		assert.NoError(t, err)

		// StateEventが設定されていることを確認（詳細は別パッケージで検証）
		assert.NotNil(t, world.Resources.Game)
	})

	t.Run("脱出ワープ", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		escape := MovementWarpEscape{}
		scope := &Scope{}

		err = escape.Validate(world, scope)
		assert.NoError(t, err)

		err = escape.Apply(world, scope)
		assert.NoError(t, err)

		assert.NotNil(t, world.Resources.Game)
	})

	t.Run("特定階層へのワープ - 検証エラー", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		// 無効な階層（0以下）
		warpInvalid := MovementWarpToFloor{Floor: 0}
		scope := &Scope{}

		err = warpInvalid.Validate(world, scope)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "階層は1以上である必要があります")
	})

	t.Run("特定階層へのワープ", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		// 階層3へのワープ
		warp := MovementWarpToFloor{Floor: 3}
		scope := &Scope{}

		err = warp.Validate(world, scope)
		assert.NoError(t, err)

		err = warp.Apply(world, scope)
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
		healing := Healing{Amount: gc.NumeralAmount{Numeral: 20}}
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
		processor.AddEffect(Healing{Amount: gc.NumeralAmount{Numeral: 10}}, nil, ecs.Entity(1))
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
		scope := &Scope{Targets: []ecs.Entity{player}}

		// 回復量がnil
		healingInvalid := Healing{Amount: nil}
		err = healingInvalid.Validate(world, scope)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "回復量が指定されていません")

		// スタミナ消費量がnil
		consumeInvalid := ConsumeStamina{Amount: nil}
		err = consumeInvalid.Validate(world, scope)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "スタミナ消費量が指定されていません")

		// スタミナ回復量がnil
		restoreInvalid := RestoreStamina{Amount: nil}
		err = restoreInvalid.Validate(world, scope)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "スタミナ回復量が指定されていません")
	})

	t.Run("ターゲットなしの検証", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		emptyScope := &Scope{Targets: []ecs.Entity{}}

		// 各エフェクトでターゲットなしエラーをテスト
		effects := []Effect{
			Healing{Amount: gc.NumeralAmount{Numeral: 10}},
			ConsumeStamina{Amount: gc.NumeralAmount{Numeral: 10}},
			RestoreStamina{Amount: gc.NumeralAmount{Numeral: 10}},
			FullRecoveryHP{},
			FullRecoverySP{},
		}

		for _, effect := range effects {
			err = effect.Validate(world, emptyScope)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "対象が指定されていません")
		}
	})
}

func TestLoggerIntegration(t *testing.T) {
	t.Run("戦闘時回復エフェクトのログ出力", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		player := createTestPlayerEntity(world, 50, 50)
		mockLogger := &MockLogger{}

		// 戦闘時回復エフェクトをLoggerとともに実行
		healing := Healing{Amount: gc.NumeralAmount{Numeral: 30}}
		scope := &Scope{
			Targets: []ecs.Entity{player},
			Logger:  mockLogger,
		}

		err = healing.Apply(world, scope)
		assert.NoError(t, err)

		// ログが記録されていることを確認
		assert.Len(t, mockLogger.Entries, 1)
		assert.Contains(t, mockLogger.Entries[0], "が30回復。")
	})

	t.Run("ダメージエフェクトのログ出力", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		player := createTestPlayerEntity(world, 100, 50)
		mockLogger := &MockLogger{}

		// ダメージエフェクトをLoggerとともに実行
		damage := CombatDamage{Amount: 25, Source: DamageSourceWeapon}
		scope := &Scope{
			Targets: []ecs.Entity{player},
			Logger:  mockLogger,
		}

		err = damage.Apply(world, scope)
		assert.NoError(t, err)

		// ダメージログが記録されていることを確認
		assert.Len(t, mockLogger.Entries, 1)
		assert.Contains(t, mockLogger.Entries[0], "に25のダメージ。")
	})

	t.Run("Logger無しの場合はログ出力なし", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		player := createTestPlayerEntity(world, 50, 50)

		// Logger無しでエフェクト実行
		healing := Healing{Amount: gc.NumeralAmount{Numeral: 20}}
		scope := &Scope{
			Targets: []ecs.Entity{player},
			Logger:  nil, // Logger無し
		}

		err = healing.Apply(world, scope)
		assert.NoError(t, err)

		// HPは回復しているが、ログ出力はない
		pools := world.Components.Pools.Get(player).(*gc.Pools)
		assert.Equal(t, 70, pools.HP.Current)
	})

	t.Run("ProcessorのAddEffectWithLoggerメソッド", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		player := createTestPlayerEntity(world, 50, 50)
		mockLogger := &MockLogger{}
		processor := NewProcessor()

		// AddEffectWithLoggerを使用してエフェクトを追加
		healing := Healing{Amount: gc.NumeralAmount{Numeral: 25}}
		processor.AddEffectWithLogger(healing, nil, mockLogger, player)

		err = processor.Execute(world)
		assert.NoError(t, err)

		// ログが記録されていることを確認
		assert.Len(t, mockLogger.Entries, 1)
		assert.Contains(t, mockLogger.Entries[0], "が25回復。")
	})

	t.Run("統合されたHealingエフェクトの確認", func(t *testing.T) {
		// 統合されたHealingエフェクトを使用
		healing := Healing{Amount: gc.NumeralAmount{Numeral: 30}}
		assert.Equal(t, "Healing({30})", healing.String())

		// 同じエフェクトが戦闘・非戦闘で共用可能
		healingForCombat := Healing{Amount: gc.NumeralAmount{Numeral: 25}}
		assert.Equal(t, "Healing({25})", healingForCombat.String())
	})

	t.Run("gamelog.BattleLogを使った戦闘ダメージログ", func(t *testing.T) {
		world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
		assert.NoError(t, err)

		// BattleLogをクリア
		gamelog.BattleLog.Flush()

		player := createTestPlayerEntity(world, 100, 50)
		processor := NewProcessor()

		// battle_command.goと同じ方法でダメージエフェクトを実行
		damageEffect := CombatDamage{Amount: 25, Source: DamageSourceWeapon}
		processor.AddEffectWithLogger(damageEffect, nil, &gamelog.BattleLog, player)

		err = processor.Execute(world)
		assert.NoError(t, err)

		// BattleLogからログを取得して確認
		logs := gamelog.BattleLog.Get()
		assert.Len(t, logs, 1)
		assert.Contains(t, logs[0], "に25のダメージ。")
	})
}
