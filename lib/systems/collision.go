package systems

import (
	"errors"
	"fmt"
	"log"
	"math"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// CollisionSystem はプレイヤーと敵の衝突を検出し、戦闘遷移を発火する
func CollisionSystem(world w.World) {
	// 既に戦闘遷移イベントが設定されている場合は処理しない
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	if gameResources.GetStateEvent() != resources.StateEventNone {
		return
	}

	// プレイヤー数のカウント
	playerCount := 0
	world.Manager.Join(
		world.Components.Position,
		world.Components.Operator,
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		playerCount++
	}))

	// 敵数のカウント
	enemyCount := 0
	world.Manager.Join(
		world.Components.Position,
		world.Components.AIMoveFSM,
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		enemyCount++
	}))

	// エンティティが存在しない場合は早期リターン
	if playerCount == 0 || enemyCount == 0 {
		return
	}

	// プレイヤー（Operatorコンポーネントを持つエンティティ）と敵の衝突をチェック
	world.Manager.Join(
		world.Components.Position,
		world.Components.Operator,
	).Visit(ecs.Visit(func(playerEntity ecs.Entity) {
		playerPos := world.Components.Position.Get(playerEntity).(*gc.Position)

		// 敵エンティティ（AIMoveFSMコンポーネントを持つ）との衝突をチェック
		world.Manager.Join(
			world.Components.Position,
			world.Components.AIMoveFSM,
		).Visit(ecs.Visit(func(enemyEntity ecs.Entity) {
			enemyPos := world.Components.Position.Get(enemyEntity).(*gc.Position)

			// 衝突判定（スプライトサイズを考慮した距離ベース）
			if checkCollisionSimple(world, playerEntity, enemyEntity, playerPos, enemyPos) {
				gameResources.SetStateEvent(resources.StateEventBattleStart)

				// 衝突した両エンティティの移動を停止
				if playerVelocity := world.Components.Velocity.Get(playerEntity); playerVelocity != nil {
					velocity := playerVelocity.(*gc.Velocity)
					velocity.ThrottleMode = gc.ThrottleModeNope
					velocity.Speed = 0
				}
				if enemyVelocity := world.Components.Velocity.Get(enemyEntity); enemyVelocity != nil {
					velocity := enemyVelocity.(*gc.Velocity)
					velocity.ThrottleMode = gc.ThrottleModeNope
					velocity.Speed = 0
				}

				return // 1回の衝突処理で終了
			}
		}))
	}))
}

// checkCollisionDistance は2つのエンティティ間の距離を返す
func checkCollisionDistance(_ w.World, _, _ ecs.Entity, pos1, pos2 *gc.Position) float64 {
	// 中心間の距離を計算
	dx := float64(pos1.X - pos2.X)
	dy := float64(pos1.Y - pos2.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

// checkCollisionSimple はスプライトサイズを考慮したシンプルな距離判定を行う
func checkCollisionSimple(world w.World, entity1, entity2 ecs.Entity, pos1, pos2 *gc.Position) bool {
	// 両エンティティのスプライトサイズを取得
	size1, err1 := getSpriteSize(world, entity1)
	if err1 != nil {
		log.Printf("Entity %v sprite size error: %v", entity1, err1)
		return false
	}

	size2, err2 := getSpriteSize(world, entity2)
	if err2 != nil {
		log.Printf("Entity %v sprite size error: %v", entity2, err2)
		return false
	}

	// 衝突判定距離を計算（両スプライトの半径の合計）
	radius1 := math.Max(float64(size1.width), float64(size1.height)) / 2
	radius2 := math.Max(float64(size2.width), float64(size2.height)) / 2
	collisionDistance := radius1 + radius2

	// 中心間の距離を計算
	distance := checkCollisionDistance(world, entity1, entity2, pos1, pos2)

	// 衝突判定
	return distance < collisionDistance
}

// spriteSize はスプライトのサイズを表す構造体
type spriteSize struct {
	width  int
	height int
}

// getSpriteSize はエンティティのスプライトサイズを取得する
// SpriteRenderコンポーネントが存在しない、またはスプライトが見つからない場合はエラーを返す
func getSpriteSize(world w.World, entity ecs.Entity) (spriteSize, error) {
	if world.Components.SpriteRender.Get(entity) == nil {
		return spriteSize{}, fmt.Errorf("entity %v does not have SpriteRender component", entity)
	}

	spriteRender := world.Components.SpriteRender.Get(entity).(*gc.SpriteRender)

	// Resourcesからスプライトシートを取得
	if world.Resources.SpriteSheets == nil {
		return spriteSize{}, errors.New("SpriteSheets resources not available")
	}
	spriteSheet, exists := (*world.Resources.SpriteSheets)[spriteRender.Name]
	if !exists {
		return spriteSize{}, fmt.Errorf("SpriteSheet %s not found", spriteRender.Name)
	}

	if len(spriteSheet.Sprites) <= spriteRender.SpriteNumber {
		return spriteSize{}, fmt.Errorf("sprite number %d is out of range (length: %d)",
			spriteRender.SpriteNumber, len(spriteSheet.Sprites))
	}

	sprite := spriteSheet.Sprites[spriteRender.SpriteNumber]
	return spriteSize{width: sprite.Width, height: sprite.Height}, nil
}
