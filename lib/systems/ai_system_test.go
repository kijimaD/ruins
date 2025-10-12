package systems

import (
	"testing"

	"github.com/kijimaD/ruins/lib/aiinput"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/testutil"
	"github.com/kijimaD/ruins/lib/turns"
)

func TestAISystem(t *testing.T) {
	t.Parallel()

	// テスト用のワールド作成
	world := testutil.InitTestWorld(t)

	// TurnManagerを初期化
	world.Resources.TurnManager = turns.NewTurnManager()

	// プレイヤーエンティティを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: gc.Tile(10), Y: gc.Tile(10)})

	// AIエンティティを作成
	aiEntity := world.Manager.NewEntity()
	aiEntity.AddComponent(world.Components.AIMoveFSM, &gc.AIMoveFSM{})
	aiEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: gc.Tile(5), Y: gc.Tile(5)})
	aiEntity.AddComponent(world.Components.AIVision, &gc.AIVision{
		ViewDistance: gc.Pixel(100), // 3タイル程度の視界
		TargetEntity: &player,
	})
	aiEntity.AddComponent(world.Components.AIRoaming, &gc.AIRoaming{
		SubState:              gc.AIRoamingWaiting,
		StartSubStateTurn:     1,
		DurationSubStateTurns: 2,
	})

	// システム実行前の位置を記録
	initialGrid := world.Components.GridElement.Get(aiEntity).(*gc.GridElement)
	initialX, initialY := int(initialGrid.X), int(initialGrid.Y)

	// AIシステムを実行（aiinputパッケージを使用）
	processor := aiinput.NewProcessor()
	processor.ProcessAllEntities(world)

	// システム実行後の位置を記録
	finalGrid := world.Components.GridElement.Get(aiEntity).(*gc.GridElement)
	finalX, finalY := int(finalGrid.X), int(finalGrid.Y)

	// 位置が変わったかどうかを確認（ランダムな動きなので移動有無は不確定）
	moved := (initialX != finalX) || (initialY != finalY)
	t.Logf("AI移動: (%d,%d) -> (%d,%d), moved: %t", initialX, initialY, finalX, finalY, moved)

	// 状態が適切に管理されているかチェック
	roaming := world.Components.AIRoaming.Get(aiEntity).(*gc.AIRoaming)
	validStates := []gc.AIRoamingSubState{gc.AIRoamingWaiting, gc.AIRoamingDriving, gc.AIRoamingChasing}
	isValidState := false
	for _, state := range validStates {
		if roaming.SubState == state {
			isValidState = true
			break
		}
	}
	if !isValidState {
		t.Errorf("AI状態が無効: %v", roaming.SubState)
	}
	t.Logf("AI状態: %v", roaming.SubState)
}
