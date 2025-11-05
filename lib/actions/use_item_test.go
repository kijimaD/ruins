package actions

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUseItemActivity_applyNutrition(t *testing.T) {
	t.Parallel()

	t.Run("満腹度が正常に増加する", func(t *testing.T) {
		t.Parallel()

		world := testutil.InitTestWorld(t)
		actor := world.Manager.NewEntity()

		// Hungerコンポーネントを追加（初期値は1000）
		hunger := gc.NewHunger()
		hunger.Current = 500 // 半分の満腹度
		actor.AddComponent(world.Components.Hunger, hunger)

		item := world.Manager.NewEntity()
		activity := &UseItemActivity{}
		act := NewActivity(activity, actor, 1)

		// 200の満腹度回復
		err := activity.applyNutrition(act, world, 200, item)
		require.NoError(t, err)

		// 満腹度が500 + 200 = 700になっているはず
		hungerComp := world.Components.Hunger.Get(actor)
		require.NotNil(t, hungerComp)
		updatedHunger := hungerComp.(*gc.Hunger)
		assert.Equal(t, 700, updatedHunger.Current, "満腹度が正しく増加していない")
	})

	t.Run("満腹度が上限を超えない", func(t *testing.T) {
		t.Parallel()

		world := testutil.InitTestWorld(t)
		actor := world.Manager.NewEntity()

		hunger := gc.NewHunger()
		hunger.Current = 950 // ほぼ満腹
		actor.AddComponent(world.Components.Hunger, hunger)

		item := world.Manager.NewEntity()
		activity := &UseItemActivity{}
		act := NewActivity(activity, actor, 1)

		// 200の満腹度回復（上限を超える）
		err := activity.applyNutrition(act, world, 200, item)
		require.NoError(t, err)

		hungerComp := world.Components.Hunger.Get(actor)
		require.NotNil(t, hungerComp)
		updatedHunger := hungerComp.(*gc.Hunger)
		assert.Equal(t, gc.DefaultMaxHunger, updatedHunger.Current, "満腹度が上限を超えている")
	})

	t.Run("満腹状態になった場合", func(t *testing.T) {
		t.Parallel()

		world := testutil.InitTestWorld(t)
		actor := world.Manager.NewEntity()
		actor.AddComponent(world.Components.Player, &gc.Player{})

		hunger := gc.NewHunger()
		hunger.Current = 850 // 85%
		actor.AddComponent(world.Components.Hunger, hunger)

		item := world.Manager.NewEntity()
		item.AddComponent(world.Components.Name, &gc.Name{Name: "パン"})

		activity := &UseItemActivity{}
		act := NewActivity(activity, actor, 1)

		// 100の満腹度回復で90%以上になる
		err := activity.applyNutrition(act, world, 100, item)
		require.NoError(t, err)

		hungerComp := world.Components.Hunger.Get(actor)
		require.NotNil(t, hungerComp)
		updatedHunger := hungerComp.(*gc.Hunger)
		assert.Equal(t, 950, updatedHunger.Current)
		assert.Equal(t, gc.HungerSatiated, updatedHunger.GetLevel(), "満腹状態になっているはず")
	})

	t.Run("Hungerコンポーネントがない場合は何もしない", func(t *testing.T) {
		t.Parallel()

		world := testutil.InitTestWorld(t)
		actor := world.Manager.NewEntity()
		// Hungerコンポーネントを追加しない

		item := world.Manager.NewEntity()
		activity := &UseItemActivity{}
		act := NewActivity(activity, actor, 1)

		// エラーにならずに完了する
		err := activity.applyNutrition(act, world, 200, item)
		assert.NoError(t, err)
	})

	t.Run("飢餓状態から回復する", func(t *testing.T) {
		t.Parallel()

		world := testutil.InitTestWorld(t)
		actor := world.Manager.NewEntity()

		hunger := gc.NewHunger()
		hunger.Current = 100 // 10% - 飢餓状態
		actor.AddComponent(world.Components.Hunger, hunger)

		item := world.Manager.NewEntity()
		activity := &UseItemActivity{}
		act := NewActivity(activity, actor, 1)

		assert.Equal(t, gc.HungerStarving, hunger.GetLevel(), "初期状態は飢餓状態")

		// 600の満腹度回復で70%になる
		err := activity.applyNutrition(act, world, 600, item)
		require.NoError(t, err)

		hungerComp := world.Components.Hunger.Get(actor)
		require.NotNil(t, hungerComp)
		updatedHunger := hungerComp.(*gc.Hunger)
		assert.Equal(t, 700, updatedHunger.Current)
		assert.Equal(t, gc.HungerNormal, updatedHunger.GetLevel(), "普通状態に回復しているはず")
	})
}
