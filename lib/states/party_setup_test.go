package states

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestPartySetupState_InitializePartyData(t *testing.T) {
	t.Parallel()

	// テスト用のワールドを作成
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// テスト用メンバーを作成
	protagonist := world.Manager.NewEntity()
	protagonist.AddComponent(world.Components.FactionAlly, &gc.FactionAllyData{})
	protagonist.AddComponent(world.Components.Name, &gc.Name{Name: "主人公"})
	protagonist.AddComponent(world.Components.InParty, &gc.InParty{})
	protagonist.AddComponent(world.Components.Attributes, &gc.Attributes{
		Vitality:  gc.Attribute{Base: 10},
		Strength:  gc.Attribute{Base: 12},
		Sensation: gc.Attribute{Base: 8},
		Dexterity: gc.Attribute{Base: 9},
		Agility:   gc.Attribute{Base: 11},
		Defense:   gc.Attribute{Base: 7},
	})
	protagonist.AddComponent(world.Components.Pools, &gc.Pools{
		Level: 1,
		XP:    0,
	})
	protagonist.AddComponent(world.Components.Player, &gc.Player{})

	ally1 := world.Manager.NewEntity()
	ally1.AddComponent(world.Components.FactionAlly, &gc.FactionAllyData{})
	ally1.AddComponent(world.Components.Name, &gc.Name{Name: "アリス"})
	ally1.AddComponent(world.Components.InParty, &gc.InParty{})
	ally1.AddComponent(world.Components.Attributes, &gc.Attributes{
		Vitality:  gc.Attribute{Base: 15},
		Strength:  gc.Attribute{Base: 18},
		Sensation: gc.Attribute{Base: 10},
		Dexterity: gc.Attribute{Base: 12},
		Agility:   gc.Attribute{Base: 14},
		Defense:   gc.Attribute{Base: 20},
	})
	ally1.AddComponent(world.Components.Pools, &gc.Pools{
		Level: 5,
		XP:    150,
	})

	ally2 := world.Manager.NewEntity()
	ally2.AddComponent(world.Components.FactionAlly, &gc.FactionAllyData{})
	ally2.AddComponent(world.Components.Name, &gc.Name{Name: "ボブ"})
	ally2.AddComponent(world.Components.Attributes, &gc.Attributes{
		Vitality:  gc.Attribute{Base: 12},
		Strength:  gc.Attribute{Base: 8},
		Sensation: gc.Attribute{Base: 16},
		Dexterity: gc.Attribute{Base: 14},
		Agility:   gc.Attribute{Base: 10},
		Defense:   gc.Attribute{Base: 5},
	})
	ally2.AddComponent(world.Components.Pools, &gc.Pools{
		Level: 3,
		XP:    75,
	})

	// PartySetupStateを作成
	state := &PartySetupState{}
	state.initializePartyData(world)

	// 主人公が正しく設定されているかテスト
	assert.Equal(t, protagonist, state.protagonistEntity)

	// パーティスロットが正しく設定されているかテスト
	assert.NotNil(t, state.currentPartySlots[0])
	assert.Equal(t, protagonist, *state.currentPartySlots[0])

	// アリスが2番目のスロットに配置されているかテスト
	assert.NotNil(t, state.currentPartySlots[1])
	assert.Equal(t, ally1, *state.currentPartySlots[1])

	// ボブはパーティに参加していないことをテスト
	assert.Nil(t, state.currentPartySlots[2])
	assert.Nil(t, state.currentPartySlots[3])
}

func TestPartySetupState_CreateMenuItems(t *testing.T) {
	t.Parallel()

	// テスト用のワールドを作成
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// テスト用メンバーを作成
	protagonist := world.Manager.NewEntity()
	protagonist.AddComponent(world.Components.FactionAlly, &gc.FactionAllyData{})
	protagonist.AddComponent(world.Components.Name, &gc.Name{Name: "主人公"})
	protagonist.AddComponent(world.Components.InParty, &gc.InParty{})
	protagonist.AddComponent(world.Components.Pools, &gc.Pools{
		Level: 1,
		XP:    0,
	})
	protagonist.AddComponent(world.Components.Attributes, &gc.Attributes{
		Vitality: gc.Attribute{Base: 10},
	})
	protagonist.AddComponent(world.Components.Player, &gc.Player{})

	ally := world.Manager.NewEntity()
	ally.AddComponent(world.Components.FactionAlly, &gc.FactionAllyData{})
	ally.AddComponent(world.Components.Name, &gc.Name{Name: "アリス"})
	ally.AddComponent(world.Components.Pools, &gc.Pools{
		Level: 5,
		XP:    150,
	})
	ally.AddComponent(world.Components.Attributes, &gc.Attributes{
		Vitality: gc.Attribute{Base: 15},
	})

	// PartySetupStateを作成
	state := &PartySetupState{
		protagonistEntity:    protagonist,
		currentPartySlots:    [4]*ecs.Entity{&protagonist, nil, nil, nil},
		selectedMemberEntity: nil,
	}

	// メニュー項目を作成
	items := state.createMenuItems(world)

	// 最低限のメニュー項目数をテスト（ヘッダー + メンバー2人 + セパレータ + 操作2つ）
	expectedMinItems := 1 + 2 + 1 + 2 // 最小でも6項目
	assert.GreaterOrEqual(t, len(items), expectedMinItems)

	// 主人公が暗転表示されているかテスト
	foundProtagonist := false
	for _, item := range items {
		if item.Label == "[主人公] 主人公" {
			foundProtagonist = true
			assert.False(t, item.Disabled, "主人公が無効化されている（選択可能でなければならない）")
			assert.True(t, item.Dimmed, "主人公が暗転表示されていない")
			break
		}
	}
	assert.True(t, foundProtagonist, "主人公がメニュー項目に含まれていない")

	// アリスが待機状態で含まれているかテスト
	foundAlice := false
	for _, item := range items {
		if item.Label == "[ 待機 ] アリス" {
			foundAlice = true
			assert.False(t, item.Disabled, "アリスが無効化されている")
			break
		}
	}
	assert.True(t, foundAlice, "アリスがメニュー項目に含まれていない")
}

func TestPartySetupState_ApplyPartyChanges(t *testing.T) {
	t.Parallel()

	// テスト用のワールドを作成
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// テスト用メンバーを作成
	protagonist := world.Manager.NewEntity()
	protagonist.AddComponent(world.Components.FactionAlly, &gc.FactionAllyData{})
	protagonist.AddComponent(world.Components.Name, &gc.Name{Name: "主人公"})
	protagonist.AddComponent(world.Components.InParty, &gc.InParty{})

	ally1 := world.Manager.NewEntity()
	ally1.AddComponent(world.Components.FactionAlly, &gc.FactionAllyData{})
	ally1.AddComponent(world.Components.Name, &gc.Name{Name: "アリス"})

	ally2 := world.Manager.NewEntity()
	ally2.AddComponent(world.Components.FactionAlly, &gc.FactionAllyData{})
	ally2.AddComponent(world.Components.Name, &gc.Name{Name: "ボブ"})
	ally2.AddComponent(world.Components.InParty, &gc.InParty{}) // 既にパーティに参加

	// PartySetupStateを作成
	state := &PartySetupState{
		protagonistEntity: protagonist,
		currentPartySlots: [4]*ecs.Entity{&protagonist, &ally1, nil, nil},
	}

	// パーティ変更を適用
	state.applyPartyChanges(world)

	// 主人公がInPartyを持っているかテスト
	assert.True(t, protagonist.HasComponent(world.Components.InParty))

	// アリスがInPartyを持っているかテスト
	assert.True(t, ally1.HasComponent(world.Components.InParty))

	// ボブがInPartyを持っていないかテスト（パーティから外された）
	assert.False(t, ally2.HasComponent(world.Components.InParty))
}

func TestPartySetupState_HandleMemberSelection_AddToParty(t *testing.T) {
	t.Parallel()

	// テスト用のワールドを作成
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// テスト用メンバーを作成
	protagonist := world.Manager.NewEntity()
	protagonist.AddComponent(world.Components.FactionAlly, &gc.FactionAllyData{})
	protagonist.AddComponent(world.Components.Name, &gc.Name{Name: "主人公"})
	protagonist.AddComponent(world.Components.Pools, &gc.Pools{
		Level: 1,
		XP:    0,
	})
	protagonist.AddComponent(world.Components.Player, &gc.Player{})

	ally := world.Manager.NewEntity()
	ally.AddComponent(world.Components.FactionAlly, &gc.FactionAllyData{})
	ally.AddComponent(world.Components.Name, &gc.Name{Name: "アリス"})
	ally.AddComponent(world.Components.Pools, &gc.Pools{
		Level: 5,
		XP:    150,
	})

	// PartySetupStateを作成
	state := &PartySetupState{
		protagonistEntity: protagonist,
		currentPartySlots: [4]*ecs.Entity{&protagonist, nil, nil, nil},
	}

	// 未加入メンバーを選択（パーティに追加）
	userData := map[string]interface{}{
		"type":       "member",
		"entity":     ally,
		"in_party":   false,
		"slot_index": -1,
	}

	// 選択前の状態をテスト
	assert.Nil(t, state.currentPartySlots[1])

	// 選択を実行
	state.handleMemberSelection(world, userData)

	// 選択後の状態をテスト
	assert.NotNil(t, state.currentPartySlots[1])
	assert.Equal(t, ally, *state.currentPartySlots[1])
}

func TestPartySetupState_HandleMemberSelection_RemoveFromParty(t *testing.T) {
	t.Parallel()

	// テスト用のワールドを作成
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// テスト用メンバーを作成
	protagonist := world.Manager.NewEntity()
	protagonist.AddComponent(world.Components.FactionAlly, &gc.FactionAllyData{})
	protagonist.AddComponent(world.Components.Name, &gc.Name{Name: "主人公"})
	protagonist.AddComponent(world.Components.Pools, &gc.Pools{
		Level: 1,
		XP:    0,
	})
	protagonist.AddComponent(world.Components.Player, &gc.Player{})

	ally := world.Manager.NewEntity()
	ally.AddComponent(world.Components.FactionAlly, &gc.FactionAllyData{})
	ally.AddComponent(world.Components.Name, &gc.Name{Name: "アリス"})
	ally.AddComponent(world.Components.Pools, &gc.Pools{
		Level: 5,
		XP:    150,
	})

	// PartySetupStateを作成
	state := &PartySetupState{
		protagonistEntity: protagonist,
		currentPartySlots: [4]*ecs.Entity{&protagonist, &ally, nil, nil},
	}

	// 加入済みメンバーを選択（パーティから外す）
	userData := map[string]interface{}{
		"type":       "member",
		"entity":     ally,
		"in_party":   true,
		"slot_index": 1,
	}

	// 選択前の状態をテスト
	assert.NotNil(t, state.currentPartySlots[1])

	// 選択を実行
	state.handleMemberSelection(world, userData)

	// 選択後の状態をテスト
	assert.Nil(t, state.currentPartySlots[1])
}

func TestPartySetupState_ProtagonistSlotProtection(t *testing.T) {
	t.Parallel()

	// テスト用のワールドを作成
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// テスト用メンバーを作成
	protagonist := world.Manager.NewEntity()
	protagonist.AddComponent(world.Components.FactionAlly, &gc.FactionAllyData{})
	protagonist.AddComponent(world.Components.Name, &gc.Name{Name: "主人公"})
	protagonist.AddComponent(world.Components.Pools, &gc.Pools{
		Level: 1,
		XP:    0,
	})
	protagonist.AddComponent(world.Components.Player, &gc.Player{})

	// PartySetupStateを作成
	state := &PartySetupState{
		protagonistEntity: protagonist,
		currentPartySlots: [4]*ecs.Entity{&protagonist, nil, nil, nil},
	}

	// 主人公を選択しようとする
	userData := map[string]interface{}{
		"type":       "member",
		"entity":     protagonist,
		"in_party":   true,
		"slot_index": 0,
	}

	// 選択前の状態をテスト
	assert.NotNil(t, state.currentPartySlots[0])
	assert.Equal(t, protagonist, *state.currentPartySlots[0])

	// 選択を実行（何も起こらないはず）
	state.handleMemberSelection(world, userData)

	// 選択後も主人公は残っているべき
	assert.NotNil(t, state.currentPartySlots[0])
	assert.Equal(t, protagonist, *state.currentPartySlots[0])
}
