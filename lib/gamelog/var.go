package gamelog

var (
	// BattleLog は戦闘用ログ
	BattleLog = NewSafeSlice(BattleLogMaxSize)
	// FieldLog はフィールド用ログ
	FieldLog = NewSafeSlice(FieldLogMaxSize)
	// SceneLog は会話シーンでステータス変化を通知する用ログ
	SceneLog = NewSafeSlice(SceneLogMaxSize)
)
