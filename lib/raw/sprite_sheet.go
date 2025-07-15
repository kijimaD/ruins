package raw

// SpriteSheet はスプライトシートの情報を保持する
type SpriteSheet struct {
	Name string
	// 戦闘時の立ち絵
	BattleBody *Image `toml:"battle_body"`
}

// Image はresourceのspriteシートから画像を特定するために必要な情報
type Image struct {
	SheetName   string
	SheetNumber *int
}
