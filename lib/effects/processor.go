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

// ProcessorHooks はプロセッサーの実行時フックを定義する
type ProcessorHooks struct {
	// BeforeEffect は各エフェクト実行前に呼ばれる
	BeforeEffect func(Effect, *Context)

	// AfterEffect は各エフェクト実行後に呼ばれる（エラーの有無に関わらず）
	AfterEffect func(Effect, *Context, error)
}

// Processor はエフェクトの実行管理を行うプロセッサー
type Processor struct {
	queue  []EffectExecution // エフェクト実行キュー
	hooks  ProcessorHooks    // 実行時フック
	logger Logger            // ログ出力
}

// NewProcessor は新しいプロセッサーを作成する
func NewProcessor() *Processor {
	return &Processor{
		queue:  make([]EffectExecution, 0),
		hooks:  ProcessorHooks{},
		logger: defaultLogger{},
	}
}

// SetHooks はプロセッサーのフックを設定する
func (p *Processor) SetHooks(hooks ProcessorHooks) {
	p.hooks = hooks
}

// SetLogger はカスタムログ出力を設定する
func (p *Processor) SetLogger(logger Logger) {
	p.logger = logger
}

// AddEffect はエフェクトをキューに追加する
func (p *Processor) AddEffect(effect Effect, creator *ecs.Entity, targets ...ecs.Entity) error {
	ctx := &Context{
		Creator: creator,
		Targets: targets,
	}

	p.queue = append(p.queue, EffectExecution{
		Effect:  effect,
		Context: ctx,
	})

	p.logger.Debug("エフェクトをキューに追加: %s (対象: %d体)", effect, len(targets))
	return nil
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

		// 実行前フック
		if p.hooks.BeforeEffect != nil {
			p.hooks.BeforeEffect(execution.Effect, execution.Context)
		}

		// エフェクトを実行
		p.logger.Debug("エフェクト実行: %s", execution.Effect)
		err := execution.Effect.Apply(world, execution.Context)

		// 実行後フック（エラーの有無に関わらず呼ぶ）
		if p.hooks.AfterEffect != nil {
			p.hooks.AfterEffect(execution.Effect, execution.Context, err)
		}

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
