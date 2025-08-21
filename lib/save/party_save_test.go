package save

import (
	"os"
	"testing"

	"github.com/kijimaD/ruins/lib/game"
	"github.com/kijimaD/ruins/lib/worldhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestPartySaveAndLoad(t *testing.T) {
	t.Parallel()
	// テスト用ディレクトリを準備
	testDir := "./test_party_save"
	defer func() {
		_ = os.RemoveAll(testDir)
	}()

	// ワールドを作成してデバッグデータを初期化
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)
	worldhelper.InitDebugData(world)

	// セーブ前のパーティ状態を確認
	originalPartyMembers := []ecs.Entity{}

	// パーティメンバーを確認
	world.Manager.Join(world.Components.InParty, world.Components.FactionAlly).Visit(ecs.Visit(func(entity ecs.Entity) {
		originalPartyMembers = append(originalPartyMembers, entity)
	}))
	require.Greater(t, len(originalPartyMembers), 0, "Should have party members")

	// パーティからメンバーを一人外す（プレイヤー以外）
	var removedMember ecs.Entity
	for _, member := range originalPartyMembers {
		if !member.HasComponent(world.Components.Player) {
			member.RemoveComponent(world.Components.InParty)
			removedMember = member
			break
		}
	}

	// 変更後のパーティ状態を確認
	modifiedPartyMembers := []ecs.Entity{}
	world.Manager.Join(world.Components.InParty, world.Components.FactionAlly).Visit(ecs.Visit(func(entity ecs.Entity) {
		modifiedPartyMembers = append(modifiedPartyMembers, entity)
	}))

	// セーブ実行
	saveManager := NewSerializationManager(testDir)
	err = saveManager.SaveWorld(world, "party_test")
	require.NoError(t, err)

	// 新しいワールドでロード
	newWorld, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	err = saveManager.LoadWorld(newWorld, "party_test")
	require.NoError(t, err)

	// ロード後のプレイヤー確認
	loadedPlayerEntity := ecs.Entity(0)
	newWorld.Manager.Join(newWorld.Components.Player).Visit(ecs.Visit(func(entity ecs.Entity) {
		loadedPlayerEntity = entity
	}))

	// ロード後のパーティメンバー確認
	loadedPartyMembers := []ecs.Entity{}
	newWorld.Manager.Join(newWorld.Components.InParty, newWorld.Components.FactionAlly).Visit(ecs.Visit(func(entity ecs.Entity) {
		loadedPartyMembers = append(loadedPartyMembers, entity)
	}))

	// 検証
	// loadedPlayerEntityが有効かチェック（プレイヤーコンポーネントがあるかで判定）
	assert.True(t, loadedPlayerEntity.HasComponent(newWorld.Components.Player), "Player should be loaded with Player component")
	assert.Equal(t, len(modifiedPartyMembers), len(loadedPartyMembers), "Party size should be preserved")

	// プレイヤーがパーティにいることを確認
	playerInParty := false
	for _, member := range loadedPartyMembers {
		if member.HasComponent(newWorld.Components.Player) {
			playerInParty = true
			break
		}
	}
	assert.True(t, playerInParty, "Player should be in party after load")

	// 外したメンバーがパーティにいないことを確認
	// (エンティティIDは変わるので名前で確認)
	if removedMember != ecs.Entity(0) && removedMember.HasComponent(world.Components.Name) {
		removedName := world.Components.Name.Get(removedMember)

		removedMemberInLoadedParty := false
		for _, member := range loadedPartyMembers {
			if member.HasComponent(newWorld.Components.Name) {
				name := newWorld.Components.Name.Get(member)
				if name == removedName {
					removedMemberInLoadedParty = true
					break
				}
			}
		}
		assert.False(t, removedMemberInLoadedParty, "Removed member should not be in party after load")
	}
}
