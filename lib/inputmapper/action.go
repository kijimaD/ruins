package inputmapper

// ActionID はアクションの一意な識別子
type ActionID string

// 移動系アクション
const (
	ActionMoveNorth     ActionID = "move_north"
	ActionMoveSouth     ActionID = "move_south"
	ActionMoveEast      ActionID = "move_east"
	ActionMoveWest      ActionID = "move_west"
	ActionMoveNorthEast ActionID = "move_north_east"
	ActionMoveNorthWest ActionID = "move_north_west"
	ActionMoveSouthEast ActionID = "move_south_east"
	ActionMoveSouthWest ActionID = "move_south_west"
	ActionWait          ActionID = "wait"
)

// UI系アクション
const (
	ActionOpenInventory   ActionID = "open_inventory"
	ActionOpenEquipment   ActionID = "open_equipment"
	ActionOpenCraft       ActionID = "open_craft"
	ActionOpenShop        ActionID = "open_shop"
	ActionOpenDungeonMenu ActionID = "open_dungeon_menu"
	ActionOpenDebugMenu   ActionID = "open_debug_menu"
	ActionCloseMenu       ActionID = "close_menu"
)

// アイテム系アクション
const (
	ActionPickup  ActionID = "pickup"
	ActionDrop    ActionID = "drop"
	ActionUseItem ActionID = "use_item"
	ActionEquip   ActionID = "equip"
	ActionUnequip ActionID = "unequip"
)

// 世界との相互作用アクション
const (
	ActionInteract ActionID = "interact" // 汎用的な相互作用（ワープ、アイテム拾得など）
)

// 戦闘系アクション
const (
	ActionAttack ActionID = "attack"
)

// メニュー操作アクション
const (
	ActionMenuUp     ActionID = "menu_up"
	ActionMenuDown   ActionID = "menu_down"
	ActionMenuLeft   ActionID = "menu_left"
	ActionMenuRight  ActionID = "menu_right"
	ActionMenuSelect ActionID = "menu_select"
	ActionMenuCancel ActionID = "menu_cancel"
)

// メッセージウィンドウ系アクション
const (
	ActionConfirm ActionID = "confirm" // メッセージ確認
	ActionSkip    ActionID = "skip"    // メッセージスキップ
)
