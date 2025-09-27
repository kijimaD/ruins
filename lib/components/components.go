package components

import (
	"fmt"
	"reflect"

	ecs "github.com/x-hgg-x/goecs/v2"
)

// EntitySpec はエンティティ作成用の仕様定義
// エンティティに付与するコンポーネントのセットを定義し、
// AddEntities関数でECSエンティティに変換される
type EntitySpec struct {
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
	AIMoveFSM    *AIMoveFSM
	AIRoaming    *AIRoaming
	AIVision     *AIVision
	AIChasing    *AIChasing
	Camera       *Camera
	Warp         *Warp
	Position     *Position
	GridElement  *GridElement
	SpriteRender *SpriteRender
	BlockView    *BlockView
	BlockPass    *BlockPass
	TurnBased    *TurnBased
	PropType     *PropType

	// member ================
	Player      *Player
	Hunger      *Hunger
	FactionType *FactionType
	Dead        *Dead

	// event ================
	EquipmentChanged *EquipmentChanged
	ProvidesHealing  *ProvidesHealing
	InflictsDamage   *InflictsDamage

	// battle ================
	CommandTable *CommandTable
	DropTable    *DropTable
}

// Components はECSコンポーネントストレージ
// 各コンポーネント型のSliceComponent/NullComponentを保持し、
// Manager.Join()でのクエリに使用される
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
	AIMoveFSM    *ecs.SliceComponent
	AIRoaming    *ecs.SliceComponent
	AIVision     *ecs.SliceComponent
	AIChasing    *ecs.SliceComponent
	Camera       *ecs.SliceComponent
	Warp         *ecs.SliceComponent
	Position     *ecs.SliceComponent
	GridElement  *ecs.SliceComponent
	SpriteRender *ecs.SliceComponent
	BlockView    *ecs.NullComponent
	BlockPass    *ecs.NullComponent
	PropType     *ecs.SliceComponent

	// member ================
	Player       *ecs.NullComponent
	Hunger       *ecs.SliceComponent
	FactionAlly  *ecs.NullComponent
	FactionEnemy *ecs.NullComponent
	Dead         *ecs.NullComponent
	TurnBased    *ecs.SliceComponent

	// event ================
	EquipmentChanged *ecs.NullComponent
	ProvidesHealing  *ecs.SliceComponent
	InflictsDamage   *ecs.SliceComponent

	// battle ================
	CommandTable *ecs.SliceComponent
	DropTable    *ecs.SliceComponent
}

// InitializeComponents はComponentInitializerインターフェースを実装する
// リフレクションを使用して自動的に全コンポーネントを初期化する
// コンポーネント追加時の手動更新が不要
func (c *Components) InitializeComponents(manager *ecs.Manager) error {
	val := reflect.ValueOf(c).Elem() // *Components から Components へ
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name

		// フィールドが設定可能かチェック
		if !field.CanSet() {
			return fmt.Errorf("field %s is not settable", fieldName)
		}

		// フィールドの型に基づいて適切なコンポーネントを作成
		switch field.Type() {
		case reflect.TypeOf((*ecs.SliceComponent)(nil)):
			// SliceComponent の初期化
			field.Set(reflect.ValueOf(manager.NewSliceComponent()))
		case reflect.TypeOf((*ecs.NullComponent)(nil)):
			// NullComponent の初期化
			field.Set(reflect.ValueOf(manager.NewNullComponent()))
		default:
			// 未対応の型はエラーとして扱う
			return fmt.Errorf("unsupported component type %v for field %s", field.Type(), fieldName)
		}
	}

	return nil
}

// Camera はカメラ
// 滑らかなズームの変更のため実際のズーム率と対象ズーム率を持つ
type Camera struct {
	Scale   float64
	ScaleTo float64
}

// Warp はワープパッド
// TODO: 接触をトリガーに何かさせたいことはよくあるので、共通の仕組みを作る
type Warp struct {
	Mode warpMode
}

// Item はキャラクターが保持できるもの。フィールド上、装備上、インベントリ上など位置状態を持ち、1スロットを消費する
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

// Player は操作対象の主人公キャラクター
type Player struct{}

// Dead はキャラクターが死亡している状態を示すマーカーコンポーネント
// 死亡時の処理(ドロップ/統計処理/ゲームログ...)を共通化するために使う
type Dead struct{}

// Pools はキャラクターのプール情報
type Pools struct {
	// 生命力 Health point
	// なくなるとゲームオーバー
	HP Pool
	// スタミナ Stamina point
	// 走ったり攻撃したら減る。自動回復する
	SP Pool
	// 電力 Electricity point
	// 機能のトグルで消費量が変わる
	EP Pool
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
