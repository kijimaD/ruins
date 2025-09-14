package systems

import (
	"github.com/kijimaD/ruins/lib/ai_input"
	w "github.com/kijimaD/ruins/lib/world"
)

// AISystem はAIエンティティの行動を処理するシステム
// ai_inputパッケージを使用してAI処理を委譲する
func AISystem(world w.World) {
	processor := ai_input.NewProcessor()
	processor.ProcessAllEntities(world)
}
