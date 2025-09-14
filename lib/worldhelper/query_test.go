package worldhelper

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestQueryOwnedMaterial(t *testing.T) {
	t.Parallel()
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// テスト用素材エンティティを作成
	materialEntity := world.Manager.NewEntity()
	materialEntity.AddComponent(world.Components.Material, &gc.Material{Amount: 5})
	materialEntity.AddComponent(world.Components.ItemLocationInBackpack, &gc.ItemLocationInBackpack)
	materialEntity.AddComponent(world.Components.Name, &gc.Name{Name: "テスト素材"})

	// 素材でないエンティティを作成（除外されることを確認）
	nonMaterialEntity := world.Manager.NewEntity()
	nonMaterialEntity.AddComponent(world.Components.ItemLocationInBackpack, &gc.ItemLocationInBackpack)
	nonMaterialEntity.AddComponent(world.Components.Name, &gc.Name{Name: "テストアイテム"})

	// クエリを実行
	var foundEntities []ecs.Entity
	QueryOwnedMaterial(func(entity ecs.Entity) {
		foundEntities = append(foundEntities, entity)
	}, world)

	// 結果を検証
	assert.Len(t, foundEntities, 1, "素材エンティティが1つだけ見つかるべき")
	assert.Equal(t, materialEntity, foundEntities[0], "正しい素材エンティティが見つかるべき")

	// クリーンアップ
	world.Manager.DeleteEntity(materialEntity)
	world.Manager.DeleteEntity(nonMaterialEntity)
}

func TestQueryPlayer(t *testing.T) {
	t.Parallel()
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

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
