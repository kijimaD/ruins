package systems

import (
	"log"
	"math"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// CollisionSystem はプレイヤーと敵の衝突を検出し、戦闘遷移を発火する
func CollisionSystem(world w.World) {
	// 既に戦闘遷移イベントが設定されている場合は処理しない
	gameResources := world.Resources.Game.(*resources.Game)
	if gameResources.StateEvent != resources.StateEventNone {
		return
	}

	// プレイヤー数のカウント
	playerCount := 0
	world.Manager.Join(
		world.Components.Position,
		world.Components.Operator,
	).Visit(ecs.Visit(func(playerEntity ecs.Entity) {
		playerCount++
	}))

	// 敵数のカウント
	enemyCount := 0
	world.Manager.Join(
		world.Components.Position,
		world.Components.FactionEnemy,
	).Visit(ecs.Visit(func(enemyEntity ecs.Entity) {
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

		// 敵エンティティ（FactionEnemyコンポーネントを持つ）との衝突をチェック
		world.Manager.Join(
			world.Components.Position,
			world.Components.FactionEnemy,
		).Visit(ecs.Visit(func(enemyEntity ecs.Entity) {
			enemyPos := world.Components.Position.Get(enemyEntity).(*gc.Position)

			// 衝突判定（スプライトサイズを考慮した距離ベース）
			if checkCollisionSimple(world, playerEntity, enemyEntity, playerPos, enemyPos) {

				// 戦闘遷移エフェクトを実行
				processor := effects.NewProcessor()
				battleEffect := &effects.BattleEncounter{
					PlayerEntity:     playerEntity,
					FieldEnemyEntity: enemyEntity, // フィールド上の敵シンボル
				}
				processor.AddEffect(battleEffect, &playerEntity)

				if err := processor.Execute(world); err != nil {
					// エラーログは残しておく（重要なエラー情報のため）
					log.Printf("戦闘遷移エラー: %v", err)
				}
			}
		}))
	}))
}

// checkCollisionDistance は2つのエンティティ間の距離を返す
func checkCollisionDistance(world w.World, entity1, entity2 ecs.Entity, pos1, pos2 *gc.Position) float64 {
	// 中心間の距離を計算
	dx := float64(pos1.X - pos2.X)
	dy := float64(pos1.Y - pos2.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

// checkCollisionSimple はスプライトサイズを考慮したシンプルな距離判定を行う
func checkCollisionSimple(world w.World, entity1, entity2 ecs.Entity, pos1, pos2 *gc.Position) bool {
	// 両エンティティのスプライトサイズを取得
	size1 := getSpriteSize(world, entity1)
	size2 := getSpriteSize(world, entity2)

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
func getSpriteSize(world w.World, entity ecs.Entity) spriteSize {
	// TODO: ハードコーディングしないようにする
	defaultSize := spriteSize{width: 32, height: 32}

	if world.Components.SpriteRender.Get(entity) != nil {
		spriteRender := world.Components.SpriteRender.Get(entity).(*gc.SpriteRender)
		if spriteRender.SpriteSheet != nil && len(spriteRender.SpriteSheet.Sprites) > spriteRender.SpriteNumber {
			sprite := spriteRender.SpriteSheet.Sprites[spriteRender.SpriteNumber]
			return spriteSize{width: sprite.Width, height: sprite.Height}
		}
	}

	return defaultSize
}
