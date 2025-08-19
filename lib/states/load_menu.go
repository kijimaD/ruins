package states

import (
	"fmt"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/save"
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	w "github.com/kijimaD/ruins/lib/world"
)

// LoadMenuState はロードスロット選択メニュー
type LoadMenuState struct {
	es.BaseState
	ui            *ebitenui.UI
	menu          *menu.Menu
	uiBuilder     *menu.UIBuilder
	keyboardInput input.KeyboardInput
	saveManager   *save.SerializationManager
}

func (st LoadMenuState) String() string {
	return "LoadMenu"
}

var _ es.State = &LoadMenuState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *LoadMenuState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *LoadMenuState) OnResume(_ w.World) {}

// OnStart はステート開始時の処理を行う
func (st *LoadMenuState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}

	// セーブマネージャーを初期化
	saveDir := "./saves"
	st.saveManager = save.NewSerializationManager(saveDir)

	st.initMenu(world)
	st.ui = st.initUI(world)
}

// OnStop はステートが停止される際に呼ばれる
func (st *LoadMenuState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *LoadMenuState) Update(_ w.World) es.Transition {
	// メニューの更新
	st.menu.Update(st.keyboardInput)
	st.ui.Update()

	return st.ConsumeTransition()
}

// Draw はスクリーンに描画する
func (st *LoadMenuState) Draw(world w.World, screen *ebiten.Image) {
	// 背景を描画
	bg := (*world.Resources.SpriteSheets)["bg_cup1"]
	screen.DrawImage(bg.Texture.Image, nil)

	st.ui.Draw(screen)
}

// initMenu はメニューコンポーネントを初期化する
func (st *LoadMenuState) initMenu(world w.World) {
	// セーブスロットの情報を取得
	saveSlots := st.getSaveSlotInfoFromDir("./saves")

	// メニュー項目の定義
	items := make([]menu.Item, 0, len(saveSlots)+1)

	for i, slot := range saveSlots {
		slotName := fmt.Sprintf("slot%d", i+1)

		items = append(items, menu.Item{
			ID:          slotName,
			Label:       slot.Label,
			Description: slot.Description,
			UserData:    slotName,
			Disabled:    !slot.Exists, // 存在しない場合は無効化
		})
	}

	// 戻るオプション
	items = append(items, menu.Item{
		ID:          "back",
		Label:       "戻る",
		Description: "前の画面に戻る",
		UserData:    "back",
	})

	// メニューの設定
	config := menu.Config{
		Items:          items,
		InitialIndex:   0,
		WrapNavigation: true,
		Orientation:    menu.Vertical,
	}

	// コールバックの設定
	callbacks := menu.Callbacks{
		OnSelect: func(_ int, item menu.Item) {
			slotName, ok := item.UserData.(string)
			if !ok {
				return
			}

			if slotName == "back" || slotName == "no_saves" {
				// メインメニューに戻る
				st.SetTransition(es.Transition{Type: es.TransPop})
				return
			}

			// ロードを実行
			err := st.saveManager.LoadWorld(world, slotName)
			if err != nil {
				// TODO: エラーハンドリング（エラーダイアログなど）
				fmt.Printf("Load failed: %v\n", err)
				return
			}

			// ロード成功後、ホームメニューに遷移
			st.SetTransition(es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewHomeMenuState}})
		},
		OnCancel: func() {
			// メインメニューに戻る
			st.SetTransition(es.Transition{Type: es.TransPop})
		},
		OnFocusChange: func(_, _ int) {
			if st.uiBuilder != nil {
				st.uiBuilder.UpdateFocus(st.menu)
			}
		},
	}

	// メニューを作成
	st.menu = menu.NewMenu(config, callbacks)

	// UIビルダーを作成
	st.uiBuilder = menu.NewUIBuilder(world)
}

// getSaveSlotInfoFromDir は指定ディレクトリからセーブスロット情報を取得する
func (st *LoadMenuState) getSaveSlotInfoFromDir(saveDir string) []SaveSlotInfo {
	slots := make([]SaveSlotInfo, 3) // 3つのセーブスロット
	for i := 0; i < 3; i++ {
		slotName := fmt.Sprintf("slot%d", i+1)

		if st.saveManager.SaveFileExists(slotName) {
			// セーブデータが存在する場合
			if timestamp, err := st.saveManager.GetSaveFileTimestamp(slotName); err == nil {
				slots[i] = SaveSlotInfo{
					Label:       fmt.Sprintf("%d [%s]", i+1, timestamp.Format("01/02 15:04")),
					Description: fmt.Sprintf("保存日時: %s", timestamp.Format("2006-01-02 15:04:05")),
					Exists:      true,
				}
			} else {
				slots[i] = SaveSlotInfo{
					Label:       fmt.Sprintf("%d [データあり]", i+1),
					Description: "セーブデータ",
					Exists:      true,
				}
			}
		} else {
			// セーブデータが存在しない場合
			slots[i] = SaveSlotInfo{
				Label:       fmt.Sprintf("%d [空]", i+1),
				Description: "空のスロット",
				Exists:      false,
			}
		}
	}

	return slots
}

// initUI はUIを初期化する
func (st *LoadMenuState) initUI(world w.World) *ebitenui.UI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// タイトル
	titleText := widget.NewText(
		widget.TextOpts.Text("読込", world.Resources.UIResources.Text.TitleFace, styles.TextColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				Padding: widget.Insets{
					Top: 100,
				},
			}),
		),
	)

	// メニューコンテナ
	menuContainer := st.uiBuilder.BuildUI(st.menu)
	wrapperContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	wrapperContainer.AddChild(menuContainer)

	rootContainer.AddChild(titleText)
	rootContainer.AddChild(wrapperContainer)

	return &ebitenui.UI{Container: rootContainer}
}
