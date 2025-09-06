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

// SafeSlice はスレッドセーフなスライス
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

// Append はログを追加する。最大サイズを超えた場合は古いものから削除する
func (s *SafeSlice) Append(value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 新しい値を追加
	s.content = append(s.content, value)

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

// Get は古い順に取り出す。副作用はない
func (s *SafeSlice) Get() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	copiedSlice := make([]string, len(s.content))
	copy(copiedSlice, s.content)

	return copiedSlice
}

// Pop は古い順に取り出す。取得した分は消える
func (s *SafeSlice) Pop() []string {
	result := s.Get()
	s.Flush()

	return result
}

// Flush はログの内容を消す
func (s *SafeSlice) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.content = []string{}
}
