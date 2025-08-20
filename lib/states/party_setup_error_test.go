package states

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPartySetupState_InitializePartyData_NoPlayerComponent(t *testing.T) {
	t.Parallel()

	// テスト用のワールドを作成
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// Playerコンポーネントを持たないエンティティのみ作成
	ally := world.Manager.NewEntity()
	ally.AddComponent(world.Components.FactionAlly, &gc.FactionAllyData{})
	ally.AddComponent(world.Components.Name, &gc.Name{Name: "味方"})
	ally.AddComponent(world.Components.InParty, &gc.InParty{})
	ally.AddComponent(world.Components.Attributes, &gc.Attributes{
		Vitality:  gc.Attribute{Base: 10},
		Strength:  gc.Attribute{Base: 12},
		Sensation: gc.Attribute{Base: 8},
		Dexterity: gc.Attribute{Base: 9},
		Agility:   gc.Attribute{Base: 11},
		Defense:   gc.Attribute{Base: 7},
	})
	ally.AddComponent(world.Components.Pools, &gc.Pools{
		Level: 1,
		XP:    0,
	})

	// PartySetupStateを作成
	state := &PartySetupState{}

	// initializePartyDataを呼び出すとpanicが発生することを確認
	assert.Panics(t, func() {
		state.initializePartyData(world)
	}, "Playerコンポーネントが見つからない場合はpanicが発生するべき")
}

func TestPartySetupState_InitializePartyData_WithPlayerComponent(t *testing.T) {
	t.Parallel()

	// テスト用のワールドを作成
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// Playerコンポーネントを持つエンティティを作成
	protagonist := world.Manager.NewEntity()
	protagonist.AddComponent(world.Components.FactionAlly, &gc.FactionAllyData{})
	protagonist.AddComponent(world.Components.Name, &gc.Name{Name: "主人公"})
	protagonist.AddComponent(world.Components.InParty, &gc.InParty{})
	protagonist.AddComponent(world.Components.Player, &gc.Player{})
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

	// PartySetupStateを作成
	state := &PartySetupState{}

	// initializePartyDataを呼び出してもpanicが発生しないことを確認
	assert.NotPanics(t, func() {
		state.initializePartyData(world)
	}, "Playerコンポーネントがある場合はpanicが発生しないべき")

	// 主人公が正しく設定されていることを確認
	assert.Equal(t, protagonist, state.protagonistEntity)
}
