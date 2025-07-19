package effects

import (
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// GameLogAppender はゲームログ出力のインターフェース
type GameLogAppender interface {
	Append(entry string)
}

// Scope はエフェクトの影響範囲を保持する
type Scope struct {
	Creator *ecs.Entity     // 効果の発動者（nilの場合もある）
	Targets []ecs.Entity    // 効果の対象エンティティ一覧
	Logger  GameLogAppender // ゲームログ出力先（nilの場合はログ出力なし）
}

// Effect はゲーム内の効果を表す核心インターフェース
type Effect interface {
	// Apply は効果を実際に適用する
	Apply(world w.World, scope *Scope) error

	// Validate は効果の適用前に妥当性を検証する
	Validate(world w.World, scope *Scope) error

	// String は効果の文字列表現を返す（ログとデバッグ用）
	String() string
}
