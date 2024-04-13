package components

import ecs "github.com/x-hgg-x/goecs/v2"

type Components struct {
	GridElement     *ecs.SliceComponent
	Player          *ecs.NullComponent
	Wall            *ecs.NullComponent
	Warp            *ecs.NullComponent
	Item            *ecs.NullComponent
	Consumable      *ecs.SliceComponent
	Name            *ecs.SliceComponent
	Description     *ecs.SliceComponent
	InBackpack      *ecs.NullComponent
	InParty         *ecs.NullComponent
	Equipped        *ecs.SliceComponent
	Member          *ecs.NullComponent
	Pools           *ecs.SliceComponent
	ProvidesHealing *ecs.SliceComponent
	InflictsDamage  *ecs.SliceComponent
	Weapon          *ecs.SliceComponent
	Material        *ecs.SliceComponent
	Recipe          *ecs.SliceComponent
	Wearable        *ecs.SliceComponent
	Attributes      *ecs.SliceComponent
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
// 一切使用できない売却専用アイテムとかはItem単独に含まれる
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

// インベントリに所持している
type InBackpack struct{}

// キャラクタが装備している
type Equipped struct {
	Owner         ecs.Entity
	EquipmentSlot EquipmentSlotNumber
}

// 武器
type Weapon struct {
	Accuracy          int            // 命中率
	BaseDamage        int            // 攻撃力
	AttackCount       int            // 攻撃回数
	EnergyConsumption int            // 消費エネルギー
	DamageAttr        DamageAttrType // 攻撃属性
	WeaponCategory    WeaponType     // 武器種別
}

type Wearable struct {
	BaseDefense       int           // 防御力
	EquipmentCategory EquipmentType // 装備部位
}

// パーティに参加している状態
type InParty struct{}

// 冒険に参加できるメンバー
type Member struct{}

// 最大値と現在値を持つようなパラメータ
type Pool struct {
	Max     int
	Current int
}

type Pools struct {
	HP    Pool
	SP    Pool
	Level int
}

type Attribute struct {
	Base     int // 固有の値
	Modifier int // 装備などで変動する値
	Total    int // 足し合わせた現在値。メモ
}

// ステータス
type Attributes struct {
	Vitality Attribute // 生命力。HPやSPに影響する
	Strength Attribute // 膂力。近接攻撃
}

// 回復する性質
type ProvidesHealing struct {
	Amount int
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
