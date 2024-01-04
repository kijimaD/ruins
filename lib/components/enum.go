package components

type warpMode string

var (
	WarpModeNext   = warpMode("NEXT")
	WarpModeEscape = warpMode("ESCAPE")
)

type TargetFactionType string

var (
	TargetFactionAlly  = TargetFactionType("ALLY")  // 味方
	TargetFactionEnemy = TargetFactionType("ENEMY") //  敵
	TargetFactionNone  = TargetFactionType("NONE")  // なし
)

type UsableSceneType string

var (
	UsableSceneBattle = UsableSceneType("BATTLE") // 戦闘
	UsableSceneField  = UsableSceneType("FIELD")  // フィールド
	UsableSceneAny    = UsableSceneType("ANY")    // いつでも
)
