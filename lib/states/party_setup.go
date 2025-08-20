package states

import (
	"fmt"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

const (
	// itemTypeMember はメンバーアイテムのタイプを表す定数
	itemTypeMember = "member"
	itemTypeAction = "action"
)

// PartySetupState はパーティ編成画面のゲームステート
type PartySetupState struct {
	es.BaseState
	ui *ebitenui.UI

	menu                  *menu.Menu
	uiBuilder             *menu.UIBuilder
	keyboardInput         input.KeyboardInput
	memberDescContainer   *widget.Container // 中央カラム：メンバー詳細
	memberStatusContainer *widget.Container // 右カラム：メンバーステータス
	rootContainer         *widget.Container

	// パーティ状態
	currentPartySlots    [4]*ecs.Entity // 現在のパーティスロット（主人公は0番固定）
	protagonistEntity    ecs.Entity     // 主人公のエンティティ
	selectedMemberEntity *ecs.Entity    // 選択中のメンバー
}

func (st PartySetupState) String() string {
	return "PartySetup"
}

// NewPartySetupState はPartySetupStateを作成する
func NewPartySetupState() es.State {
	return &PartySetupState{}
}

var _ es.State = &PartySetupState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *PartySetupState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *PartySetupState) OnResume(_ w.World) {
	// フォーカス状態を更新
	if st.uiBuilder != nil && st.menu != nil {
		st.uiBuilder.UpdateFocus(st.menu)
	}
}

// OnStart はステートが開始される際に呼ばれる
func (st *PartySetupState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}

	st.initializePartyData(world)
	st.ui = st.initUI(world)

	// 初期フォーカス状態を設定
	if st.uiBuilder != nil && st.menu != nil {
		st.uiBuilder.UpdateFocus(st.menu)
	}
}

// OnStop はステートが停止される際に呼ばれる
func (st *PartySetupState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *PartySetupState) Update(_ w.World) es.Transition {
	st.menu.Update(st.keyboardInput)
	st.ui.Update()

	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *PartySetupState) Draw(_ w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// initializePartyData はパーティデータを初期化する
func (st *PartySetupState) initializePartyData(world w.World) {
	// 主人公を特定（Playerコンポーネントを持つエンティティが必須）
	var protagonist ecs.Entity
	found := false

	// Playerコンポーネントを持つエンティティを探す
	world.Manager.Join(
		world.Components.Player,
		world.Components.FactionAlly,
		world.Components.Name,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		if !found {
			protagonist = entity
			found = true
		}
	}))

	// Playerコンポーネントを持つエンティティが見つからない場合はエラーで停止
	if !found {
		panic("パーティ編成画面：Playerコンポーネントを持つ主人公エンティティが見つかりません。ゲームデータに問題があります。")
	}

	st.protagonistEntity = protagonist

	// 現在のパーティ構成を取得
	st.currentPartySlots = [4]*ecs.Entity{&protagonist, nil, nil, nil}
	slotIndex := 1

	worldhelper.QueryInPartyMember(world, func(entity ecs.Entity) {
		if entity != protagonist && slotIndex < 4 {
			st.currentPartySlots[slotIndex] = &entity
			slotIndex++
		}
	})
}

// initUI はUIを初期化する
func (st *PartySetupState) initUI(world w.World) *ebitenui.UI {
	return st.initUIWithFocus(world, -1)
}

// initUIWithFocus はUIを初期化し、指定されたインデックスにフォーカスを設定する
func (st *PartySetupState) initUIWithFocus(world w.World, focusIndex int) *ebitenui.UI {
	// メニュー項目を作成
	items := st.createMenuItems(world)

	// 初期インデックスを決定
	initialIndex := 0
	if focusIndex >= 0 && focusIndex < len(items) {
		initialIndex = focusIndex
	}

	// メニューの設定
	config := menu.Config{
		Items:          items,
		InitialIndex:   initialIndex,
		WrapNavigation: true,
		Orientation:    menu.Vertical,
	}

	// コールバックの設定
	callbacks := menu.Callbacks{
		OnSelect: func(_ int, item menu.Item) {
			st.handleItemSelection(world, item)
		},
		OnCancel: func() {
			// ESCでホームメニューに戻る
			st.SetTransition(es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewHomeMenuState}})
		},
		OnFocusChange: func(_, newIndex int) {
			if newIndex >= 0 && newIndex < len(items) {
				st.handleItemChange(world, items[newIndex])
			}
			// フォーカス変更時にUIを更新
			if st.uiBuilder != nil {
				st.uiBuilder.UpdateFocus(st.menu)
			}
		},
	}

	// メニューを作成
	st.menu = menu.NewMenu(config, callbacks)

	// UIビルダーを作成してメニューをbuild
	if st.uiBuilder == nil {
		st.uiBuilder = menu.NewUIBuilder(world)
	}
	menuContainer := st.uiBuilder.BuildUI(st.menu)

	// 中央と右カラムのコンテナを作成
	st.memberDescContainer = eui.NewVerticalContainer(
		widget.ContainerOpts.BackgroundImage(world.Resources.UIResources.Panel.ImageTrans),
	)
	st.memberStatusContainer = eui.NewVerticalContainer(
		widget.ContainerOpts.BackgroundImage(world.Resources.UIResources.Panel.ImageTrans),
	)

	// 初期表示を更新（最初のメンバーを選択）
	if len(items) > 0 {
		// 最初のメンバーアイテムを見つける
		firstMemberIndex := -1
		for i, item := range items {
			if userData, ok := item.UserData.(map[string]interface{}); ok {
				if userData["type"] == itemTypeMember {
					firstMemberIndex = i
					break
				}
			}
		}

		// メンバーが見つかった場合、そのメンバーを選択
		if firstMemberIndex >= 0 {
			st.handleItemChange(world, items[firstMemberIndex])
		}
	}

	// 3カラムレイアウトのルートコンテナ
	st.rootContainer = eui.NewItemGridContainer(
		widget.ContainerOpts.BackgroundImage(world.Resources.UIResources.Panel.ImageTrans),
	)

	// レイアウト構成
	st.rootContainer.AddChild(eui.NewTitleText("パーティ編成", world))
	st.rootContainer.AddChild(widget.NewContainer()) // 空の要素
	st.rootContainer.AddChild(widget.NewContainer()) // 空の要素

	st.rootContainer.AddChild(menuContainer)            // 左カラム
	st.rootContainer.AddChild(st.memberDescContainer)   // 中央カラム
	st.rootContainer.AddChild(st.memberStatusContainer) // 右カラム

	return &ebitenui.UI{Container: st.rootContainer}
}

// createMenuItems はメニュー項目を作成する
func (st *PartySetupState) createMenuItems(world w.World) []menu.Item {
	items := []menu.Item{}

	allMembers := st.getAllMembers(world)

	// パーティメンバーセクション
	items = append(items, menu.Item{
		ID:       "party_header",
		Label:    "パーティ",
		Disabled: true,
	})

	// パーティメンバーを追加
	for _, member := range allMembers {
		inParty, slotIndex := st.getMemberPartyStatus(member)
		if !inParty && member != st.protagonistEntity {
			continue // パーティにいないメンバーはスキップ
		}

		memberName := getDisplayName(world, member)
		var label string

		if member == st.protagonistEntity {
			label = fmt.Sprintf("%s(固定)", memberName)
		} else {
			label = memberName
		}

		items = append(items, menu.Item{
			ID:       fmt.Sprintf("member_%d", member),
			Label:    label,
			Disabled: false, // 全メンバー選択可能（主人公は操作のみ不可）
			UserData: map[string]interface{}{
				"type":           itemTypeMember,
				"entity":         member,
				"in_party":       inParty,
				"slot_index":     slotIndex,
				"is_protagonist": member == st.protagonistEntity,
			},
		})
	}

	// 待機メンバーセクション
	items = append(items, menu.Item{
		ID:       "waiting_header",
		Label:    "待機",
		Disabled: true,
	})

	// 待機メンバーを追加
	hasWaitingMembers := false
	for _, member := range allMembers {
		inParty, slotIndex := st.getMemberPartyStatus(member)
		if inParty || member == st.protagonistEntity {
			continue // パーティにいるメンバーや主人公はスキップ
		}

		hasWaitingMembers = true
		memberName := getDisplayName(world, member)
		label := memberName

		items = append(items, menu.Item{
			ID:       fmt.Sprintf("member_%d", member),
			Label:    label,
			Disabled: false,
			UserData: map[string]interface{}{
				"type":           itemTypeMember,
				"entity":         member,
				"in_party":       inParty,
				"slot_index":     slotIndex,
				"is_protagonist": false,
			},
		})
	}

	// 待機メンバーがいない場合は表示
	if !hasWaitingMembers {
		items = append(items, menu.Item{
			ID:       "no_waiting",
			Label:    "  (なし)",
			Disabled: true,
		})
	}

	// 操作説明
	items = append(items, menu.Item{
		ID:       "separator",
		Label:    "操作",
		Disabled: true,
	})
	items = append(items, menu.Item{
		ID:    "apply",
		Label: "編成を確定",
		UserData: map[string]interface{}{
			"type":   itemTypeAction,
			"action": "apply",
		},
	})
	items = append(items, menu.Item{
		ID:    "cancel",
		Label: "キャンセル",
		UserData: map[string]interface{}{
			"type":   itemTypeAction,
			"action": "cancel",
		},
	})

	return items
}

// getAllMembers は全ての味方メンバーを取得する
func (st *PartySetupState) getAllMembers(world w.World) []ecs.Entity {
	members := []ecs.Entity{}

	// 主人公を最初に追加
	members = append(members, st.protagonistEntity)

	// その他の味方を追加
	world.Manager.Join(
		world.Components.FactionAlly,
		world.Components.Name,
		world.Components.Attributes,
		world.Components.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		if entity != st.protagonistEntity {
			members = append(members, entity)
		}
	}))

	return members
}

// getMemberPartyStatus はメンバーのパーティ参加状況を取得する
func (st *PartySetupState) getMemberPartyStatus(entity ecs.Entity) (bool, int) {
	for i, member := range st.currentPartySlots {
		if member != nil && *member == entity {
			return true, i
		}
	}
	return false, -1
}

// getDisplayName は表示名を取得する
func getDisplayName(world w.World, entity ecs.Entity) string {
	return world.Components.Name.Get(entity).(*gc.Name).Name
}

// getJobName は職業名を取得する
func getJobName(world w.World, entity ecs.Entity) string {
	if entity.HasComponent(world.Components.Job) {
		return world.Components.Job.Get(entity).(*gc.Job).Job
	}
	return ""
}

// handleItemSelection はアイテム選択時の処理
func (st *PartySetupState) handleItemSelection(world w.World, item menu.Item) {
	userData, ok := item.UserData.(map[string]interface{})
	if !ok {
		return
	}

	switch userData["type"].(string) {
	case itemTypeMember:
		st.handleMemberSelection(world, userData)
	case "action":
		st.handleActionSelection(world, userData)
	}
}

// handleItemChange はアイテム変更時の処理（カーソル移動）
func (st *PartySetupState) handleItemChange(world w.World, item menu.Item) {
	userData, ok := item.UserData.(map[string]interface{})
	if !ok {
		st.selectedMemberEntity = nil
		st.updateMemberDisplay(world)
		return
	}

	var targetEntity *ecs.Entity

	switch userData["type"].(string) {
	case itemTypeMember:
		entity := userData["entity"].(ecs.Entity)
		targetEntity = &entity
	}

	st.selectedMemberEntity = targetEntity
	st.updateMemberDisplay(world)
}

// handleMemberSelection はメンバー選択時の処理（シンプルな切り替え）
func (st *PartySetupState) handleMemberSelection(world w.World, userData map[string]interface{}) {
	entity := userData["entity"].(ecs.Entity)
	inParty := userData["in_party"].(bool)

	// 主人公は操作不可（ステータス確認のみ）
	isProtagonist, exists := userData["is_protagonist"]
	if exists && isProtagonist.(bool) {
		return
	}

	if inParty {
		// パーティから外す
		st.removeMemberFromParty(entity)
	} else {
		// パーティに追加
		st.addMemberToParty(entity)
	}

	// メニューを再構築
	st.rebuildMenu(world)
}

// addMemberToParty はメンバーをパーティに追加する
func (st *PartySetupState) addMemberToParty(entity ecs.Entity) {
	// 空きスロットを見つけて配置
	for i := 1; i < 4; i++ { // 0番は主人公なのでスキップ
		if st.currentPartySlots[i] == nil {
			st.currentPartySlots[i] = &entity
			return
		}
	}
}

// removeMemberFromParty はメンバーをパーティから外す
func (st *PartySetupState) removeMemberFromParty(entity ecs.Entity) {
	for i := 1; i < 4; i++ { // 0番は主人公なのでスキップ
		if st.currentPartySlots[i] != nil && *st.currentPartySlots[i] == entity {
			st.currentPartySlots[i] = nil
			return
		}
	}
}

// handleActionSelection はアクション選択時の処理
func (st *PartySetupState) handleActionSelection(world w.World, userData map[string]interface{}) {
	action := userData["action"].(string)

	switch action {
	case "apply":
		st.applyPartyChanges(world)
		st.SetTransition(es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewHomeMenuState}})
	case "cancel":
		st.SetTransition(es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewHomeMenuState}})
	}
}

// applyPartyChanges はパーティ変更を適用する
func (st *PartySetupState) applyPartyChanges(world w.World) {
	// 全メンバーからInPartyコンポーネントを削除
	world.Manager.Join(
		world.Components.FactionAlly,
		world.Components.InParty,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		entity.RemoveComponent(world.Components.InParty)
	}))

	// 新しいパーティメンバーにInPartyコンポーネントを追加
	for _, member := range st.currentPartySlots {
		if member != nil {
			(*member).AddComponent(world.Components.InParty, &gc.InParty{})
		}
	}
}

// rebuildMenu はメニューを再構築する
func (st *PartySetupState) rebuildMenu(world w.World) {
	// 現在フォーカスされているアイテムのIDを保存
	var currentFocusedID string
	if st.menu != nil {
		currentItems := st.menu.GetItems()
		currentIndex := st.menu.GetFocusedIndex()
		if currentIndex >= 0 && currentIndex < len(currentItems) {
			currentFocusedID = currentItems[currentIndex].ID
		}
	}

	// 新しいメニュー項目を取得
	items := st.createMenuItems(world)

	// IDベースでフォーカス位置を見つける
	targetFocusIndex := -1
	if currentFocusedID != "" {
		for i, item := range items {
			if item.ID == currentFocusedID {
				targetFocusIndex = i
				break
			}
		}
	}

	// IDが見つからない場合は、最初の有効なアイテムにフォーカス
	if targetFocusIndex == -1 && len(items) > 0 {
		for i, item := range items {
			if !item.Disabled {
				targetFocusIndex = i
				break
			}
		}
	}

	// UIを再構築(フォーカス位置を指定)
	if st.uiBuilder != nil {
		st.ui = st.initUIWithFocus(world, targetFocusIndex)

		// 現在フォーカスされているアイテムの詳細表示を更新
		if targetFocusIndex >= 0 && targetFocusIndex < len(items) {
			st.handleItemChange(world, items[targetFocusIndex])
		}
	}
}

// updateMemberDisplay はメンバー表示を更新する
func (st *PartySetupState) updateMemberDisplay(world w.World) {
	// UIが初期化されていない場合は何もしない
	if st.memberDescContainer == nil || st.memberStatusContainer == nil {
		return
	}

	// 中央カラムをクリア
	st.memberDescContainer.RemoveChildren()

	// 右カラムをクリア
	st.memberStatusContainer.RemoveChildren()

	if st.selectedMemberEntity == nil {
		// 選択なしの場合
		st.memberDescContainer.AddChild(eui.NewBodyText("メンバーを選択してください", styles.TextColor, world))
		return
	}

	entity := *st.selectedMemberEntity

	// 中央カラム：基本情報
	// 職業を表示
	jobName := getJobName(world, entity)
	if jobName != "" {
		st.memberDescContainer.AddChild(eui.NewSubtitleText(fmt.Sprintf("職業: %s", jobName), world))
	}

	if entity.HasComponent(world.Components.Pools) {
		pools := world.Components.Pools.Get(entity).(*gc.Pools)
		st.memberDescContainer.AddChild(eui.NewBodyText(fmt.Sprintf("レベル: %d", pools.Level), styles.TextColor, world))
		st.memberDescContainer.AddChild(eui.NewBodyText(fmt.Sprintf("経験値: %d", pools.XP), styles.TextColor, world))
	}

	// 計算値（攻撃力、防御力など）を表示
	if entity.HasComponent(world.Components.Attributes) {
		attrs := world.Components.Attributes.Get(entity).(*gc.Attributes)
		st.memberDescContainer.AddChild(eui.NewBodyText(fmt.Sprintf("攻撃力: %d", attrs.Strength.Total), styles.TextColor, world))
		st.memberDescContainer.AddChild(eui.NewBodyText(fmt.Sprintf("防御力: %d", attrs.Defense.Total), styles.TextColor, world))
	}

	// 装備情報
	st.memberDescContainer.AddChild(eui.NewBodyText("", styles.TextColor, world)) // 空行
	st.memberDescContainer.AddChild(eui.NewBodyText("装備:", styles.TextColor, world))

	// 防具装備
	wearSlots := worldhelper.GetWearEquipments(world, entity)
	for i, slot := range wearSlots {
		if slot != nil {
			itemName := world.Components.Name.Get(*slot).(*gc.Name).Name
			st.memberDescContainer.AddChild(eui.NewBodyText(fmt.Sprintf("防具%d: %s", i+1, itemName), styles.TextColor, world))
		} else {
			st.memberDescContainer.AddChild(eui.NewBodyText(fmt.Sprintf("防具%d: (なし)", i+1), styles.TextColor, world))
		}
	}

	// 手札装備
	cardSlots := worldhelper.GetCardEquipments(world, entity)
	for i, slot := range cardSlots {
		if slot != nil {
			itemName := world.Components.Name.Get(*slot).(*gc.Name).Name
			st.memberDescContainer.AddChild(eui.NewBodyText(fmt.Sprintf("手札%d: %s", i+1, itemName), styles.TextColor, world))
		} else {
			st.memberDescContainer.AddChild(eui.NewBodyText(fmt.Sprintf("手札%d: (なし)", i+1), styles.TextColor, world))
		}
	}

	// 右カラム：詳細ステータス
	views.AddMemberBar(world, st.memberStatusContainer, entity)

	if entity.HasComponent(world.Components.Attributes) {
		attrs := world.Components.Attributes.Get(entity).(*gc.Attributes)
		st.memberStatusContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", consts.VitalityLabel, attrs.Vitality.Total, attrs.Vitality.Modifier), styles.TextColor, world))
		st.memberStatusContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", consts.StrengthLabel, attrs.Strength.Total, attrs.Strength.Modifier), styles.TextColor, world))
		st.memberStatusContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", consts.SensationLabel, attrs.Sensation.Total, attrs.Sensation.Modifier), styles.TextColor, world))
		st.memberStatusContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", consts.DexterityLabel, attrs.Dexterity.Total, attrs.Dexterity.Modifier), styles.TextColor, world))
		st.memberStatusContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", consts.AgilityLabel, attrs.Agility.Total, attrs.Agility.Modifier), styles.TextColor, world))
		st.memberStatusContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", consts.DefenseLabel, attrs.Defense.Total, attrs.Defense.Modifier), styles.TextColor, world))
	}
}
