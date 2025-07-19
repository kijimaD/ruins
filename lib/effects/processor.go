package effects

import (
	"fmt"

	"github.com/kijimaD/ruins/lib/logger"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// EffectExecution はキューに蓄積される実行単位
type EffectExecution struct {
	Effect  Effect   // 実行する効果
	Context *Context // 実行コンテキスト
}

// Processor はエフェクトの実行管理を行うプロセッサー
type Processor struct {
	queue  []EffectExecution // エフェクト実行キュー
	logger *logger.Logger    // ログ出力
}

// NewProcessor は新しいプロセッサーを作成する
func NewProcessor() *Processor {
	return &Processor{
		queue:  make([]EffectExecution, 0),
		logger: logger.New(logger.CategoryEffect),
	}
}


// AddEffect はエフェクトをキューに追加する
func (p *Processor) AddEffect(effect Effect, creator *ecs.Entity, targets ...ecs.Entity) {
	ctx := &Context{
		Creator: creator,
		Targets: targets,
	}

	p.queue = append(p.queue, EffectExecution{
		Effect:  effect,
		Context: ctx,
	})

	p.logger.Debug("エフェクトをキューに追加", "effect", effect.String(), "targets", len(targets))
}


// Execute はキュー内のすべてのエフェクトを順次実行する
func (p *Processor) Execute(world w.World) error {
	p.logger.Debug("エフェクトキュー実行開始", "queue_size", len(p.queue))

	executed := 0
	for len(p.queue) > 0 {
		execution := p.queue[0]
		p.queue = p.queue[1:]

		// エフェクトの妥当性を検証
		if err := execution.Effect.Validate(world, execution.Context); err != nil {
			p.logger.Error("エフェクト検証失敗", "effect", execution.Effect.String(), "error", err)
			return fmt.Errorf("エフェクト検証失敗 %s: %w", execution.Effect, err)
		}

		// エフェクトを実行
		p.logger.Debug("エフェクト実行", "effect", execution.Effect.String())
		err := execution.Effect.Apply(world, execution.Context)

		if err != nil {
			p.logger.Error("エフェクト実行失敗", "effect", execution.Effect.String(), "error", err)
			return fmt.Errorf("エフェクト実行失敗 %s: %w", execution.Effect, err)
		}

		executed++
	}

	p.logger.Debug("エフェクトキュー実行完了", "executed", executed)
	return nil
}

// Clear はキューをクリアする（テストやリセット用）
func (p *Processor) Clear() {
	cleared := len(p.queue)
	p.queue = p.queue[:0]
	p.logger.Debug("エフェクトキューをクリア", "cleared", cleared)
}

// QueueSize は現在のキューサイズを返す
func (p *Processor) QueueSize() int {
	return len(p.queue)
}

// IsEmpty はキューが空かどうかを判定する
func (p *Processor) IsEmpty() bool {
	return len(p.queue) == 0
}

// QueuedEffects はキュー内のエフェクト一覧を返す（デバッグ用）
func (p *Processor) QueuedEffects() []Effect {
	effects := make([]Effect, len(p.queue))
	for i, execution := range p.queue {
		effects[i] = execution.Effect
	}
	return effects
}
