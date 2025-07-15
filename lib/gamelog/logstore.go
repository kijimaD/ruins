package gamelog

import (
	"sync"
)

var (
	// BattleLog は戦闘用ログ
	BattleLog SafeSlice
	// FieldLog はフィールド用ログ
	FieldLog SafeSlice
	// SceneLog は会話シーンでステータス変化を通知する用ログ
	SceneLog SafeSlice
)

// SafeSlice はスレッドセーフなスライス
// TODO: 無限に追加される可能性があるので、最大の長さを設定する
type SafeSlice struct {
	content []string
	mu      sync.Mutex
}

// Append はログを追加する
func (s *SafeSlice) Append(value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.content = append(s.content, value)
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
