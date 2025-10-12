package systems

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/testutil"
	"github.com/stretchr/testify/assert"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestDeadCleanupSystem(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// テスト用エンティティを作成

	// 1. 通常の敵（AI）エンティティ - 削除されるべき
	enemy := world.Manager.NewEntity()
	enemy.AddComponent(world.Components.Name, &gc.Name{Name: "テスト敵"})
	enemy.AddComponent(world.Components.AIMoveFSM, &gc.AIMoveFSM{})
	enemy.AddComponent(world.Components.Dead, &gc.Dead{})

	// 2. プレイヤーエンティティ - 削除されないべき
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Name, &gc.Name{Name: "プレイヤー"})
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.Dead, &gc.Dead{})

	// 3. その他のDeadエンティティ - 削除されるべき
	otherDead := world.Manager.NewEntity()
	otherDead.AddComponent(world.Components.Name, &gc.Name{Name: "その他"})
	otherDead.AddComponent(world.Components.Dead, &gc.Dead{})

	// 4. 生きているエンティティ - 削除されないべき
	alive := world.Manager.NewEntity()
	alive.AddComponent(world.Components.Name, &gc.Name{Name: "生きている敵"})
	alive.AddComponent(world.Components.AIMoveFSM, &gc.AIMoveFSM{})

	// DeadCleanupSystemを実行
	DeadCleanupSystem(world)

	// 結果を検証

	// 敵エンティティは削除されているべき（Nameコンポーネントも削除される）
	assert.False(t, enemy.HasComponent(world.Components.Name), "敵エンティティは削除されるべき")

	// プレイヤーエンティティは削除されていないべき
	assert.True(t, player.HasComponent(world.Components.Name), "プレイヤーエンティティは削除されないべき")
	assert.True(t, player.HasComponent(world.Components.Dead), "プレイヤーのDeadコンポーネントは残るべき")

	// その他のDeadエンティティは削除されているべき
	assert.False(t, otherDead.HasComponent(world.Components.Name), "その他のDeadエンティティは削除されるべき")

	// 生きているエンティティは削除されていないべき
	assert.True(t, alive.HasComponent(world.Components.Name), "生きているエンティティは削除されないべき")
	assert.False(t, alive.HasComponent(world.Components.Dead), "生きているエンティティにDeadコンポーネントはないべき")

	// クリーンアップ
	world.Manager.DeleteEntity(player)
	world.Manager.DeleteEntity(alive)
}

func TestDeadCleanupSystem_NoDeadEntities(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// Deadエンティティが存在しない状態でテスト
	alive1 := world.Manager.NewEntity()
	alive1.AddComponent(world.Components.Name, &gc.Name{Name: "生きている1"})

	alive2 := world.Manager.NewEntity()
	alive2.AddComponent(world.Components.Name, &gc.Name{Name: "生きている2"})
	alive2.AddComponent(world.Components.AIMoveFSM, &gc.AIMoveFSM{})

	// DeadCleanupSystemを実行
	DeadCleanupSystem(world)

	// すべてのエンティティが残っているべき
	assert.True(t, alive1.HasComponent(world.Components.Name), "生きているエンティティ1は残るべき")
	assert.True(t, alive2.HasComponent(world.Components.Name), "生きているエンティティ2は残るべき")

	// クリーンアップ
	world.Manager.DeleteEntity(alive1)
	world.Manager.DeleteEntity(alive2)
}

func TestDeadCleanupSystem_EmptyWorld(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// エンティティが存在しない状態でテスト
	// パニックやエラーが発生しないことを確認
	DeadCleanupSystem(world)

	// エンティティ数が0であることを確認
	count := 0
	world.Manager.Join().Visit(ecs.Visit(func(_ ecs.Entity) {
		count++
	}))
	assert.Equal(t, 0, count, "空のworldではエンティティ数は0であるべき")
}
