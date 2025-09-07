package gamelog

import "image/color"

// LogFragment は色付きテキストの断片
type LogFragment struct {
	Color color.RGBA `json:"color"`
	Text  string     `json:"text"`
}

// LogKind はログの種類
type LogKind int

const (
	LogKindField LogKind = iota
	LogKindBattle
	LogKindScene
)

// String はLogKindの文字列表現を返す
func (lk LogKind) String() string {
	switch lk {
	case LogKindField:
		return "Field"
	case LogKindBattle:
		return "Battle"
	case LogKindScene:
		return "Scene"
	default:
		return "Unknown"
	}
}
