package gamelog

import "fmt"

// この例は、Loggerのメソッドチェーンの使い方を示しています
func ExampleLogger_methodChaining() {
	// メソッドチェーンでのログ作成
	New().
		NPCName("Goblin").
		Append(" attacks you for ").
		Damage(15).
		Append(" damage!").
		Log(LogKindField)

	// 色指定も同様
	New().
		Color(0, 255, 255). // Cyan
		Append("John").
		ColorRGBA(ColorWhite).
		Append(" considers attacking ").
		ColorRGBA(ColorCyan).
		Append("Orc").
		Log(LogKindField)

	// アイテムとプレイヤー名
	New().
		PlayerName("Hero").
		Append(" picks up ").
		ItemName("Iron Sword").
		Append(".").
		Log(LogKindField)

	fmt.Println("Logs created successfully!")
	// Output: Logs created successfully!
}

// 戦闘ログの例
func ExampleLogger_battleLog() {
	BattleLog.Clear()

	// 複雑な戦闘ログ
	New().
		NPCName("Skeleton Warrior").
		Append(" swings ").
		ItemName("Rusty Sword").
		Append(" at ").
		PlayerName("Hero").
		Append(" for ").
		Damage(12).
		Append(" damage!").
		Log(LogKindBattle)

	// ログの取得と表示
	messages := BattleLog.GetRecent(1)
	fmt.Println(messages[0])
	// Output: Skeleton Warrior swings Rusty Sword at Hero for 12 damage!
}

// カスタム色の例
func ExampleLogger_customColors() {
	FieldLog.Clear()

	New().
		ColorRGBA(ColorPurple).
		Append("Magic spell ").
		ColorRGBA(ColorOrange).
		Append("Fire Bolt").
		ColorRGBA(ColorWhite).
		Append(" hits for ").
		Damage(20).
		Append(" damage!").
		Log(LogKindField)

	// 色付きエントリの取得
	entries := FieldLog.GetRecentEntries(1)
	fmt.Printf("Entry has %d colored fragments\n", len(entries[0].Fragments))
	// Output: Entry has 5 colored fragments
}

// 連続攻撃のログ例
func ExampleLogger_chainedAttack() {
	New().
		NPCName("Orc").
		Append(" attacks ").
		Append("->").
		NPCName("Player").
		Append(" for ").
		Damage(8).
		Log(LogKindField)

	fmt.Println("Chained attack log created!")
	// Output: Chained attack log created!
}
