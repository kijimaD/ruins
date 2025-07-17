// Package logger はゲーム内の構造化ログを提供する
//
// 機能別にカテゴリを分けて、デバッグしやすいログ出力を実現する。
// 環境変数でログレベルやカテゴリ別の出力制御が可能。
//
// 使用例:
//
//	var log = logger.New(logger.CategoryBattle)
//	log.Info("戦闘開始", "enemy", "スライム", "hp", 100)
//	log.Debug("ダメージ計算", "base", 10, "bonus", 5)
//
// 環境変数:
//   - LOG_LEVEL: デフォルトログレベル (debug, info, warn, error, fatal)
//   - LOG_CATEGORIES: カテゴリ別ログレベル (例: "battle=debug,render=warn")
package logger
