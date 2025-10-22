package components

import "fmt"

// TriggerType はトリガーの種類を表す
type TriggerType string

const (
	// TriggerTypeWarp はワープホール
	TriggerTypeWarp = TriggerType("WARP")
	// TriggerTypeDoor はドア
	TriggerTypeDoor = TriggerType("DOOR")
)

// TriggerData はトリガーのデータインターフェース
type TriggerData interface {
	TriggerType() TriggerType
}

// Trigger は接触で発動するイベント
type Trigger struct {
	Detail          TriggerData
	ActivationRange ActivationRange // 発動範囲
	ActivationMode  ActivationMode  // 発動方式
}

// WarpNextTrigger は次の階層へワープするトリガー
type WarpNextTrigger struct{}

// TriggerType はトリガータイプを返す
func (t WarpNextTrigger) TriggerType() TriggerType {
	return TriggerTypeWarp
}

// WarpEscapeTrigger は脱出ワープするトリガー
type WarpEscapeTrigger struct{}

// TriggerType はトリガータイプを返す
func (t WarpEscapeTrigger) TriggerType() TriggerType {
	return TriggerTypeWarp
}

// DoorTrigger はドアのトリガー
type DoorTrigger struct{}

// TriggerType はトリガータイプを返す
func (t DoorTrigger) TriggerType() TriggerType {
	return TriggerTypeDoor
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

// ActivationMode はトリガーの発動方式を表す
type ActivationMode string

const (
	// ActivationModeAuto は自動発動（接触時に即座に発動）
	ActivationModeAuto ActivationMode = "AUTO"
	// ActivationModeManual は手動発動（Enterキーなどで発動）
	ActivationModeManual ActivationMode = "MANUAL"
)

// Valid はActivationModeの値が有効かを検証する
func (enum ActivationMode) Valid() error {
	switch enum {
	case ActivationModeAuto, ActivationModeManual:
		return nil
	}
	return fmt.Errorf("get %s: %w", enum, ErrInvalidEnumType)
}
