package states

// BaseState は共通のtransition管理を持つベース構造体
type BaseState struct {
	trans *Transition
}

// SetTransition は遷移を設定する
func (bs *BaseState) SetTransition(trans Transition) {
	bs.trans = &trans
}

// GetTransition は現在の遷移を取得する
func (bs *BaseState) GetTransition() *Transition {
	return bs.trans
}

// ClearTransition は遷移をクリアする
func (bs *BaseState) ClearTransition() {
	bs.trans = nil
}

// ConsumeTransition は遷移を消費して返す
func (bs *BaseState) ConsumeTransition() Transition {
	if bs.trans != nil {
		next := *bs.trans
		bs.trans = nil
		return next
	}
	return Transition{Type: TransNone}
}

// StateWithTransition はtransition管理機能を持つstateのインターフェース
type StateWithTransition interface {
	State
	GetTransition() *Transition
	ClearTransition()
}