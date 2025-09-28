package components

// PropType は置物のタイプを表す
type PropType string

const (
	// PropTypeTable はテーブル
	PropTypeTable PropType = "table"
	// PropTypeChair は椅子
	PropTypeChair PropType = "chair"
	// PropTypeBookshelf は本棚
	PropTypeBookshelf PropType = "bookshelf"
	// PropTypeBed はベッド
	PropTypeBed PropType = "bed"
	// PropTypeBarrel は樽
	PropTypeBarrel PropType = "barrel"
	// PropTypeCrate は木箱
	PropTypeCrate PropType = "crate"
)
