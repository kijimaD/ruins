package gamelog

import (
	"fmt"

	"github.com/kijimaD/ruins/lib/colors"
)

// この例は、Loggerのメソッドチェーンの使い方を示しています
func ExampleLogger_methodChaining() {
	// ローカルログストアを作成
	testLog := NewSafeSlice(FieldLogMaxSize)

	// メソッドチェーンでのログ作成
	New(testLog).
		NPCName("Goblin").
		Append(" attacks you for ").
		Damage(15).
		Append(" damage!").
		Log()

	// 色指定も同様
	New(testLog).
		ColorRGBA(colors.ColorCyan). // Cyan
		Append("John").
		ColorRGBA(colors.ColorWhite).
		Append(" considers attacking ").
		ColorRGBA(colors.ColorCyan).
		Append("Orc").
		Log()

	// アイテムとプレイヤー名
	New(testLog).
		PlayerName("Hero").
		Append(" picks up ").
		ItemName("Iron Sword").
		Append(".").
		Log()

	fmt.Println("Logs created successfully!")
	// Output: Logs created successfully!
}

// 戦闘ログの例
func ExampleLogger_battleLog() {
	// ローカル戦闘ログストアを作成
	testBattleLog := NewSafeSlice(BattleLogMaxSize)
	testBattleLog.Clear()

	// 複雑な戦闘ログ
	New(testBattleLog).
		NPCName("Skeleton Warrior").
		Append(" swings ").
		ItemName("Rusty Sword").
		Append(" at ").
		PlayerName("Hero").
		Append(" for ").
		Damage(12).
		Append(" damage!").
		Log()

	// ログの取得と表示
	messages := testBattleLog.GetRecent(1)
	fmt.Println(messages[0])
	// Output: Skeleton Warrior swings Rusty Sword at Hero for 12 damage!
}

// カスタム色の例
func ExampleLogger_customColors() {
	// ローカルフィールドログストアを作成
	testFieldLog := NewSafeSlice(FieldLogMaxSize)
	testFieldLog.Clear()

	New(testFieldLog).
		ColorRGBA(colors.ColorPurple).
		Append("Magic spell ").
		ColorRGBA(colors.ColorOrange).
		Append("Fire Bolt").
		ColorRGBA(colors.ColorWhite).
		Append(" hits for ").
		Damage(20).
		Append(" damage!").
		Log()

	// 色付きエントリの取得
	entries := testFieldLog.GetRecentEntries(1)
	fmt.Printf("Entry has %d colored fragments\n", len(entries[0].Fragments))
	// Output: Entry has 5 colored fragments
}

// 連続攻撃のログ例
func ExampleLogger_chainedAttack() {
	// ローカルフィールドログストアを作成
	testFieldLog := NewSafeSlice(FieldLogMaxSize)

	New(testFieldLog).
		NPCName("Orc").
		Append(" attacks ").
		Append("->").
		NPCName("Player").
		Append(" for ").
		Damage(8).
		Log()

	fmt.Println("Chained attack log created!")
	// Output: Chained attack log created!
}
