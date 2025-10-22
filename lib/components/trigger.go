package components

// TriggerType はトリガーの種類を表す
type TriggerType string

const (
	// TriggerTypeWarp はワープホール
	TriggerTypeWarp = TriggerType("WARP")
)

// TriggerData はトリガーのデータインターフェース
type TriggerData interface {
	TriggerType() TriggerType
}

// Trigger は接触で発動するイベント
type Trigger struct {
	Detail      TriggerData
	AutoExecute bool // 接触時に自動実行するか（false=手動でEnterキー必要）
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
