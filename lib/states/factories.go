package states

import (
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/mapbuilder"
	"github.com/kijimaD/ruins/lib/messagedata"
	w "github.com/kijimaD/ruins/lib/world"
)

// 各ステートのファクトリー関数を集約したファイル

// NewHomeMenuState は新しいHomeMenuStateインスタンスを作成するファクトリー関数
func NewHomeMenuState() es.State[w.World] {
	return &HomeMenuState{}
}

// NewDungeonSelectState は新しいDungeonSelectStateインスタンスを作成するファクトリー関数
func NewDungeonSelectState() es.State[w.World] {
	return &DungeonSelectState{}
}

// NewCraftMenuState は新しいCraftMenuStateインスタンスを作成するファクトリー関数
func NewCraftMenuState() es.State[w.World] {
	return &CraftMenuState{}
}

// NewInventoryMenuState は新しいInventoryMenuStateインスタンスを作成するファクトリー関数
func NewInventoryMenuState() es.State[w.World] {
	return &InventoryMenuState{}
}

// NewEquipMenuState は新しいEquipMenuStateインスタンスを作成するファクトリー関数
func NewEquipMenuState() es.State[w.World] {
	return &EquipMenuState{}
}

// NewDebugMenuState は新しいDebugMenuStateインスタンスを作成するファクトリー関数
func NewDebugMenuState() es.State[w.World] {
	return &DebugMenuState{}
}

// NewDungeonMenuState は新しいDungeonMenuStateインスタンスを作成するファクトリー関数
func NewDungeonMenuState() es.State[w.World] {
	return &DungeonMenuState{}
}

// NewDungeonStateWithDepth は指定されたDepthでDungeonStateインスタンスを作成するファクトリー関数
func NewDungeonStateWithDepth(depth int) es.StateFactory[w.World] {
	return func() es.State[w.World] {
		return &DungeonState{Depth: depth, BuilderType: mapbuilder.BuilderTypeRandom}
	}
}

// NewDungeonStateWithSeed は指定されたDepthとSeedでDungeonStateインスタンスを作成するファクトリー関数
func NewDungeonStateWithSeed(depth int, seed uint64) es.StateFactory[w.World] {
	return func() es.State[w.World] {
		return &DungeonState{Depth: depth, Seed: seed, BuilderType: mapbuilder.BuilderTypeRandom}
	}
}

// NewDungeonStateWithBuilder は指定されたBuilderTypeでDungeonStateインスタンスを作成するファクトリー関数
func NewDungeonStateWithBuilder(depth int, builderType mapbuilder.BuilderType) es.StateFactory[w.World] {
	return func() es.State[w.World] {
		return &DungeonState{Depth: depth, BuilderType: builderType}
	}
}

// NewMainMenuState は新しいMainMenuStateインスタンスを作成するファクトリー関数
func NewMainMenuState() es.State[w.World] {
	return &MainMenuState{}
}

// NewGameOverState は新しいGameOverStateインスタンスを作成するファクトリー関数
func NewGameOverState() es.State[w.World] {
	return &GameOverState{}
}

// NewSaveMenuState は新しいSaveMenuStateインスタンスを作成するファクトリー関数
func NewSaveMenuState() es.State[w.World] {
	return &SaveMenuState{}
}

// NewLoadMenuState は新しいLoadMenuStateインスタンスを作成するファクトリー関数
func NewLoadMenuState() es.State[w.World] {
	return &LoadMenuState{}
}

// NewMessageWindowState はメッセージデータを受け取って新しいMessageWindowStateを作成するファクトリー関数
func NewMessageWindowState(messageData *messagedata.MessageData) es.State[w.World] {
	return &MessageWindowState{
		messageData: messageData,
	}
}
