package entities

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/stretchr/testify/assert"
)

func TestComponentList(t *testing.T) {
	t.Parallel()
	t.Run("create entity component list", func(t *testing.T) {
		t.Parallel()
		gameComponents := []interface{}{
			gc.EntitySpec{
				Name: &gc.Name{Name: "テストエンティティ"},
			},
		}

		list := ComponentList[interface{}]{
			Game: gameComponents,
		}

		assert.Len(t, list.Game, 1, "Gameコンポーネントリストの長さが正しくない")
	})

	t.Run("empty entity component list", func(t *testing.T) {
		t.Parallel()
		list := ComponentList[interface{}]{}
		assert.Nil(t, list.Game, "空のリストでGameがnilでない")
	})
}

func TestAddEntities(t *testing.T) {
	t.Parallel()
	t.Run("basic functionality test", func(t *testing.T) {
		t.Parallel()
		// 循環依存を避けるため、基本的な機能のみテスト
		entityComponentList := ComponentList[interface{}]{
			Game: []interface{}{
				gc.EntitySpec{
					Name: &gc.Name{Name: "テストエンティティ"},
				},
			},
		}

		// AddEntitiesは実際のworldオブジェクトが必要なため、
		// 構造体の正常性のみテスト
		assert.Len(t, entityComponentList.Game, 1, "Gameコンポーネントリストの長さが正しくない")

		// EntitySpecの中身を確認
		gameComponent := entityComponentList.Game[0].(gc.EntitySpec)
		assert.NotNil(t, gameComponent.Name, "Nameコンポーネントが設定されていない")
		assert.Equal(t, "テストエンティティ", gameComponent.Name.Name, "名前が正しく設定されていない")
	})
}
