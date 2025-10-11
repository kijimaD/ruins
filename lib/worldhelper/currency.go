package worldhelper

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// AddCurrency はエンティティに所持金を追加する
func AddCurrency(world w.World, entity ecs.Entity, amount int) {
	wallet := world.Components.Wallet.Get(entity)
	if wallet == nil {
		return
	}
	w := wallet.(*gc.Wallet)
	w.Currency += amount
}

// GetCurrency はエンティティの所持金を取得する
func GetCurrency(world w.World, entity ecs.Entity) int {
	wallet := world.Components.Wallet.Get(entity)
	if wallet == nil {
		return 0
	}
	return wallet.(*gc.Wallet).Currency
}

// SetCurrency はエンティティの所持金を設定する
func SetCurrency(world w.World, entity ecs.Entity, amount int) {
	wallet := world.Components.Wallet.Get(entity)
	if wallet == nil {
		return
	}
	w := wallet.(*gc.Wallet)
	w.Currency = amount
}

// HasCurrency は指定額以上の所持金を持っているか確認
func HasCurrency(world w.World, entity ecs.Entity, amount int) bool {
	return GetCurrency(world, entity) >= amount
}

// ConsumeCurrency はエンティティの所持金を消費する
// 所持金が足りない場合はfalseを返す
func ConsumeCurrency(world w.World, entity ecs.Entity, amount int) bool {
	if !HasCurrency(world, entity, amount) {
		return false
	}
	wallet := world.Components.Wallet.Get(entity)
	if wallet == nil {
		return false
	}
	w := wallet.(*gc.Wallet)
	w.Currency -= amount
	return true
}
