package worldhelper

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/testutil"
	"github.com/stretchr/testify/assert"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestQueryOwnedStackable(t *testing.T) {
	t.Parallel()
	world := testutil.InitTestWorld(t)

	// テスト用スタック可能アイテムエンティティを作成
	stackableEntity := world.Manager.NewEntity()
	stackableEntity.AddComponent(world.Components.Stackable, &gc.Stackable{Count: 5})
	stackableEntity.AddComponent(world.Components.ItemLocationInBackpack, &gc.ItemLocationInBackpack)
	stackableEntity.AddComponent(world.Components.Name, &gc.Name{Name: "テストスタック可能アイテム"})

	// スタック不可アイテムを作成（除外されることを確認）
	nonStackableEntity := world.Manager.NewEntity()
	nonStackableEntity.AddComponent(world.Components.ItemLocationInBackpack, &gc.ItemLocationInBackpack)
	nonStackableEntity.AddComponent(world.Components.Name, &gc.Name{Name: "テストアイテム"})

	// クエリを実行
	var foundEntities []ecs.Entity
	QueryOwnedStackable(world, func(entity ecs.Entity) {
		foundEntities = append(foundEntities, entity)
	})

	// 結果を検証
	assert.Len(t, foundEntities, 1, "スタック可能アイテムが1つだけ見つかるべき")
	assert.Equal(t, stackableEntity, foundEntities[0], "正しいスタック可能アイテムが見つかるべき")

	// クリーンアップ
	world.Manager.DeleteEntity(stackableEntity)
	world.Manager.DeleteEntity(nonStackableEntity)
}

func TestQueryPlayer(t *testing.T) {
	t.Parallel()
	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.FactionAlly, &gc.FactionAlly)
	player.AddComponent(world.Components.Name, &gc.Name{Name: "プレイヤー"})

	// 敵を作成（除外されることを確認）
	enemy := world.Manager.NewEntity()
	enemy.AddComponent(world.Components.FactionEnemy, &gc.FactionEnemy)
	enemy.AddComponent(world.Components.Name, &gc.Name{Name: "敵"})

	// クエリを実行
	var foundEntities []ecs.Entity
	QueryPlayer(world, func(entity ecs.Entity) {
		foundEntities = append(foundEntities, entity)
	})

	// 結果を検証
	assert.Len(t, foundEntities, 1, "プレイヤーが1つだけ見つかるべき")
	assert.Equal(t, player, foundEntities[0], "正しいプレイヤーが見つかるべき")

	// クリーンアップ
	world.Manager.DeleteEntity(player)
	world.Manager.DeleteEntity(enemy)
}
