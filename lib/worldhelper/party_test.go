package worldhelper

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewParty(t *testing.T) {
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// 味方キャラクターを作成
	ally1 := world.Manager.NewEntity()
	ally1.AddComponent(world.Components.FactionAlly, &gc.FactionAlly)
	ally1.AddComponent(world.Components.Pools, &gc.Pools{
		HP:    gc.Pool{Current: 100, Max: 100},
		Level: 1,
	})
	ally1.AddComponent(world.Components.Attributes, &gc.Attributes{
		Vitality: gc.Attribute{Base: 10, Total: 10},
	})

	ally2 := world.Manager.NewEntity()
	ally2.AddComponent(world.Components.FactionAlly, &gc.FactionAlly)
	ally2.AddComponent(world.Components.Pools, &gc.Pools{
		HP:    gc.Pool{Current: 0, Max: 50}, // 死亡状態
		Level: 1,
	})
	ally2.AddComponent(world.Components.Attributes, &gc.Attributes{
		Vitality: gc.Attribute{Base: 8, Total: 8},
	})

	// パーティを作成
	party, err := NewParty(world, gc.FactionAlly)
	assert.NoError(t, err, "パーティ作成でエラーが発生してはいけない")

	// パーティの状態を検証
	assert.Len(t, party.members, 2, "パーティのメンバー数が正しくない")
	assert.Len(t, party.lives, 2, "パーティの生存状況配列の長さが正しくない")
	assert.Equal(t, 1, party.LivesLen(), "生存メンバー数が正しくない")

	// 現在選択されているメンバーが生存していることを確認
	currentMember := party.Value()
	assert.NotNil(t, currentMember, "現在のメンバーがnilであってはいけない")
	currentPools := world.Components.Pools.Get(*currentMember).(*gc.Pools)
	assert.Greater(t, currentPools.HP.Current, 0, "現在のメンバーは生存していなければならない")

	// クリーンアップ
	world.Manager.DeleteEntity(ally1)
	world.Manager.DeleteEntity(ally2)
}

func TestNewByEntity(t *testing.T) {
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// 味方エンティティを作成
	ally := world.Manager.NewEntity()
	ally.AddComponent(world.Components.FactionAlly, &gc.FactionAlly)
	ally.AddComponent(world.Components.Pools, &gc.Pools{
		HP:    gc.Pool{Current: 100, Max: 100},
		Level: 1,
	})
	ally.AddComponent(world.Components.Attributes, &gc.Attributes{
		Vitality: gc.Attribute{Base: 10, Total: 10},
	})

	// エンティティからパーティを作成
	party, err := NewByEntity(world, ally)
	assert.NoError(t, err, "エンティティからのパーティ作成でエラーが発生してはいけない")
	assert.NotNil(t, party.Value(), "パーティが正しく作成されていない")

	// 無効なエンティティでテスト
	invalidEntity := world.Manager.NewEntity()
	_, err = NewByEntity(world, invalidEntity)
	assert.Error(t, err, "無効なエンティティでエラーが発生するべき")

	// クリーンアップ
	world.Manager.DeleteEntity(ally)
	world.Manager.DeleteEntity(invalidEntity)
}

func TestPartyNavigation(t *testing.T) {
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// 3人の味方キャラクターを作成（1人は死亡状態）
	ally1 := world.Manager.NewEntity()
	ally1.AddComponent(world.Components.FactionAlly, &gc.FactionAlly)
	ally1.AddComponent(world.Components.Pools, &gc.Pools{
		HP:    gc.Pool{Current: 100, Max: 100},
		Level: 1,
	})
	ally1.AddComponent(world.Components.Attributes, &gc.Attributes{
		Vitality: gc.Attribute{Base: 10, Total: 10},
	})

	ally2 := world.Manager.NewEntity()
	ally2.AddComponent(world.Components.FactionAlly, &gc.FactionAlly)
	ally2.AddComponent(world.Components.Pools, &gc.Pools{
		HP:    gc.Pool{Current: 0, Max: 50}, // 死亡状態
		Level: 1,
	})
	ally2.AddComponent(world.Components.Attributes, &gc.Attributes{
		Vitality: gc.Attribute{Base: 8, Total: 8},
	})

	ally3 := world.Manager.NewEntity()
	ally3.AddComponent(world.Components.FactionAlly, &gc.FactionAlly)
	ally3.AddComponent(world.Components.Pools, &gc.Pools{
		HP:    gc.Pool{Current: 75, Max: 75},
		Level: 1,
	})
	ally3.AddComponent(world.Components.Attributes, &gc.Attributes{
		Vitality: gc.Attribute{Base: 9, Total: 9},
	})

	// パーティを作成
	party, err := NewParty(world, gc.FactionAlly)
	assert.NoError(t, err)

	// 初期状態の確認
	currentMember := party.Value()
	assert.NotNil(t, currentMember, "現在のメンバーがnilであってはいけない")

	// 次のメンバーに移動
	err = party.Next()
	assert.NoError(t, err, "次のメンバーへの移動でエラーが発生してはいけない")
	nextMember := party.Value()
	assert.NotEqual(t, currentMember, nextMember, "次のメンバーは現在のメンバーと異なるべき")

	// 生存しているメンバーのみがナビゲーションの対象であることを確認
	nextPools := world.Components.Pools.Get(*nextMember).(*gc.Pools)
	assert.Greater(t, nextPools.HP.Current, 0, "ナビゲーション対象は生存メンバーでなければならない")

	// 前のメンバーに戻る
	err = party.Prev()
	assert.NoError(t, err, "前のメンバーへの移動でエラーが発生してはいけない")
	prevMember := party.Value()
	assert.Equal(t, currentMember, prevMember, "前のメンバーは元のメンバーと同じであるべき")

	// クリーンアップ
	world.Manager.DeleteEntity(ally1)
	world.Manager.DeleteEntity(ally2)
	world.Manager.DeleteEntity(ally3)
}
