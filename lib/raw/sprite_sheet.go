package raw

type SpriteSheet struct {
	Name string
	// 戦闘時の立ち絵
	BattleBody *Image `toml:"battle_body"`
}

// resource の sprite sheet から画像を特定するために必要な情報
type Image struct {
	SheetName   string
	SheetNumber *int
}
