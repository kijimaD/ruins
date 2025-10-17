package states

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/testutil"
	"github.com/kijimaD/ruins/lib/turns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecuteMoveAction(t *testing.T) {
	t.Parallel()

	t.Run("正常な移動", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// プレイヤーを作成
		player := world.Manager.NewEntity()
		player.AddComponent(world.Components.Player, &gc.Player{})
		player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
		player.AddComponent(world.Components.TurnBased, &gc.TurnBased{})

		// 移動前の座標を確認
		gridBefore := world.Components.GridElement.Get(player).(*gc.GridElement)
		initialX := int(gridBefore.X)
		initialY := int(gridBefore.Y)

		// 北に移動
		ExecuteMoveAction(world, gc.DirectionUp)

		// 移動後の座標を確認
		gridAfter := world.Components.GridElement.Get(player).(*gc.GridElement)
		assert.Equal(t, initialX, int(gridAfter.X), "X座標は変化しないべき")
		assert.Equal(t, initialY-1, int(gridAfter.Y), "Y座標が1減るべき")
	})

	t.Run("プレイヤーが存在しない場合", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// プレイヤーなしで移動を試みる（パニックしないことを確認）
		ExecuteMoveAction(world, gc.DirectionUp)
		// エラーにならず何も起きないべき
	})

	t.Run("GridElementがない場合", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// GridElementなしのプレイヤーを作成
		player := world.Manager.NewEntity()
		player.AddComponent(world.Components.Player, &gc.Player{})

		// 移動を試みる（パニックしないことを確認）
		ExecuteMoveAction(world, gc.DirectionUp)
		// エラーにならず何も起きないべき
	})

	t.Run("8方向の移動", func(t *testing.T) {
		t.Parallel()

		directions := []struct {
			name      string
			direction gc.Direction
			expectedX int
			expectedY int
		}{
			{"北", gc.DirectionUp, 10, 9},
			{"南", gc.DirectionDown, 10, 11},
			{"東", gc.DirectionRight, 11, 10},
			{"西", gc.DirectionLeft, 9, 10},
			{"北東", gc.DirectionUpRight, 11, 9},
			{"北西", gc.DirectionUpLeft, 9, 9},
			{"南東", gc.DirectionDownRight, 11, 11},
			{"南西", gc.DirectionDownLeft, 9, 11},
		}

		for _, tt := range directions {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				world := testutil.InitTestWorld(t)

				player := world.Manager.NewEntity()
				player.AddComponent(world.Components.Player, &gc.Player{})
				player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
				player.AddComponent(world.Components.TurnBased, &gc.TurnBased{})

				ExecuteMoveAction(world, tt.direction)

				gridAfter := world.Components.GridElement.Get(player).(*gc.GridElement)
				assert.Equal(t, tt.expectedX, int(gridAfter.X), "X座標が正しく移動するべき")
				assert.Equal(t, tt.expectedY, int(gridAfter.Y), "Y座標が正しく移動するべき")
			})
		}
	})
}

func TestExecuteWaitAction(t *testing.T) {
	t.Parallel()

	t.Run("待機アクションの実行", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// プレイヤーを作成
		player := world.Manager.NewEntity()
		player.AddComponent(world.Components.Player, &gc.Player{})
		player.AddComponent(world.Components.TurnBased, &gc.TurnBased{})

		// 待機アクションを実行（パニックしないことを確認）
		ExecuteWaitAction(world)

		// プレイヤーエンティティが存在することを確認
		assert.True(t, player.HasComponent(world.Components.Player), "プレイヤーが存在するべき")
	})

	t.Run("プレイヤーが存在しない場合", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// プレイヤーなしで待機を試みる（パニックしないことを確認）
		ExecuteWaitAction(world)
		// エラーにならず何も起きないべき
	})
}

func TestExecuteEnterAction(t *testing.T) {
	t.Parallel()

	t.Run("アイテムがある場合", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// プレイヤーを作成
		player := world.Manager.NewEntity()
		player.AddComponent(world.Components.Player, &gc.Player{})
		player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
		player.AddComponent(world.Components.TurnBased, &gc.TurnBased{})

		// 同じ位置にアイテムを作成
		item := world.Manager.NewEntity()
		item.AddComponent(world.Components.Item, &gc.Item{})
		item.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
		item.AddComponent(world.Components.ItemLocationOnField, &gc.ItemLocationOnField)
		item.AddComponent(world.Components.Name, &gc.Name{Name: "テストアイテム"})

		// Enterアクションを実行
		ExecuteEnterAction(world)

		// Enterアクションが実行されることを確認（パニックしない）
		assert.True(t, player.HasComponent(world.Components.Player), "プレイヤーが存在するべき")
	})

	t.Run("ワープホールがある場合", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// プレイヤーを作成
		player := world.Manager.NewEntity()
		player.AddComponent(world.Components.Player, &gc.Player{})
		player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
		player.AddComponent(world.Components.TurnBased, &gc.TurnBased{})

		// 同じ位置にワープホールを作成
		warp := world.Manager.NewEntity()
		warp.AddComponent(world.Components.Warp, &gc.Warp{Mode: gc.WarpModeNext})
		warp.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

		// Enterアクションを実行（ワープ処理が呼ばれることを期待）
		ExecuteEnterAction(world)

		// ワープ処理が実行されたかは、実装によって検証方法が異なる
		// ここではパニックしないことを確認
	})

	t.Run("プレイヤーが存在しない場合", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// プレイヤーなしでEnterを試みる（パニックしないことを確認）
		ExecuteEnterAction(world)
		// エラーにならず何も起きないべき
	})

	t.Run("GridElementがない場合", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// GridElementなしのプレイヤーを作成
		player := world.Manager.NewEntity()
		player.AddComponent(world.Components.Player, &gc.Player{})

		// Enterを試みる（パニックしないことを確認）
		ExecuteEnterAction(world)
		// エラーにならず何も起きないべき
	})

	t.Run("何もない場所でEnter", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// プレイヤーを作成
		player := world.Manager.NewEntity()
		player.AddComponent(world.Components.Player, &gc.Player{})
		player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
		player.AddComponent(world.Components.TurnBased, &gc.TurnBased{})

		// Enterアクションを実行
		ExecuteEnterAction(world)

		// 何も起きないことを確認（パニックしない）
		assert.True(t, player.HasComponent(world.Components.Player), "プレイヤーが存在するべき")
	})
}

func TestExecuteMoveActionWithEnemy(t *testing.T) {
	t.Parallel()

	t.Run("敵がいる位置への移動は攻撃になる", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// ターンマネージャーを初期化
		turnManager := turns.NewTurnManager()
		world.Resources.TurnManager = turnManager

		// プレイヤーを作成
		player := world.Manager.NewEntity()
		player.AddComponent(world.Components.Player, &gc.Player{})
		player.AddComponent(world.Components.FactionAlly, &gc.FactionAlly)
		player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
		player.AddComponent(world.Components.TurnBased, &gc.TurnBased{})

		// 北隣に敵を作成
		enemy := world.Manager.NewEntity()
		enemy.AddComponent(world.Components.FactionEnemy, &gc.FactionEnemy)
		enemy.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 9})
		enemy.AddComponent(world.Components.TurnBased, &gc.TurnBased{})
		enemy.AddComponent(world.Components.Pools, &gc.Pools{
			HP: gc.Pool{Current: 100, Max: 100},
		})

		initialPlayerX := int(world.Components.GridElement.Get(player).(*gc.GridElement).X)
		initialPlayerY := int(world.Components.GridElement.Get(player).(*gc.GridElement).Y)

		// 北に移動（敵がいる方向）
		ExecuteMoveAction(world, gc.DirectionUp)

		// プレイヤーが移動していないことを確認（攻撃したため）
		gridAfter := world.Components.GridElement.Get(player).(*gc.GridElement)
		assert.Equal(t, initialPlayerX, int(gridAfter.X), "攻撃時はX座標が変化しないべき")
		assert.Equal(t, initialPlayerY, int(gridAfter.Y), "攻撃時はY座標が変化しないべき")

		// 敵エンティティが存在することを確認（攻撃が実行された）
		assert.True(t, enemy.HasComponent(world.Components.Pools), "敵が存在するべき")
	})
}

func TestCheckTileEvents(t *testing.T) {
	t.Parallel()

	t.Run("プレイヤーエンティティの場合のみイベントチェック", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// プレイヤーを作成
		player := world.Manager.NewEntity()
		player.AddComponent(world.Components.Player, &gc.Player{})
		player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

		// checkTileEventsを呼び出し（パニックしないことを確認）
		checkTileEvents(world, player, 10, 10)
	})

	t.Run("非プレイヤーエンティティの場合はイベントチェックしない", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// 敵を作成
		enemy := world.Manager.NewEntity()
		enemy.AddComponent(world.Components.FactionEnemy, &gc.FactionEnemy)
		enemy.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

		// checkTileEventsを呼び出し（パニックしないことを確認）
		checkTileEvents(world, enemy, 10, 10)
	})
}

func TestFindEnemyAtPosition(t *testing.T) {
	t.Parallel()

	t.Run("指定位置に敵がいる場合", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// プレイヤーを作成
		player := world.Manager.NewEntity()
		player.AddComponent(world.Components.Player, &gc.Player{})
		player.AddComponent(world.Components.FactionAlly, &gc.FactionAlly)
		player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

		// 指定位置に敵を作成
		enemy := world.Manager.NewEntity()
		enemy.AddComponent(world.Components.FactionEnemy, &gc.FactionEnemy)
		enemy.AddComponent(world.Components.GridElement, &gc.GridElement{X: 11, Y: 10})

		foundEnemy := findEnemyAtPosition(world, player, 11, 10)
		assert.Equal(t, enemy, foundEnemy, "敵が見つかるべき")
	})

	t.Run("指定位置に敵がいない場合", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// プレイヤーを作成
		player := world.Manager.NewEntity()
		player.AddComponent(world.Components.Player, &gc.Player{})
		player.AddComponent(world.Components.FactionAlly, &gc.FactionAlly)
		player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

		foundEnemy := findEnemyAtPosition(world, player, 11, 10)
		assert.Equal(t, 0, int(foundEnemy), "敵が見つからないべき")
	})

	t.Run("死亡している敵は無視される", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// プレイヤーを作成
		player := world.Manager.NewEntity()
		player.AddComponent(world.Components.Player, &gc.Player{})
		player.AddComponent(world.Components.FactionAlly, &gc.FactionAlly)
		player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

		// 死亡した敵を作成
		enemy := world.Manager.NewEntity()
		enemy.AddComponent(world.Components.FactionEnemy, &gc.FactionEnemy)
		enemy.AddComponent(world.Components.GridElement, &gc.GridElement{X: 11, Y: 10})
		enemy.AddComponent(world.Components.Dead, &gc.Dead{})

		foundEnemy := findEnemyAtPosition(world, player, 11, 10)
		assert.Equal(t, 0, int(foundEnemy), "死亡した敵は無視されるべき")
	})

	t.Run("味方は敵として扱われない", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// プレイヤーを作成
		player := world.Manager.NewEntity()
		player.AddComponent(world.Components.Player, &gc.Player{})
		player.AddComponent(world.Components.FactionAlly, &gc.FactionAlly)
		player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

		// 味方を作成
		ally := world.Manager.NewEntity()
		ally.AddComponent(world.Components.FactionAlly, &gc.FactionAlly)
		ally.AddComponent(world.Components.GridElement, &gc.GridElement{X: 11, Y: 10})

		foundEnemy := findEnemyAtPosition(world, player, 11, 10)
		assert.Equal(t, 0, int(foundEnemy), "味方は敵として扱われないべき")
	})
}

func TestIsHostileFaction(t *testing.T) {
	t.Parallel()

	t.Run("プレイヤー側と敵側は敵対関係", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		player := world.Manager.NewEntity()
		player.AddComponent(world.Components.FactionAlly, &gc.FactionAlly)

		enemy := world.Manager.NewEntity()
		enemy.AddComponent(world.Components.FactionEnemy, &gc.FactionEnemy)

		require.True(t, isHostileFaction(world, player, enemy), "プレイヤーと敵は敵対関係であるべき")
		require.True(t, isHostileFaction(world, enemy, player), "敵とプレイヤーは敵対関係であるべき")
	})

	t.Run("プレイヤー側同士は敵対関係ではない", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		player1 := world.Manager.NewEntity()
		player1.AddComponent(world.Components.FactionAlly, &gc.FactionAlly)

		player2 := world.Manager.NewEntity()
		player2.AddComponent(world.Components.FactionAlly, &gc.FactionAlly)

		assert.False(t, isHostileFaction(world, player1, player2), "プレイヤー側同士は敵対関係ではないべき")
	})

	t.Run("敵側同士は敵対関係ではない", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		enemy1 := world.Manager.NewEntity()
		enemy1.AddComponent(world.Components.FactionEnemy, &gc.FactionEnemy)

		enemy2 := world.Manager.NewEntity()
		enemy2.AddComponent(world.Components.FactionEnemy, &gc.FactionEnemy)

		assert.False(t, isHostileFaction(world, enemy1, enemy2), "敵側同士は敵対関係ではないべき")
	})
}

func TestFindClosedDoorAtPosition(t *testing.T) {
	t.Parallel()

	t.Run("指定位置に閉じたドアがある場合", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// 閉じたドアを作成
		door := world.Manager.NewEntity()
		door.AddComponent(world.Components.Door, &gc.Door{IsOpen: false, Orientation: gc.DoorOrientationHorizontal})
		door.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

		foundDoor := findClosedDoorAtPosition(world, 10, 10)
		assert.Equal(t, door, foundDoor, "閉じたドアが見つかるべき")
	})

	t.Run("指定位置にドアがない場合", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		foundDoor := findClosedDoorAtPosition(world, 10, 10)
		assert.Equal(t, 0, int(foundDoor), "ドアが見つからないべき")
	})

	t.Run("ドアが開いている場合は無視される", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// 開いたドアを作成
		door := world.Manager.NewEntity()
		door.AddComponent(world.Components.Door, &gc.Door{IsOpen: true, Orientation: gc.DoorOrientationVertical})
		door.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

		foundDoor := findClosedDoorAtPosition(world, 10, 10)
		assert.Equal(t, 0, int(foundDoor), "開いたドアは見つからないべき")
	})

	t.Run("複数のドアがある場合", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// 閉じたドアを作成
		door1 := world.Manager.NewEntity()
		door1.AddComponent(world.Components.Door, &gc.Door{IsOpen: false, Orientation: gc.DoorOrientationHorizontal})
		door1.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

		// 別の位置に開いたドアを作成
		door2 := world.Manager.NewEntity()
		door2.AddComponent(world.Components.Door, &gc.Door{IsOpen: true, Orientation: gc.DoorOrientationVertical})
		door2.AddComponent(world.Components.GridElement, &gc.GridElement{X: 11, Y: 10})

		foundDoor := findClosedDoorAtPosition(world, 10, 10)
		assert.Equal(t, door1, foundDoor, "閉じたドアのみが見つかるべき")

		foundDoor2 := findClosedDoorAtPosition(world, 11, 10)
		assert.Equal(t, 0, int(foundDoor2), "開いたドアは見つからないべき")
	})
}
