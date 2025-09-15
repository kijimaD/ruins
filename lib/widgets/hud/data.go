package hud

// Data はすべてのHUDウィジェットが必要とするデータを統合する
type Data struct {
	GameInfo     GameInfoData
	MinimapData  MinimapData
	DebugOverlay DebugOverlayData
	MessageData  MessageData
}

// GameInfoData はゲーム基本情報のデータ
type GameInfoData struct {
	FloorNumber int // フロア番号
	TurnNumber  int // ターン番号
	PlayerMoves int // プレイヤーの残り移動ポイント
}

// MinimapData はミニマップ描画に必要なデータ
type MinimapData struct {
	PlayerTileX      int                      // プレイヤーのタイル座標X
	PlayerTileY      int                      // プレイヤーのタイル座標Y
	ExploredTiles    map[string]bool          // 探索済みタイル
	TileColors       map[string]TileColorInfo // タイル色情報
	MinimapConfig    MinimapConfig            // ミニマップ設定
	ScreenDimensions ScreenDimensions         // 画面サイズ
}

// TileColorInfo はタイルの色情報
type TileColorInfo struct {
	R, G, B, A uint8
}

// MinimapConfig はミニマップの設定
type MinimapConfig struct {
	Width  int // ミニマップ幅
	Height int // ミニマップ高さ
	Scale  int // スケール（1タイルのピクセル数）
}

// ScreenDimensions は画面サイズ
type ScreenDimensions struct {
	Width  int
	Height int
}

// DebugOverlayData はAIデバッグ情報のデータ
type DebugOverlayData struct {
	Enabled          bool              // デバッグ表示有効フラグ
	AIStates         []AIStateInfo     // AI状態情報
	VisionRanges     []VisionRangeInfo // 視界範囲情報
	ScreenDimensions ScreenDimensions  // 画面サイズ
}

// AIStateInfo はAI状態の情報
type AIStateInfo struct {
	ScreenX   float64 // 画面上のX座標
	ScreenY   float64 // 画面上のY座標
	StateText string  // 状態テキスト
}

// VisionRangeInfo は視界範囲の情報
type VisionRangeInfo struct {
	ScreenX      float64 // 中心の画面X座標
	ScreenY      float64 // 中心の画面Y座標
	ScaledRadius float32 // スケール済み半径
}

// MessageData はメッセージ表示に必要なデータ
type MessageData struct {
	Messages         []string          // 表示するメッセージ一覧
	ScreenDimensions ScreenDimensions  // 画面サイズ
	Config           MessageAreaConfig // メッセージエリア設定
}
