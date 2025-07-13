package loader

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	c "github.com/kijimaD/ruins/lib/engine/components"
	"github.com/stretchr/testify/assert"
)

func TestEngineComponentList(t *testing.T) {
	t.Run("create engine component list", func(t *testing.T) {
		spriteRender := &c.SpriteRender{
			SpriteNumber: 5,
		}

		list := EngineComponentList{
			SpriteRender: spriteRender,
		}

		assert.NotNil(t, list.SpriteRender, "SpriteRenderが設定されない")
		assert.Equal(t, 5, list.SpriteRender.SpriteNumber, "SpriteNumberが正しく設定されない")
	})

	t.Run("empty engine component list", func(t *testing.T) {
		list := EngineComponentList{}
		assert.Nil(t, list.SpriteRender, "空のリストでSpriteRenderがnilでない")
	})
}

func TestEntityComponentList(t *testing.T) {
	t.Run("create entity component list", func(t *testing.T) {
		gameComponents := []interface{}{
			gc.GameComponentList{
				Name: &gc.Name{Name: "テストエンティティ"},
			},
		}

		list := EntityComponentList{
			Game: gameComponents,
		}

		assert.Len(t, list.Game, 1, "Gameコンポーネントリストの長さが正しくない")
	})

	t.Run("empty entity component list", func(t *testing.T) {
		list := EntityComponentList{}
		assert.Nil(t, list.Game, "空のリストでGameがnilでない")
	})
}

func TestAddEntities(t *testing.T) {
	t.Run("basic functionality test", func(t *testing.T) {
		// 循環依存を避けるため、基本的な機能のみテスト
		entityComponentList := EntityComponentList{
			Game: []interface{}{
				gc.GameComponentList{
					Name: &gc.Name{Name: "テストエンティティ"},
				},
			},
		}

		// AddEntitiesは実際のworldオブジェクトが必要なため、
		// 構造体の正常性のみテスト
		assert.Len(t, entityComponentList.Game, 1, "Gameコンポーネントリストの長さが正しくない")

		// GameComponentListの中身を確認
		gameComponent := entityComponentList.Game[0].(gc.GameComponentList)
		assert.NotNil(t, gameComponent.Name, "Nameコンポーネントが設定されていない")
		assert.Equal(t, "テストエンティティ", gameComponent.Name.Name, "名前が正しく設定されていない")
	})
}

func TestProcessSpriteRenderData(t *testing.T) {
	t.Run("nil sprite render data", func(t *testing.T) {
		// processSpriteRenderDataは非公開関数だが、nilの場合の処理をテスト
		// 実際のworldオブジェクトなしでテスト可能な部分

		// fillDataの構造体作成テスト
		fillData := fillData{
			Width:  32,
			Height: 32,
			Color:  [4]uint8{255, 0, 0, 255}, // 赤色
		}

		assert.Equal(t, 32, fillData.Width, "fillDataの幅が正しくない")
		assert.Equal(t, 32, fillData.Height, "fillDataの高さが正しくない")
		assert.Equal(t, [4]uint8{255, 0, 0, 255}, fillData.Color, "fillDataの色が正しくない")

		// spriteRenderDataの構造体作成テスト
		spriteRenderData := spriteRenderData{
			Fill:            &fillData,
			SpriteSheetName: "",
			SpriteNumber:    0,
		}

		assert.NotNil(t, spriteRenderData.Fill, "fillDataが設定されていない")
		assert.Equal(t, "", spriteRenderData.SpriteSheetName, "SpriteSheetNameが正しくない")
		assert.Equal(t, 0, spriteRenderData.SpriteNumber, "SpriteNumberが正しくない")
	})
}
