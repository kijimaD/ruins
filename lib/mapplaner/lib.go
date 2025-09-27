// Package mapplanner はマップ生成機能を提供する
// 参考: https://bfnightly.bracketproductions.com
package mapplanner

import (
	"fmt"
	"log"
	"time"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// WarpPortalType はワープポータルの種別
type WarpPortalType uint8

const (
	// WarpPortalNext は次の階に向かうワープポータル
	WarpPortalNext WarpPortalType = iota
	// WarpPortalEscape は脱出用ワープポータル
	WarpPortalEscape
)

// WarpPortal はワープポータルエンティティの配置情報
type WarpPortal struct {
	X    int            // X座標
	Y    int            // Y座標
	Type WarpPortalType // ポータルの種別
}

// MetaPlan は階層のタイルを作る元になる概念の集合体
type MetaPlan struct {
	// 階層情報
	Level resources.Level
	// 階層を構成するタイル群。長さはステージの大きさで決まる
	Tiles []Tile
	// 部屋群。部屋は長方形の移動可能な空間のことをいう。
	// 部屋はタイルの集合体である
	Rooms []gc.Rect
	// 廊下群。廊下は部屋と部屋をつなぐ移動可能な空間のことをいう。
	// 廊下はタイルの集合体である
	Corridors [][]resources.TileIdx
	// RandomSource はシード値による再現可能なランダム生成を提供する（内部プランナーのみアクセス可能）
	RandomSource *RandomSource
	// WarpPortals は配置予定のワープポータルリスト（内部プランナーからアクセス可能、外部からは読み取り専用）
	WarpPortals []WarpPortal
	// NPCs は配置予定のNPCリスト（内部プランナーからアクセス可能、外部からは読み取り専用）
	NPCs []NPCSpec
	// Items は配置予定のアイテムリスト（内部プランナーからアクセス可能、外部からは読み取り専用）
	Items []ItemSpec
	// Props は配置予定のPropsリスト（内部プランナーからアクセス可能、外部からは読み取り専用）
	Props []PropsSpec
}

// GetLevel は階層情報を読み取り専用で返す（外部からのアクセス用）
func (bm *MetaPlan) GetLevel() *resources.Level {
	return &bm.Level
}

// GetTiles はタイル配列を読み取り専用で返す（外部からのアクセス用）
func (bm MetaPlan) GetTiles() []Tile {
	return bm.Tiles
}

// GetRooms は部屋リストを読み取り専用で返す（外部からのアクセス用）
func (bm MetaPlan) GetRooms() []gc.Rect {
	return bm.Rooms
}

// GetCorridors は廊下リストを読み取り専用で返す（外部からのアクセス用）
func (bm MetaPlan) GetCorridors() [][]resources.TileIdx {
	return bm.Corridors
}

// GetWarpPortals はワープポータルリストを読み取り専用で返す（外部からのアクセス用）
func (bm MetaPlan) GetWarpPortals() []WarpPortal {
	return bm.WarpPortals
}

// GetNPCs はNPCリストを読み取り専用で返す（外部からのアクセス用）
func (bm MetaPlan) GetNPCs() []NPCSpec {
	return bm.NPCs
}

// GetItems はアイテムリストを読み取り専用で返す（外部からのアクセス用）
func (bm MetaPlan) GetItems() []ItemSpec {
	return bm.Items
}

// GetProps はPropsリストを読み取り専用で返す（外部からのアクセス用）
func (bm MetaPlan) GetProps() []PropsSpec {
	return bm.Props
}

// SetTiles はタイル配列を設定する（テスト用）
func (bm *MetaPlan) SetTiles(tiles []Tile) {
	bm.Tiles = tiles
}

// SetWarpPortals はワープポータルリストを設定する（テスト用）
func (bm *MetaPlan) SetWarpPortals(portals []WarpPortal) {
	bm.WarpPortals = portals
}

// AddWarpPortal はワープポータルを追加する（テスト用）
func (bm *MetaPlan) AddWarpPortal(portal WarpPortal) {
	bm.WarpPortals = append(bm.WarpPortals, portal)
}

// IsSpawnableTile は指定タイル座標がスポーン可能かを返す
// スポーンチェックは地図生成時にしか使わないだろう
func (bm MetaPlan) IsSpawnableTile(_ w.World, tx gc.Tile, ty gc.Tile) bool {
	idx := bm.Level.XYTileIndex(tx, ty)
	tile := bm.Tiles[idx]
	if !tile.Walkable {
		return false
	}

	// planning段階では、MetaPlan内の計画済みエンティティをチェック
	if bm.existPlannedEntityOnTile(int(tx), int(ty)) {
		return false
	}

	return true
}

// existPlannedEntityOnTile は指定座標に計画済みエンティティがあるかをチェック
func (bm MetaPlan) existPlannedEntityOnTile(x, y int) bool {
	// ワープポータルをチェック
	for _, portal := range bm.WarpPortals {
		if portal.X == x && portal.Y == y {
			return true
		}
	}

	// NPCをチェック
	for _, npc := range bm.NPCs {
		if npc.X == x && npc.Y == y {
			return true
		}
	}

	// アイテムをチェック
	for _, item := range bm.Items {
		if item.X == x && item.Y == y {
			return true
		}
	}

	// Propsをチェック
	for _, prop := range bm.Props {
		if prop.X == x && prop.Y == y {
			return true
		}
	}

	return false
}

// UpTile は上にあるタイルを調べる
func (bm MetaPlan) UpTile(idx resources.TileIdx) Tile {
	targetIdx := resources.TileIdx(int(idx) - int(bm.Level.TileWidth))
	if targetIdx < 0 {
		return TileEmpty
	}

	return bm.Tiles[targetIdx]
}

// DownTile は下にあるタイルを調べる
func (bm MetaPlan) DownTile(idx resources.TileIdx) Tile {
	targetIdx := int(idx) + int(bm.Level.TileWidth)
	if targetIdx > len(bm.Tiles)-1 {
		return TileEmpty
	}

	return bm.Tiles[targetIdx]
}

// LeftTile は左にあるタイルを調べる
func (bm MetaPlan) LeftTile(idx resources.TileIdx) Tile {
	targetIdx := idx - 1
	if targetIdx < 0 {
		return TileEmpty
	}

	return bm.Tiles[targetIdx]
}

// RightTile は右にあるタイルを調べる
func (bm MetaPlan) RightTile(idx resources.TileIdx) Tile {
	targetIdx := idx + 1
	if int(targetIdx) > len(bm.Tiles)-1 {
		return TileEmpty
	}

	return bm.Tiles[targetIdx]
}

// AdjacentOrthoAnyFloor は直交する近傍4タイルに移動可能タイルがあるか判定する
func (bm MetaPlan) AdjacentOrthoAnyFloor(idx resources.TileIdx) bool {
	return bm.UpTile(idx).Walkable ||
		bm.DownTile(idx).Walkable ||
		bm.RightTile(idx).Walkable ||
		bm.LeftTile(idx).Walkable
}

// AdjacentAnyFloor は直交・斜めを含む近傍8タイルに床があるか判定する
func (bm MetaPlan) AdjacentAnyFloor(idx resources.TileIdx) bool {
	x, y := bm.Level.XYTileCoord(idx)
	width := int(bm.Level.TileWidth)
	height := int(bm.Level.TileHeight)

	// 8方向の隣接タイル座標をチェック
	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1}, // 上段
		{0, -1}, {0, 1}, // 中段（中心を除く）
		{1, -1}, {1, 0}, {1, 1}, // 下段
	}

	for _, dir := range directions {
		nx, ny := int(x)+dir[0], int(y)+dir[1]

		// 境界チェック
		if nx < 0 || nx >= width || ny < 0 || ny >= height {
			continue
		}

		neighborIdx := bm.Level.XYTileIndex(gc.Tile(nx), gc.Tile(ny))
		tile := bm.Tiles[neighborIdx]

		if tile.Walkable {
			return true
		}
	}

	return false
}

// GetWallType は近傍パターンから適切な壁タイプを判定する
func (bm MetaPlan) GetWallType(idx resources.TileIdx) WallType {
	// 4方向の隣接タイルの床状況をチェック
	upFloor := bm.isFloorOrWarp(bm.UpTile(idx))
	downFloor := bm.isFloorOrWarp(bm.DownTile(idx))
	leftFloor := bm.isFloorOrWarp(bm.LeftTile(idx))
	rightFloor := bm.isFloorOrWarp(bm.RightTile(idx))

	// 単純なケース：一方向のみに床がある場合
	if singleWallType := bm.checkSingleDirectionWalls(upFloor, downFloor, leftFloor, rightFloor); singleWallType != WallTypeGeneric {
		return singleWallType
	}

	// 角のケース：2方向に床がある場合
	if cornerWallType := bm.checkCornerWalls(upFloor, downFloor, leftFloor, rightFloor); cornerWallType != WallTypeGeneric {
		return cornerWallType
	}

	// 複雑なパターンまたは判定不可の場合
	return WallTypeGeneric
}

// checkSingleDirectionWalls は単一方向に床がある場合の壁タイプを返す
func (bm MetaPlan) checkSingleDirectionWalls(upFloor, downFloor, leftFloor, rightFloor bool) WallType {
	if downFloor && !upFloor && !leftFloor && !rightFloor {
		return WallTypeTop // 下に床がある → 上壁
	}
	if upFloor && !downFloor && !leftFloor && !rightFloor {
		return WallTypeBottom // 上に床がある → 下壁
	}
	if rightFloor && !upFloor && !downFloor && !leftFloor {
		return WallTypeLeft // 右に床がある → 左壁
	}
	if leftFloor && !upFloor && !downFloor && !rightFloor {
		return WallTypeRight // 左に床がある → 右壁
	}
	return WallTypeGeneric
}

// checkCornerWalls は2方向に床がある場合の壁タイプを返す
func (bm MetaPlan) checkCornerWalls(upFloor, downFloor, leftFloor, rightFloor bool) WallType {
	if downFloor && rightFloor && !upFloor && !leftFloor {
		return WallTypeTopLeft // 下右に床 → 左上角
	}
	if downFloor && leftFloor && !upFloor && !rightFloor {
		return WallTypeTopRight // 下左に床 → 右上角
	}
	if upFloor && rightFloor && !downFloor && !leftFloor {
		return WallTypeBottomLeft // 上右に床 → 左下角
	}
	if upFloor && leftFloor && !downFloor && !rightFloor {
		return WallTypeBottomRight // 上左に床 → 右下角
	}
	return WallTypeGeneric
}

// isFloorOrWarp は移動可能タイルかを判定する
func (bm MetaPlan) isFloorOrWarp(tile Tile) bool {
	return tile.Walkable
}

// PlannerChain は階層データMetaPlanに対して適用する生成ロジックを保持する構造体
type PlannerChain struct {
	Starter  *InitialMapPlanner
	Planners []MetaMapPlanner
	PlanData MetaPlan
}

// NewPlannerChain はシード値を指定してプランナーチェーンを作成する
// シードが0の場合はランダムなシードを生成する
func NewPlannerChain(width gc.Tile, height gc.Tile, seed uint64) *PlannerChain {
	tileCount := int(width) * int(height)
	tiles := make([]Tile, tileCount)

	// シードが0の場合はランダムなシードを生成
	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
	}

	return &PlannerChain{
		Starter:  nil,
		Planners: []MetaMapPlanner{},
		PlanData: MetaPlan{
			Level: resources.Level{
				TileWidth:  width,
				TileHeight: height,
				Entities:   make([]ecs.Entity, tileCount),
			},
			Tiles:        tiles,
			Rooms:        []gc.Rect{},
			Corridors:    [][]resources.TileIdx{},
			RandomSource: NewRandomSource(seed),
			WarpPortals:  []WarpPortal{},
			NPCs:         []NPCSpec{},
			Items:        []ItemSpec{},
			Props:        []PropsSpec{},
		},
	}
}

// StartWith は初期プランナーを設定する
func (b *PlannerChain) StartWith(initialMapPlanner InitialMapPlanner) {
	b.Starter = &initialMapPlanner
}

// With はメタプランナーを追加する
func (b *PlannerChain) With(metaMapPlanner MetaMapPlanner) {
	b.Planners = append(b.Planners, metaMapPlanner)
}

// Build はプランナーチェーンを実行してマップを生成する
func (b *PlannerChain) Build() {
	if b.Starter == nil {
		log.Fatal("empty starter planner!")
	}
	(*b.Starter).BuildInitial(&b.PlanData)

	for _, meta := range b.Planners {
		meta.BuildMeta(&b.PlanData)
	}
}

// ValidateConnectivity はマップの接続性を検証する
// プレイヤーのスタート位置からワープ/脱出ポータルへの到達可能性をチェック
func (b *PlannerChain) ValidateConnectivity(playerStartX, playerStartY int) MapConnectivityResult {
	pf := NewPathFinder(&b.PlanData)
	return pf.ValidateMapConnectivity(playerStartX, playerStartY)
}

// BuildPlanFromTiles はMetaPlanからEntityPlanを構築する
// MetaPlanは生成過程で使用される中間データ、EntityPlanは最終的な配置計画
func (bm *MetaPlan) BuildPlanFromTiles() (*EntityPlan, error) {
	plan := NewEntityPlan(int(bm.Level.TileWidth), int(bm.Level.TileHeight))

	// プレイヤー開始位置を設定（タイル配列ベースの場合は中央付近）
	width := int(bm.Level.TileWidth)
	height := int(bm.Level.TileHeight)
	centerX := width / 2
	centerY := height / 2

	// スポーン可能な位置を探す
	playerX, playerY := centerX, centerY
	found := false

	// 複数の候補位置を試す
	attempts := []struct{ x, y int }{
		{width / 2, height / 2},         // 中央
		{width / 4, height / 4},         // 左上寄り
		{3 * width / 4, height / 4},     // 右上寄り
		{width / 4, 3 * height / 4},     // 左下寄り
		{3 * width / 4, 3 * height / 4}, // 右下寄り
	}

	// 最適な位置を探す
	for _, pos := range attempts {
		tileIdx := bm.Level.XYTileIndex(gc.Tile(pos.x), gc.Tile(pos.y))
		if int(tileIdx) < len(bm.Tiles) && bm.Tiles[tileIdx].Walkable {
			playerX, playerY = pos.x, pos.y
			found = true
			break
		}
	}

	// 見つからない場合は全体をスキャン
	if !found {
		for _i, tile := range bm.Tiles {
			if tile.Walkable {
				i := resources.TileIdx(_i)
				x, y := bm.Level.XYTileCoord(i)
				playerX, playerY = int(x), int(y)
				found = true
				break
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("プレイヤー配置可能な床タイルが見つかりません")
	}

	// プレイヤー位置を設定
	plan.SetPlayerStartPosition(playerX, playerY)

	// タイルを走査してEntityPlanを構築
	for _i, tile := range bm.Tiles {
		i := resources.TileIdx(_i)
		x, y := bm.Level.XYTileCoord(i)

		switch tile.Type {
		case TileTypeFloor:
			plan.AddFloor(int(x), int(y))

		case TileTypeWall:
			// 近傍8タイル（直交・斜め）にフロアがあるときだけ壁にする
			if bm.AdjacentAnyFloor(i) {
				// 壁タイプを判定（スプライト番号はmapspawnerで決定）
				wallType := bm.GetWallType(i)
				plan.AddWallWithType(int(x), int(y), wallType)
			}

		case TileTypeEmpty:
			// 空のタイルはエンティティを生成しない
			continue

		default:
			return nil, fmt.Errorf("未知のタイルタイプ: %d", tile.Type)
		}
	}

	// ワープポータルエンティティをEntityPlanに追加
	for _, portal := range bm.WarpPortals {
		switch portal.Type {
		case WarpPortalNext:
			plan.AddWarpNext(portal.X, portal.Y)
		case WarpPortalEscape:
			plan.AddWarpEscape(portal.X, portal.Y)
		}
	}

	// NPCエンティティをEntityPlanに追加
	for _, npc := range bm.NPCs {
		plan.AddNPC(npc.X, npc.Y, npc.NPCType)
	}

	// アイテムエンティティをEntityPlanに追加
	for _, item := range bm.Items {
		plan.AddItem(item.X, item.Y, item.ItemName)
	}

	// PropsエンティティをEntityPlanに追加
	for _, prop := range bm.Props {
		plan.AddProp(prop.X, prop.Y, prop.PropType)
	}

	return plan, nil
}

// InitialMapPlanner は初期マップをプランするインターフェース
// タイルへの描画は行わず、構造体フィールドの値を初期化するだけ
type InitialMapPlanner interface {
	BuildInitial(*MetaPlan)
}

// MetaMapPlanner はメタ情報をプランするインターフェース
type MetaMapPlanner interface {
	BuildMeta(*MetaPlan)
}

// NewSmallRoomPlanner はシンプルな小部屋プランナーを作成する
func NewSmallRoomPlanner(width gc.Tile, height gc.Tile, seed uint64) *PlannerChain {
	chain := NewPlannerChain(width, height, seed)
	chain.StartWith(RectRoomPlanner{})
	chain.With(NewFillAll(TileWall))      // 全体を壁で埋める
	chain.With(RoomDraw{})                // 部屋を描画
	chain.With(LineCorridorPlanner{})     // 廊下を作成
	chain.With(NewBoundaryWall(TileWall)) // 最外周を壁で囲む

	return chain
}

// NewTownPlanner は街の固定マッププランナーを作成する
func NewTownPlanner(width gc.Tile, height gc.Tile, seed uint64) *PlannerChain {
	// 新しい文字列ベースの街プランナーを使用
	return NewStringTownPlanner(width, height, seed)
}

// NewBigRoomPlanner は大部屋プランナーを作成する
// ランダムにバリエーションを適用する統合版
func NewBigRoomPlanner(width gc.Tile, height gc.Tile, seed uint64) *PlannerChain {
	chain := NewPlannerChain(width, height, seed)
	chain.StartWith(BigRoomPlanner{})
	chain.With(NewFillAll(TileWall))      // 全体を壁で埋める
	chain.With(BigRoomDraw{})             // 大部屋を描画（バリエーション込み）
	chain.With(NewBoundaryWall(TileWall)) // 最外周を壁で囲む

	return chain
}

// PlannerType はマップ生成の設定を表す構造体
type PlannerType struct {
	// プランナー名
	Name string
	// 敵をスポーンするか
	SpawnEnemies bool
	// アイテムをスポーンするか
	SpawnItems bool
	// ポータル位置を固定するか
	UseFixedPortalPos bool
	// プランナー関数
	PlannerFunc func(width gc.Tile, height gc.Tile, seed uint64) *PlannerChain
}

var (
	// PlannerTypeRandom はランダム選択用のプランナータイプ
	PlannerTypeRandom = PlannerType{
		Name:              "ランダム",
		SpawnEnemies:      true,
		SpawnItems:        true,
		UseFixedPortalPos: false,
	}

	// PlannerTypeSmallRoom は小部屋ダンジョンのプランナータイプ
	PlannerTypeSmallRoom = PlannerType{
		Name:              "小部屋",
		SpawnEnemies:      true,
		SpawnItems:        true,
		UseFixedPortalPos: false,
		PlannerFunc:       NewSmallRoomPlanner,
	}

	// PlannerTypeBigRoom は大部屋ダンジョンのプランナータイプ
	PlannerTypeBigRoom = PlannerType{
		Name:              "大部屋",
		SpawnEnemies:      true,
		SpawnItems:        true,
		UseFixedPortalPos: false,
		PlannerFunc:       NewBigRoomPlanner,
	}

	// PlannerTypeCave は洞窟ダンジョンのプランナータイプ
	PlannerTypeCave = PlannerType{
		Name:              "洞窟",
		SpawnEnemies:      true,
		SpawnItems:        true,
		UseFixedPortalPos: false,
		PlannerFunc:       NewCavePlanner,
	}

	// PlannerTypeRuins は遺跡ダンジョンのプランナータイプ
	PlannerTypeRuins = PlannerType{
		Name:              "遺跡",
		SpawnEnemies:      true,
		SpawnItems:        true,
		UseFixedPortalPos: false,
		PlannerFunc:       NewRuinsPlanner,
	}

	// PlannerTypeForest は森ダンジョンのプランナータイプ
	PlannerTypeForest = PlannerType{
		Name:              "森",
		SpawnEnemies:      true,
		SpawnItems:        true,
		UseFixedPortalPos: false,
		PlannerFunc:       NewForestPlanner,
	}

	// PlannerTypeTown は市街地のプランナータイプ
	PlannerTypeTown = PlannerType{
		Name:              "市街地",
		SpawnEnemies:      false, // 街では敵をスポーンしない
		SpawnItems:        false, // 街ではフィールドアイテムをスポーンしない
		UseFixedPortalPos: true,  // ポータル位置を固定
		PlannerFunc:       NewTownPlanner,
	}
)

// NewRandomPlanner はシード値を使用してランダムにプランナーを選択し作成する
func NewRandomPlanner(width gc.Tile, height gc.Tile, seed uint64) *PlannerChain {
	// シードが0の場合はランダムなシードを生成する。後続のビルダーに渡される
	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
	}

	// シード値からランダムソースを作成（ビルダー選択用）
	rs := NewRandomSource(seed)

	// ランダム選択対象のプランナータイプ（街は除外）
	candidateTypes := []PlannerType{
		PlannerTypeSmallRoom,
		PlannerTypeBigRoom,
		PlannerTypeCave,
		PlannerTypeRuins,
		PlannerTypeForest,
	}

	// ランダムに選択
	selectedType := candidateTypes[rs.Intn(len(candidateTypes))]

	return selectedType.PlannerFunc(width, height, seed)
}

// GetPlanners は登録されているプランナーのスライスを返す
func (b *PlannerChain) GetPlanners() []MetaMapPlanner {
	return b.Planners
}

// BuildPlan はPlannerChainを実行してEntityPlanを生成する
func BuildPlan(chain *PlannerChain) (*EntityPlan, error) {
	// プランナーチェーンを実行
	chain.Build()

	// PlanDataからEntityPlanを構築
	plan, err := chain.PlanData.BuildPlanFromTiles()
	if err != nil {
		return nil, fmt.Errorf("EntityPlan構築エラー: %w", err)
	}

	return plan, nil
}
