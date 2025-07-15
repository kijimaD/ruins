package worldhelper

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/stretchr/testify/assert"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestQueryOwnedMaterial(t *testing.T) {
	world := game.InitWorld(960, 720)
	gameComponents := world.Components.Game

	// テスト用素材エンティティを作成
	materialEntity := world.Manager.NewEntity()
	materialEntity.AddComponent(gameComponents.Material, &gc.Material{Amount: 5})
	materialEntity.AddComponent(gameComponents.ItemLocationInBackpack, &gc.ItemLocationInBackpack)
	materialEntity.AddComponent(gameComponents.Name, &gc.Name{Name: "テスト素材"})

	// 素材でないエンティティを作成（除外されることを確認）
	nonMaterialEntity := world.Manager.NewEntity()
	nonMaterialEntity.AddComponent(gameComponents.ItemLocationInBackpack, &gc.ItemLocationInBackpack)
	nonMaterialEntity.AddComponent(gameComponents.Name, &gc.Name{Name: "テストアイテム"})

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

func TestQueryInPartyMember(t *testing.T) {
	world := game.InitWorld(960, 720)
	gameComponents := world.Components.Game

	// パーティメンバーを作成
	partyMember := world.Manager.NewEntity()
	partyMember.AddComponent(gameComponents.FactionAlly, &gc.FactionAlly)
	partyMember.AddComponent(gameComponents.InParty, &gc.InParty{})
	partyMember.AddComponent(gameComponents.Name, &gc.Name{Name: "パーティメンバー"})

	// 味方だがパーティにいないメンバーを作成
	allyMember := world.Manager.NewEntity()
	allyMember.AddComponent(gameComponents.FactionAlly, &gc.FactionAlly)
	allyMember.AddComponent(gameComponents.Name, &gc.Name{Name: "味方メンバー"})

	// クエリを実行
	var foundEntities []ecs.Entity
	QueryInPartyMember(world, func(entity ecs.Entity) {
		foundEntities = append(foundEntities, entity)
	})

	// 結果を検証
	assert.Len(t, foundEntities, 1, "パーティメンバーが1つだけ見つかるべき")
	assert.Equal(t, partyMember, foundEntities[0], "正しいパーティメンバーが見つかるべき")

	// クリーンアップ
	world.Manager.DeleteEntity(partyMember)
	world.Manager.DeleteEntity(allyMember)
}
