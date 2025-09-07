package hud

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// Minimap はHUDのミニマップエリア
type Minimap struct {
	enabled bool
}

// NewMinimap は新しいHUDMinimapを作成する
func NewMinimap() *Minimap {
	return &Minimap{
		enabled: true,
	}
}

// SetEnabled は有効/無効を設定する
func (minimap *Minimap) SetEnabled(enabled bool) {
	minimap.enabled = enabled
}

// Update はミニマップを更新する
func (minimap *Minimap) Update(_ w.World) {
	// 現在は更新処理なし
}

// Draw はミニマップを描画する
func (minimap *Minimap) Draw(world w.World, screen *ebiten.Image) {
	if !minimap.enabled {
		return
	}

	// プレイヤー位置を取得
	var playerPos *gc.Position

	world.Manager.Join(
		world.Components.Position,
		world.Components.Operator,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerPos = world.Components.Position.Get(entity).(*gc.Position)
	}))

	if playerPos == nil {
		return // プレイヤーが見つからない場合は描画しない
	}

	// Dungeonリソースから探索済みマップとミニマップ設定を取得
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)

	// 探索済みタイルがない場合でも、ミニマップの枠だけは表示する
	if len(gameResources.ExploredTiles) == 0 {
		// 空のミニマップを描画
		minimap.drawEmpty(world, screen)
		return
	}

	// ミニマップの設定
	minimapWidth := gameResources.Minimap.Width
	minimapHeight := gameResources.Minimap.Height
	minimapScale := gameResources.Minimap.Scale // 1タイルをscaleピクセルで表現
	screenWidth := world.Resources.ScreenDimensions.Width
	minimapX := screenWidth - minimapWidth - 10 // 画面右端から10ピクセル内側
	minimapY := 10                              // 画面上端から10ピクセル下

	// ミニマップの背景を描画（半透明の黒い四角）
	minimapBg := ebiten.NewImage(minimapWidth, minimapHeight)
	minimapBg.Fill(color.RGBA{0, 0, 0, 128}) // 半透明の黒
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(minimapX), float64(minimapY))
	screen.DrawImage(minimapBg, op)

	// プレイヤーの現在位置をタイル座標に変換
	tileSize := 32 // タイルサイズ
	playerTileX := int(playerPos.X) / tileSize
	playerTileY := int(playerPos.Y) / tileSize

	// ミニマップの中心をプレイヤー位置に合わせる
	centerX := minimapX + minimapWidth/2
	centerY := minimapY + minimapHeight/2

	// 探索済みタイルを描画
	for tileKey := range gameResources.ExploredTiles {
		var tileX, tileY int
		if _, err := fmt.Sscanf(tileKey, "%d,%d", &tileX, &tileY); err != nil {
			continue
		}

		// プレイヤー位置からの相対位置を計算
		relativeX := tileX - playerTileX
		relativeY := tileY - playerTileY

		// ミニマップ上の座標を計算
		mapX := float32(centerX + relativeX*minimapScale)
		mapY := float32(centerY + relativeY*minimapScale)

		// ミニマップの範囲内かチェック
		if mapX >= float32(minimapX) && mapX <= float32(minimapX+minimapWidth-minimapScale) &&
			mapY >= float32(minimapY) && mapY <= float32(minimapY+minimapHeight-minimapScale) {

			// タイルのタイプに応じて色を決定
			tileColor := minimap.getTileColor(world, tileX, tileY)

			// 小さな四角形でタイルを表現
			vector.DrawFilledRect(screen, mapX, mapY, float32(minimapScale), float32(minimapScale), tileColor, false)
		}
	}

	// プレイヤーの位置を赤い点で表示
	playerMapX := float32(centerX)
	playerMapY := float32(centerY)
	vector.DrawFilledCircle(screen, playerMapX, playerMapY, 2, color.RGBA{255, 0, 0, 255}, false)

	// ミニマップの枠を描画
	minimap.drawFrame(screen, minimapX, minimapY, minimapWidth, minimapHeight)
}

// getTileColor はタイルの種類に応じてミニマップ上の色を返す
func (minimap *Minimap) getTileColor(world w.World, tileX, tileY int) color.RGBA {
	// そのタイル位置に実際にエンティティが存在するかチェック
	hasWall := false
	hasFloor := false

	// GridElement を持つエンティティをチェック
	world.Manager.Join(
		world.Components.GridElement,
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		grid := world.Components.GridElement.Get(entity).(*gc.GridElement)

		// グリッドの座標がタイル座標と一致するかチェック
		if int(grid.Row) == tileX && int(grid.Col) == tileY {
			// このタイルにエンティティが存在する
			if entity.HasComponent(world.Components.BlockView) {
				hasWall = true
			} else {
				hasFloor = true
			}
		}
	}))

	// 実際にエンティティが存在する場合のみ描画
	if hasWall {
		return color.RGBA{100, 100, 100, 255} // 壁は灰色
	} else if hasFloor {
		return color.RGBA{200, 200, 200, 128} // 床は薄い灰色
	}

	// 何もない場所は描画しない（透明）
	return color.RGBA{0, 0, 0, 0} // 透明
}

// drawEmpty は空のミニマップ（枠のみ）を描画する
func (minimap *Minimap) drawEmpty(world w.World, screen *ebiten.Image) {
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	minimapWidth := gameResources.Minimap.Width
	minimapHeight := gameResources.Minimap.Height
	screenWidth := world.Resources.ScreenDimensions.Width
	minimapX := screenWidth - minimapWidth - 10
	minimapY := 10

	// ミニマップの背景を描画（半透明の黒い四角）
	minimapBg := ebiten.NewImage(minimapWidth, minimapHeight)
	minimapBg.Fill(color.RGBA{0, 0, 0, 128})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(minimapX), float64(minimapY))
	screen.DrawImage(minimapBg, op)

	// ミニマップの枠を描画
	minimap.drawFrame(screen, minimapX, minimapY, minimapWidth, minimapHeight)

	// 中央に"No Data"テキストを表示
	ebitenutil.DebugPrintAt(screen, "No Data", minimapX+50, minimapY+70)
}

// drawFrame はミニマップの枠を描画する
func (minimap *Minimap) drawFrame(screen *ebiten.Image, x, y, width, height int) {
	whiteColor := color.RGBA{255, 255, 255, 255}

	// 枠線を描画
	vector.DrawFilledRect(screen, float32(x-1), float32(y-1), 1, float32(height+2), whiteColor, false)     // 左
	vector.DrawFilledRect(screen, float32(x+width), float32(y-1), 1, float32(height+2), whiteColor, false) // 右
	vector.DrawFilledRect(screen, float32(x-1), float32(y-1), float32(width+2), 1, whiteColor, false)      // 上
	vector.DrawFilledRect(screen, float32(x-1), float32(y+height), float32(width+2), 1, whiteColor, false) // 下
}
