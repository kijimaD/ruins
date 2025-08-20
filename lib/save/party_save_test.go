package save

import (
	"fmt"
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
	fmt.Println("=== セーブ前のパーティ状態 ===")

	originalPartyMembers := []ecs.Entity{}

	// プレイヤーを探す
	world.Manager.Join(world.Components.Player).Visit(ecs.Visit(func(entity ecs.Entity) {
		fmt.Printf("Player: %d", entity)
		if entity.HasComponent(world.Components.Name) {
			name := world.Components.Name.Get(entity)
			fmt.Printf(" (%+v)", name)
		}
		fmt.Println()
	}))

	// パーティメンバーを確認
	world.Manager.Join(world.Components.InParty, world.Components.FactionAlly).Visit(ecs.Visit(func(entity ecs.Entity) {
		originalPartyMembers = append(originalPartyMembers, entity)
		fmt.Printf("Party member: %d", entity)
		if entity.HasComponent(world.Components.Name) {
			name := world.Components.Name.Get(entity)
			fmt.Printf(" (%+v)", name)
		}
		if entity.HasComponent(world.Components.Player) {
			fmt.Printf(" [PLAYER]")
		}
		fmt.Println()
	}))

	fmt.Printf("Original party size: %d\n", len(originalPartyMembers))
	require.Greater(t, len(originalPartyMembers), 0, "Should have party members")

	// パーティからメンバーを一人外す（プレイヤー以外）
	var removedMember ecs.Entity
	for _, member := range originalPartyMembers {
		if !member.HasComponent(world.Components.Player) {
			member.RemoveComponent(world.Components.InParty)
			removedMember = member
			fmt.Printf("Removed from party: %d\n", member)
			break
		}
	}

	// 変更後のパーティ状態を確認
	fmt.Println("\n=== パーティ変更後の状態 ===")
	modifiedPartyMembers := []ecs.Entity{}
	world.Manager.Join(world.Components.InParty, world.Components.FactionAlly).Visit(ecs.Visit(func(entity ecs.Entity) {
		modifiedPartyMembers = append(modifiedPartyMembers, entity)
		fmt.Printf("Party member: %d", entity)
		if entity.HasComponent(world.Components.Name) {
			name := world.Components.Name.Get(entity)
			fmt.Printf(" (%+v)", name)
		}
		if entity.HasComponent(world.Components.Player) {
			fmt.Printf(" [PLAYER]")
		}
		fmt.Println()
	}))
	fmt.Printf("Modified party size: %d\n", len(modifiedPartyMembers))

	// セーブ実行
	saveManager := NewSerializationManager(testDir)
	err = saveManager.SaveWorld(world, "party_test")
	require.NoError(t, err)

	// 新しいワールドでロード
	fmt.Println("\n=== ロード後の状態 ===")
	newWorld, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	err = saveManager.LoadWorld(newWorld, "party_test")
	require.NoError(t, err)

	// ロード後のプレイヤー確認
	loadedPlayerEntity := ecs.Entity(0)
	newWorld.Manager.Join(newWorld.Components.Player).Visit(ecs.Visit(func(entity ecs.Entity) {
		loadedPlayerEntity = entity
		fmt.Printf("Loaded Player: %d", entity)
		if entity.HasComponent(newWorld.Components.Name) {
			name := newWorld.Components.Name.Get(entity)
			fmt.Printf(" (%+v)", name)
		}
		fmt.Println()
	}))

	// ロード後のパーティメンバー確認
	loadedPartyMembers := []ecs.Entity{}
	newWorld.Manager.Join(newWorld.Components.InParty, newWorld.Components.FactionAlly).Visit(ecs.Visit(func(entity ecs.Entity) {
		loadedPartyMembers = append(loadedPartyMembers, entity)
		fmt.Printf("Loaded party member: %d", entity)
		if entity.HasComponent(newWorld.Components.Name) {
			name := newWorld.Components.Name.Get(entity)
			fmt.Printf(" (%+v)", name)
		}
		if entity.HasComponent(newWorld.Components.Player) {
			fmt.Printf(" [PLAYER]")
		}
		fmt.Println()
	}))
	fmt.Printf("Loaded party size: %d\n", len(loadedPartyMembers))

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

	fmt.Println("\n✓ パーティ編成の保存・復元が正常に完了")
}
