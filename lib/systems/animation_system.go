package systems

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// グローバルアニメーションカウンタ
// TODO(kijima): ほかでも使うようになったらworldで保持したほうがいいかもしれない
var globalAnimationCounter int64

const (
	// AnimationFrameInterval はアニメーションフレーム切替間隔（フレーム数）
	// 30フレームごとに切り替え（60FPSで0.5秒）
	AnimationFrameInterval = 30
)

// AnimationSystem は全エンティティのスプライトアニメーションを更新する
// グローバルカウンタを使用して、全エンティティが同期してアニメーションする
func AnimationSystem(world w.World) {
	globalAnimationCounter++

	world.Manager.Join(
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		spriteRender := world.Components.SpriteRender.Get(entity).(*gc.SpriteRender)

		// AnimKeysが空ならアニメーションなし
		if len(spriteRender.AnimKeys) == 0 {
			return
		}

		// フレームインデックスを計算
		frameIndex := (globalAnimationCounter / AnimationFrameInterval) % int64(len(spriteRender.AnimKeys))

		// SpriteKeyを更新
		spriteRender.SpriteKey = spriteRender.AnimKeys[frameIndex]
	}))
}
