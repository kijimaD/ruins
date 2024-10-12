package components

import (
	ec "github.com/kijimaD/ruins/lib/engine/components"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// コンポーネントのリストが格納されたオブジェクト。
// コンポーネントの型に応じて、対応するECSコンポーネントを取得するために使用する。
type GameComponentList struct {
	// general ================
	Player      *Player
	Camera      *Camera
	Warp        *Warp
	Item        *Item
	Name        *Name
	Description *Description
	InBackpack  *InBackpack
	Equipped    *Equipped
	Consumable  *Consumable
	InParty     *InParty
	Member      *Member
	// Enemy            *Enemy
	Pools            *Pools
	ProvidesHealing  *ProvidesHealing
	InflictsDamage   *InflictsDamage
	Attack           *Attack
	Material         *Material
	Recipe           *Recipe
	Wearable         *Wearable
	Attributes       *Attributes
	EquipmentChanged *EquipmentChanged
	Card             *Card
	FactionType      *FactionType

	// field ================
	Position     *Position
	GridElement  *GridElement
	SpriteRender *ec.SpriteRender
	BlockView    *BlockView
	BlockPass    *BlockPass

	// battle ================
	BattleCommand *BattleCommand
}

// componentsを溜めるスライス群
// Join時はこのフィールドでクエリする
type Components struct {
	Player      *ecs.NullComponent
	Camera      *ecs.SliceComponent
	Wall        *ecs.NullComponent
	Warp        *ecs.SliceComponent
	Item        *ecs.NullComponent
	Consumable  *ecs.SliceComponent
	Name        *ecs.SliceComponent
	Description *ecs.SliceComponent
	InBackpack  *ecs.NullComponent
	InParty     *ecs.NullComponent
	Equipped    *ecs.SliceComponent
	Member      *ecs.NullComponent
	// Enemy            *ecs.NullComponent
	FactionTypeEnemy *ecs.NullComponent
	Pools            *ecs.SliceComponent
	ProvidesHealing  *ecs.SliceComponent
	InflictsDamage   *ecs.SliceComponent
	Attack           *ecs.SliceComponent
	Material         *ecs.SliceComponent
	Recipe           *ecs.SliceComponent
	Wearable         *ecs.SliceComponent
	Attributes       *ecs.SliceComponent
	EquipmentChanged *ecs.NullComponent
	Card             *ecs.SliceComponent

	Position     *ecs.SliceComponent
	GridElement  *ecs.SliceComponent
	SpriteRender *ecs.SliceComponent
	BlockView    *ecs.NullComponent
	BlockPass    *ecs.NullComponent

	BattleCommand *ecs.SliceComponent
}

// フィールドで操作対象となる対象
// operatorとかのほうがよさそうか?
type Player struct{}

// カメラ
type Camera struct {
	Scale   float64
	ScaleTo float64
}

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

// インベントリに入っている状態
type InBackpack struct{}

// キャラクタが装備している状態。InBackpackとは排反
type Equipped struct {
	Owner         ecs.Entity
	EquipmentSlot EquipmentSlotNumber
}

// 装備品。キャラクタが装備することでパラメータを変更できる
type Wearable struct {
	Defense           int           // 防御力
	EquipmentCategory EquipmentType // 装備部位
	EquipBonus        EquipBonus    // ステータスへのボーナス
}

// 冒険パーティに参加している状態
type InParty struct{}

// 冒険に参加できるメンバー
// TODO: allyにしたほうがいいかも
type Member struct{}

// 敵
// 「派閥」コンポーネントとして、味方と敵を排反にしたほうがいいのかもしれない
// そうするとJoinだけではだめなので面倒だな
// ただ排反を表現できればいいのだが
// コンポーネントでチェックを入れるようにすればいいのか
// type Enemy struct{}

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
