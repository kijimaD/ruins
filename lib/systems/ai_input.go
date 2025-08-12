package systems

import (
	"math"
	"math/rand/v2"
	"time"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// AIInputSystem は AI制御されたエンティティの入力処理を行う
func AIInputSystem(world w.World) {
	// まずプレイヤーのエンティティと位置を取得
	var playerPos *gc.Position
	world.Manager.Join(
		world.Components.Operator,
		world.Components.Position,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerPos = world.Components.Position.Get(entity).(*gc.Position)
	}))

	// AI制御エンティティの処理
	world.Manager.Join(
		world.Components.Velocity,
		world.Components.Position,
		world.Components.AIMoveFSM,
		world.Components.AIRoaming,
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		roaming := world.Components.AIRoaming.Get(entity).(*gc.AIRoaming)
		velocity := world.Components.Velocity.Get(entity).(*gc.Velocity)
		position := world.Components.Position.Get(entity).(*gc.Position)

		// 視界コンポーネントがある場合、プレイヤーとの距離をチェック
		if entity.HasComponent(world.Components.AIVision) && playerPos != nil {
			vision := world.Components.AIVision.Get(entity).(*gc.AIVision)

			// プレイヤーとの距離を計算
			dx := float64(playerPos.X - position.X)
			dy := float64(playerPos.Y - position.Y)
			distance := math.Sqrt(dx*dx + dy*dy)

			// 視界内にプレイヤーがいる場合
			if distance <= vision.ViewDistance {
				// 追跡状態に移行
				roaming.SubState = gc.AIRoamingChasing
				roaming.StartSubState = time.Now()
				roaming.DurationSubState = time.Second * 2 // 追跡時間

				// 追跡ターゲット情報を更新または作成
				if entity.HasComponent(world.Components.AIChasing) {
					chasing := world.Components.AIChasing.Get(entity).(*gc.AIChasing)
					chasing.TargetX = float64(playerPos.X)
					chasing.TargetY = float64(playerPos.Y)
					chasing.LastSeen = time.Now()
				} else {
					// AIChasingコンポーネントを新規作成
					world.Components.AIChasing.Set(entity, &gc.AIChasing{
						TargetX:  float64(playerPos.X),
						TargetY:  float64(playerPos.Y),
						LastSeen: time.Now(),
					})
				}

				// プレイヤーへの角度を計算
				angle := math.Atan2(dy, dx) * 180 / math.Pi
				velocity.Angle = angle
				velocity.ThrottleMode = gc.ThrottleModeFront
				return
			}
		}

		// 追跡状態の処理
		if roaming.SubState == gc.AIRoamingChasing {
			if entity.HasComponent(world.Components.AIChasing) {
				chasing := world.Components.AIChasing.Get(entity).(*gc.AIChasing)

				// 最後の視認から時間が経過したら通常の徘徊に戻る
				if time.Since(chasing.LastSeen).Seconds() > 3 {
					roaming.SubState = gc.AIRoamingWaiting
					roaming.StartSubState = time.Now()
					roaming.DurationSubState = time.Second * time.Duration(rand.IntN(3))
				} else {
					// ターゲット位置へ向かう
					dx := chasing.TargetX - float64(position.X)
					dy := chasing.TargetY - float64(position.Y)
					angle := math.Atan2(dy, dx) * 180 / math.Pi
					velocity.Angle = angle
					velocity.ThrottleMode = gc.ThrottleModeFront
					return
				}
			}
		}

		// 通常の徘徊処理
		if time.Since(roaming.StartSubState).Seconds() > roaming.DurationSubState.Seconds() {
			roaming.StartSubState = time.Now()
			roaming.DurationSubState = time.Second * time.Duration(rand.IntN(3))

			var subState gc.AIRoamingSubState
			switch rand.IntN(2) {
			case 0:
				subState = gc.AIRoamingWaiting
			case 1:
				subState = gc.AIRoamingDriving
			}
			roaming.SubState = subState

			switch subState {
			case gc.AIRoamingWaiting:
				velocity.ThrottleMode = gc.ThrottleModeNope
				velocity.Angle += float64(rand.IntN(91))
			case gc.AIRoamingDriving:
				velocity.ThrottleMode = gc.ThrottleModeFront
			}
		}
	}))
}
