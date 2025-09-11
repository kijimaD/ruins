// Package typewriter provides a simple typewriter-style text display system.
//
// このパッケージは、RPG風のメッセージ表示に必要な基本機能を提供します：
// - 1文字ずつの表示（タイプライター効果）
// - スキップ機能
// - カスタマイズ可能な表示速度
// - UI作成機能（MessageUIBuilder）
//
// 使用例（基本的な使い方）:
//
//	config := typewriter.DialogConfig()
//	tw := typewriter.New(config)
//
//	tw.OnComplete(func() {
//	    fmt.Println("表示完了")
//	})
//
// 使用例（UI付きの使い方）:
//
//	// MessageHandlerとUIBuilderを使用
//	handler := typewriter.NewMessageHandler(typewriter.DialogConfig(), keyboardInput)
//	uiConfig := typewriter.DefaultUIConfig()
//	uiConfig.TextFace = yourTextFace
//	uiConfig.TextColor = yourTextColor
//	uiBuilder := typewriter.NewMessageUIBuilder(handler, uiConfig)
//
//	// 更新ループ内で
//	handler.Update()
//	uiBuilder.Update()
//
//	// 描画
//	uiBuilder.GetUI().Draw(screen)
//
//	tw.Start("こんにちは、世界！")
//
//	for tw.IsTyping() {
//	    tw.Update()
//	    time.Sleep(16 * time.Millisecond) // 60FPS
//	}
package typewriter
