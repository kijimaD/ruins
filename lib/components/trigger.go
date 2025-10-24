package components

import "fmt"

// Trigger は接触で発動するイベント
type Trigger struct {
	Data TriggerData
}

// TriggerConfig はトリガーの設定
type TriggerConfig struct {
	ActivationRange ActivationRange // 発動範囲
	ActivationWay   ActivationWay   // 発動方式
}

// TriggerData はトリガーのデータインターフェース
type TriggerData interface {
	Config() TriggerConfig
}

// WarpNextTrigger は次の階層へワープするトリガー
type WarpNextTrigger struct{}

// Config はトリガー設定を返す
func (t WarpNextTrigger) Config() TriggerConfig {
	return TriggerConfig{
		ActivationRange: ActivationRangeSameTile,
		ActivationWay:   ActivationWayManual,
	}
}

// WarpEscapeTrigger は脱出ワープするトリガー
type WarpEscapeTrigger struct{}

// Config はトリガー設定を返す
func (t WarpEscapeTrigger) Config() TriggerConfig {
	return TriggerConfig{
		ActivationRange: ActivationRangeSameTile,
		ActivationWay:   ActivationWayManual,
	}
}

// DoorTrigger はドアのトリガー
type DoorTrigger struct{}

// Config はトリガー設定を返す
func (t DoorTrigger) Config() TriggerConfig {
	return TriggerConfig{
		ActivationRange: ActivationRangeAdjacent,
		ActivationWay:   ActivationWayOnCollision,
	}
}

// TalkTrigger は会話のトリガー
type TalkTrigger struct{}

// Config はトリガー設定を返す
func (t TalkTrigger) Config() TriggerConfig {
	return TriggerConfig{
		ActivationRange: ActivationRangeAdjacent,
		ActivationWay:   ActivationWayOnCollision,
	}
}

// ItemTrigger はアイテム拾得のトリガー
type ItemTrigger struct{}

// Config はトリガー設定を返す
func (t ItemTrigger) Config() TriggerConfig {
	return TriggerConfig{
		ActivationRange: ActivationRangeSameTile,
		ActivationWay:   ActivationWayManual,
	}
}

// ActivationRange はトリガーの発動範囲を表す
type ActivationRange string

const (
	// ActivationRangeSameTile は直上（同じタイル）で発動
	ActivationRangeSameTile ActivationRange = "SAME_TILE"
	// ActivationRangeAdjacent は隣接タイルで発動
	ActivationRangeAdjacent ActivationRange = "ADJACENT"
)

// Valid はActivationRangeの値が有効かを検証する
func (enum ActivationRange) Valid() error {
	switch enum {
	case ActivationRangeSameTile, ActivationRangeAdjacent:
		return nil
	}
	return fmt.Errorf("get %s: %w", enum, ErrInvalidEnumType)
}

// ================

// ActivationWay はトリガーの発動方式を表す
type ActivationWay string

const (
	// ActivationWayAuto は自動発動（範囲内に入ったら即座に発動）
	ActivationWayAuto ActivationWay = "AUTO"
	// ActivationWayManual は手動発動（Enterキーやアクションメニューで発動）
	ActivationWayManual ActivationWay = "MANUAL"
	// ActivationWayOnCollision は移動先衝突時に自動発動（移動先として指定された時に発動）
	ActivationWayOnCollision ActivationWay = "ON_COLLISION"
)

// Valid はActivationWayの値が有効かを検証する
func (enum ActivationWay) Valid() error {
	switch enum {
	case ActivationWayAuto, ActivationWayManual, ActivationWayOnCollision:
		return nil
	}
	return fmt.Errorf("get %s: %w", enum, ErrInvalidEnumType)
}
