package components

import ecs "github.com/x-hgg-x/goecs/v2"

type BattleCommand struct {
	// 攻撃主体(死んだらこの攻撃は実行されない}
	Owner ecs.Entity
	// 攻撃対象
	Target ecs.Entity
	// 攻撃方法
	Way Card
}
