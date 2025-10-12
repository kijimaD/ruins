// Package testutil はテストで共通して使用するユーティリティ関数を提供する。
//
// # Overview
//
// testutilパッケージは、テスト全体で共有される便利な関数やヘルパーを提供します。
// 主に、並行テスト実行時の競合状態を避けるためのユーティリティを含みます。
//
// # 並行テストにおけるWorld初期化
//
// maingame.InitWorld()は共有リソース（フォント、スプライトシート、Rawsなど）を読み込むため、
// 並行テスト実行時に競合状態（concurrent map writes）が発生する可能性があります。
//
// この問題を解決するため、InitTestWorld()関数は内部でmutexを使ってInitWorld()の呼び出しを
// 排他制御します。これにより、テストはt.Parallel()を使って並行実行できますが、
// InitWorld()自体は順次実行されるため、競合状態を回避できます。
//
// # 使用例
//
// 並行テストでWorldを初期化する場合:
//
//	func TestSomething(t *testing.T) {
//	    t.Parallel() // 並行実行可能
//	    world := testutil.InitTestWorld(t)
//
//	    // テストコード
//	}
//
// # 設計上の利点
//
//   - テストの並行実行を維持（テスト実行速度の最適化）
//   - 共有リソースへのアクセスを保護（競合状態の回避）
//   - 各テストが完全に独立したWorldインスタンスを取得
//   - 既存のテストコードへの影響を最小限に抑える
//
// # 歴史的背景
//
// 当初、equip_menu_sort_test.goで"fatal error: concurrent map writes"が発生しました。
// これは、複数のテストが並行してInitWorld()を呼び出すことで、内部のリソースローダーが
// 同時にマップに書き込むことが原因でした。
//
// 最初の修正案としてt.Parallel()を削除して順次実行にしましたが、テスト速度が低下するため、
// より良い解決策として、共通のヘルパー関数InitTestWorld()を作成し、mutex保護を行いました。
//
// これにより、以下のテストファイルが安全に並行実行できるようになりました：
//   - lib/states/equip_menu_sort_test.go
//   - lib/states/craft_menu_sort_test.go
//   - その他、InitWorldを使用する全てのテスト
package testutil
