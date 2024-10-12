package components

import (
	"fmt"

	ec "github.com/kijimaD/ruins/lib/engine/components"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// コンポーネントのリストが格納されたオブジェクト。
// この構造体を元にエンティティに対してコンポーネントを作成する。
// フィールドの型や値に応じて、対応するECSコンポーネントを取得する。
type GameComponentList struct {
	// general ================
	Name        *Name
	Description *Description

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
	Player       *Player
	Camera       *Camera
	Wall         *Wall
	Warp         *Warp
	Position     *Position
	GridElement  *GridElement
	SpriteRender *ec.SpriteRender
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
}

// componentsを溜めるスライス群
// Join時はこのフィールドでクエリする
type Components struct {
	// general ================
	Name        *ecs.SliceComponent
	Description *ecs.SliceComponent

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
	Player       *ecs.NullComponent
	Camera       *ecs.SliceComponent
	Wall         *ecs.NullComponent
	Warp         *ecs.SliceComponent
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
}

// フィールドで操作対象となる対象
// operatorとかのほうがよさそうか?
type Player struct{}

// カメラ
type Camera struct {
	Scale   float64
	ScaleTo float64
}

// 壁
type Wall struct{}

// ワープパッド
// TODO: 接触をトリガーに何かさせたいことはよくあるので、共通の仕組みを作る
type Warp struct {
	Mode warpMode
}

// キャラクターが保持できるもの
// 装備品、カード、回復アイテム、売却アイテム
type Item struct{}

// 消耗品。一度使うとなくなる
type Consumable struct {
	UsableScene UsableSceneType
	TargetType  TargetType
}

// 表示名
type Name struct {
	Name string
}

// 説明
type Description struct {
	Description string
}

// 装備品。キャラクタが装備することでパラメータを変更できる
type Wearable struct {
	Defense           int           // 防御力
	EquipmentCategory EquipmentType // 装備部位
	EquipBonus        EquipBonus    // ステータスへのボーナス
}

// 冒険パーティに参加している状態
type InParty struct{}

type Pools struct {
	HP    Pool
	SP    Pool
	Level int
}

// エンティティが持つステータス値。各種計算式で使う
type Attributes struct {
	Vitality  Attribute // 体力。丈夫さ、持久力、しぶとさ。HPやSPに影響する
	Strength  Attribute // 筋力。主に近接攻撃のダメージに影響する
	Sensation Attribute // 感覚。主に射撃攻撃のダメージに影響する
	Dexterity Attribute // 器用。攻撃時の命中率に影響する
	Agility   Attribute // 敏捷。回避率、行動の速さに影響する
	Defense   Attribute // 防御。被弾ダメージを軽減させる
}

// 回復する性質
// 直接的な数値が作用し、ステータスなどは考慮されない
type ProvidesHealing struct {
	Amount Amounter
}

// ダメージを与える性質
// 直接的な数値が作用し、ステータスなどは考慮されない
type InflictsDamage struct {
	Amount int
}

// 合成素材
// アイテムとの違い:
// - 個々のインスタンスで性能の違いはなく、単に数量だけを見る
// - 複数の単位で扱うのでAmountを持つ。x3を合成で使ったりする
type Material struct {
	Amount int
}

// 合成に必要な素材
type Recipe struct {
	Inputs []RecipeInput
}

// 装備変更が行われたことを示すダーティーフラグ
type EquipmentChanged struct{}

// カードは戦闘中に選択するコマンド
// 攻撃、防御、回復など、人に影響を及ぼすものをアクションカードという
type Card struct {
	TargetType TargetType
	Cost       int
}

// 攻撃の性質。攻撃毎にこの数値と作用対象のステータスを加味して、最終的なダメージ量を決定する
type Attack struct {
	Accuracy       int         // 命中率
	Damage         int         // 攻撃力
	AttackCount    int         // 攻撃回数
	Element        ElementType // 攻撃属性
	AttackCategory AttackType  // 攻撃種別
}

// 派閥
type FactionType string

const (
	// 味方
	FactionAlly FactionType = "FactionAlly"
	// 敵
	FactionEnemy FactionType = "FactionEnemy"
)

// ================
// アイテムの場所
type ItemLocationType fmt.Stringer

var (
	// バックパック内
	ItemLocationInBackpack ItemLocationType = LocationInBackpack{}
	// 味方が装備中
	ItemLocationEquipped ItemLocationType = LocationEquipped{}
	// フィールド上
	ItemLocationOnField ItemLocationType = LocationOnField{}
	// いずれにも存在しない。マスター用
	ItemLocationNone ItemLocationType = LocationNone{}
)

type LocationInBackpack struct{}

func (c LocationInBackpack) String() string {
	return "ItemLocationInBackpack"
}

type LocationEquipped struct {
	Owner         ecs.Entity
	EquipmentSlot EquipmentSlotNumber
}

func (c LocationEquipped) String() string {
	return "ItemLocationEquipped"
}

type LocationOnField struct{}

func (c LocationOnField) String() string {
	return "ItemLocationOnField"
}

type LocationNone struct{}

func (c LocationNone) String() string {
	return "ItemLocationNone"
}
