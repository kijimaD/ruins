package gamelog

import (
	"sync"
)

var (
	// 戦闘用
	BattleLog SafeSlice
	// フィールド用
	FieldLog SafeSlice
	// 会話シーンでステータス変化を通知する用
	SceneLog SafeSlice
)

// TODO: 無限に追加される可能性があるので、最大の長さを設定する
type SafeSlice struct {
	content []string
	mu      sync.Mutex
}

// ログを追加する
func (s *SafeSlice) Append(value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.content = append(s.content, value)
}

// 古い順に取り出す。副作用はない
func (s *SafeSlice) Get() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	copiedSlice := make([]string, len(s.content))
	copy(copiedSlice, s.content)

	return copiedSlice
}

// 古い順に取り出す。取得した分は消える
func (s *SafeSlice) Pop() []string {
	result := s.Get()
	s.Flush()

	return result
}

// ログの内容を消す
func (s *SafeSlice) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.content = []string{}
}
