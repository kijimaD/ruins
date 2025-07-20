package testing

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/mocks"
)

// テストヘルパーの使用例
func TestEntityBuilder(t *testing.T) {
	t.Parallel()
	// 単体エンティティの作成
	player := NewEntityBuilder(t).
		WithName("テストプレイヤー").
		WithPosition(100, 100).
		WithHealth(100, 100).
		WithStats(10, 10, 10, 10, 10, 5).
		AsPlayer().
		WithRender("player").
		Build()

	// アサーションの使用例
	AssertName(t, player.Name, "テストプレイヤー")
	AssertPosition(t, player.Position, 100, 100)
	AssertPools(t, player.Pools, 100, 100)
	AssertHasComponent(t, player, "Player")
	AssertNotHasComponent(t, player, "Enemy")
}

func TestMultiEntityBuilder(t *testing.T) {
	t.Parallel()
	// 複数エンティティの作成
	entities := NewMultiEntityBuilder(t).
		AddBuilder(func(b *EntityBuilder) *EntityBuilder {
			return b.WithName("プレイヤー").AsPlayer()
		}).
		AddBuilder(func(b *EntityBuilder) *EntityBuilder {
			return b.WithName("スライム").AsEnemy()
		}).
		AddBuilder(func(b *EntityBuilder) *EntityBuilder {
			return b.WithName("剣").AsWeapon(10, 80)
		}).
		Build()

	if len(entities) != 3 {
		t.Errorf("エンティティ数が期待値と異なります: 期待値3, 実際値%d", len(entities))
	}

	// 各エンティティの検証
	AssertName(t, entities[0].Name, "プレイヤー")
	AssertName(t, entities[1].Name, "スライム")
	AssertName(t, entities[2].Name, "剣")
}

func TestStandardEntities(t *testing.T) {
	t.Parallel()
	// 標準エンティティの作成
	player := CreateStandardPlayer(t)
	enemy := CreateStandardEnemy(t, "ゴブリン")
	weapon := CreateStandardWeapon(t, "鉄の剣", 15, 85)
	potion := CreateStandardPotion(t, "回復薬", 50)

	// 検証
	AssertName(t, player.Name, "プレイヤー")
	AssertName(t, enemy.Name, "ゴブリン")
	AssertName(t, weapon.Name, "鉄の剣")
	AssertName(t, potion.Name, "回復薬")

	AssertHasComponent(t, player, "Player")
	AssertHasComponent(t, enemy, "Enemy")
	AssertHasComponent(t, weapon, "Weapon")
	AssertHasComponent(t, potion, "Consumable")
}

func TestMockInput(t *testing.T) {
	t.Parallel()
	// モック入力の使用例
	mockInput := mocks.NewMockInputHandler()

	// 初期状態の確認
	if mockInput.IsKeyPressed(ebiten.KeyW) {
		t.Error("キーが押されていないはずです")
	}

	// キーを押す
	mockInput.PressKey(ebiten.KeyW)
	if !mockInput.IsKeyPressed(ebiten.KeyW) {
		t.Error("キーが押されているはずです")
	}
	if !mockInput.IsKeyJustPressed(ebiten.KeyW) {
		t.Error("キーが今フレームで押されているはずです")
	}

	// フレーム終了処理
	mockInput.EndFrame()
	if mockInput.IsKeyJustPressed(ebiten.KeyW) {
		t.Error("フレーム終了後はJustPressedがfalseになるはずです")
	}
	if !mockInput.IsKeyPressed(ebiten.KeyW) {
		t.Error("キーはまだ押されているはずです")
	}

	// キーを離す
	mockInput.ReleaseKey(ebiten.KeyW)
	if mockInput.IsKeyPressed(ebiten.KeyW) {
		t.Error("キーが離されているはずです")
	}
}

func TestMockRandom(t *testing.T) {
	t.Parallel()
	// モック乱数の使用例
	mockRandom := mocks.NewMockRandomGenerator()

	// 値を設定
	mockRandom.SetFloat64Values(0.1, 0.5, 0.9)
	mockRandom.SetIntValues(1, 3, 7)

	// 値の取得
	if got := mockRandom.Float64(); got != 0.1 {
		t.Errorf("期待値0.1, 実際値%f", got)
	}
	if got := mockRandom.Float64(); got != 0.5 {
		t.Errorf("期待値0.5, 実際値%f", got)
	}
	if got := mockRandom.Float64(); got != 0.9 {
		t.Errorf("期待値0.9, 実際値%f", got)
	}

	// 循環して最初に戻る
	if got := mockRandom.Float64(); got != 0.1 {
		t.Errorf("期待値0.1（循環）, 実際値%f", got)
	}

	// 整数値のテスト
	if got := mockRandom.Intn(10); got != 1 {
		t.Errorf("期待値1, 実際値%d", got)
	}
}

func TestBattleScenarioUsage(t *testing.T) {
	t.Parallel()
	// バトルシナリオの使用例
	scenario := CreateTestBattleScenario(t)

	// プレイヤーの確認
	AssertName(t, scenario.Player.Name, "テストプレイヤー")
	AssertHasComponent(t, scenario.Player, "Player")

	// 敵の確認
	if len(scenario.Enemies) != 2 {
		t.Errorf("敵の数が期待値と異なります: 期待値2, 実際値%d", len(scenario.Enemies))
	}
	AssertName(t, scenario.Enemies[0].Name, "スライム")
	AssertName(t, scenario.Enemies[1].Name, "ゴブリン")

	// アイテムの確認
	if len(scenario.Items) != 2 {
		t.Errorf("アイテムの数が期待値と異なります: 期待値2, 実際値%d", len(scenario.Items))
	}
	AssertName(t, scenario.Items[0].Name, "テスト剣")
	AssertName(t, scenario.Items[1].Name, "テスト薬草")
}
