package typewriter

import "time"

// Config はタイプライターの動作設定
type Config struct {
	CharDelay        time.Duration // 基本的な1文字あたりの待機時間
	PunctuationDelay time.Duration // 句読点後の追加待機時間
	CommaDelay       time.Duration // カンマ後の追加待機時間
	NewlineDelay     time.Duration // 改行後の追加待機時間

	// 特殊設定
	SkipEnabled       bool          // スキップ機能の有効/無効
	PauseEnabled      bool          // 一時停止機能の有効/無効
	AutoComplete      bool          // 自動完了（時間経過で完了）
	AutoCompleteDelay time.Duration // 自動完了までの時間
}

// FastConfig は高速表示用設定
func FastConfig() Config {
	return Config{
		CharDelay:        30 * time.Millisecond,
		PunctuationDelay: 100 * time.Millisecond,
		CommaDelay:       50 * time.Millisecond,
		NewlineDelay:     200 * time.Millisecond,
		SkipEnabled:      true,
		PauseEnabled:     false,
		AutoComplete:     false,
	}
}

// SlowConfig はゆっくり表示用設定
func SlowConfig() Config {
	return Config{
		CharDelay:         150 * time.Millisecond,
		PunctuationDelay:  600 * time.Millisecond,
		CommaDelay:        300 * time.Millisecond,
		NewlineDelay:      800 * time.Millisecond,
		SkipEnabled:       true,
		PauseEnabled:      true,
		AutoComplete:      true,
		AutoCompleteDelay: 10 * time.Second,
	}
}

// DialogConfig は会話用設定
func DialogConfig() Config {
	return Config{
		CharDelay:         60 * time.Millisecond,
		PunctuationDelay:  400 * time.Millisecond,
		CommaDelay:        200 * time.Millisecond,
		NewlineDelay:      600 * time.Millisecond,
		SkipEnabled:       true,
		PauseEnabled:      true,
		AutoComplete:      true,
		AutoCompleteDelay: 8 * time.Second,
	}
}
