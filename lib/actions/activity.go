package actions

import (
	"fmt"
	"time"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/logger"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// ActivityType は継続的なアクション（アクティビティ）の種別を表す
// CDDAのactivity_idを参考にした設計
type ActivityType int

const (
	// ActivityNull は無効なアクティビティを表す
	ActivityNull ActivityType = iota
	// ActivityMove は移動アクティビティを表す
	ActivityMove
	// ActivityAttack は攻撃アクティビティを表す
	ActivityAttack
	// ActivityPickup はアイテム拾得アクティビティを表す
	ActivityPickup
	// ActivityWarp はワープアクティビティを表す
	ActivityWarp
	// ActivityRest は休息アクティビティを表す
	ActivityRest
	// ActivityRead は読書アクティビティを表す
	ActivityRead
	// ActivityCraft はクラフトアクティビティを表す
	ActivityCraft
	// ActivityWait は待機アクティビティを表す
	ActivityWait
)

// ActivityState はアクティビティの実行状態を表す
type ActivityState int

const (
	// ActivityStateRunning はアクティビティが実行中であることを表す
	ActivityStateRunning ActivityState = iota
	// ActivityStatePaused はアクティビティが一時停止中であることを表す
	ActivityStatePaused
	// ActivityStateCompleted はアクティビティが完了したことを表す
	ActivityStateCompleted
	// ActivityStateCanceled はアクティビティがキャンセルされたことを表す
	ActivityStateCanceled
)

// TimingMode は時間消費の計算方法を表す
// CDDAのbased_on_typeを参考
type TimingMode int

const (
	// TimingModeTime は実時間ベース（キャラクター速度に関係なく固定時間）
	TimingModeTime TimingMode = iota
	// TimingModeSpeed は速度ベース（速いキャラクターは早く完了）
	TimingModeSpeed
	// TimingModeCustom はカスタムタイミング
	TimingModeCustom
)

// ActivityInfo はアクティビティのメタデータを保持する
type ActivityInfo struct {
	Type             ActivityType // アクティビティの種別
	Name             string       // 表示名
	Description      string       // 説明文
	Interruptible    bool         // 中断可能か
	Resumable        bool         // 中断後の再開可能か
	TimingMode       TimingMode   // 時間計算方法
	ActionPointCost  int          // 1ターン毎のアクションポイントコスト
	TotalRequiredAP  int          // アクティビティ完了に必要な総AP量
	RequiresTarget   bool         // ターゲットが必要か
	RequiresPosition bool         // 位置が必要か
}

// ActivityInterface はアクティビティの実行を担当するインターフェース
// CDDAのactivity_actorを参考にした設計
type ActivityInterface interface {
	Validate(act *Activity, world w.World) error
	Start(act *Activity, world w.World) error
	DoTurn(act *Activity, world w.World) error
	Finish(act *Activity, world w.World) error
	Canceled(act *Activity, world w.World) error
}

// Activity は継続的なアクション（アクティビティ）のデータを表す
type Activity struct {
	Type         ActivityType  // アクティビティの種別
	State        ActivityState // 実行状態
	TurnsTotal   int           // 総必要ターン数
	TurnsLeft    int           // 残りターン数
	Actor        ecs.Entity    // 実行者
	Target       *ecs.Entity   // 対象エンティティ（nilの場合もある）
	Position     *gc.Position  // 対象位置（nilの場合もある）
	Message      string        // 進行状況メッセージ
	StartTime    time.Time     // 開始時刻
	PauseTime    *time.Time    // 一時停止時刻（nilは実行中）
	CancelReason string        // キャンセル理由

	Logger *logger.Logger
}

// アクティビティ情報テーブル
var activityInfos = map[ActivityType]ActivityInfo{
	ActivityNull: {
		Type:             ActivityNull,
		Name:             "",
		Description:      "無効なアクティビティ",
		Interruptible:    false,
		Resumable:        false,
		TimingMode:       TimingModeTime,
		ActionPointCost:  0,
		TotalRequiredAP:  0,
		RequiresTarget:   false,
		RequiresPosition: false,
	},
	ActivityMove: {
		Type:             ActivityMove,
		Name:             "移動",
		Description:      "隣接するタイルに移動する",
		Interruptible:    false,
		Resumable:        false,
		TimingMode:       TimingModeSpeed,
		ActionPointCost:  100, // 初期AP相当（基本アクション）
		TotalRequiredAP:  100,
		RequiresTarget:   false,
		RequiresPosition: true, // 移動先が必要
	},
	ActivityAttack: {
		Type:             ActivityAttack,
		Name:             "攻撃",
		Description:      "敵を攻撃する",
		Interruptible:    false,
		Resumable:        false,
		TimingMode:       TimingModeSpeed,
		ActionPointCost:  100, // 初期AP相当（基本アクション）
		TotalRequiredAP:  100,
		RequiresTarget:   true, // 攻撃対象が必要
		RequiresPosition: false,
	},
	ActivityPickup: {
		Type:             ActivityPickup,
		Name:             "拾得",
		Description:      "アイテムを拾得する",
		Interruptible:    false,
		Resumable:        false,
		TimingMode:       TimingModeSpeed,
		ActionPointCost:  50, // 初期APの半分（素早いアクション）
		TotalRequiredAP:  50,
		RequiresTarget:   false,
		RequiresPosition: false,
	},
	ActivityWarp: {
		Type:             ActivityWarp,
		Name:             "ワープ",
		Description:      "ワープホールを使用する",
		Interruptible:    false,
		Resumable:        false,
		TimingMode:       TimingModeTime,
		ActionPointCost:  0, // 時間を消費しない（瞬間移動）
		TotalRequiredAP:  0,
		RequiresTarget:   false,
		RequiresPosition: false,
	},
	ActivityRest: {
		Type:             ActivityRest,
		Name:             "休息",
		Description:      "体力を回復するために休息する",
		Interruptible:    true,
		Resumable:        true,
		TimingMode:       TimingModeTime,
		ActionPointCost:  100,  // 初期AP相当（継続アクション毎ターン）
		TotalRequiredAP:  1000, // AP100のプレイヤーで10ターン
		RequiresTarget:   false,
		RequiresPosition: false,
	},
	ActivityRead: {
		Type:             ActivityRead,
		Name:             "読書",
		Description:      "本を読んでスキルを習得する",
		Interruptible:    true,
		Resumable:        true,
		TimingMode:       TimingModeTime,
		ActionPointCost:  100,  // 初期AP相当（継続アクション毎ターン）
		TotalRequiredAP:  2000, // AP100のプレイヤーで20ターン
		RequiresTarget:   true, // 本が対象
		RequiresPosition: false,
	},
	ActivityCraft: {
		Type:             ActivityCraft,
		Name:             "クラフト",
		Description:      "アイテムを作成する",
		Interruptible:    true,
		Resumable:        false, // 一度中断すると材料が無駄になる
		TimingMode:       TimingModeSpeed,
		ActionPointCost:  100,  // 初期AP相当（継続アクション毎ターン）
		TotalRequiredAP:  1500, // AP100のプレイヤーで15ターン
		RequiresTarget:   false,
		RequiresPosition: true, // 作業場所が必要
	},
	ActivityWait: {
		Type:             ActivityWait,
		Name:             "待機",
		Description:      "指定した時間だけ待機する",
		Interruptible:    true,
		Resumable:        true,
		TimingMode:       TimingModeTime,
		ActionPointCost:  100, // 初期AP相当（意図的な時間消費）
		TotalRequiredAP:  500, // AP100のプレイヤーで5ターン
		RequiresTarget:   false,
		RequiresPosition: false,
	},
}

// アクティビティアクターのレジストリ
var activityActors = map[ActivityType]ActivityInterface{
	// 初期化は各アクティビティファイルで行う
}

// RegisterActivityActor はアクティビティアクターを登録する
func RegisterActivityActor(activityType ActivityType, actor ActivityInterface) {
	activityActors[activityType] = actor
}

// GetActivityActor は指定されたアクティビティのアクターを取得する
func GetActivityActor(activityType ActivityType) ActivityInterface {
	return activityActors[activityType]
}

// GetActivityInfo は指定されたアクティビティの情報を取得する
func GetActivityInfo(activityType ActivityType) ActivityInfo {
	if info, exists := activityInfos[activityType]; exists {
		return info
	}
	// 未知のアクティビティに対してはNullActivityを返す
	return activityInfos[ActivityNull]
}

// NewActivity は新しいアクティビティを作成する
func NewActivity(activityType ActivityType, actor ecs.Entity, duration int) *Activity {
	// durationは0以下の値は許可しない（呼び出し側で適切な値を指定する必要がある）
	if duration <= 0 {
		duration = 1 // 最低1ターンは必要
	}

	return &Activity{
		Type:       activityType,
		State:      ActivityStateRunning,
		TurnsTotal: duration,
		TurnsLeft:  duration,
		Actor:      actor,
		StartTime:  time.Now(),
		Logger:     logger.New(logger.CategoryAction),
	}
}

// CalculateRequiredTurns はキャラクターのAP量に基づいて必要ターン数を計算する
func CalculateRequiredTurns(activityType ActivityType, characterAP int) int {
	info := GetActivityInfo(activityType)

	// AP積み上げ方式の場合
	if info.TotalRequiredAP > 0 && characterAP > 0 {
		// 必要総AP ÷ キャラクターのAP = 必要ターン数（切り上げ）
		return (info.TotalRequiredAP + characterAP - 1) / characterAP
	}

	// 即座実行アクションや特殊なアクションの場合は1ターン固定
	return 1
}

// CanInterrupt はアクティビティが中断可能かを返す
func (a *Activity) CanInterrupt() bool {
	info := GetActivityInfo(a.Type)
	return info.Interruptible && a.State == ActivityStateRunning
}

// CanResume はアクティビティが再開可能かを返す
func (a *Activity) CanResume() bool {
	info := GetActivityInfo(a.Type)
	return info.Resumable && a.State == ActivityStatePaused
}

// Interrupt はアクティビティを中断する
func (a *Activity) Interrupt(reason string) error {
	if !a.CanInterrupt() {
		return fmt.Errorf("アクティビティ '%s' は中断できません", a.GetDisplayName())
	}

	a.State = ActivityStatePaused
	now := time.Now()
	a.PauseTime = &now
	a.CancelReason = reason

	a.Logger.Debug("アクティビティ中断",
		"type", a.Type.String(),
		"actor", a.Actor,
		"reason", reason,
		"turns_left", a.TurnsLeft)

	return nil
}

// Resume はアクティビティを再開する
func (a *Activity) Resume() error {
	if !a.CanResume() {
		return fmt.Errorf("アクティビティ '%s' は再開できません", a.GetDisplayName())
	}

	a.State = ActivityStateRunning
	a.PauseTime = nil
	a.CancelReason = ""

	a.Logger.Debug("アクティビティ再開",
		"type", a.Type.String(),
		"actor", a.Actor,
		"turns_left", a.TurnsLeft)

	return nil
}

// Cancel はアクティビティをキャンセルする
func (a *Activity) Cancel(reason string) {
	a.State = ActivityStateCanceled
	a.CancelReason = reason

	a.Logger.Debug("アクティビティキャンセル",
		"type", a.Type.String(),
		"actor", a.Actor,
		"reason", reason)
}

// Complete はアクティビティを完了状態にする
func (a *Activity) Complete() {
	a.State = ActivityStateCompleted
	a.TurnsLeft = 0

	a.Logger.Debug("アクティビティ完了",
		"type", a.Type.String(),
		"actor", a.Actor,
		"duration", a.TurnsTotal)
}

// IsActive はアクティビティがアクティブかを返す
func (a *Activity) IsActive() bool {
	return a.State == ActivityStateRunning
}

// IsCompleted はアクティビティが完了しているかを返す
func (a *Activity) IsCompleted() bool {
	return a.State == ActivityStateCompleted || a.TurnsLeft <= 0
}

// IsCanceled はアクティビティがキャンセルされているかを返す
func (a *Activity) IsCanceled() bool {
	return a.State == ActivityStateCanceled
}

// GetProgressPercent は進捗率を0-100の値で返す
func (a *Activity) GetProgressPercent() float64 {
	if a.TurnsTotal <= 0 {
		return 100.0
	}
	completed := float64(a.TurnsTotal - a.TurnsLeft)
	return (completed / float64(a.TurnsTotal)) * 100.0
}

// GetDisplayName は表示用の名前を返す
func (a *Activity) GetDisplayName() string {
	info := GetActivityInfo(a.Type)
	return info.Name
}

// GetDisplayMessage は現在の状況を表示するメッセージを返す
func (a *Activity) GetDisplayMessage() string {
	if a.Message != "" {
		return a.Message
	}

	switch a.State {
	case ActivityStateRunning:
		percent := a.GetProgressPercent()
		return fmt.Sprintf("%s中... (%0.1f%%)", a.GetDisplayName(), percent)
	case ActivityStatePaused:
		return fmt.Sprintf("%s (一時停止中)", a.GetDisplayName())
	case ActivityStateCompleted:
		return fmt.Sprintf("%s完了", a.GetDisplayName())
	case ActivityStateCanceled:
		return fmt.Sprintf("%sキャンセル: %s", a.GetDisplayName(), a.CancelReason)
	default:
		return a.GetDisplayName()
	}
}

// String はActivityTypeの文字列表現を返す
func (t ActivityType) String() string {
	switch t {
	case ActivityMove:
		return "Move"
	case ActivityAttack:
		return "Attack"
	case ActivityPickup:
		return "Pickup"
	case ActivityWarp:
		return "Warp"
	case ActivityRest:
		return "Rest"
	case ActivityRead:
		return "Read"
	case ActivityCraft:
		return "Craft"
	case ActivityWait:
		return "Wait"
	case ActivityNull:
		return "Null"
	default:
		return fmt.Sprintf("ActivityType(%d)", int(t))
	}
}

// String はActivityStateの文字列表現を返す
func (s ActivityState) String() string {
	switch s {
	case ActivityStateRunning:
		return "Running"
	case ActivityStatePaused:
		return "Paused"
	case ActivityStateCompleted:
		return "Completed"
	case ActivityStateCanceled:
		return "Canceled"
	default:
		return fmt.Sprintf("ActivityState(%d)", int(s))
	}
}

// isPlayerActivity はアクティビティの実行者がプレイヤーかを判定する
func isPlayerActivity(act *Activity, world w.World) bool {
	return act.Actor.HasComponent(world.Components.Player)
}
