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
	messageState := &MessageState{}

	// ダンジョン選択メッセージを作成
	messageData := messagedata.NewSystemMessage("ダンジョン選択").
		WithChoice("森の遺跡", func(_ w.World) {
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{NewDungeonStateWithDepth(1)}})
		}).
		WithChoice("山の遺跡", func(_ w.World) {
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{NewDungeonStateWithDepth(1)}})
		}).
		WithChoice("塔の遺跡", func(_ w.World) {
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{NewDungeonStateWithDepth(1)}})
		}).
		WithChoice("拠点メニューに戻る", func(_ w.World) {
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewHomeMenuState}})
		})

	// MessageStateにMessageDataを設定
	messageState.messageData = messageData

	return messageState
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
	messageState := &MessageState{}

	// ダンジョンメニューのメッセージを作成
	messageData := messagedata.NewSystemMessage("どうしますか？").
		WithChoice("脱出", func(_ w.World) {
			// MessageStateで直接HomeMenuStateに遷移
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewHomeMenuState}})
		}).
		WithChoice("閉じる", func(_ w.World) {
			// 何も処理しない（デフォルトのTransPopが適用される）
		})

	// MessageStateにMessageDataを設定
	messageState.messageData = messageData

	return messageState
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

// NewGameOverMessageState はゲームオーバー用のMessageStateを作成するファクトリー関数
func NewGameOverMessageState() es.State[w.World] {
	// MessageStateインスタンスを作成
	messageState := &MessageState{}

	// ゲームオーバーメッセージを作成（選択肢付き）
	messageData := messagedata.NewSystemMessage("死亡した。").
		WithChoice("メインメニューに戻る", func(_ w.World) {
			// メインメニューに遷移
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewMainMenuState}})
		})

	// MessageStateにMessageDataを設定
	messageState.messageData = messageData

	return messageState
}

// NewSaveMenuState は新しいSaveMenuStateインスタンスを作成するファクトリー関数
func NewSaveMenuState() es.State[w.World] {
	return &SaveMenuState{}
}

// NewLoadMenuState は新しいLoadMenuStateインスタンスを作成するファクトリー関数
func NewLoadMenuState() es.State[w.World] {
	return &LoadMenuState{}
}

// NewMessageState はメッセージデータを受け取って新しいMessageStateを作成するファクトリー関数
func NewMessageState(messageData *messagedata.MessageData) es.State[w.World] {
	return &MessageState{
		messageData: messageData,
	}
}
