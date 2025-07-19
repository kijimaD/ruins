package worldhelper

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAmount(t *testing.T) {
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// テスト用素材エンティティを作成
	materialEntity := world.Manager.NewEntity()
	materialEntity.AddComponent(world.Components.Material, &gc.Material{Amount: 10})
	materialEntity.AddComponent(world.Components.ItemLocationInBackpack, &gc.ItemLocationInBackpack)
	materialEntity.AddComponent(world.Components.Name, &gc.Name{Name: "鉄"})

	// 素材の数量を取得
	amount := GetAmount("鉄", world)
	assert.Equal(t, 10, amount, "素材の数量が正しく取得できない")

	// 存在しない素材の数量を取得
	amount = GetAmount("存在しない素材", world)
	assert.Equal(t, 0, amount, "存在しない素材の数量は0であるべき")

	// クリーンアップ
	world.Manager.DeleteEntity(materialEntity)
}

func TestPlusMinusAmount(t *testing.T) {
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// テスト用素材エンティティを作成
	materialEntity := world.Manager.NewEntity()
	materialEntity.AddComponent(world.Components.Material, &gc.Material{Amount: 10})
	materialEntity.AddComponent(world.Components.ItemLocationInBackpack, &gc.ItemLocationInBackpack)
	materialEntity.AddComponent(world.Components.Name, &gc.Name{Name: "鉄"})

	// 数量を増加
	PlusAmount("鉄", 5, world)
	amount := GetAmount("鉄", world)
	assert.Equal(t, 15, amount, "数量増加が正しく動作しない")

	// 数量を減少
	MinusAmount("鉄", 3, world)
	amount = GetAmount("鉄", world)
	assert.Equal(t, 12, amount, "数量減少が正しく動作しない")

	// 上限テスト（999を超えない）
	PlusAmount("鉄", 1000, world)
	amount = GetAmount("鉄", world)
	assert.Equal(t, 999, amount, "数量は999を超えてはいけない")

	// 下限テスト（0未満にならない）
	MinusAmount("鉄", 1500, world)
	amount = GetAmount("鉄", world)
	assert.Equal(t, 0, amount, "数量は0未満になってはいけない")

	// クリーンアップ
	world.Manager.DeleteEntity(materialEntity)
}
