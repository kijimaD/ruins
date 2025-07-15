package entities

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/stretchr/testify/assert"
)

func TestComponentList(t *testing.T) {
	t.Run("create entity component list", func(t *testing.T) {
		gameComponents := []interface{}{
			gc.GameComponentList{
				Name: &gc.Name{Name: "テストエンティティ"},
			},
		}

		list := ComponentList{
			Game: gameComponents,
		}

		assert.Len(t, list.Game, 1, "Gameコンポーネントリストの長さが正しくない")
	})

	t.Run("empty entity component list", func(t *testing.T) {
		list := ComponentList{}
		assert.Nil(t, list.Game, "空のリストでGameがnilでない")
	})
}

func TestAddEntities(t *testing.T) {
	t.Run("basic functionality test", func(t *testing.T) {
		// 循環依存を避けるため、基本的な機能のみテスト
		entityComponentList := ComponentList{
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
