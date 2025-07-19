package effects

import (
	"fmt"

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
	logger Logger            // ログ出力
}

// NewProcessor は新しいプロセッサーを作成する
func NewProcessor() *Processor {
	return &Processor{
		queue:  make([]EffectExecution, 0),
		logger: defaultLogger{},
	}
}

// SetLogger はカスタムログ出力を設定する
func (p *Processor) SetLogger(logger Logger) {
	p.logger = logger
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

	p.logger.Debug("エフェクトをキューに追加: %s (対象: %d体)", effect, len(targets))
}


// Execute はキュー内のすべてのエフェクトを順次実行する
func (p *Processor) Execute(world w.World) error {
	p.logger.Debug("エフェクトキュー実行開始 (キューサイズ: %d)", len(p.queue))

	executed := 0
	for len(p.queue) > 0 {
		execution := p.queue[0]
		p.queue = p.queue[1:]

		// エフェクトの妥当性を検証
		if err := execution.Effect.Validate(world, execution.Context); err != nil {
			p.logger.Error("エフェクト検証失敗: %s - %v", execution.Effect, err)
			return fmt.Errorf("エフェクト検証失敗 %s: %w", execution.Effect, err)
		}

		// エフェクトを実行
		p.logger.Debug("エフェクト実行: %s", execution.Effect)
		err := execution.Effect.Apply(world, execution.Context)

		if err != nil {
			p.logger.Error("エフェクト実行失敗: %s - %v", execution.Effect, err)
			return fmt.Errorf("エフェクト実行失敗 %s: %w", execution.Effect, err)
		}

		executed++
	}

	p.logger.Debug("エフェクトキュー実行完了 (実行数: %d)", executed)
	return nil
}

// Clear はキューをクリアする（テストやリセット用）
func (p *Processor) Clear() {
	cleared := len(p.queue)
	p.queue = p.queue[:0]
	p.logger.Debug("エフェクトキューをクリア (クリア数: %d)", cleared)
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
