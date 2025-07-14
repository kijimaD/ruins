package components

// AIRoamingSubState はAI徘徊行動のサブ状態を表す
type AIRoamingSubState string

const (
	// AIRoamingWaiting はAI徘徊における待機状態
	AIRoamingWaiting = AIRoamingSubState("WAIT")
	// AIRoamingDriving はAI徘徊における移動状態
	AIRoamingDriving = AIRoamingSubState("DRIVING")
)
