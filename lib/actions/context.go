package actions

import (
	gc "github.com/kijimaD/ruins/lib/components"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// Context はアクション実行時のコンテキストを表す
type Context struct {
	Actor  ecs.Entity   // アクションを実行するエンティティ
	Target *ecs.Entity  // 対象エンティティ（nilの場合もある）
	Dest   *gc.Position // アクションの対象位置
}

// Result はアクション実行結果を表す
type Result struct {
	Success  bool     // 実行成功/失敗
	ActionID ActionID // 実行されたアクション
	Message  string   // 結果メッセージ
}
