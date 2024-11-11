package gamelog

import (
	"sync"

	"github.com/kijimaD/ruins/lib/utils/mathutil"
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

// 新しい順にログを取り出す。副作用はない
func (s *SafeSlice) Latest(num int) []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	copiedSlice := make([]string, len(s.content))
	l := int(mathutil.Min(len(s.content), num))
	copy(copiedSlice, s.content)

	return copiedSlice[len(s.content)-l:]
}

// ログの内容を消す
func (s *SafeSlice) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.content = []string{}
}
