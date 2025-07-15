package testing

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
)

// EntityBuilder はテスト用のエンティティビルダー
type EntityBuilder struct {
	t          *testing.T
	components gc.GameComponentList
}

// NewEntityBuilder は新しいエンティティビルダーを作成する
func NewEntityBuilder(t *testing.T) *EntityBuilder {
	t.Helper()
	return &EntityBuilder{
		t:          t,
		components: gc.GameComponentList{},
	}
}

// WithName は名前を設定する
func (b *EntityBuilder) WithName(name string) *EntityBuilder {
	b.t.Helper()
	b.components.Name = &gc.Name{Name: name}
	return b
}

// WithPosition は位置を設定する
func (b *EntityBuilder) WithPosition(x, y float64) *EntityBuilder {
	b.t.Helper()
	b.components.Position = &gc.Position{X: gc.Pixel(x), Y: gc.Pixel(y)}
	return b
}

// WithHealth は体力を設定する
func (b *EntityBuilder) WithHealth(current, maxVal int) *EntityBuilder {
	b.t.Helper()
	b.components.Pools = &gc.Pools{
		HP:    gc.Pool{Current: current, Max: maxVal},
		SP:    gc.Pool{Current: maxVal / 2, Max: maxVal / 2},
		XP:    0,
		Level: 1,
	}
	return b
}

// WithStats は基本ステータスを設定する
func (b *EntityBuilder) WithStats(vitality, strength, sensation, dexterity, agility, defense int) *EntityBuilder {
	b.t.Helper()
	b.components.Attributes = &gc.Attributes{
		Vitality:  gc.Attribute{Base: vitality, Modifier: 0, Total: vitality},
		Strength:  gc.Attribute{Base: strength, Modifier: 0, Total: strength},
		Sensation: gc.Attribute{Base: sensation, Modifier: 0, Total: sensation},
		Dexterity: gc.Attribute{Base: dexterity, Modifier: 0, Total: dexterity},
		Agility:   gc.Attribute{Base: agility, Modifier: 0, Total: agility},
		Defense:   gc.Attribute{Base: defense, Modifier: 0, Total: defense},
	}
	return b
}

// AsPlayer はプレイヤーとして設定する
func (b *EntityBuilder) AsPlayer() *EntityBuilder {
	b.t.Helper()
	b.components.InParty = &gc.InParty{}
	return b
}

// AsEnemy は敵として設定する
func (b *EntityBuilder) AsEnemy() *EntityBuilder {
	b.t.Helper()
	b.components.FactionType = &gc.FactionEnemy
	return b
}

// AsItem はアイテムとして設定する
func (b *EntityBuilder) AsItem() *EntityBuilder {
	b.t.Helper()
	b.components.Item = &gc.Item{}
	b.components.ItemLocationType = &gc.ItemLocationInBackpack
	return b
}

// AsWeapon は武器として設定する
func (b *EntityBuilder) AsWeapon(damage, accuracy int) *EntityBuilder {
	b.t.Helper()
	b.AsItem()
	b.components.Wearable = &gc.Wearable{
		Defense:           0,
		EquipmentCategory: gc.EquipmentTorso,
	}
	b.components.Attack = &gc.Attack{
		Accuracy:       accuracy,
		Damage:         damage,
		AttackCount:    1,
		Element:        gc.ElementTypeNone,
		AttackCategory: gc.AttackSword,
	}
	return b
}

// AsConsumable は消耗品として設定する
func (b *EntityBuilder) AsConsumable(healAmount int) *EntityBuilder {
	b.t.Helper()
	b.AsItem()
	b.components.Consumable = &gc.Consumable{
		UsableScene: gc.UsableSceneBattle,
		TargetType: gc.TargetType{
			TargetGroup: gc.TargetGroupAlly,
			TargetNum:   gc.TargetSingle,
		},
	}
	b.components.ProvidesHealing = &gc.ProvidesHealing{
		Amount: gc.NumeralAmount{Numeral: healAmount},
	}
	return b
}

// AsMaterial は素材として設定する
func (b *EntityBuilder) AsMaterial(amount int) *EntityBuilder {
	b.t.Helper()
	b.AsItem()
	b.components.Material = &gc.Material{Amount: amount}
	return b
}

// WithRender はレンダリングを設定する（簡略版）
func (b *EntityBuilder) WithRender(_ string) *EntityBuilder {
	b.t.Helper()
	// レンダリングの設定は複雑なので、実際のテストでは適切に設定する必要がある
	// ここでは簡単な設定に留める
	return b
}

// WithDescription は説明を設定する
func (b *EntityBuilder) WithDescription(description string) *EntityBuilder {
	b.t.Helper()
	b.components.Description = &gc.Description{Description: description}
	return b
}

// Build は最終的なコンポーネントリストを返す
func (b *EntityBuilder) Build() gc.GameComponentList {
	b.t.Helper()
	return b.components
}

// MultiEntityBuilder は複数のエンティティを作成するためのヘルパー
type MultiEntityBuilder struct {
	t        *testing.T
	entities []gc.GameComponentList
}

// NewMultiEntityBuilder は複数エンティティビルダーを作成する
func NewMultiEntityBuilder(t *testing.T) *MultiEntityBuilder {
	t.Helper()
	return &MultiEntityBuilder{
		t:        t,
		entities: []gc.GameComponentList{},
	}
}

// Add はエンティティを追加する
func (mb *MultiEntityBuilder) Add(entity gc.GameComponentList) *MultiEntityBuilder {
	mb.t.Helper()
	mb.entities = append(mb.entities, entity)
	return mb
}

// AddBuilder はビルダーを使ってエンティティを追加する
func (mb *MultiEntityBuilder) AddBuilder(builderFunc func(*EntityBuilder) *EntityBuilder) *MultiEntityBuilder {
	mb.t.Helper()
	builder := NewEntityBuilder(mb.t)
	entity := builderFunc(builder).Build()
	mb.entities = append(mb.entities, entity)
	return mb
}

// Build は最終的なエンティティリストを返す
func (mb *MultiEntityBuilder) Build() []gc.GameComponentList {
	mb.t.Helper()
	return mb.entities
}

// よく使われるエンティティパターンのヘルパー関数

// CreateStandardPlayer は標準的なプレイヤーを作成する
func CreateStandardPlayer(t *testing.T) gc.GameComponentList {
	t.Helper()
	return NewEntityBuilder(t).
		WithName("プレイヤー").
		WithPosition(100, 100).
		WithHealth(100, 100).
		WithStats(10, 10, 10, 10, 10, 5).
		AsPlayer().
		WithRender("player").
		Build()
}

// CreateStandardEnemy は標準的な敵を作成する
func CreateStandardEnemy(t *testing.T, name string) gc.GameComponentList {
	t.Helper()
	return NewEntityBuilder(t).
		WithName(name).
		WithPosition(200, 200).
		WithHealth(50, 50).
		WithStats(8, 8, 8, 8, 8, 3).
		AsEnemy().
		WithRender("enemy").
		Build()
}

// CreateStandardWeapon は標準的な武器を作成する
func CreateStandardWeapon(t *testing.T, name string, damage, accuracy int) gc.GameComponentList {
	t.Helper()
	return NewEntityBuilder(t).
		WithName(name).
		WithDescription("テスト用の武器").
		AsWeapon(damage, accuracy).
		Build()
}

// CreateStandardPotion は標準的な回復アイテムを作成する
func CreateStandardPotion(t *testing.T, name string, healAmount int) gc.GameComponentList {
	t.Helper()
	return NewEntityBuilder(t).
		WithName(name).
		WithDescription("テスト用の回復アイテム").
		AsConsumable(healAmount).
		Build()
}
