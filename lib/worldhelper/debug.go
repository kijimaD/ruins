package worldhelper

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// InitDebugData はデバッグ用の初期データを設定する
// パーティメンバーが存在しない場合のみ実行される
// テスト、VRT、デバッグで使用される共通のエンティティセットを生成する
func InitDebugData(world w.World) {
	// 既にパーティメンバーが存在するかチェック
	memberCount := 0
	world.Manager.Join(
		world.Components.Player,
		world.Components.FactionAlly,
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		memberCount++
	}))

	// 既にメンバーがいる場合は何もしない
	if memberCount > 0 {
		return
	}

	// 基本アイテムの生成
	card1, _ := SpawnItem(world, "木刀", gc.ItemLocationInBackpack)
	card2, _ := SpawnItem(world, "ハンドガン", gc.ItemLocationInBackpack)
	card3, _ := SpawnItem(world, "M72 LAW", gc.ItemLocationInBackpack)
	armor, _ := SpawnItem(world, "西洋鎧", gc.ItemLocationInBackpack)
	_, _ = SpawnItem(world, "作業用ヘルメット", gc.ItemLocationInBackpack)
	_, _ = SpawnItem(world, "革のブーツ", gc.ItemLocationInBackpack)
	_, _ = SpawnItem(world, "ルビー原石", gc.ItemLocationInBackpack)
	_, _ = SpawnItem(world, "レイガン", gc.ItemLocationInBackpack)
	_, _ = SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)
	_, _ = SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)
	_, _ = SpawnItem(world, "回復スプレー", gc.ItemLocationInBackpack)
	_, _ = SpawnItem(world, "回復スプレー", gc.ItemLocationInBackpack)
	_, _ = SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
	_, _ = SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
	_, _ = SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
	_, _ = SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)

	// 大量アイテム生成（デバッグ用）
	for i := 0; i < 10; i++ {
		_, _ = SpawnItem(world, "木刀", gc.ItemLocationInBackpack)
		_, _ = SpawnItem(world, "ハンドガン", gc.ItemLocationInBackpack)
		_, _ = SpawnItem(world, "レイガン", gc.ItemLocationInBackpack)
	}
	for i := 0; i < 10; i++ {
		_, _ = SpawnItem(world, "西洋鎧", gc.ItemLocationInBackpack)
		_, _ = SpawnItem(world, "作業用ヘルメット", gc.ItemLocationInBackpack)
		_, _ = SpawnItem(world, "革のブーツ", gc.ItemLocationInBackpack)
	}
	for i := 0; i < 10; i++ {
		_, _ = SpawnItem(world, "ルビー原石", gc.ItemLocationInBackpack)
	}
	for i := 0; i < 10; i++ {
		_, _ = SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)
		_, _ = SpawnItem(world, "回復スプレー", gc.ItemLocationInBackpack)
		_, _ = SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
	}

	// プレイヤー生成（一人のみ）
	celestine, _ := SpawnPlayer(world, "セレスティン")

	// マテリアルとレシピ
	_ = SpawnAllMaterials(world)
	PlusAmount("鉄", 40, world)
	PlusAmount("鉄くず", 4, world)
	PlusAmount("緑ハーブ", 2, world)
	PlusAmount("黄ハーブ", 1, world)
	PlusAmount("木の棒", 1, world)
	PlusAmount("フェライトコア", 30, world)
	PlusAmount("銀の欠片", 1, world)
	PlusAmount("古い歯車", 1, world)
	PlusAmount("水晶の粉", 1, world)
	PlusAmount("黒曜石", 1, world)
	PlusAmount("血赤石", 1, world)
	PlusAmount("月光草", 1, world)
	PlusAmount("古代の骨", 1, world)
	PlusAmount("星鉄", 1, world)
	PlusAmount("霧の結晶", 1, world)
	PlusAmount("古布の切れ端", 1, world)
	PlusAmount("琥珀", 1, world)
	PlusAmount("深海の塩", 1, world)
	PlusAmount("雷光石", 1, world)
	PlusAmount("灰", 1, world)
	PlusAmount("蒼鉛", 1, world)
	PlusAmount("獣の毛皮", 1, world)
	PlusAmount("朽ちた羊皮紙", 1, world)
	PlusAmount("紫水晶", 1, world)
	PlusAmount("錆びた鎖", 1, world)
	PlusAmount("聖油", 1, world)
	_ = SpawnAllRecipes(world)
	_ = SpawnAllCards(world)

	// 装備
	Equip(world, card1, celestine, gc.EquipmentSlotNumber(0))
	Equip(world, card2, celestine, gc.EquipmentSlotNumber(1))
	Equip(world, card3, celestine, gc.EquipmentSlotNumber(2))
	Equip(world, armor, celestine, gc.EquipmentSlotNumber(0))
}
