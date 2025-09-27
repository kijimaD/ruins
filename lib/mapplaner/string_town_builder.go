// Package mapplaner の文字列ベース街ビルダー
// 文字列で直接街のレイアウトを定義するシステム
package mapplaner

import (
	"fmt"
	"strings"

	gc "github.com/kijimaD/ruins/lib/components"
)

// NewStringTownPlanner は文字列ベースの街プランナーを作成する
func NewStringTownPlanner(width gc.Tile, height gc.Tile, seed uint64) *PlannerChain {
	// 50x50の大規模な街レイアウトを文字列で定義
	tileMap, entityMap := GetTownLayout()

	// レイアウトの整合性を検証
	if err := ValidateTownLayout(tileMap, entityMap); err != nil {
		// 連結性エラーの場合は警告のみ（ゲームは継続可能）
		if strings.Contains(err.Error(), "連結性検証エラー") {
			fmt.Printf("警告: %v\n", err)
		} else {
			// その他の重要なエラー（サイズ不整合、ワープホール数など）は致命的
			panic(fmt.Sprintf("街レイアウト検証エラー: %v", err))
		}
	}

	planner := &StringMapPlanner{
		TileMap:       tileMap,
		EntityMap:     entityMap,
		TileMapping:   getDefaultTileMapping(),
		EntityMapping: getDefaultEntityMapping(),
	}

	// 実際のマップサイズは文字列から自動検出される（50x50）
	chain := NewPlannerChain(width, height, seed)
	chain.StartWith(planner)
	return chain
}

// GetTownLayout は街のタイルとエンティティレイアウトを返す
// TODO: tileMap と entityMap を明確に分ける(両方で使えるものがあって使い分けがわかりにくい)
// TODO: 建物をプレハブ式にする
func GetTownLayout() ([]string, []string) {
	// 50x50の街レイアウト（幅3の道路と5x5以上の建物）
	tileMap := []string{
		"##################################################",
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##", // 北の境界道路（幅3）
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##",
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##",
		"#rrr#########f#########rrr##########f#########r###",
		"#rrr#fffffffffffffffffr#rrr#ffffffffffffffffff#r##", // 北区域の大きな建物
		"#rrr#fffffffffffffffffr#rrr#ffffffffffffffffff#r##",
		"#rrr#fffffffffffffffffr#rrr#ffffffffffffffffff#r##",
		"#rrr#fffffffffffffffffr#rrr#ffffffffffffffffff#r##",
		"#rrr#fffffffffffffffffr#rrr#ffffffffffffffffff#r##",
		"#rrr#########f#########rrr##########f#########r###",
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##", // 北から中央への道路（幅3）
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##",
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##",
		"#rrr####f###rrr####f###rrr####f###rrr####f####r###", // 住宅区域
		"#rrr#ffffff#rrr#ffffff#rrr#ffffff#rrr#fffffff#r###", // 5x5の家
		"#rrrfffffff#rrr#ffffff#rrr#ffffff#rrr#fffffff#r###",
		"#rrr#ffffff#rrrfffffff#rrr#ffffff#rrr#fffffff#r###",
		"#rrr#ffffff#rrr#ffffff#rrrfffffff#rrr#fffffff#r###",
		"#rrr#ffffff#rrr#ffffff#rrr#ffffff#rrrfffffffr#r###",
		"#rrr########rrr########rrr########rrr#########r###",
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##", // 中央大通り（幅3）
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##",
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##",
		"#rrr#########f#########rrr##########f#########r###", // 中央広場エリア
		"#rrr#fffffffffffffffrrfrrr#ffffffffffffffffff#r###", // 広場（床のみ、ワープホールはentityMapで定義）
		"#rrr#ffffffffffffffffrrrrr#ffffffffffffffffff#r###",
		"#rrr#ffffffffffffffffrrrrr#ffffffffffffffffff#r###",
		"#rrr#ffffffffffffffffrrrrr#ffffffffffffffffff#r###",
		"#rrr#ffffffffffffffffrrrrr#ffffffffffffffffff#r###",
		"#rrr#########f#########rrr##########f#########r###",
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##", // 広場から南への道路（幅3）
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##",
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##",
		"#rrr####f###rrr####f###rrr####f###rrr####f####r###", // 南区域の住宅
		"#rrr#ffffff#rrr#ffffff#rrr#ffffff#rrr#fffffff#r###", // 5x5の家
		"#rrrfffffff#rrr#ffffff#rrr#ffffff#rrr#fffffff#r###",
		"#rrr#ffffff#rrrfffffff#rrr#ffffff#rrr#fffffff#r###",
		"#rrr#ffffff#rrr#ffffff#rrrfffffff#rrr#fffffff#r###",
		"#rrr#ffffff#rrr#ffffff#rrr#ffffff#rrrfffffffr#r###",
		"#rrr########rrr########rrr########rrr#########r###",
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##", // 南区域の道路（幅3）
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##",
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##",
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##",
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##", // 南の大きな建物
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##",
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##",
		"#rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr##",
		"##################################################",
	}

	// エンティティ配置（建物、家具、NPCなど） 50x50（幅3道路対応）
	entityMap := []string{
		"..................................................",
		"..................................................", // 北の境界道路（幅3）
		"..................................................",
		"..................................................",
		"......................rrr.........................",
		"....&.CT..BH.....MS..r....SSSS...HHHH...MMMM......", // 北区域の大きな建物
		".....................r............................",
		"..........HH.........r.........OOO................",
		".....................r............................",
		"..MS.................r.TTT........................",
		"......................rrr.........................",
		"..................................................", // 北から中央への道路（幅3）
		"..................................................",
		"..................................................",
		"....#######....#######....#######....#########....", // 住宅区域
		"....&.TCH......&.BTS......&.HCT......&.SBH........", // 5x5の家
		"...............................................r..",
		"....M..........A..........O..........R............",
		"..................................................",
		"..................................................",
		"....#######....#######....#######....#########....",
		"..................................................", // 中央大通り（幅3）
		"..................................................",
		"..................................................",
		"......................rrr.........................", // 中央広場エリア
		"...........................&......................", // 広場（NPC）
		"..................@...............................", // プレイヤー開始位置
		"........HH......OOO......MM.......................",
		"..................................................",
		"..................................................",
		"......................rrr.........................",
		"..................................................", // 広場から南への道路（幅3）
		"..................................................",
		"..................................................",
		"....#######....#######....#######....#########....", // 南区域の住宅
		"....&.BCT......&.HSM......&.TAO......&.RHB........", // 5x5の家
		"..................................................",
		"....H..........C..........S..........A............",
		"..................................................",
		"..................................................",
		"....#######....#######....#######....#########....",
		"..................................................", // 南区域の道路（幅3）
		"..................................................",
		"..................................................",
		"..................................................",
		"..................................................", // 広場
		"..................................................",
		"..................................................", // ワープホール（下の広場の真ん中）
		"..................................................",
		"........................w.........................",
	}

	return tileMap, entityMap
}

// ValidateTownLayout は街レイアウトの整合性を検証する
func ValidateTownLayout(tileMap, entityMap []string) error {
	if len(tileMap) == 0 || len(entityMap) == 0 {
		return fmt.Errorf("マップが空です")
	}

	height := len(tileMap)
	width := len(tileMap[0])

	// サイズ一致確認
	if len(entityMap) != height {
		return fmt.Errorf("タイルマップとエンティティマップの行数が一致しません: %d vs %d", len(tileMap), len(entityMap))
	}

	// 各行の長さ確認
	for i, row := range tileMap {
		if len(row) != width {
			return fmt.Errorf("タイルマップの行 %d の長さが不一致: 期待値 %d, 実際 %d", i, width, len(row))
		}
	}

	for i, row := range entityMap {
		if len(row) != width {
			return fmt.Errorf("エンティティマップの行 %d の長さが不一致: 期待値 %d, 実際 %d", i, width, len(row))
		}
	}

	// 境界壁の確認
	for x := 0; x < width; x++ {
		if tileMap[0][x] != '#' || tileMap[height-1][x] != '#' {
			return fmt.Errorf("上下の境界が壁でありません: 位置 (%d, 0) または (%d, %d)", x, x, height-1)
		}
	}

	for y := 0; y < height; y++ {
		if tileMap[y][0] != '#' || tileMap[y][width-1] != '#' {
			return fmt.Errorf("左右の境界が壁でありません: 位置 (0, %d) または (%d, %d)", y, width-1, y)
		}
	}

	// ワープホールの確認（tileMapとentityMapの両方をチェック）
	warpCount := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if tileMap[y][x] == 'w' {
				warpCount++
			}
		}
	}

	// entityMapでもワープホールをカウント
	for y := 0; y < height && y < len(entityMap); y++ {
		for x := 0; x < width && x < len(entityMap[y]); x++ {
			if entityMap[y][x] == 'w' {
				warpCount++
			}
		}
	}

	if warpCount != 1 {
		return fmt.Errorf("ワープホールが正確に1つある必要があります: 実際 %d", warpCount)
	}

	// 連結性の確認（簡易版）
	if err := validateConnectivity(tileMap, width, height); err != nil {
		return fmt.Errorf("連結性検証エラー: %w", err)
	}

	return nil
}

// validateConnectivity は簡易的な連結性チェックを行う
func validateConnectivity(tileMap []string, width, height int) error {
	// 歩行可能なタイルを探索
	visited := make([][]bool, height)
	for i := range visited {
		visited[i] = make([]bool, width)
	}

	// 最初の歩行可能な位置を見つける
	startX, startY := -1, -1
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if isWalkable(tileMap[y][x]) {
				startX, startY = x, y
				break
			}
		}
		if startX >= 0 {
			break
		}
	}

	if startX < 0 {
		return fmt.Errorf("歩行可能なタイルが見つかりません")
	}

	// 深度優先探索で到達可能な全タイルをマーク
	reachableCount := dfs(tileMap, visited, startX, startY, width, height)

	// 全歩行可能タイルの数をカウント
	totalWalkable := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if isWalkable(tileMap[y][x]) {
				totalWalkable++
			}
		}
	}

	if reachableCount != totalWalkable {
		return fmt.Errorf("一部のフロアに到達できません: 到達可能 %d, 総フロア数 %d", reachableCount, totalWalkable)
	}

	return nil
}

// isWalkable は歩行可能なタイルかを判定する
func isWalkable(tile byte) bool {
	return tile == 'f' || tile == 'r' || tile == 'w' // 床、道路、ワープホール
}

// dfs は深度優先探索で到達可能なタイル数を返す
func dfs(tileMap []string, visited [][]bool, x, y, width, height int) int {
	if x < 0 || x >= width || y < 0 || y >= height || visited[y][x] || !isWalkable(tileMap[y][x]) {
		return 0
	}

	visited[y][x] = true
	count := 1

	// 4方向を探索
	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	for _, dir := range directions {
		count += dfs(tileMap, visited, x+dir[0], y+dir[1], width, height)
	}

	return count
}
