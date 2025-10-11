package worldhelper

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/maingame"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddCurrency(t *testing.T) {
	t.Parallel()
	world, err := maingame.InitWorld(960, 720)
	require.NoError(t, err)

	// プレイヤーを作成してWalletを追加
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Wallet, &gc.Wallet{Currency: 100})

	// 通貨を追加
	AddCurrency(world, player, 50)

	// 結果を検証
	currency := GetCurrency(world, player)
	assert.Equal(t, 150, currency, "通貨が150になるべき")

	// クリーンアップ
	world.Manager.DeleteEntity(player)
}

func TestGetCurrency(t *testing.T) {
	t.Parallel()
	world, err := maingame.InitWorld(960, 720)
	require.NoError(t, err)

	// プレイヤーを作成してWalletを追加
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Wallet, &gc.Wallet{Currency: 200})

	// 通貨を取得
	currency := GetCurrency(world, player)
	assert.Equal(t, 200, currency, "通貨が200であるべき")

	// Walletがない場合
	playerWithoutWallet := world.Manager.NewEntity()
	currency = GetCurrency(world, playerWithoutWallet)
	assert.Equal(t, 0, currency, "Walletがない場合は0を返すべき")

	// クリーンアップ
	world.Manager.DeleteEntity(player)
	world.Manager.DeleteEntity(playerWithoutWallet)
}

func TestSetCurrency(t *testing.T) {
	t.Parallel()
	world, err := maingame.InitWorld(960, 720)
	require.NoError(t, err)

	// プレイヤーを作成してWalletを追加
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Wallet, &gc.Wallet{Currency: 100})

	// 通貨を設定
	SetCurrency(world, player, 500)

	// 結果を検証
	currency := GetCurrency(world, player)
	assert.Equal(t, 500, currency, "通貨が500に設定されるべき")

	// クリーンアップ
	world.Manager.DeleteEntity(player)
}

func TestHasCurrency(t *testing.T) {
	t.Parallel()
	world, err := maingame.InitWorld(960, 720)
	require.NoError(t, err)

	// プレイヤーを作成してWalletを追加
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Wallet, &gc.Wallet{Currency: 100})

	// 通貨チェック
	assert.True(t, HasCurrency(world, player, 50), "50以上持っているべき")
	assert.True(t, HasCurrency(world, player, 100), "100以上持っているべき")
	assert.False(t, HasCurrency(world, player, 101), "101以上は持っていない")

	// クリーンアップ
	world.Manager.DeleteEntity(player)
}

func TestConsumeCurrency(t *testing.T) {
	t.Parallel()
	world, err := maingame.InitWorld(960, 720)
	require.NoError(t, err)

	// プレイヤーを作成してWalletを追加
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Wallet, &gc.Wallet{Currency: 100})

	// 通貨を消費（成功）
	success := ConsumeCurrency(world, player, 50)
	assert.True(t, success, "消費が成功するべき")
	assert.Equal(t, 50, GetCurrency(world, player), "残り50になるべき")

	// 通貨を消費（失敗：足りない）
	success = ConsumeCurrency(world, player, 100)
	assert.False(t, success, "消費が失敗するべき")
	assert.Equal(t, 50, GetCurrency(world, player), "残り50のまま変わらないべき")

	// 通貨を消費（成功：ちょうど）
	success = ConsumeCurrency(world, player, 50)
	assert.True(t, success, "消費が成功するべき")
	assert.Equal(t, 0, GetCurrency(world, player), "残り0になるべき")

	// クリーンアップ
	world.Manager.DeleteEntity(player)
}

func TestCurrencyOperationsWithoutWallet(t *testing.T) {
	t.Parallel()
	world, err := maingame.InitWorld(960, 720)
	require.NoError(t, err)

	// Walletを持たないエンティティ
	entity := world.Manager.NewEntity()

	// 各操作がエラーなく動作することを確認
	AddCurrency(world, entity, 100)
	assert.Equal(t, 0, GetCurrency(world, entity), "Walletがないので0")

	SetCurrency(world, entity, 200)
	assert.Equal(t, 0, GetCurrency(world, entity), "Walletがないので0")

	assert.False(t, HasCurrency(world, entity, 1), "Walletがないのでfalse")
	assert.False(t, ConsumeCurrency(world, entity, 1), "Walletがないのでfalse")

	// クリーンアップ
	world.Manager.DeleteEntity(entity)
}
