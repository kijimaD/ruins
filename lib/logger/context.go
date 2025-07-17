package logger

// Context は機能別のログコンテキストを表す型
type Context string

const (
	// ステート管理関連
	ContextTransition Context = "transition" // ステート遷移

	// ゲームシステム関連
	ContextBattle Context = "battle" // 戦闘処理
	ContextMove   Context = "move"   // 移動処理
	ContextInput  Context = "input"  // 入力処理
	ContextRender Context = "render" // 描画処理
	ContextVision Context = "vision" // 視界計算
	ContextHUD    Context = "hud"    // HUD表示

	// エフェクト関連
	ContextEffect Context = "effect" // エフェクト全般
	ContextDamage Context = "damage" // ダメージ処理
	ContextHeal   Context = "heal"   // 回復処理

	// リソース関連
	ContextResource Context = "resource" // リソース管理
	ContextLoad     Context = "load"     // ファイル読み込み
	ContextCache    Context = "cache"    // キャッシュ管理

	// ワールド・エンティティ関連
	ContextWorld     Context = "world"     // ワールド管理
	ContextEntity    Context = "entity"    // エンティティ操作
	ContextComponent Context = "component" // コンポーネント操作

	// デバッグ・パフォーマンス関連
	ContextPerf  Context = "perf"  // パフォーマンス計測
	ContextDebug Context = "debug" // デバッグ全般
)
