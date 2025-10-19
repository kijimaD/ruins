package states

import (
	"fmt"

	"github.com/kijimaD/ruins/lib/actions"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/config"
	es "github.com/kijimaD/ruins/lib/engine/states"
	mapplanner "github.com/kijimaD/ruins/lib/mapplanner"
	"github.com/kijimaD/ruins/lib/messagedata"
	"github.com/kijimaD/ruins/lib/save"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
)

// 各ステートのファクトリー関数を集約したファイル

// NewDungeonMenuState は新しいDungeonMenuStateインスタンスを作成するファクトリー関数
func NewDungeonMenuState() es.State[w.World] {
	persistentState := NewPersistentMessageState(nil)

	persistentState.messageData = messagedata.NewSystemMessage("").
		WithChoice("合成", func(_ w.World) error {
			persistentState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewCraftMenuState}})
			return nil
		}).
		WithChoice("所持", func(_ w.World) error {
			persistentState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewInventoryMenuState}})
			return nil
		}).
		WithChoice("装備", func(_ w.World) error {
			persistentState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewEquipMenuState}})
			return nil
		}).
		WithChoice("店", func(_ w.World) error {
			persistentState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewShopMenuState}})
			return nil
		}).
		WithChoice("書込", func(_ w.World) error {
			persistentState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewSaveMenuState}})
			return nil
		}).
		WithChoice("終了", func(_ w.World) error {
			persistentState.SetTransition(es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewMainMenuState}})
			return nil
		}).
		WithChoice("閉じる", func(_ w.World) error {
			persistentState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
			return nil
		})

	return persistentState
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

	messageState.messageData = messagedata.NewSystemMessage("デバッグメニュー").
		WithChoice("回復薬スポーン(インベントリ)", func(world w.World) error {
			_, err := worldhelper.SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)
			if err != nil {
				return fmt.Errorf("error spawning item: %w", err)
			}
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
			return nil
		}).
		WithChoice("手榴弾スポーン(インベントリ)", func(world w.World) error {
			_, err := worldhelper.SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
			if err != nil {
				return fmt.Errorf("error spawning item: %w", err)
			}
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
			return nil
		}).
		WithChoice("ゲームオーバー", func(_ w.World) error {
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewGameOverMessageState}})
			return nil
		}).
		WithChoice("ダンジョン開始(大部屋)", func(_ w.World) error {
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{
				NewDungeonState(1, WithBuilderType(mapplanner.PlannerTypeBigRoom)),
			}})
			return nil
		}).
		WithChoice("ダンジョン開始(小部屋)", func(_ w.World) error {
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{
				NewDungeonState(1, WithBuilderType(mapplanner.PlannerTypeSmallRoom)),
			}})
			return nil
		}).
		WithChoice("ダンジョン開始(洞窟)", func(_ w.World) error {
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{
				NewDungeonState(1, WithBuilderType(mapplanner.PlannerTypeCave)),
			}})
			return nil
		}).
		WithChoice("ダンジョン開始(廃墟)", func(_ w.World) error {
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{
				NewDungeonState(1, WithBuilderType(mapplanner.PlannerTypeRuins)),
			}})
			return nil
		}).
		WithChoice("ダンジョン開始(森)", func(_ w.World) error {
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{
				NewDungeonState(1, WithBuilderType(mapplanner.PlannerTypeForest)),
			}})
			return nil
		}).
		WithChoice("市街地開始", func(_ w.World) error {
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{
				NewDungeonState(1, WithBuilderType(mapplanner.PlannerTypeTown)),
			}})
			return nil
		}).
		WithChoice("メッセージ表示テスト", func(_ w.World) error {
			testMessageData := messagedata.NewSystemMessage("ゲームが自動保存されました。\n\n進行状況は安全に記録されています。")
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{func() es.State[w.World] { return NewMessageState(testMessageData) }}})
			return nil
		}).
		WithChoice("アイテム入手イベント", func(world w.World) error {
			// アイテムを実際にインベントリに追加
			_ = worldhelper.AddStackableCount(world, "鉄", 1)
			_ = worldhelper.AddStackableCount(world, "木の棒", 1)
			_ = worldhelper.AddStackableCount(world, "フェライトコア", 2)

			// アイテム入手完了後の表示用メッセージを生成
			messageText := "宝箱を発見した。\n\n" +
				"鉄を手に入れた。\n" +
				"木の棒を手に入れた。\n" +
				"フェライトコアを2個手に入れた。\n"

			itemMessageData := &messagedata.MessageData{
				Speaker: "",
			}
			itemMessageData.AddText(messageText)
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{func() es.State[w.World] { return NewMessageState(itemMessageData) }}})
			return nil
		}).
		WithChoice("長いメッセージテスト", func(_ w.World) error {
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
			return nil
		}).
		WithChoice("連鎖メッセージテスト", func(_ w.World) error {
			chainMessageData := messagedata.NewSystemMessage("戦闘開始。").
				SystemMessage("剣と剣がぶつかり合う。").
				SystemMessage("勝利した。")

			messageState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{func() es.State[w.World] { return NewMessageState(chainMessageData) }}})
			return nil
		}).
		WithChoice("選択肢分岐メッセージテスト", func(_ w.World) error {
			battleMessage := messagedata.NewSystemMessage("戦闘した。")
			negotiateMessage := messagedata.NewSystemMessage("交渉した。")
			escapeMessage := messagedata.NewSystemMessage("逃走した。")

			choiceMessageData := messagedata.NewDialogMessage("敵に遭遇した。", "").
				WithChoiceMessage("戦う", battleMessage).
				WithChoiceMessage("交渉する", negotiateMessage).
				WithChoiceMessage("逃走する", escapeMessage)

			messageState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{func() es.State[w.World] { return NewMessageState(choiceMessageData) }}})
			return nil
		}).
		WithChoice("選択肢処理テスト", func(_ w.World) error {
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

			testMessageData.Choices[0].Action = func(_ w.World) error { choiceAction1(); return nil }
			testMessageData.Choices[1].Action = func(_ w.World) error { choiceAction2(); return nil }

			messageState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{func() es.State[w.World] { return NewMessageState(testMessageData) }}})
			return nil
		}).
		WithChoice("デバッグ表示切り替え", func(_ w.World) error {
			cfg := config.Get()
			cfg.ShowAIDebug = !cfg.ShowAIDebug
			cfg.NoEncounter = !cfg.NoEncounter
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
			return nil
		}).
		WithChoice("ゲーム開始メッセージ", func(_ w.World) error {
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewGameStartMessageState}})
			return nil
		}).
		WithChoice("閉じる", func(_ w.World) error {
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
			return nil
		})

	return messageState
}

// DungeonStateOption はDungeonStateのオプション設定関数
type DungeonStateOption func(*DungeonState)

// WithSeed はシード値を設定するオプション
func WithSeed(seed uint64) DungeonStateOption {
	return func(ds *DungeonState) {
		ds.Seed = seed
	}
}

// WithBuilderType はマップビルダータイプを設定するオプション
func WithBuilderType(builderType mapplanner.PlannerType) DungeonStateOption {
	return func(ds *DungeonState) {
		ds.BuilderType = builderType
	}
}

// NewDungeonState はDungeonStateインスタンスを作成するファクトリー関数
// デフォルトではBuilderTypeはPlannerTypeRandomになる
func NewDungeonState(depth int, opts ...DungeonStateOption) es.StateFactory[w.World] {
	return func() es.State[w.World] {
		ds := &DungeonState{
			Depth:       depth,
			BuilderType: mapplanner.PlannerTypeRandom,
		}
		for _, opt := range opts {
			opt(ds)
		}
		return ds
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
		WithChoice("メインメニューに戻る", func(_ w.World) error {
			// メインメニューに遷移
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewMainMenuState}})
			return nil
		})

	// MessageStateにMessageDataを設定
	messageState.messageData = messageData

	return messageState
}

// NewGameStartMessageState はゲーム開始時の目的を説明するMessageStateを作成するファクトリー関数
func NewGameStartMessageState() es.State[w.World] {
	// MessageStateインスタンスを作成
	messageState := &MessageState{}

	// メッセージを作成
	messageData := messagedata.NewDialogMessage(`「あんた、遺跡の『珠狙い』だろ?
外からこの街に来る異常に若い連中はみんなそうさ。
向こう見ずで破滅的で、...
どうしようもない事情を持ってる。」

「あんたは...、そうか、母親が...。
言っちゃ悪いが、そういう奴らはここでは珍しくない。
どんな事情があるにせよ、遺跡で辿る結末は1つさ。」`, "老兵")

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

		messageData = messageData.WithChoice(label, func(world w.World) error {
			if err := saveManager.SaveWorld(world, slotName); err != nil {
				return fmt.Errorf("save failed: %w", err)
			}
			// セーブ後は同じセーブメニューを再作成してメニューを維持
			messageState.SetTransition(es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewSaveMenuState}})
			return nil
		})
	}

	messageData = messageData.WithChoice("戻る", func(_ w.World) error {
		messageState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
		return nil
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

			messageData = messageData.WithChoice(label, func(world w.World) error {
				// ロードを実行
				err := saveManager.LoadWorld(world, slotName)
				if err != nil {
					println("Load failed:", err.Error())
					messageState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
					return err
				}
				// 遷移（街マップを生成してプレイヤーを配置）
				stateFactory := NewDungeonState(1, WithBuilderType(mapplanner.PlannerTypeTown))
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{stateFactory}})
				return nil
			})
		}
	}

	// 有効なセーブデータが存在しない場合の処理
	if !hasValidSlot {
		messageData = messageData.WithChoice("セーブデータがありません", func(_ w.World) error {
			// 何もしない（選択不可を示すためのダミー選択肢）
			return nil
		})
	}

	messageData = messageData.WithChoice("戻る", func(_ w.World) error {
		messageState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
		return nil
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

// NewShopMenuState は新しいShopMenuStateインスタンスを作成するファクトリー関数
func NewShopMenuState() es.State[w.World] {
	return &ShopMenuState{}
}

// NewInteractionMenuState はインタラクションメニューStateを作成する
func NewInteractionMenuState(world w.World) es.State[w.World] {
	messageState := &MessageState{}

	// プレイヤー周辺の実行可能なアクションを取得
	interactionActions := GetInteractionActions(world)

	if len(interactionActions) == 0 {
		// アクションがない場合
		messageState.messageData = messagedata.NewSystemMessage("実行可能なアクションがありません。")
		return messageState
	}

	// アクションメニューを構築
	messageState.messageData = messagedata.NewSystemMessage("")

	for _, action := range interactionActions {
		// クロージャで変数をキャプチャ
		capturedAction := action
		messageState.messageData = messageState.messageData.WithChoice(capturedAction.Label, func(world w.World) error {
			// アクションを実行
			playerEntity, err := worldhelper.GetPlayerEntity(world)
			if err != nil {
				messageState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
				return err
			}

			params := actions.ActionParams{
				Actor:  playerEntity,
				Target: &capturedAction.Target,
			}
			executeActivity(world, capturedAction.Activity, params)

			messageState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
			return nil
		})
	}

	// キャンセル用の「閉じる」選択肢を追加
	messageState.messageData = messageState.messageData.WithChoice("キャンセル", func(_ w.World) error {
		messageState.SetTransition(es.Transition[w.World]{Type: es.TransPop})
		return nil
	})

	return messageState
}
