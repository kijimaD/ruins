package mapplaner

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
)

// PropPlacer は部屋に置物を配置するビルダー
type PropPlacer struct {
	// PropDensity は部屋あたりの置物配置密度（0.0-1.0）
	PropDensity float64
	// PropTypes は配置可能な置物タイプのリスト
	PropTypes []gc.PropType
}

// NewPropPlacer は新しいPropPlacerを作成する
func NewPropPlacer(density float64, propTypes []gc.PropType) PropPlacer {
	// デフォルトの置物タイプを設定
	if len(propTypes) == 0 {
		propTypes = []gc.PropType{
			gc.PropTypeTable,
			gc.PropTypeChair,
			gc.PropTypeBookshelf,
			gc.PropTypeBarrel,
			gc.PropTypeCrate,
		}
	}

	return PropPlacer{
		PropDensity: density,
		PropTypes:   propTypes,
	}
}

// BuildMeta は部屋に置物を配置する
func (pp PropPlacer) BuildMeta(buildData *BuilderMap) {
	// 各部屋に置物を配置
	for _, room := range buildData.Rooms {
		pp.placePropsInRoom(buildData, room)
	}
}

// placePropsInRoom は指定された部屋に置物を配置する
func (pp PropPlacer) placePropsInRoom(buildData *BuilderMap, room gc.Rect) {
	// 部屋のサイズに基づいて配置する置物の数を決定
	roomArea := int(room.X2-room.X1-2) * int(room.Y2-room.Y1-2) // 壁を除いた内部面積
	if roomArea <= 0 {
		return // 面積が0以下の場合は配置しない
	}

	// 密度に基づいて配置数を計算（最低1個、最大で面積の半分）
	maxProps := max(1, int(float64(roomArea)*pp.PropDensity))
	propCount := 1 + buildData.RandomSource.Intn(maxProps)

	// 部屋内の有効な位置を特定
	validPositions := pp.getValidPositions(buildData, room)
	if len(validPositions) == 0 {
		return // 配置可能な位置がない場合
	}

	// 実際に配置する数を有効位置数で制限
	propCount = min(propCount, len(validPositions))

	// ランダムに位置を選んで置物を配置
	for i := 0; i < propCount; i++ {
		// 残りの位置からランダムに選択
		posIndex := buildData.RandomSource.Intn(len(validPositions))
		pos := validPositions[posIndex]

		// 選択した位置を除去（重複を防ぐ）
		validPositions = append(validPositions[:posIndex], validPositions[posIndex+1:]...)

		// ランダムに置物タイプを選択
		propType := pp.PropTypes[buildData.RandomSource.Intn(len(pp.PropTypes))]

		// 置物を配置（エラーは無視して続行）
		pp.placeProp(buildData, propType, pos.X, pos.Y)
	}
}

// getValidPositions は部屋内の配置可能な位置を取得する
func (pp PropPlacer) getValidPositions(buildData *BuilderMap, room gc.Rect) []gc.GridElement {
	var validPositions []gc.GridElement

	// 部屋の内部をスキャン（壁を除く）
	for y := room.Y1 + 1; y < room.Y2-1; y++ {
		for x := room.X1 + 1; x < room.X2-1; x++ {
			idx := buildData.Level.XYTileIndex(x, y)

			// 床タイルかチェック
			if buildData.Tiles[idx] == TileFloor {
				validPositions = append(validPositions, gc.GridElement{X: x, Y: y})
			}
		}
	}

	return validPositions
}

// placeProp は指定位置に置物を配置する（実際のワールドには配置せず、ログのみ出力）
func (pp PropPlacer) placeProp(_ *BuilderMap, propType gc.PropType, x gc.Tile, y gc.Tile) {
	// 注意: この段階ではまだワールドが存在しないため、実際の配置は行わない
	// マップ生成完了後にSpawnPropを呼び出す必要がある
	fmt.Printf("置物配置予定: %s at (%d, %d)\n", propType, x, y)
}
