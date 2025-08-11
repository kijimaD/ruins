package states

import (
	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/world"
)

// 各ステートのファクトリー関数を集約したファイル

// NewIntroState は新しいIntroStateインスタンスを作成するファクトリー関数
func NewIntroState() es.State {
	return &IntroState{}
}

// NewBattleState は新しいBattleStateインスタンスを作成するファクトリー関数
func NewBattleState() es.State {
	return &BattleState{}
}

// NewBattleStateWithEnemies は指定された敵でBattleStateを作成するファクトリー関数
func NewBattleStateWithEnemies(enemies []string) es.StateFactory {
	return func() es.State {
		return &BattleState{
			FixedEnemies: enemies,
		}
	}
}

// NewHomeMenuState は新しいHomeMenuStateインスタンスを作成するファクトリー関数
func NewHomeMenuState() es.State {
	return &HomeMenuState{}
}

// NewDungeonSelectState は新しいDungeonSelectStateインスタンスを作成するファクトリー関数
func NewDungeonSelectState() es.State {
	return &DungeonSelectState{}
}

// NewCraftMenuState は新しいCraftMenuStateインスタンスを作成するファクトリー関数
func NewCraftMenuState() es.State {
	return &CraftMenuState{}
}

// NewInventoryMenuState は新しいInventoryMenuStateインスタンスを作成するファクトリー関数
func NewInventoryMenuState() es.State {
	return &InventoryMenuState{}
}

// NewEquipMenuState は新しいEquipMenuStateインスタンスを作成するファクトリー関数
func NewEquipMenuState() es.State {
	return &EquipMenuState{}
}

// NewDebugMenuState は新しいDebugMenuStateインスタンスを作成するファクトリー関数
func NewDebugMenuState() es.State {
	return &DebugMenuState{}
}

// NewDungeonMenuState は新しいDungeonMenuStateインスタンスを作成するファクトリー関数
func NewDungeonMenuState() es.State {
	return &DungeonMenuState{}
}

// NewDungeonStateWithDepth は指定されたDepthでDungeonStateインスタンスを作成するファクトリー関数
func NewDungeonStateWithDepth(depth int) es.StateFactory {
	return func() es.State {
		return &DungeonState{Depth: depth}
	}
}

// NewMessageStateWithText は指定されたテキストでMessageStateインスタンスを作成するファクトリー関数
func NewMessageStateWithText(text string) es.StateFactory {
	return func() es.State {
		return &MessageState{
			text: text,
		}
	}
}

// NewExecStateWithFunc は指定された関数でExecStateインスタンスを作成するファクトリー関数
func NewExecStateWithFunc(f func(w.World)) es.StateFactory {
	return func() es.State {
		return NewExecState(f)
	}
}

// NewMainMenuState は新しいMainMenuStateインスタンスを作成するファクトリー関数
func NewMainMenuState() es.State {
	return &MainMenuState{}
}

// NewGameOverState は新しいGameOverStateインスタンスを作成するファクトリー関数
func NewGameOverState() es.State {
	return &GameOverState{}
}
