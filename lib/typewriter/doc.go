// Package typewriter provides a simple typewriter-style text display system.
//
// このパッケージは、RPG風のメッセージ表示に必要な基本機能を提供します：
// - 1文字ずつの表示（タイプライター効果）
// - スキップ機能
// - カスタマイズ可能な表示速度
//
// フォーカスしないこと:
// - UIは提供しない
//
// 使用例:
//
//	config := typewriter.BattleConfig()
//	tw := typewriter.New(config)
//
//	tw.OnComplete(func() {
//	    fmt.Println("表示完了")
//	})
//
//	tw.Start("こんにちは、世界！")
//
//	for tw.IsTyping() {
//	    tw.Update()
//	    time.Sleep(16 * time.Millisecond) // 60FPS
//	}
package typewriter
