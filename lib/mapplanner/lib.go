// Package mapplanner はマップ生成機能を提供する
// 参考: https://bfnightly.bracketproductions.com
package mapplanner

import (
	"log"
	"time"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/raw"
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
	// 部屋群。部屋は長方形の移動可能な空間のことをいう。
	// 部屋はタイルの集合体である
	Rooms []gc.Rect
	// 廊下群。廊下は部屋と部屋をつなぐ移動可能な空間のことをいう。
	// 廊下はタイルの集合体である
	Corridors [][]resources.TileIdx
	// RandomSource はシード値による再現可能なランダム生成を提供する
	RandomSource *RandomSource
	// 階層を構成するタイル群。長さはステージの大きさで決まる
	// 通行可能かを判定するための情報を保持している必要がある
	Tiles []raw.TileRaw
	// WarpPortals は配置予定のワープポータルリスト
	WarpPortals []WarpPortal
	// NPCs は配置予定のNPCリスト
	NPCs []NPCSpec
	// Items は配置予定のアイテムリスト
	Items []ItemSpec
	// Props は配置予定のPropsリスト
	Props []PropsSpec
	// RawMaster はタイル生成に使用するマスターデータ
	RawMaster *raw.Master
}

// IsSpawnableTile は指定タイル座標がスポーン可能かを返す
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
func (bm MetaPlan) UpTile(idx resources.TileIdx) raw.TileRaw {
	targetIdx := resources.TileIdx(int(idx) - int(bm.Level.TileWidth))
	if targetIdx < 0 {
		return bm.GenerateTile("Empty")
	}

	return bm.Tiles[targetIdx]
}

// DownTile は下にあるタイルを調べる
func (bm MetaPlan) DownTile(idx resources.TileIdx) raw.TileRaw {
	targetIdx := int(idx) + int(bm.Level.TileWidth)
	if targetIdx > len(bm.Tiles)-1 {
		return bm.GenerateTile("Empty")
	}

	return bm.Tiles[targetIdx]
}

// LeftTile は左にあるタイルを調べる
func (bm MetaPlan) LeftTile(idx resources.TileIdx) raw.TileRaw {
	targetIdx := idx - 1
	if targetIdx < 0 {
		return bm.GenerateTile("Empty")
	}

	return bm.Tiles[targetIdx]
}

// RightTile は右にあるタイルを調べる
func (bm MetaPlan) RightTile(idx resources.TileIdx) raw.TileRaw {
	targetIdx := idx + 1
	if int(targetIdx) > len(bm.Tiles)-1 {
		return bm.GenerateTile("Empty")
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
func (bm MetaPlan) isFloorOrWarp(tile raw.TileRaw) bool {
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
	tiles := make([]raw.TileRaw, tileCount)

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

// Plan はプランナーチェーンを実行してマップを生成する
func (b *PlannerChain) Plan() {
	if b.Starter == nil {
		log.Fatal("empty starter planner!")
	}
	if err := (*b.Starter).PlanInitial(&b.PlanData); err != nil {
		log.Fatalf("PlanInitial failed: %v", err)
	}

	for _, meta := range b.Planners {
		meta.PlanMeta(&b.PlanData)
	}
}

// ValidateConnectivity はマップの接続性を検証する
// プレイヤーのスタート位置からワープポータルへの到達可能性をチェックし、問題があればエラーを返す
func (b *PlannerChain) ValidateConnectivity(playerStartX, playerStartY int) error {
	pf := NewPathFinder(&b.PlanData)
	return pf.ValidateConnectivity(playerStartX, playerStartY)
}

// InitialMapPlanner は初期マップをプランするインターフェース
// タイルへの描画は行わず、構造体フィールドの値を初期化するだけ
type InitialMapPlanner interface {
	PlanInitial(*MetaPlan) error
}

// MetaMapPlanner はメタ情報をプランするインターフェース
type MetaMapPlanner interface {
	PlanMeta(*MetaPlan)
}

// NewSmallRoomPlanner はシンプルな小部屋プランナーを作成する
func NewSmallRoomPlanner(width gc.Tile, height gc.Tile, seed uint64) *PlannerChain {
	chain := NewPlannerChain(width, height, seed)
	chain.StartWith(RectRoomPlanner{})
	chain.With(NewFillAll("Wall"))      // 全体を壁で埋める
	chain.With(RoomDraw{})              // 部屋を描画
	chain.With(LineCorridorPlanner{})   // 廊下を作成
	chain.With(NewBoundaryWall("Wall")) // 最外周を壁で囲む

	return chain
}

// NewBigRoomPlanner は大部屋プランナーを作成する
// ランダムにバリエーションを適用する統合版
func NewBigRoomPlanner(width gc.Tile, height gc.Tile, seed uint64) *PlannerChain {
	chain := NewPlannerChain(width, height, seed)
	chain.StartWith(BigRoomPlanner{})
	chain.With(NewFillAll("Wall"))      // 全体を壁で埋める
	chain.With(BigRoomDraw{})           // 大部屋を描画（バリエーション込み）
	chain.With(NewBoundaryWall("Wall")) // 最外周を壁で囲む

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

// GenerateTile は指定されたタイルを生成する
// TOMLからの生成に失敗した場合はパニックする
// TODO: 消して直接呼び出せばよい
func (bm *MetaPlan) GenerateTile(name string) raw.TileRaw {
	if bm.RawMaster == nil {
		panic("RawMasterが設定されていない。TOMLからのタイル生成が必須である")
	}
	// RawMaster.GenerateTileは内部でpanicするため、そのまま呼び出す
	return bm.RawMaster.GenerateTile(name)
}

// GetPlayerStartPosition はプレイヤーの開始位置を取得する
func (bm *MetaPlan) GetPlayerStartPosition() (int, int, bool) {
	// 適切な開始位置を探す（SpawnFromMetaPlanと同じロジック）
	width := int(bm.Level.TileWidth)
	height := int(bm.Level.TileHeight)

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
			return pos.x, pos.y, true
		}
	}

	// 見つからない場合は全体をスキャン
	for _i, tile := range bm.Tiles {
		if tile.Walkable {
			i := resources.TileIdx(_i)
			x, y := bm.Level.XYTileCoord(i)
			return int(x), int(y), true
		}
	}

	// 見つからない場合
	return 0, 0, false
}
