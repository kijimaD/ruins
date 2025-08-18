package save

import (
	"os"
	"path/filepath"
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestSaveLoadXPandLevel(t *testing.T) {
	t.Parallel()
	// 一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "save_test_")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// ワールドを作成
	w, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// キャラクターエンティティを作成（経験値とレベルを設定）
	character := w.Manager.NewEntity()
	character.AddComponent(w.Components.Name, &gc.Name{Name: "テストキャラ"})
	character.AddComponent(w.Components.FactionAlly, &gc.FactionAllyData{})
	character.AddComponent(w.Components.InParty, &gc.InParty{})
	character.AddComponent(w.Components.Pools, &gc.Pools{
		HP:      gc.Pool{Max: 120, Current: 100},
		SP:      gc.Pool{Max: 50, Current: 50},
		XP:      75,  // 75XP保持している状態
		Level:   3,   // レベル3
	})

	// 別のキャラクターも作成（レベルアップ直後の状態をシミュレート）
	character2 := w.Manager.NewEntity()
	character2.AddComponent(w.Components.Name, &gc.Name{Name: "テストキャラ2"})
	character2.AddComponent(w.Components.FactionAlly, &gc.FactionAllyData{})
	character2.AddComponent(w.Components.InParty, &gc.InParty{})
	character2.AddComponent(w.Components.Pools, &gc.Pools{
		HP:      gc.Pool{Max: 100, Current: 80},
		SP:      gc.Pool{Max: 40, Current: 40},
		XP:      0,   // レベルアップでリセットされた状態
		Level:   2,   // レベル2
	})

	// セーブマネージャーを作成
	sm := NewSerializationManager(tempDir)

	// ワールドを保存
	err = sm.SaveWorld(w, "test_xp")
	require.NoError(t, err)

	// セーブファイルが存在することを確認
	saveFile := filepath.Join(tempDir, "test_xp.json")
	_, err = os.Stat(saveFile)
	require.NoError(t, err)

	// 新しいワールドを作成してロード
	newWorld, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	err = sm.LoadWorld(newWorld, "test_xp")
	require.NoError(t, err)

	// キャラクターの経験値とレベルが正しく復元されることを確認
	characterCount := 0
	newWorld.Manager.Join(
		newWorld.Components.FactionAlly,
		newWorld.Components.InParty,
		newWorld.Components.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := newWorld.Components.Name.Get(entity).(*gc.Name)
		pools := newWorld.Components.Pools.Get(entity).(*gc.Pools)
		
		switch name.Name {
		case "テストキャラ":
			assert.Equal(t, 100, pools.HP.Current, "HPが正しく復元されていない")
			assert.Equal(t, 120, pools.HP.Max, "MaxHPが正しく復元されていない")
			assert.Equal(t, 75, pools.XP, "XPが正しく復元されていない")
			assert.Equal(t, 3, pools.Level, "Levelが正しく復元されていない")
			characterCount++
		case "テストキャラ2":
			assert.Equal(t, 80, pools.HP.Current, "キャラ2のHPが正しく復元されていない")
			assert.Equal(t, 100, pools.HP.Max, "キャラ2のMaxHPが正しく復元されていない")
			assert.Equal(t, 0, pools.XP, "キャラ2のXPが正しく復元されていない")
			assert.Equal(t, 2, pools.Level, "キャラ2のLevelが正しく復元されていない")
			characterCount++
		}
	}))
	assert.Equal(t, 2, characterCount, "キャラクターが正しくロードされていない")
}