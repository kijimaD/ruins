package worldhelper

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestSetMaxHPSP(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		vitality    int
		strength    int
		sensation   int
		dexterity   int
		agility     int
		level       int
		expectedHP  int
		expectedSP  int
		description string
	}{
		{
			name:        "レベル1基本ステータス",
			vitality:    10,
			strength:    8,
			sensation:   7,
			dexterity:   6,
			agility:     9,
			level:       1,
			expectedHP:  int(30 + float64(10*8+8+7)*1.0), // 30 + 95 = 125
			expectedSP:  int(float64(10*2+6+9) * 1.0),    // 35
			description: "レベル1での基本的なHP/SP計算",
		},
		{
			name:        "レベル5でのステータス",
			vitality:    15,
			strength:    12,
			sensation:   10,
			dexterity:   8,
			agility:     11,
			level:       5,
			expectedHP:  189, // 30 + 142 * 1.12 = 189.04 → 189
			expectedSP:  52,  // 49 * 1.08 = 52.92 → 52
			description: "レベル5でのHP/SP計算",
		},
		{
			name:        "高ステータス",
			vitality:    20,
			strength:    18,
			sensation:   15,
			dexterity:   14,
			agility:     16,
			level:       1,
			expectedHP:  int(30 + float64(20*8+18+15)*1.0), // 30 + 193 = 223
			expectedSP:  int(float64(20*2+14+16) * 1.0),    // 70
			description: "高ステータスでのHP/SP計算",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// 独立したworldを作成
			world, err := game.InitWorld(960, 720)
			require.NoError(t, err)

			// エンティティを作成
			entity := world.Manager.NewEntity()

			// Attributesコンポーネントを追加（BaseとTotalを0に設定してsetMaxHPSPの初期化をテスト）
			entity.AddComponent(world.Components.Attributes, &gc.Attributes{
				Vitality:  gc.Attribute{Base: tt.vitality, Total: 0},
				Strength:  gc.Attribute{Base: tt.strength, Total: 0},
				Sensation: gc.Attribute{Base: tt.sensation, Total: 0},
				Dexterity: gc.Attribute{Base: tt.dexterity, Total: 0},
				Agility:   gc.Attribute{Base: tt.agility, Total: 0},
				Defense:   gc.Attribute{Base: 5, Total: 0},
			})

			// Poolsコンポーネントを追加
			entity.AddComponent(world.Components.Pools, &gc.Pools{
				Level: tt.level,
				HP:    gc.Pool{Current: 0, Max: 0},
				SP:    gc.Pool{Current: 0, Max: 0},
			})

			// 関数を実行
			err = setMaxHPSP(world, entity)
			require.NoError(t, err)

			// 結果を検証
			pools := world.Components.Pools.Get(entity).(*gc.Pools)
			attrs := world.Components.Attributes.Get(entity).(*gc.Attributes)

			// Totalが正しく初期化されたことを確認
			assert.Equal(t, tt.vitality, attrs.Vitality.Total, "体力のTotal値が正しく初期化されていない")
			assert.Equal(t, tt.strength, attrs.Strength.Total, "力のTotal値が正しく初期化されていない")
			assert.Equal(t, tt.sensation, attrs.Sensation.Total, "感覚のTotal値が正しく初期化されていない")
			assert.Equal(t, tt.dexterity, attrs.Dexterity.Total, "器用さのTotal値が正しく初期化されていない")
			assert.Equal(t, tt.agility, attrs.Agility.Total, "素早さのTotal値が正しく初期化されていない")

			// HP/SPが正しく計算されたことを確認
			assert.Equal(t, tt.expectedHP, pools.HP.Max, "最大HPの計算が正しくない: %s", tt.description)
			assert.Equal(t, tt.expectedHP, pools.HP.Current, "現在HPが最大HPと同じでない: %s", tt.description)
			assert.Equal(t, tt.expectedSP, pools.SP.Max, "最大SPの計算が正しくない: %s", tt.description)
			assert.Equal(t, tt.expectedSP, pools.SP.Current, "現在SPが最大SPと同じでない: %s", tt.description)

			// クリーンアップ
			world.Manager.DeleteEntity(entity)
		})
	}
}

func TestSetMaxHPSP_WithoutComponents(t *testing.T) {
	t.Parallel()
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// 必要なコンポーネントがないエンティティ
	entity := world.Manager.NewEntity()

	// 関数を実行してエラーが発生することを確認
	err = setMaxHPSP(world, entity)
	require.Error(t, err, "必要なコンポーネントがない場合はエラーを返すべき")
	assert.Contains(t, err.Error(), "does not have required components", "エラーメッセージが適切であるべき")

	// クリーンアップ
	world.Manager.DeleteEntity(entity)
}

func TestFullRecover(t *testing.T) {
	t.Parallel()
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// テスト用エンティティを作成
	entity := world.Manager.NewEntity()
	entity.AddComponent(world.Components.Attributes, &gc.Attributes{
		Vitality:  gc.Attribute{Base: 10, Total: 0},
		Strength:  gc.Attribute{Base: 8, Total: 0},
		Sensation: gc.Attribute{Base: 7, Total: 0},
		Dexterity: gc.Attribute{Base: 6, Total: 0},
		Agility:   gc.Attribute{Base: 9, Total: 0},
		Defense:   gc.Attribute{Base: 5, Total: 0},
	})
	entity.AddComponent(world.Components.Pools, &gc.Pools{
		Level: 1,
		HP:    gc.Pool{Current: 0, Max: 0},
		SP:    gc.Pool{Current: 0, Max: 0},
	})

	// fullRecoverを実行
	fullRecover(world, entity)

	// 結果を検証
	pools := world.Components.Pools.Get(entity).(*gc.Pools)
	attrs := world.Components.Attributes.Get(entity).(*gc.Attributes)

	// 属性のTotalが正しく設定されたことを確認
	assert.Equal(t, 10, attrs.Vitality.Total, "体力のTotal値が正しく設定されていない")
	assert.Equal(t, 8, attrs.Strength.Total, "力のTotal値が正しく設定されていない")

	// HP/SPが正しく計算されたことを確認
	expectedHP := int(30 + float64(10*8+8+7)*1.0) // 30 + 95 = 125
	expectedSP := int(float64(10*2+6+9) * 1.0)    // 35
	assert.Equal(t, expectedHP, pools.HP.Max, "最大HPが正しく計算されていない")
	assert.Equal(t, expectedHP, pools.HP.Current, "現在HPが最大HPと一致していない")
	assert.Equal(t, expectedSP, pools.SP.Max, "最大SPが正しく計算されていない")
	assert.Equal(t, expectedSP, pools.SP.Current, "現在SPが最大SPと一致していない")

	// クリーンアップ
	world.Manager.DeleteEntity(entity)
}

func TestSpawnNPCHasAIMoveFSM(t *testing.T) {
	t.Parallel()
	// NPCが敵として認識されるAIMoveFSMコンポーネントを持つことを確認
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// SpriteSheetsを初期化
	spriteSheets := make(map[string]gc.SpriteSheet)
	spriteSheets["field"] = gc.SpriteSheet{
		Sprites: []gc.Sprite{
			{Width: 32, Height: 32}, // インデックス0
			{Width: 32, Height: 32}, // インデックス1
			{Width: 32, Height: 32}, // インデックス2
			{Width: 32, Height: 32}, // インデックス3
			{Width: 32, Height: 32}, // インデックス4
			{Width: 32, Height: 32}, // インデックス5
			{Width: 32, Height: 32}, // インデックス6 (NPCのスプライト)
		},
	}
	world.Resources.SpriteSheets = &spriteSheets

	// NPCを生成
	require.NoError(t, SpawnNPC(world, 100, 100))

	// AIMoveFSMコンポーネントを持つエンティティが存在することを確認
	enemyFound := false
	world.Manager.Join(
		world.Components.Position,
		world.Components.AIMoveFSM,
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		enemyFound = true
	}))

	assert.True(t, enemyFound, "SpawnNPCで生成されたエンティティはAIMoveFSMコンポーネントを持つべき")
}
