package logger

// Category は機能別のログカテゴリを表す型
type Category string

const (
	// ステート管理関連
	CategoryTransition Category = "transition" // ステート遷移

	// ゲームシステム関連
	CategoryBattle Category = "battle" // 戦闘処理
	CategoryMove   Category = "move"   // 移動処理
	CategoryInput  Category = "input"  // 入力処理
	CategoryRender Category = "render" // 描画処理
	CategoryVision Category = "vision" // 視界計算
	CategoryHUD    Category = "hud"    // HUD表示

	// エフェクト関連
	CategoryEffect Category = "effect" // エフェクト全般
	CategoryDamage Category = "damage" // ダメージ処理
	CategoryHeal   Category = "heal"   // 回復処理

	// リソース関連
	CategoryResource Category = "resource" // リソース管理
	CategoryLoad     Category = "load"     // ファイル読み込み
	CategoryCache    Category = "cache"    // キャッシュ管理

	// ワールド・エンティティ関連
	CategoryWorld     Category = "world"     // ワールド管理
	CategoryEntity    Category = "entity"    // エンティティ操作
	CategoryComponent Category = "component" // コンポーネント操作

	// デバッグ・パフォーマンス関連
	CategoryPerf  Category = "perf"  // パフォーマンス計測
	CategoryDebug Category = "debug" // デバッグ全般
)
