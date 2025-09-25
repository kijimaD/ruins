package states

import (
	"fmt"
	"log"

	gc "github.com/kijimaD/ruins/lib/components"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/mapbuilder"
	"github.com/kijimaD/ruins/lib/messagedata"
	"github.com/kijimaD/ruins/lib/save"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
)

// 各ステートのファクトリー関数を集約したファイル

// NewHomeMenuState は新しいHomeMenuStateインスタンスを作成するファクトリー関数
func NewHomeMenuState() es.State[w.World] {
	messageState := &MessageState{}

	messageData := messagedata.NewSystemMessage("拠点メニュー")

	// ホームメニューの選択肢を生成
	homeActions := []struct {
		label  string
		action func(w.World)
	}{
		{
			"出発",
			func(_ w.World) {
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewDungeonSelectState}})
			},
		},
		{
			"合成",
			func(_ w.World) {
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewCraftMenuState}})
			},
		},
		{
			"所持",
			func(_ w.World) {
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewInventoryMenuState}})
			},
		},
		{
			"装備",
			func(_ w.World) {
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewEquipMenuState}})
			},
		},
		{
			"書込",
			func(_ w.World) {
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewSaveMenuState}})
			},
		},
		{
			"終了",
			func(_ w.World) {
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewMainMenuState}})
			},
		},
	}

	// 各選択肢を追加
	for _, homeAction := range homeActions {
		actionCopy := homeAction.action // クロージャキャプチャ対策
		messageData = messageData.WithChoice(homeAction.label, actionCopy)
	}

	// MessageStateにMessageDataを設定
	messageState.messageData = messageData

	return messageState
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
	messageState := &MessageState{}

	messageData := messagedata.NewSystemMessage("デバッグメニュー")

	// デバッグ機能の一覧を生成
	debugActions := []struct {
		label  string
		action func(w.World)
	}{
		{
			"回復薬スポーン(インベントリ)",
			func(world w.World) {
				_, err := worldhelper.SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)
				if err != nil {
					log.Fatal("Error spawning item:", err.Error())
				}
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
			},
		},
		{
			"手榴弾スポーン(インベントリ)",
			func(world w.World) {
				_, err := worldhelper.SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
				if err != nil {
					log.Fatal("Error spawning item:", err.Error())
				}
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
			},
		},
		{
			"ゲームオーバー",
			func(_ w.World) {
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewGameOverMessageState}})
			},
		},
		{
			"ダンジョン開始(大部屋)",
			func(_ w.World) {
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{
					NewDungeonStateWithBuilder(1, mapbuilder.BuilderTypeBigRoom),
				}})
			},
		},
		{
			"ダンジョン開始(小部屋)",
			func(_ w.World) {
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{
					NewDungeonStateWithBuilder(1, mapbuilder.BuilderTypeSmallRoom),
				}})
			},
		},
		{
			"ダンジョン開始(洞窟)",
			func(_ w.World) {
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{
					NewDungeonStateWithBuilder(1, mapbuilder.BuilderTypeCave),
				}})
			},
		},
		{
			"ダンジョン開始(廃墟)",
			func(_ w.World) {
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{
					NewDungeonStateWithBuilder(1, mapbuilder.BuilderTypeRuins),
				}})
			},
		},
		{
			"ダンジョン開始(森)",
			func(_ w.World) {
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{
					NewDungeonStateWithBuilder(1, mapbuilder.BuilderTypeForest),
				}})
			},
		},
		{
			"メッセージ表示テスト",
			func(_ w.World) {
				testMessageData := messagedata.NewSystemMessage("ゲームが自動保存されました。\n\n進行状況は安全に記録されています。")
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{func() es.State[w.World] { return NewMessageState(testMessageData) }}})
			},
		},
		{
			"アイテム入手イベント",
			func(world w.World) {
				// アイテムを実際にインベントリに追加
				worldhelper.PlusAmount("鉄", 1, world)
				worldhelper.PlusAmount("木の棒", 1, world)
				worldhelper.PlusAmount("フェライトコア", 2, world)

				// アイテム入手完了後の表示用メッセージを生成
				messageText := "宝箱を発見した。\n\n" +
					"鉄を手に入れた。\n" +
					"木の棒を手に入れた。\n" +
					"フェライトコアを2個手に入れた。\n"

				itemMessageData := &messagedata.MessageData{
					Text:    messageText,
					Speaker: "",
				}
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{func() es.State[w.World] { return NewMessageState(itemMessageData) }}})
			},
		},
		{
			"長いメッセージテスト",
			func(_ w.World) {
				longText := `これは非常に長いメッセージのテストです。

メッセージウィンドウは自動的にサイズを調整し、
長いテキストでも適切に表示されることを確認しています。

複数行のテキストと改行が正しく処理されること、
そしてウィンドウの背景やボーダーが適切に描画されることを
このテストで検証できます。

日本語のテキストも問題なく表示されるはずです。
句読点、記号、数字123なども含めて確認してみましょう。`

				longMessageData := messagedata.NewSystemMessage(longText)
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{func() es.State[w.World] { return NewMessageState(longMessageData) }}})
			},
		},
		{
			"連鎖メッセージテスト",
			func(_ w.World) {
				chainMessageData := messagedata.NewSystemMessage("戦闘開始。").
					SystemMessage("剣と剣がぶつかり合う。").
					SystemMessage("勝利した。")

				messageState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{func() es.State[w.World] { return NewMessageState(chainMessageData) }}})
			},
		},
		{
			"選択肢分岐メッセージテスト",
			func(_ w.World) {
				battleMessage := messagedata.NewSystemMessage("戦闘した。")
				negotiateMessage := messagedata.NewSystemMessage("交渉した。")
				escapeMessage := messagedata.NewSystemMessage("逃走した。")

				choiceMessageData := messagedata.NewDialogMessage("敵に遭遇した。", "").
					WithChoiceMessage("戦う", battleMessage).
					WithChoiceMessage("交渉する", negotiateMessage).
					WithChoiceMessage("逃走する", escapeMessage)

				messageState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{func() es.State[w.World] { return NewMessageState(choiceMessageData) }}})
			},
		},
		{
			"選択肢処理テスト",
			func(_ w.World) {
				choiceAction1 := func() {
					println("実行: 1")
				}
				choiceAction2 := func() {
					println("実行: 2")
				}

				onCompleteAction := func() {
					println("Complete Action")
				}

				result1 := messagedata.NewSystemMessage("選択肢1を選びました。").
					SystemMessage("何かの処理が実行されました。").
					WithOnComplete(onCompleteAction)

				result2 := messagedata.NewSystemMessage("選択肢2を選びました。").
					SystemMessage("別の処理が実行されました。").
					WithOnComplete(onCompleteAction)

				testMessageData := messagedata.NewDialogMessage("処理のテストです。選択肢を選んでください。", "システム").
					WithChoiceMessage("処理1を実行", result1).
					WithChoiceMessage("処理2を実行", result2)

				testMessageData.Choices[0].Action = func(_ w.World) { choiceAction1() }
				testMessageData.Choices[1].Action = func(_ w.World) { choiceAction2() }

				messageState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{func() es.State[w.World] { return NewMessageState(testMessageData) }}})
			},
		},
	}

	// 各デバッグ機能の選択肢を追加
	for _, debugAction := range debugActions {
		actionCopy := debugAction.action // クロージャキャプチャ対策
		messageData = messageData.WithChoice(debugAction.label, actionCopy)
	}

	// 閉じる選択肢を追加
	messageData = messageData.WithChoice("閉じる", func(_ w.World) {
		messageState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
	})

	// MessageStateにMessageDataを設定
	messageState.messageData = messageData

	return messageState
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
	messageState := &MessageState{}

	// セーブマネージャーで現在のスロット状態を取得
	saveManager := save.NewSerializationManager("./saves")

	messageData := messagedata.NewSystemMessage("どのスロットにセーブしますか？")

	// 各スロットの状態を確認して選択肢を動的に生成
	for i := 1; i <= 3; i++ {
		slotName := fmt.Sprintf("slot%d", i)
		var label string

		if saveManager.SaveFileExists(slotName) {
			if timestamp, err := saveManager.GetSaveFileTimestamp(slotName); err == nil {
				label = fmt.Sprintf("スロット%d [%s]", i, timestamp.Format("01/02 15:04"))
			} else {
				label = fmt.Sprintf("スロット%d [データあり]", i)
			}
		} else {
			label = fmt.Sprintf("スロット%d [空]", i)
		}

		slotNameCopy := slotName // クロージャキャプチャ対策
		messageData = messageData.WithChoice(label, func(world w.World) {
			if err := saveManager.SaveWorld(world, slotNameCopy); err != nil {
				log.Fatal("Save failed:", err.Error())
			}
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
		})
	}

	messageData = messageData.WithChoice("戻る", func(_ w.World) {
		messageState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
	})

	messageState.messageData = messageData
	return messageState
}

// NewLoadMenuState は新しいLoadMenuStateインスタンスを作成するファクトリー関数
func NewLoadMenuState() es.State[w.World] {
	messageState := &MessageState{}

	// セーブマネージャーで現在のスロット状態を取得
	saveManager := save.NewSerializationManager("./saves")

	messageData := messagedata.NewSystemMessage("どのスロットから読み込みますか？")

	// 各スロットの状態を確認して選択肢を動的に生成
	hasValidSlot := false
	for i := 1; i <= 3; i++ {
		slotName := fmt.Sprintf("slot%d", i)
		var label string

		if saveManager.SaveFileExists(slotName) {
			hasValidSlot = true
			if timestamp, err := saveManager.GetSaveFileTimestamp(slotName); err == nil {
				label = fmt.Sprintf("スロット%d [%s]", i, timestamp.Format("01/02 15:04"))
			} else {
				label = fmt.Sprintf("スロット%d [データあり]", i)
			}

			slotNameCopy := slotName // クロージャキャプチャ対策
			messageData = messageData.WithChoice(label, func(world w.World) {
				// ロードを実行
				err := saveManager.LoadWorld(world, slotNameCopy)
				if err != nil {
					println("Load failed:", err.Error())
					messageState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
					return
				}
				// ロード成功後、ホームメニューに遷移
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewHomeMenuState}})
			})
		}
	}

	// 有効なセーブデータが存在しない場合の処理
	if !hasValidSlot {
		messageData = messageData.WithChoice("セーブデータがありません", func(_ w.World) {
			// 何もしない（選択不可を示すためのダミー選択肢）
		})
	}

	messageData = messageData.WithChoice("戻る", func(_ w.World) {
		messageState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
	})

	// MessageStateにMessageDataを設定
	messageState.messageData = messageData

	return messageState
}

// NewMessageState はメッセージデータを受け取って新しいMessageStateを作成するファクトリー関数
func NewMessageState(messageData *messagedata.MessageData) es.State[w.World] {
	return &MessageState{
		messageData: messageData,
	}
}
