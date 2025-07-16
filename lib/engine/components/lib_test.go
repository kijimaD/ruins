package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// テスト用のコンポーネント
type TestComponents struct {
	TestSlice *ecs.SliceComponent
	TestNull  *ecs.NullComponent
}

func (t *TestComponents) InitializeComponents(manager *ecs.Manager) error {
	t.TestSlice = manager.NewSliceComponent()
	t.TestNull = manager.NewNullComponent()
	return nil
}

func TestInitComponents(t *testing.T) {
	t.Run("正常にコンポーネントを初期化できる", func(t *testing.T) {
		manager := ecs.NewManager()
		gameComponents := &TestComponents{}

		components, err := InitComponents(manager, gameComponents)

		assert.NoError(t, err)
		assert.NotNil(t, components)
		assert.NotNil(t, components.Game)
		assert.NotNil(t, components.Game.TestSlice)
		assert.NotNil(t, components.Game.TestNull)
	})

	t.Run("型安全性が保たれている", func(t *testing.T) {
		manager := ecs.NewManager()
		gameComponents := &TestComponents{}

		components, err := InitComponents(manager, gameComponents)

		assert.NoError(t, err)
		// 型アサーションが不要で、直接アクセスできる
		assert.IsType(t, &TestComponents{}, components.Game)
		assert.IsType(t, &ecs.SliceComponent{}, components.Game.TestSlice)
		assert.IsType(t, &ecs.NullComponent{}, components.Game.TestNull)
	})
}

func TestComponentInitializer(t *testing.T) {
	t.Run("ComponentInitializerインターフェースを実装している", func(_ *testing.T) {
		var _ ComponentInitializer = &TestComponents{}
	})
}
