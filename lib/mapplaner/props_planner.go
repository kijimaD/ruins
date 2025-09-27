// Package mapplanner のProps配置プランナー - 責務分離によりmapspawnerから移動
package mapplanner

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
)

// PropsSpec はProps配置仕様を表す
type PropsSpec struct {
	X        int         // X座標
	Y        int         // Y座標
	PropType gc.PropType // Propsタイプ
}

// PropsPlanner はProps配置を担当するプランナー
type PropsPlanner struct {
	world       w.World
	plannerType PlannerType
}

// NewPropsPlanner はPropsプランナーを作成する
func NewPropsPlanner(world w.World, plannerType PlannerType) *PropsPlanner {
	return &PropsPlanner{
		world:       world,
		plannerType: plannerType,
	}
}

// PlanMeta はProps配置情報をMetaPlanに追加する
func (p *PropsPlanner) PlanMeta(planData *MetaPlan) {
	// 町タイプの場合は固定Props配置を追加
	if p.plannerType.Name == PlannerTypeTown.Name {
		p.addTownProps(planData)
	}

	// ダンジョンタイプの場合は固定Props配置を追加（必要に応じて実装）
	// TODO: 必要に応じて他のタイプも実装
}

// addTownProps は町用の固定Props配置をMetaPlanに追加する
func (p *PropsPlanner) addTownProps(planData *MetaPlan) {
	centerX := int(planData.Level.TileWidth) / 2
	centerY := int(planData.Level.TileHeight) / 2

	// 図書館の家具配置
	libraryProps := []struct {
		propType gc.PropType
		offsetX  int
		offsetY  int
	}{
		{gc.PropTypeBookshelf, -8, -19}, // 北壁沿い
		{gc.PropTypeBookshelf, -6, -19}, // 北壁沿い
		{gc.PropTypeTable, -7, -17},     // 閲覧机
		{gc.PropTypeChair, -7, -16},     // 閲覧用椅子
		{gc.PropTypeTable, -4, -15},     // 学習机
	}

	for _, prop := range libraryProps {
		x := centerX + prop.offsetX
		y := centerY + prop.offsetY
		if p.isValidPropPosition(planData, gc.Tile(x), gc.Tile(y)) {
			planData.Props = append(planData.Props, PropsSpec{
				X:        x,
				Y:        y,
				PropType: prop.propType,
			})
		}
	}

	// 学校の家具配置
	schoolProps := []struct {
		propType gc.PropType
		offsetX  int
		offsetY  int
	}{
		{gc.PropTypeBookshelf, 6, -19},  // 北壁
		{gc.PropTypeBookshelf, 8, -19},  // 北壁
		{gc.PropTypeBookshelf, 10, -19}, // 北壁
		{gc.PropTypeBookshelf, 5, -16},  // 西壁
		{gc.PropTypeBookshelf, 11, -16}, // 東壁
		{gc.PropTypeTable, 8, -15},      // 教卓
		{gc.PropTypeChair, 8, -14},      // 教師用椅子
	}

	for _, prop := range schoolProps {
		x := centerX + prop.offsetX
		y := centerY + prop.offsetY
		if p.isValidPropPosition(planData, gc.Tile(x), gc.Tile(y)) {
			planData.Props = append(planData.Props, PropsSpec{
				X:        x,
				Y:        y,
				PropType: prop.propType,
			})
		}
	}

	// 住民の家1の家具配置
	house1Props := []struct {
		propType gc.PropType
		offsetX  int
		offsetY  int
	}{
		{gc.PropTypeBed, 13, -7},   // 寝室
		{gc.PropTypeTable, 15, -5}, // 食事台
		{gc.PropTypeChair, 15, -4}, // 食事用椅子
		{gc.PropTypeChair, 16, -5}, // 食事用椅子
	}

	for _, prop := range house1Props {
		x := centerX + prop.offsetX
		y := centerY + prop.offsetY
		if p.isValidPropPosition(planData, gc.Tile(x), gc.Tile(y)) {
			planData.Props = append(planData.Props, PropsSpec{
				X:        x,
				Y:        y,
				PropType: prop.propType,
			})
		}
	}

	// 住民の家2の家具配置
	house2Props := []struct {
		propType gc.PropType
		offsetX  int
		offsetY  int
	}{
		{gc.PropTypeBed, 14, 2},   // 寝室
		{gc.PropTypeTable, 16, 4}, // 食事台
		{gc.PropTypeChair, 16, 5}, // 食事用椅子
	}

	for _, prop := range house2Props {
		x := centerX + prop.offsetX
		y := centerY + prop.offsetY
		if p.isValidPropPosition(planData, gc.Tile(x), gc.Tile(y)) {
			planData.Props = append(planData.Props, PropsSpec{
				X:        x,
				Y:        y,
				PropType: prop.propType,
			})
		}
	}

	// 公民館の座席配置
	hallProps := []struct {
		propType gc.PropType
		offsetX  int
		offsetY  int
	}{
		{gc.PropTypeChair, -6, 12}, // 集会用座席
		{gc.PropTypeChair, -4, 12}, // 集会用座席
		{gc.PropTypeChair, 4, 12},  // 集会用座席
		{gc.PropTypeChair, 6, 12},  // 集会用座席
	}

	for _, prop := range hallProps {
		x := centerX + prop.offsetX
		y := centerY + prop.offsetY
		if p.isValidPropPosition(planData, gc.Tile(x), gc.Tile(y)) {
			planData.Props = append(planData.Props, PropsSpec{
				X:        x,
				Y:        y,
				PropType: prop.propType,
			})
		}
	}

	// 事務所の家具配置
	officeProps := []struct {
		propType gc.PropType
		offsetX  int
		offsetY  int
	}{
		{gc.PropTypeBed, 12, 13},       // 休憩用ベッド
		{gc.PropTypeTable, 14, 15},     // 事務机
		{gc.PropTypeChair, 14, 16},     // 事務用椅子
		{gc.PropTypeBookshelf, 18, 14}, // 書類棚
	}

	for _, prop := range officeProps {
		x := centerX + prop.offsetX
		y := centerY + prop.offsetY
		if p.isValidPropPosition(planData, gc.Tile(x), gc.Tile(y)) {
			planData.Props = append(planData.Props, PropsSpec{
				X:        x,
				Y:        y,
				PropType: prop.propType,
			})
		}
	}

	// 市場の露店（簡略化して一部のみ配置）
	marketProps := []struct {
		propType gc.PropType
		offsetX  int
		offsetY  int
	}{
		{gc.PropTypeTable, -12, 5}, // 露店1
		{gc.PropTypeTable, -9, 5},  // 露店2
		{gc.PropTypeTable, -6, 5},  // 露店3
	}

	for _, prop := range marketProps {
		x := centerX + prop.offsetX
		y := centerY + prop.offsetY
		if p.isValidPropPosition(planData, gc.Tile(x), gc.Tile(y)) {
			planData.Props = append(planData.Props, PropsSpec{
				X:        x,
				Y:        y,
				PropType: prop.propType,
			})
		}
	}
}

// isValidPropPosition はProp配置に適した位置かチェックする
func (p *PropsPlanner) isValidPropPosition(planData *MetaPlan, x, y gc.Tile) bool {
	// 範囲チェック
	if x < 0 || x >= planData.Level.TileWidth || y < 0 || y >= planData.Level.TileHeight {
		return false
	}

	tileIdx := planData.Level.XYTileIndex(x, y)
	if int(tileIdx) >= len(planData.Tiles) {
		return false
	}

	tile := planData.Tiles[tileIdx]
	// 床タイルにのみ配置可能
	return tile == TileFloor
}
