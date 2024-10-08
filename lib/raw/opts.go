package raw

type SpawnType string

var (
	// どの場所にも属さない。マスタとして使う
	SpawnInNone SpawnType = "NONE"
	// バックパック内に生成する
	SpawnInBackpack SpawnType = "IN_BACKPACK"
)
