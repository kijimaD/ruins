package gamelog

import (
	"sync"
)

const (
	// DefaultMaxLogSize はデフォルトの最大ログサイズ
	DefaultMaxLogSize = 1000
	// FieldLogMaxSize はフィールドログの最大サイズ
	FieldLogMaxSize = 100
	// BattleLogMaxSize は戦闘ログの最大サイズ
	BattleLogMaxSize = 500
	// SceneLogMaxSize はシーンログの最大サイズ
	SceneLogMaxSize = 200
)

var (
	// BattleLog は戦闘用ログ
	BattleLog = NewSafeSlice(BattleLogMaxSize)
	// FieldLog はフィールド用ログ
	FieldLog = NewSafeSlice(FieldLogMaxSize)
	// SceneLog は会話シーンでステータス変化を通知する用ログ
	SceneLog = NewSafeSlice(SceneLogMaxSize)
)

// SafeSlice はスレッドセーフなログストレージ
type SafeSlice struct {
	content []string
	maxSize int
	mu      sync.Mutex
}

// NewSafeSlice は指定されたサイズの新しいSafeSliceを作成する
func NewSafeSlice(maxSize int) *SafeSlice {
	if maxSize <= 0 {
		maxSize = DefaultMaxLogSize
	}
	return &SafeSlice{
		content: make([]string, 0, maxSize),
		maxSize: maxSize,
	}
}

// Push は新しいログを追加
func (s *SafeSlice) Push(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 新しい値を追加
	s.content = append(s.content, message)

	// 最大サイズを超えた場合、古いものから削除（FIFO）
	if len(s.content) > s.maxSize {
		// 効率的な削除：最新のmaxSize分だけを保持
		keepCount := s.maxSize
		if keepCount > len(s.content) {
			keepCount = len(s.content)
		}

		// 新しいスライスに最新の要素をコピー
		newContent := make([]string, keepCount)
		copy(newContent, s.content[len(s.content)-keepCount:])
		s.content = newContent
	}
}

// GetRecent は最新N行を表示順で取得
// 結果: [..., 3番目に新しい, 2番目に新しい, 最新]
func (s *SafeSlice) GetRecent(lines int) []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.content) == 0 {
		return []string{}
	}

	// 要求された行数とコンテンツの長さの小さい方を使用
	actualLines := lines
	if actualLines > len(s.content) {
		actualLines = len(s.content)
	}

	// 最新のactualLines行を取得（表示順）
	startIndex := len(s.content) - actualLines
	result := make([]string, actualLines)
	copy(result, s.content[startIndex:])

	return result
}

// GetHistory は全履歴を表示順で取得
// 結果: [最古, ..., 2番目に新しい, 最新]
func (s *SafeSlice) GetHistory() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]string, len(s.content))
	copy(result, s.content)

	return result
}

// Clear は全ログを削除
func (s *SafeSlice) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.content = []string{}
}

// Count は現在のログ行数
func (s *SafeSlice) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.content)
}

// MaxHistory は履歴の最大保持行数
func (s *SafeSlice) MaxHistory() int {
	return s.maxSize
}
