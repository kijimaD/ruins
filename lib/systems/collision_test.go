package systems

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/stretchr/testify/assert"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestCheckCollisionSimple(t *testing.T) {
	world := createTestWorldForCollision(t)

	tests := []struct {
		name        string
		pos1        *gc.Position
		pos2        *gc.Position
		size1       spriteSize
		size2       spriteSize
		expected    bool
		description string
	}{
		{
			name:        "同じ位置での衝突",
			pos1:        &gc.Position{X: 100, Y: 100},
			pos2:        &gc.Position{X: 100, Y: 100},
			size1:       spriteSize{width: 32, height: 32},
			size2:       spriteSize{width: 32, height: 32},
			expected:    true,
			description: "同じ位置にある同サイズのスプライト",
		},
		{
			name:        "近接での衝突",
			pos1:        &gc.Position{X: 100, Y: 100},
			pos2:        &gc.Position{X: 120, Y: 120},
			size1:       spriteSize{width: 32, height: 32},
			size2:       spriteSize{width: 32, height: 32},
			expected:    true,
			description: "スプライト半径内での衝突",
		},
		{
			name:        "離れた位置で非衝突",
			pos1:        &gc.Position{X: 100, Y: 100},
			pos2:        &gc.Position{X: 200, Y: 200},
			size1:       spriteSize{width: 32, height: 32},
			size2:       spriteSize{width: 32, height: 32},
			expected:    false,
			description: "十分離れた位置での非衝突",
		},
		{
			name:        "大きなスプライトでの衝突",
			pos1:        &gc.Position{X: 100, Y: 100},
			pos2:        &gc.Position{X: 130, Y: 130}, // 距離を短縮
			size1:       spriteSize{width: 64, height: 64},
			size2:       spriteSize{width: 64, height: 64},
			expected:    true,
			description: "大きなスプライト同士の衝突",
		},
		{
			name:        "異なるサイズでの衝突",
			pos1:        &gc.Position{X: 100, Y: 100},
			pos2:        &gc.Position{X: 120, Y: 120}, // 距離を短縮
			size1:       spriteSize{width: 64, height: 64},
			size2:       spriteSize{width: 16, height: 16},
			expected:    true,
			description: "異なるサイズのスプライトでの衝突",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用のエンティティを作成
			createEntityWithSpriteSize(t, world, float64(tt.pos1.X), float64(tt.pos1.Y), tt.size1, true)
			createEntityWithSpriteSize(t, world, float64(tt.pos2.X), float64(tt.pos2.Y), tt.size2, false)

			// エンティティを取得
			var entity1, entity2 ecs.Entity
			playerCount := 0
			world.Manager.Join(world.Components.Position, world.Components.Operator).Visit(ecs.Visit(func(e ecs.Entity) {
				if playerCount == 0 {
					entity1 = e
				}
				playerCount++
			}))

			enemyCount := 0
			world.Manager.Join(world.Components.Position, world.Components.FactionEnemy).Visit(ecs.Visit(func(e ecs.Entity) {
				if enemyCount == 0 {
					entity2 = e
				}
				enemyCount++
			}))

			// 衝突判定をテスト
			if playerCount > 0 && enemyCount > 0 {
				result := checkCollisionSimple(world, entity1, entity2, tt.pos1, tt.pos2)
				assert.Equal(t, tt.expected, result, tt.description)
			}

			// エンティティをクリーンアップ
			world.Manager.Join(world.Components.Position).Visit(ecs.Visit(func(e ecs.Entity) {
				world.Manager.DeleteEntity(e)
			}))
		})
	}
}

func TestGetSpriteSize(t *testing.T) {
	world := createTestWorldForCollision(t)

	tests := []struct {
		name     string
		width    int
		height   int
		expected spriteSize
	}{
		{
			name:     "標準サイズのスプライト",
			width:    32,
			height:   32,
			expected: spriteSize{width: 32, height: 32},
		},
		{
			name:     "大きなスプライト",
			width:    64,
			height:   48,
			expected: spriteSize{width: 64, height: 48},
		},
		{
			name:     "小さなスプライト",
			width:    16,
			height:   16,
			expected: spriteSize{width: 16, height: 16},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// スプライト付きエンティティを作成
			createEntityWithSprite(t, world, 100, 100, tt.width, tt.height, true)

			// エンティティを取得
			var testEntity ecs.Entity
			world.Manager.Join(world.Components.Position, world.Components.Operator).Visit(ecs.Visit(func(e ecs.Entity) {
				testEntity = e
			}))

			// スプライトサイズを取得
			size := getSpriteSize(world, testEntity)
			assert.Equal(t, tt.expected, size)

			// クリーンアップ
			world.Manager.DeleteEntity(testEntity)
		})
	}

	// スプライトがない場合のデフォルトサイズテスト
	t.Run("スプライトがない場合はデフォルトサイズ", func(t *testing.T) {
		createPlayerEntity(t, world, 100, 100)

		var testEntity ecs.Entity
		world.Manager.Join(world.Components.Position, world.Components.Operator).Visit(ecs.Visit(func(e ecs.Entity) {
			testEntity = e
		}))

		size := getSpriteSize(world, testEntity)
		expected := spriteSize{width: 32, height: 32}
		assert.Equal(t, expected, size)

		world.Manager.DeleteEntity(testEntity)
	})
}

func TestCollisionSystemWithMultipleEntities(t *testing.T) {
	world := createTestWorldForCollision(t)

	// 複数のプレイヤーと敵を作成
	createPlayerEntity(t, world, 100, 100)
	createEnemyEntity(t, world, 110, 110) // 接触
	createEnemyEntity(t, world, 200, 200) // 非接触
	createEnemyEntity(t, world, 105, 105) // 接触

	// システムが正常に動作することを確認
	assert.NotPanics(t, func() {
		CollisionSystem(world)
	})
}

func TestCollisionSystemWithNoEntities(t *testing.T) {
	world := createTestWorldForCollision(t)

	// エンティティが存在しない状態でシステムを実行
	assert.NotPanics(t, func() {
		CollisionSystem(world)
	})
}
