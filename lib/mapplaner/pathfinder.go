package mapplanner

import gc "github.com/kijimaD/ruins/lib/components"

// PathFinder はパスファインディング機能を提供する
type PathFinder struct {
	buildData *MetaPlan
}

// NewPathFinder はPathFinderを作成する
func NewPathFinder(buildData *MetaPlan) *PathFinder {
	return &PathFinder{buildData: buildData}
}

// IsWalkable は指定座標が歩行可能かを判定する
func (pf *PathFinder) IsWalkable(x, y int) bool {
	width := int(pf.buildData.Level.TileWidth)
	height := int(pf.buildData.Level.TileHeight)

	// 境界チェック
	if x < 0 || x >= width || y < 0 || y >= height {
		return false
	}

	idx := pf.buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
	tile := pf.buildData.Tiles[idx]

	return tile.Walkable
}

// FindPath はBFSを使ってスタート地点からゴールまでのパスを探索する
// 上下左右の4方向移動のみサポート
func (pf *PathFinder) FindPath(startX, startY, goalX, goalY int) []Position {
	width := int(pf.buildData.Level.TileWidth)
	height := int(pf.buildData.Level.TileHeight)

	// スタートまたはゴールが歩行不可能な場合は空のパスを返す
	if !pf.IsWalkable(startX, startY) || !pf.IsWalkable(goalX, goalY) {
		return []Position{}
	}

	// 訪問済みマップ
	visited := make([][]bool, width)
	for i := range visited {
		visited[i] = make([]bool, height)
	}

	// 親ポイントマップ（パス復元用）
	parent := make([][]Position, width)
	for i := range parent {
		parent[i] = make([]Position, height)
		for j := range parent[i] {
			parent[i][j] = Position{X: -1, Y: -1} // 無効値で初期化
		}
	}

	// BFS用のキュー
	queue := []Position{{X: startX, Y: startY}}
	visited[startX][startY] = true

	// 4方向の移動方向
	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	// BFS実行
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// ゴールに到達した場合
		if current.X == goalX && current.Y == goalY {
			// パスを復元
			return pf.reconstructPath(parent, startX, startY, goalX, goalY)
		}

		// 隣接する4方向をチェック
		for _, dir := range directions {
			nextX := current.X + dir[0]
			nextY := current.Y + dir[1]

			// 境界チェックと歩行可能性チェック
			if nextX >= 0 && nextX < width && nextY >= 0 && nextY < height &&
				!visited[nextX][nextY] && pf.IsWalkable(nextX, nextY) {

				visited[nextX][nextY] = true
				parent[nextX][nextY] = Position{X: current.X, Y: current.Y}
				queue = append(queue, Position{X: nextX, Y: nextY})
			}
		}
	}

	// パスが見つからなかった場合は空のスライスを返す
	return []Position{}
}

// Position は座標を表す構造体
type Position struct {
	X int
	Y int
}

// reconstructPath は親ポイントマップからパスを復元する
func (pf *PathFinder) reconstructPath(parent [][]Position, startX, startY, goalX, goalY int) []Position {
	var path []Position
	current := Position{X: goalX, Y: goalY}

	// ゴールからスタートまで逆順にたどる
	for current.X != -1 && current.Y != -1 {
		path = append(path, current)
		if current.X == startX && current.Y == startY {
			break
		}
		current = parent[current.X][current.Y]
	}

	// パスを反転（スタートからゴールの順序にする）
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

// IsReachable はスタート地点からゴール地点まで到達可能かを判定する
func (pf *PathFinder) IsReachable(startX, startY, goalX, goalY int) bool {
	path := pf.FindPath(startX, startY, goalX, goalY)
	return len(path) > 0
}

// ValidateMapConnectivity はマップの接続性を検証する
// プレイヤーのスタート位置からワープポータル、脱出ポータルへの到達可能性をチェックする
func (pf *PathFinder) ValidateMapConnectivity(playerStartX, playerStartY int) MapConnectivityResult {
	width := int(pf.buildData.Level.TileWidth)
	height := int(pf.buildData.Level.TileHeight)

	result := MapConnectivityResult{
		PlayerStartReachable: pf.IsWalkable(playerStartX, playerStartY),
		WarpPortals:          []PortalReachability{},
		EscapePortals:        []PortalReachability{},
	}

	// 全タイルをスキャンしてワープポータルと脱出ポータルを見つける
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			idx := pf.buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
			tile := pf.buildData.Tiles[idx]

			switch tile {
			case TileWarpNext:
				reachable := pf.IsReachable(playerStartX, playerStartY, x, y)
				result.WarpPortals = append(result.WarpPortals, PortalReachability{
					X:         x,
					Y:         y,
					Reachable: reachable,
				})

			case TileWarpEscape:
				reachable := pf.IsReachable(playerStartX, playerStartY, x, y)
				result.EscapePortals = append(result.EscapePortals, PortalReachability{
					X:         x,
					Y:         y,
					Reachable: reachable,
				})
			}
		}
	}

	return result
}

// MapConnectivityResult はマップ接続性検証の結果
type MapConnectivityResult struct {
	// PlayerStartReachable はプレイヤーのスタート位置が歩行可能か
	PlayerStartReachable bool
	// WarpPortals はワープポータルの到達可能性リスト
	WarpPortals []PortalReachability
	// EscapePortals は脱出ポータルの到達可能性リスト
	EscapePortals []PortalReachability
}

// PortalReachability はポータルの到達可能性情報
type PortalReachability struct {
	X         int  // ポータルのX座標
	Y         int  // ポータルのY座標
	Reachable bool // プレイヤーのスタート地点から到達可能か
}

// HasReachableWarpPortal は到達可能なワープポータルがあるかを返す
func (result *MapConnectivityResult) HasReachableWarpPortal() bool {
	for _, portal := range result.WarpPortals {
		if portal.Reachable {
			return true
		}
	}
	return false
}

// HasReachableEscapePortal は到達可能な脱出ポータルがあるかを返す
func (result *MapConnectivityResult) HasReachableEscapePortal() bool {
	for _, portal := range result.EscapePortals {
		if portal.Reachable {
			return true
		}
	}
	return false
}

// IsFullyConnected はプレイヤーがすべての重要なポータルに到達可能かを返す
func (result *MapConnectivityResult) IsFullyConnected() bool {
	if !result.PlayerStartReachable {
		return false
	}

	// ワープポータルがある場合は、少なくとも1つは到達可能である必要がある
	if len(result.WarpPortals) > 0 && !result.HasReachableWarpPortal() {
		return false
	}

	// 脱出ポータルがある場合は、少なくとも1つは到達可能である必要がある
	if len(result.EscapePortals) > 0 && !result.HasReachableEscapePortal() {
		return false
	}

	return true
}
