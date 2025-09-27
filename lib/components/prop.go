package components

// PropType は置物のタイプを表す
type PropType string

// String はPropTypeの文字列表現を返す
func (pt PropType) String() string {
	return "PropType"
}

const (
	// PropTypeTable はテーブル
	PropTypeTable PropType = "table"
	// PropTypeChair は椅子
	PropTypeChair PropType = "chair"
	// PropTypeBookshelf は本棚
	PropTypeBookshelf PropType = "bookshelf"
	// PropTypeBed はベッド
	PropTypeBed PropType = "bed"
	// PropTypeChest は宝箱
	PropTypeChest PropType = "chest"
	// PropTypeDoor はドア
	PropTypeDoor PropType = "door"
	// PropTypeTorch はたいまつ
	PropTypeTorch PropType = "torch"
	// PropTypeAltar は祭壇
	PropTypeAltar PropType = "altar"
	// PropTypeBarrel は樽
	PropTypeBarrel PropType = "barrel"
	// PropTypeCrate は木箱
	PropTypeCrate PropType = "crate"
)
