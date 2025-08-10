package resources

// StateEvent はフィールド上でのイベント。ステート遷移が発生する
type StateEvent string

const (
	// StateEventNone はイベントなしを表す
	StateEventNone = StateEvent("NONE")
	// StateEventWarpNext は次の階層への移動を表す
	StateEventWarpNext = StateEvent("WARP_NEXT")
	// StateEventWarpEscape は脱出を表す
	StateEventWarpEscape = StateEvent("WARP_ESCAPE")
	// StateEventBattleStart は戦闘開始を表す
	StateEventBattleStart = StateEvent("BATTLE_START")
)
