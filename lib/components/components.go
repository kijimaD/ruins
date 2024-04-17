package components

import ecs "github.com/x-hgg-x/goecs/v2"

type Components struct {
	GridElement      *ecs.SliceComponent
	Player           *ecs.NullComponent
	Wall             *ecs.NullComponent
	Warp             *ecs.NullComponent
	Item             *ecs.NullComponent
	Consumable       *ecs.SliceComponent
	Name             *ecs.SliceComponent
	Description      *ecs.SliceComponent
	InBackpack       *ecs.NullComponent
	InParty          *ecs.NullComponent
	Equipped         *ecs.SliceComponent
	Member           *ecs.NullComponent
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
}

type GridElement struct {
	Line int
	Col  int
}

// フィールドでの移動体
type Player struct{}

// 壁
type Wall struct{}

// ワープパッド
type Warp struct {
	Mode warpMode
}

// アイテム枠に入るもの
// 使用できない売却専用アイテムなどはItem単独に含まれる
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

// インベントリに所持している状態
type InBackpack struct{}

// キャラクタが装備している状態
type Equipped struct {
	Owner         ecs.Entity
	EquipmentSlot EquipmentSlotNumber
}

// 防具
type Wearable struct {
	Defense           int           // 防御力
	EquipmentCategory EquipmentType // 装備部位
	EquipBonus        EquipBonus    // ステータスへのボーナス
}

// パーティに参加している状態
type InParty struct{}

// 冒険に参加できるメンバー
type Member struct{}

type Pools struct {
	HP    Pool
	SP    Pool
	Level int
}

// エンティティが持つステータス。各種計算式で使う
type Attributes struct {
	Vitality  Attribute // 体力。丈夫さ、持久力、しぶとさ。HPやSPに影響する
	Strength  Attribute // 筋力。主に近接攻撃のダメージに影響する
	Sensation Attribute // 感覚。主に射撃攻撃のダメージに影響する
	Dexterity Attribute // 器用。攻撃時の命中率に影響する
	Agility   Attribute // 敏捷。回避率、行動の速さに影響する
	Defense   Attribute // 防御。被弾ダメージを軽減させる
}

// 回復する性質
type ProvidesHealing struct {
	Amount Amounter
}

// ダメージを与える性質
type InflictsDamage struct {
	Amount int
}

// 合成素材。
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

// 装備変更直後を示すダーティーフラグ
type EquipmentChanged struct{}

// カードは戦闘中に選択するためのコマンド
// 攻撃や防御など、人に影響を及ぼすものをアクションカードという
// 他カードをターゲットとするものをブーストカードという
type Card struct {
	TargetType TargetType
	Cost       int
}

type Attack struct {
	Accuracy       int         // 命中率
	Damage         int         // 攻撃力
	AttackCount    int         // 攻撃回数
	Element        ElementType // 攻撃属性
	AttackCategory AttackType  // 攻撃種別
	EquipBonus     EquipBonus  // ステータスへのボーナス
}
