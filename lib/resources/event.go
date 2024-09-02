package resources

// フィールド上でのイベント
type StateEvent string

const (
	StateEventNone       = StateEvent("NONE")
	StateEventWarpNext   = StateEvent("WARP_NEXT")
	StateEventWarpEscape = StateEvent("WARP_ESCAPE")
)
