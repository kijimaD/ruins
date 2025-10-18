package states

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEquipMenuWearSortIntegration(t *testing.T) {
	t.Parallel()
	world := testutil.InitTestWorld(t)

	state := &EquipMenuState{}

	// テスト用防具エンティティを作成（名前順ではない順序で作成）
	wearable1 := world.Manager.NewEntity()
	wearable1.AddComponent(world.Components.Name, &gc.Name{Name: "Shield"})
	wearable1.AddComponent(world.Components.Item, &gc.Item{})
	wearable1.AddComponent(world.Components.ItemLocationInBackpack, gc.LocationInBackpack{})
	wearable1.AddComponent(world.Components.Wearable, &gc.Wearable{})

	wearable2 := world.Manager.NewEntity()
	wearable2.AddComponent(world.Components.Name, &gc.Name{Name: "Armor"})
	wearable2.AddComponent(world.Components.Item, &gc.Item{})
	wearable2.AddComponent(world.Components.ItemLocationInBackpack, gc.LocationInBackpack{})
	wearable2.AddComponent(world.Components.Wearable, &gc.Wearable{})

	wearable3 := world.Manager.NewEntity()
	wearable3.AddComponent(world.Components.Name, &gc.Name{Name: "Helmet"})
	wearable3.AddComponent(world.Components.Item, &gc.Item{})
	wearable3.AddComponent(world.Components.ItemLocationInBackpack, gc.LocationInBackpack{})
	wearable3.AddComponent(world.Components.Wearable, &gc.Wearable{})

	// queryMenuWearのテスト
	wearables := state.queryMenuWear(world)
	require.Len(t, wearables, 3, "防具が3つ見つかるべき")

	// ソート順を確認（名前順）
	name1 := world.Components.Name.Get(wearables[0]).(*gc.Name)
	name2 := world.Components.Name.Get(wearables[1]).(*gc.Name)
	name3 := world.Components.Name.Get(wearables[2]).(*gc.Name)

	assert.Equal(t, "Armor", name1.Name, "1番目の防具名が正しくない")
	assert.Equal(t, "Helmet", name2.Name, "2番目の防具名が正しくない")
	assert.Equal(t, "Shield", name3.Name, "3番目の防具名が正しくない")

	// クリーンアップ
	world.Manager.DeleteEntity(wearable1)
	world.Manager.DeleteEntity(wearable2)
	world.Manager.DeleteEntity(wearable3)
}

func TestEquipMenuWeaponSortIntegration(t *testing.T) {
	t.Parallel()
	world := testutil.InitTestWorld(t)

	state := &EquipMenuState{}

	// テスト用武器エンティティを作成（名前順ではない順序で作成）
	weapon1 := world.Manager.NewEntity()
	weapon1.AddComponent(world.Components.Name, &gc.Name{Name: "Thunder Weapon"})
	weapon1.AddComponent(world.Components.Item, &gc.Item{})
	weapon1.AddComponent(world.Components.ItemLocationInBackpack, gc.LocationInBackpack{})
	weapon1.AddComponent(world.Components.Weapon, &gc.Weapon{})

	weapon2 := world.Manager.NewEntity()
	weapon2.AddComponent(world.Components.Name, &gc.Name{Name: "Fire Weapon"})
	weapon2.AddComponent(world.Components.Item, &gc.Item{})
	weapon2.AddComponent(world.Components.ItemLocationInBackpack, gc.LocationInBackpack{})
	weapon2.AddComponent(world.Components.Weapon, &gc.Weapon{})

	weapon3 := world.Manager.NewEntity()
	weapon3.AddComponent(world.Components.Name, &gc.Name{Name: "Ice Weapon"})
	weapon3.AddComponent(world.Components.Item, &gc.Item{})
	weapon3.AddComponent(world.Components.ItemLocationInBackpack, gc.LocationInBackpack{})
	weapon3.AddComponent(world.Components.Weapon, &gc.Weapon{})

	// queryMenuWeaponのテスト
	weapons := state.queryMenuWeapon(world)
	require.Len(t, weapons, 3, "武器が3つ見つかるべき")

	// ソート順を確認（名前順）
	name1 := world.Components.Name.Get(weapons[0]).(*gc.Name)
	name2 := world.Components.Name.Get(weapons[1]).(*gc.Name)
	name3 := world.Components.Name.Get(weapons[2]).(*gc.Name)

	assert.Equal(t, "Fire Weapon", name1.Name, "1番目の武器名が正しくない")
	assert.Equal(t, "Ice Weapon", name2.Name, "2番目の武器名が正しくない")
	assert.Equal(t, "Thunder Weapon", name3.Name, "3番目の武器名が正しくない")

	// クリーンアップ
	world.Manager.DeleteEntity(weapon1)
	world.Manager.DeleteEntity(weapon2)
	world.Manager.DeleteEntity(weapon3)
}
