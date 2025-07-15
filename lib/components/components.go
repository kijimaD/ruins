package components

import (
	"fmt"
	"time"

	ecs "github.com/x-hgg-x/goecs/v2"
)

// GameComponentList はコンポーネントのリストが格納されたオブジェクト。
// この構造体を元にエンティティに対してコンポーネントを作成する。
// フィールドの型や値に応じて、対応するECSコンポーネントを取得する。
type GameComponentList struct {
	// general ================
	Name        *Name
	Description *Description
	Render      *Render

	// item ================
	Item             *Item
	Consumable       *Consumable
	Pools            *Pools
	Attack           *Attack
	Material         *Material
	Recipe           *Recipe
	Wearable         *Wearable
	Attributes       *Attributes
	Card             *Card
	ItemLocationType *ItemLocationType

	// field ================
	Operator     *Operator
	AIMoveFSM    *AIMoveFSM
	AIRoaming    *AIRoaming
	Camera       *Camera
	Wall         *Wall
	Warp         *Warp
	Velocity     *Velocity
	Position     *Position
	GridElement  *GridElement
	SpriteRender *SpriteRender
	BlockView    *BlockView
	BlockPass    *BlockPass

	// member ================
	InParty     *InParty
	FactionType *FactionType

	// event ================
	BattleCommand    *BattleCommand
	EquipmentChanged *EquipmentChanged
	ProvidesHealing  *ProvidesHealing
	InflictsDamage   *InflictsDamage

	// battle ================
	CommandTable *CommandTable
	DropTable    *DropTable
}

// Components はcomponentsを溜めるスライス群
// Join時はこのフィールドでクエリする
type Components struct {
	// general ================
	Name        *ecs.SliceComponent
	Description *ecs.SliceComponent
	Render      *ecs.SliceComponent

	// item ================
	Item                   *ecs.NullComponent
	Consumable             *ecs.SliceComponent
	Pools                  *ecs.SliceComponent
	Attack                 *ecs.SliceComponent
	Material               *ecs.SliceComponent
	Recipe                 *ecs.SliceComponent
	Wearable               *ecs.SliceComponent
	Attributes             *ecs.SliceComponent
	Card                   *ecs.SliceComponent
	ItemLocationInBackpack *ecs.NullComponent
	ItemLocationEquipped   *ecs.SliceComponent
	ItemLocationOnField    *ecs.NullComponent
	ItemLocationNone       *ecs.NullComponent

	// field ================
	Operator     *ecs.NullComponent
	AIMoveFSM    *ecs.SliceComponent
	AIRoaming    *ecs.SliceComponent
	Camera       *ecs.SliceComponent
	Wall         *ecs.NullComponent
	Warp         *ecs.SliceComponent
	Velocity     *ecs.SliceComponent
	Position     *ecs.SliceComponent
	GridElement  *ecs.SliceComponent
	SpriteRender *ecs.SliceComponent
	BlockView    *ecs.NullComponent
	BlockPass    *ecs.NullComponent

	// member ================
	InParty      *ecs.NullComponent
	FactionAlly  *ecs.NullComponent
	FactionEnemy *ecs.NullComponent

	// event ================
	BattleCommand    *ecs.SliceComponent
	EquipmentChanged *ecs.NullComponent
	ProvidesHealing  *ecs.SliceComponent
	InflictsDamage   *ecs.SliceComponent

	// battle ================
	CommandTable *ecs.SliceComponent
	DropTable    *ecs.SliceComponent
}

// Operator はフィールドでの操作対象
type Operator struct{}

// AIMoveFSM はAI移動の有限状態マシン
type AIMoveFSM struct {
	LastStateChange time.Time
}

// AIRoaming はAI移動で歩き回り状態
type AIRoaming struct {
	SubState AIRoamingSubState
	// サブステートの開始時間
	StartSubState time.Time
	// サブステートの持続時間
	DurationSubState time.Duration
}

// Camera はカメラ
type Camera struct {
	Scale   float64
	ScaleTo float64
}

// Wall は壁
type Wall struct{}

// Warp はワープパッド
// TODO: 接触をトリガーに何かさせたいことはよくあるので、共通の仕組みを作る
type Warp struct {
	Mode warpMode
}

// Item はキャラクターが保持できるもの
// 装備品、カード、回復アイテム、売却アイテム
type Item struct{}

// Consumable は消耗品。一度使うとなくなる
type Consumable struct {
	UsableScene UsableSceneType
	TargetType  TargetType
}

// Name は表示名
type Name struct {
	Name string
}

// Description は説明
type Description struct {
	Description string
}

// Wearable は装備品。キャラクタが装備することでパラメータを変更できる
type Wearable struct {
	Defense           int           // 防御力
	EquipmentCategory EquipmentType // 装備部位
	EquipBonus        EquipBonus    // ステータスへのボーナス
}

// InParty は冒険パーティに参加している状態
type InParty struct{}

// Pools はキャラクターのプール情報
type Pools struct {
	// 生命
	HP Pool
	// 特殊行動力
	SP Pool
	// 経験値
	XP int
	// レベル
	Level int
}

// Attributes はエンティティが持つステータス値。各種計算式で使う
type Attributes struct {
	Vitality  Attribute // 体力。丈夫さ、持久力、しぶとさ。HPやSPに影響する
	Strength  Attribute // 筋力。主に近接攻撃のダメージに影響する
	Sensation Attribute // 感覚。主に射撃攻撃のダメージに影響する
	Dexterity Attribute // 器用。攻撃時の命中率に影響する
	Agility   Attribute // 敏捷。回避率、行動の速さに影響する
	Defense   Attribute // 防御。被弾ダメージを軽減させる
}

// ProvidesHealing は回復する性質
// 直接的な数値が作用し、ステータスなどは考慮されない
type ProvidesHealing struct {
	Amount Amounter
}

// InflictsDamage はダメージを与える性質
// 直接的な数値が作用し、ステータスなどは考慮されない
type InflictsDamage struct {
	Amount int
}

// Material は合成素材
// アイテムとの違い:
// - 個々のインスタンスで性能の違いはなく、単に数量だけを見る
// - 複数の単位で扱うのでAmountを持つ。x3を合成で使ったりする
type Material struct {
	Amount int
}

// Recipe は合成に必要な素材
type Recipe struct {
	Inputs []RecipeInput
}

// EquipmentChanged は装備変更が行われたことを示すダーティーフラグ
type EquipmentChanged struct{}

// Card はカードは戦闘中に選択するコマンド
// 攻撃、防御、回復など、人に影響を及ぼすものをアクションカードという
type Card struct {
	TargetType TargetType
	Cost       int
}

// Attack は攻撃の性質。攻撃毎にこの数値と作用対象のステータスを加味して、最終的なダメージ量を決定する
type Attack struct {
	Accuracy       int         // 命中率
	Damage         int         // 攻撃力
	AttackCount    int         // 攻撃回数
	Element        ElementType // 攻撃属性
	AttackCategory AttackType  // 攻撃種別
}

// CommandTable はAI用の、戦闘コマンドテーブル名
type CommandTable struct {
	Name string
}

// DropTable はドロップテーブル名
type DropTable struct {
	Name string
}

// Render は描画対象物。
// FIXME: 表示する方法を統一する...
// resource のスプライトシートから画像を特定するために必要な情報。
// Renderはスプライトシートの情報はresourceに持たせ、特定に必要な
// 情報だけ保持している。
// また、場面ごとにフィールドを分けて保持している。
//
// SpriteRenderはSpriteを内部に持っていて初期化が面倒な面がある。
type Render struct {
	// 戦闘中の立ち絵
	BattleBody *SheetImage
}

// SheetImage はシート画像情報
type SheetImage struct {
	SheetName   string
	SheetNumber *int
}

// FactionType は所属派閥。絶対的な指定
type FactionType fmt.Stringer

var (
	// FactionAlly は味方(プレイヤー側)
	FactionAlly FactionType = FactionAllyData{}
	// FactionEnemy は敵性(プレイヤーと敵対)
	FactionEnemy FactionType = FactionEnemyData{}
)

// FactionAllyData は味方派閥データ
type FactionAllyData struct{}

func (c FactionAllyData) String() string {
	return "FactionAlly"
}

// FactionEnemyData は敵性派閥データ
type FactionEnemyData struct{}

func (c FactionEnemyData) String() string {
	return "FactionEnemy"
}

// ItemLocationType はアイテムの場所
type ItemLocationType fmt.Stringer

var (
	// ItemLocationInBackpack はバックパック内
	ItemLocationInBackpack ItemLocationType = LocationInBackpack{}
	// ItemLocationEquipped は味方が装備中
	ItemLocationEquipped ItemLocationType = LocationEquipped{}
	// ItemLocationOnField はフィールド上
	ItemLocationOnField ItemLocationType = LocationOnField{}
	// ItemLocationNone はいずれにも存在しない。マスター用
	ItemLocationNone ItemLocationType = LocationNone{}
)

// LocationInBackpack はバックパック内位置
type LocationInBackpack struct{}

func (c LocationInBackpack) String() string {
	return "ItemLocationInBackpack"
}

// LocationEquipped は装備中位置
type LocationEquipped struct {
	Owner         ecs.Entity
	EquipmentSlot EquipmentSlotNumber
}

func (c LocationEquipped) String() string {
	return "ItemLocationEquipped"
}

// LocationOnField はフィールド上位置
type LocationOnField struct{}

func (c LocationOnField) String() string {
	return "ItemLocationOnField"
}

// LocationNone は位置なし
type LocationNone struct{}

func (c LocationNone) String() string {
	return "ItemLocationNone"
}
