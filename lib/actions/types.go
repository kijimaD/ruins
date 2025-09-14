package actions

import "fmt"

// ActionID はゲーム内のアクション種別を表す
// CDDAのaction_id enumを参考にした設計
type ActionID int

const (
	// ActionNull は無効なアクションを表す
	ActionNull ActionID = iota
	// ActionMove は隣接するタイルに移動するアクションを表す
	ActionMove
	// ActionWait は何もせずに時間を過ごすアクションを表す
	ActionWait
	// ActionAttack は敵を攻撃するアクションを表す
	ActionAttack
	// ActionPickupItem はアイテムを拾得するアクションを表す
	ActionPickupItem
	// ActionWarp はワープホールを使用するアクションを表す
	ActionWarp
)

// ActionInfo はアクションのメタデータを保持する
type ActionInfo struct {
	ID            ActionID // アクションID
	Name          string   // 表示名
	Description   string   // 説明文
	RequiresTurn  bool     // ターン消費が必要か
	MoveCost      int      // 移動コスト（ターンベース戦闘用）
	Interruptable bool     // 中断可能か（継続アクション用）
}

// actionInfos はすべてのアクション情報を保持するテーブル
var actionInfos = map[ActionID]ActionInfo{
	ActionNull:       {ActionNull, "", "無効なアクション", false, 0, false},
	ActionMove:       {ActionMove, "移動", "隣接するタイルに移動する", true, 100, false},
	ActionWait:       {ActionWait, "待機", "何もせずに時間を過ごす", true, 100, false},
	ActionAttack:     {ActionAttack, "攻撃", "敵を攻撃する", true, 100, false},
	ActionPickupItem: {ActionPickupItem, "拾得", "アイテムを拾得する", true, 100, false},
	ActionWarp:       {ActionWarp, "ワープ", "ワープホールを使用する", true, 100, false},
}

// GetActionInfo は指定されたアクションの情報を取得する
func GetActionInfo(id ActionID) ActionInfo {
	if info, exists := actionInfos[id]; exists {
		return info
	}
	return actionInfos[ActionNull]
}

// String はActionIDの文字列表現を返す
func (id ActionID) String() string {
	switch id {
	case ActionMove:
		return "Move"
	case ActionWait:
		return "Wait"
	case ActionAttack:
		return "Attack"
	case ActionPickupItem:
		return "PickupItem"
	case ActionWarp:
		return "Warp"
	case ActionNull:
		return "Null"
	default:
		return fmt.Sprintf("ActionID(%d)", int(id))
	}
}

// RequiresTurn はアクションがターン消費を必要とするかを返す
func (id ActionID) RequiresTurn() bool {
	return GetActionInfo(id).RequiresTurn
}

// MoveCost はアクションの移動コストを返す
func (id ActionID) MoveCost() int {
	return GetActionInfo(id).MoveCost
}

// IsInterruptable はアクションが中断可能かを返す
func (id ActionID) IsInterruptable() bool {
	return GetActionInfo(id).Interruptable
}

// GetAllActions は定義されたすべてのアクションIDを返す
func GetAllActions() []ActionID {
	return []ActionID{ActionMove, ActionWait, ActionAttack}
}
