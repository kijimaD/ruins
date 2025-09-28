package mapplanner

import "fmt"

// ValidationError は計画検証エラー
type ValidationError struct {
	Message string
	X, Y    int
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: (%d, %d)", e.Message, e.X, e.Y)
}

// NewValidationError は新しいValidationErrorを作成する
func NewValidationError(message string, x, y int) ValidationError {
	return ValidationError{
		Message: message,
		X:       x,
		Y:       y,
	}
}

// Validate は計画の妥当性と接続性をチェックする
func (mp *EntityPlan) Validate() error {
	// 座標範囲チェック
	for _, tile := range mp.Tiles {
		if tile.X < 0 || tile.X >= mp.Width || tile.Y < 0 || tile.Y >= mp.Height {
			return NewValidationError("タイル座標が範囲外", tile.X, tile.Y)
		}
	}

	for _, entity := range mp.Entities {
		if entity.X < 0 || entity.X >= mp.Width || entity.Y < 0 || entity.Y >= mp.Height {
			return NewValidationError("エンティティ座標が範囲外", entity.X, entity.Y)
		}
	}

	// 接続性チェック
	// TODO: そもそもPlayerPosがないのは異常なのでエラーを返すべき
	if mp.HasPlayerPos {
		if err := mp.validateConnectivity(); err != nil {
			return err
		}
	}

	return nil
}

// validateConnectivity はEntityPlan内の情報から接続性を検証する
func (mp *EntityPlan) validateConnectivity() error {
	// ワープポータルの位置を収集
	warpPortals := mp.collectWarpPortals()
	if len(warpPortals) == 0 {
		// ワープポータルがない場合は接続性チェック不要
		return nil
	}

	// EntityPlanからwalkableマップを構築
	walkableMap := mp.buildWalkableMap()

	// プレイヤー位置から各ワープポータルへの到達可能性をBFSで検証
	for _, portal := range warpPortals {
		if !mp.isReachableBFS(walkableMap, mp.PlayerStartX, mp.PlayerStartY, portal.X, portal.Y) {
			return ErrConnectivity
		}
	}

	return nil
}

// portalPosition はワープポータルの位置を表す
type portalPosition struct {
	X, Y int
	Type EntityType
}

// collectWarpPortals はEntityPlan内のワープポータルを収集する
func (mp *EntityPlan) collectWarpPortals() []portalPosition {
	var portals []portalPosition
	for _, entity := range mp.Entities {
		if entity.EntityType == EntityTypeWarpNext || entity.EntityType == EntityTypeWarpEscape {
			portals = append(portals, portalPosition{
				X:    entity.X,
				Y:    entity.Y,
				Type: entity.EntityType,
			})
		}
	}
	return portals
}

// buildWalkableMap はEntityPlanから歩行可能マップを構築する
func (mp *EntityPlan) buildWalkableMap() [][]bool {
	walkable := make([][]bool, mp.Height)
	for y := range walkable {
		walkable[y] = make([]bool, mp.Width)
	}

	// 全て歩行不可として初期化
	for y := 0; y < mp.Height; y++ {
		for x := 0; x < mp.Width; x++ {
			walkable[y][x] = false
		}
	}

	// TileSpecからタイルの歩行可能性を設定
	for _, tileSpec := range mp.Tiles {
		walkable[tileSpec.Y][tileSpec.X] = tileSpec.TileType.Walkable
	}

	// エンティティで歩行可能性を上書き
	// TODO: TileSpecのように通行情報をEntitySpecにもたせたほうがいいかもしれない
	for _, entity := range mp.Entities {
		switch entity.EntityType {
		case EntityTypeFloor:
			// 床エンティティは歩行可能
			walkable[entity.Y][entity.X] = true
		case EntityTypeWall:
			// 壁エンティティは歩行不可
			walkable[entity.Y][entity.X] = false
		case EntityTypeWarpNext, EntityTypeWarpEscape:
			// ワープポータルは歩行可能
			walkable[entity.Y][entity.X] = true
		case EntityTypeProp, EntityTypeNPC, EntityTypeItem:
			// 置物、NPC、アイテムは床の上に配置されるが歩行不可
			walkable[entity.Y][entity.X] = false
		}
	}

	return walkable
}

// isReachableBFS はBFSを使って到達可能性を判定する
func (mp *EntityPlan) isReachableBFS(walkable [][]bool, startX, startY, targetX, targetY int) bool {
	if startX == targetX && startY == targetY {
		return true
	}

	if !walkable[startY][startX] || !walkable[targetY][targetX] {
		return false
	}

	visited := make([][]bool, mp.Height)
	for y := range visited {
		visited[y] = make([]bool, mp.Width)
	}

	queue := []struct{ x, y int }{{startX, startY}}
	visited[startY][startX] = true

	directions := []struct{ dx, dy int }{
		{0, 1}, {0, -1}, {1, 0}, {-1, 0},
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, dir := range directions {
			nx, ny := current.x+dir.dx, current.y+dir.dy

			if nx < 0 || nx >= mp.Width || ny < 0 || ny >= mp.Height {
				continue
			}

			if visited[ny][nx] || !walkable[ny][nx] {
				continue
			}

			if nx == targetX && ny == targetY {
				return true
			}

			visited[ny][nx] = true
			queue = append(queue, struct{ x, y int }{nx, ny})
		}
	}

	return false
}