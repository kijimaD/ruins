package components

import ecs "github.com/x-hgg-x/goecs/v2"

// BattleCommand は戦闘における行動コマンドを表す
type BattleCommand struct {
	// 行動主体(死んだらこの攻撃は実行されない}
	Owner ecs.Entity
	// 行動対象
	Target ecs.Entity
	// 行動方法(カードEntity)
	Way ecs.Entity
}
