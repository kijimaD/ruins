package messagedata

// DialogueTable は会話データのテーブル
var DialogueTable = map[string]func(speakerName string) *MessageData{
	"old_soldier_greeting": func(speakerName string) *MessageData {
		// 1ページ目
		msg1 := NewDialogMessage("", speakerName).
			AddText("「あんた、").
			AddKeyword("遺跡").
			AddText("の").
			AddKeyword("珠狙い").
			AddText("だろ?").
			AddNewLine().
			AddText("外からこの街に来る異様に若い連中はみんなそうさ。").
			AddNewLine().
			AddText("向こう見ずで破滅的で、...").
			AddNewLine().
			AddText("どうしようもない事情を持ってる。」")

		// 2ページ目
		msg2 := NewDialogMessage("", speakerName).
			AddText("「あんたは...、そうか、").
			AddText("母親が...。").
			AddNewLine().
			AddText("言っちゃ悪いが、そういう奴らはここでは珍しくない。").
			AddNewLine().
			AddText("どんな事情があるにせよ、").
			AddKeyword("遺跡").
			AddText("で辿る結末は1つさ。」")

		msg1.NextMessages = append(msg1.NextMessages, msg2)
		return msg1
	},
}

// GetDialogue は指定されたキーに対応する会話データを取得する
func GetDialogue(key string, speakerName string) *MessageData {
	if dialogueFunc, ok := DialogueTable[key]; ok {
		return dialogueFunc(speakerName)
	}
	// デフォルトメッセージ
	return NewDialogMessage("...", speakerName)
}
