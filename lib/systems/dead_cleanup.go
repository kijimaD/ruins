package systems

import (
	"github.com/kijimaD/ruins/lib/logger"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// DeadCleanupSystem はDeadコンポーネントを持つ敵エンティティを削除する
func DeadCleanupSystem(world w.World) {
	logger := logger.New(logger.CategoryEntity)

	// Deadコンポーネントを持つエンティティを検索
	var toDelete []ecs.Entity
	world.Manager.Join(
		world.Components.Dead,
		world.Components.Player.Not(),
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		toDelete = append(toDelete, entity)
	}))

	// エンティティを削除
	for _, entity := range toDelete {
		world.Manager.DeleteEntity(entity)
	}

	if len(toDelete) > 0 {
		logger.Debug("Dead cleanup completed", "deleted_count", len(toDelete))
	}
}
