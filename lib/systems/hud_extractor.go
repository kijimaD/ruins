package systems

import (
	"fmt"
	"image/color"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/config"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/widgets/hud"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// ExtractHUDData はworldから全てのHUDデータを抽出する
func ExtractHUDData(world w.World) hud.HUDData {
	return hud.HUDData{
		GameInfo:     extractGameInfo(world),
		MinimapData:  extractMinimapData(world),
		DebugOverlay: extractDebugOverlay(world),
		MessageData:  extractMessageData(world),
	}
}

// extractGameInfo はゲーム基本情報を抽出する
func extractGameInfo(world w.World) hud.GameInfoData {
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	floorNumber := gameResources.Depth

	// プレイヤーの速度情報を取得
	var playerSpeed float64
	world.Manager.Join(
		world.Components.Velocity,
		world.Components.Position,
		world.Components.Operator,
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		velocity := world.Components.Velocity.Get(entity).(*gc.Velocity)
		playerSpeed = velocity.Speed
	}))

	return hud.GameInfoData{
		FloorNumber: floorNumber,
		PlayerSpeed: playerSpeed,
	}
}

// extractMinimapData はミニマップデータを抽出する
func extractMinimapData(world w.World) hud.MinimapData {
	// プレイヤー位置を取得
	var playerPos *gc.Position
	world.Manager.Join(
		world.Components.Position,
		world.Components.Operator,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerPos = world.Components.Position.Get(entity).(*gc.Position)
	}))

	if playerPos == nil {
		return hud.MinimapData{} // プレイヤーが見つからない場合は空データ
	}

	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	screenDimensions := hud.ScreenDimensions{
		Width:  world.Resources.ScreenDimensions.Width,
		Height: world.Resources.ScreenDimensions.Height,
	}

	// プレイヤーのタイル座標
	tileSize := 32
	playerTileX := int(playerPos.X) / tileSize
	playerTileY := int(playerPos.Y) / tileSize

	// タイル色情報を抽出
	tileColors := make(map[string]TileColorInfo)
	for tileKey := range gameResources.ExploredTiles {
		var tileX, tileY int
		if _, err := fmt.Sscanf(tileKey, "%d,%d", &tileX, &tileY); err != nil {
			continue
		}

		tileColor := getTileColorForMinimap(world, tileX, tileY)
		tileColors[tileKey] = TileColorInfo{
			R: tileColor.R,
			G: tileColor.G,
			B: tileColor.B,
			A: tileColor.A,
		}
	}

	return hud.MinimapData{
		PlayerTileX:   playerTileX,
		PlayerTileY:   playerTileY,
		ExploredTiles: gameResources.ExploredTiles,
		TileColors:    tileColors,
		MinimapConfig: hud.MinimapConfig{
			Width:  gameResources.Minimap.Width,
			Height: gameResources.Minimap.Height,
			Scale:  gameResources.Minimap.Scale,
		},
		ScreenDimensions: screenDimensions,
	}
}

// TileColorInfo はタイル色情報の内部型
type TileColorInfo = hud.TileColorInfo

// extractDebugOverlay はデバッグオーバーレイデータを抽出する
func extractDebugOverlay(world w.World) hud.DebugOverlayData {
	cfg := config.Get()
	if !cfg.ShowAIDebug {
		return hud.DebugOverlayData{Enabled: false}
	}

	// カメラ情報を取得
	var cameraPos gc.Position
	var cameraScale float64
	world.Manager.Join(
		world.Components.Camera,
		world.Components.Position,
	).Visit(ecs.Visit(func(camEntity ecs.Entity) {
		cameraPos = *world.Components.Position.Get(camEntity).(*gc.Position)
		camera := world.Components.Camera.Get(camEntity).(*gc.Camera)
		cameraScale = camera.Scale
	}))

	screenDimensions := hud.ScreenDimensions{
		Width:  world.Resources.ScreenDimensions.Width,
		Height: world.Resources.ScreenDimensions.Height,
	}

	// AI状態情報を抽出
	var aiStates []hud.AIStateInfo
	world.Manager.Join(
		world.Components.Position,
		world.Components.AIMoveFSM,
		world.Components.AIRoaming,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		position := world.Components.Position.Get(entity).(*gc.Position)
		roaming := world.Components.AIRoaming.Get(entity).(*gc.AIRoaming)

		// AIの現在の状態を判定
		var stateText string
		if entity.HasComponent(world.Components.AIChasing) {
			stateText = "CHASING"
		} else {
			switch roaming.SubState {
			case gc.AIRoamingWaiting:
				stateText = "WAITING"
			case gc.AIRoamingDriving:
				stateText = "ROAMING"
			case gc.AIRoamingChasing:
				stateText = "CHASING"
			default:
				stateText = "UNKNOWN"
			}
		}

		// 画面座標に変換
		screenX := (float64(position.X)-float64(cameraPos.X))*cameraScale + float64(screenDimensions.Width)/2
		screenY := (float64(position.Y)-float64(cameraPos.Y))*cameraScale + float64(screenDimensions.Height)/2

		aiStates = append(aiStates, hud.AIStateInfo{
			ScreenX:   screenX,
			ScreenY:   screenY,
			StateText: stateText,
		})
	}))

	// 視界範囲情報を抽出
	var visionRanges []hud.VisionRangeInfo
	world.Manager.Join(
		world.Components.Position,
		world.Components.AIVision,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		position := world.Components.Position.Get(entity).(*gc.Position)
		vision := world.Components.AIVision.Get(entity).(*gc.AIVision)

		screenX := (float64(position.X)-float64(cameraPos.X))*cameraScale + float64(screenDimensions.Width)/2
		screenY := (float64(position.Y)-float64(cameraPos.Y))*cameraScale + float64(screenDimensions.Height)/2
		scaledRadius := float32(float64(vision.ViewDistance) * cameraScale)

		visionRanges = append(visionRanges, hud.VisionRangeInfo{
			ScreenX:      screenX,
			ScreenY:      screenY,
			ScaledRadius: scaledRadius,
		})
	}))

	// 移動方向情報を抽出
	var movementDirs []hud.MovementDirectionInfo
	world.Manager.Join(
		world.Components.Position,
		world.Components.Velocity,
		world.Components.AIMoveFSM,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		position := world.Components.Position.Get(entity).(*gc.Position)
		velocity := world.Components.Velocity.Get(entity).(*gc.Velocity)

		if velocity.Speed > 0 && velocity.ThrottleMode == gc.ThrottleModeFront {
			screenX := (float64(position.X)-float64(cameraPos.X))*cameraScale + float64(screenDimensions.Width)/2
			screenY := (float64(position.Y)-float64(cameraPos.Y))*cameraScale + float64(screenDimensions.Height)/2

			movementDirs = append(movementDirs, hud.MovementDirectionInfo{
				ScreenX:     screenX,
				ScreenY:     screenY,
				Angle:       velocity.Angle,
				Speed:       velocity.Speed,
				CameraScale: cameraScale,
			})
		}
	}))

	return hud.DebugOverlayData{
		Enabled:            true,
		AIStates:           aiStates,
		VisionRanges:       visionRanges,
		MovementDirections: movementDirs,
		ScreenDimensions:   screenDimensions,
	}
}

// extractMessageData はメッセージデータを抽出する
func extractMessageData(world w.World) hud.MessageData {
	messages := extractMessagesFromGameLog()

	screenDimensions := hud.ScreenDimensions{
		Width:  world.Resources.ScreenDimensions.Width,
		Height: world.Resources.ScreenDimensions.Height,
	}

	// デフォルト設定を使用
	config := hud.DefaultMessageAreaConfig()

	return hud.MessageData{
		Messages:         messages,
		ScreenDimensions: screenDimensions,
		Config:           config,
	}
}

// extractMessagesFromGameLog はゲームログからメッセージを抽出する
func extractMessagesFromGameLog() []string {
	// gamelog.FieldLogから実際のメッセージを取得
	if gamelog.FieldLog == nil {
		return []string{}
	}

	// SafeSliceから全履歴のメッセージを取得
	return gamelog.FieldLog.GetHistory()
}

// getTileColorForMinimap はタイルの種類に応じてミニマップ上の色を返す
func getTileColorForMinimap(world w.World, tileX, tileY int) color.RGBA {
	hasWall := false
	hasFloor := false

	world.Manager.Join(
		world.Components.GridElement,
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		grid := world.Components.GridElement.Get(entity).(*gc.GridElement)

		if int(grid.Row) == tileX && int(grid.Col) == tileY {
			if entity.HasComponent(world.Components.BlockView) {
				hasWall = true
			} else {
				hasFloor = true
			}
		}
	}))

	if hasWall {
		return color.RGBA{100, 100, 100, 255} // 壁は灰色
	} else if hasFloor {
		return color.RGBA{200, 200, 200, 128} // 床は薄い灰色
	}

	return color.RGBA{0, 0, 0, 0} // 透明
}
