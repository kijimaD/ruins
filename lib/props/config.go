// Package props は置物（Prop）の設定とマネージャーを提供する
// 置物の種類、属性、スプライト情報などを管理し、
// ゲーム世界への置物配置をサポートする
package props

import (
	gc "github.com/kijimaD/ruins/lib/components"
)

// PropConfig は置物の設定情報を保持する
type PropConfig struct {
	Name             string      // 置物の名前
	Type             gc.PropType // 置物のタイプ
	SpriteNumber     int         // スプライト番号
	BlocksMovement   bool        // 移動を阻害するか
	BlocksVisibility bool        // 視界を遮るか
	SpawnWeight      int         // ランダム配置での重み
	Description      string      // 説明文
}

// DefaultPropConfigs はデフォルトの置物設定を返す
func DefaultPropConfigs() map[gc.PropType]PropConfig {
	return map[gc.PropType]PropConfig{
		gc.PropTypeTable: {
			Name:             "テーブル",
			Type:             gc.PropTypeTable,
			SpriteNumber:     19, // Table sprite
			BlocksMovement:   true,
			BlocksVisibility: false,
			SpawnWeight:      10,
			Description:      "頑丈な木製のテーブル",
		},
		gc.PropTypeChair: {
			Name:             "椅子",
			Type:             gc.PropTypeChair,
			SpriteNumber:     20, // Chair sprite
			BlocksMovement:   true,
			BlocksVisibility: false,
			SpawnWeight:      15,
			Description:      "座り心地の良さそうな椅子",
		},
		gc.PropTypeBookshelf: {
			Name:             "本棚",
			Type:             gc.PropTypeBookshelf,
			SpriteNumber:     21, // Bookshelf sprite
			BlocksMovement:   true,
			BlocksVisibility: true,
			SpawnWeight:      5,
			Description:      "たくさんの本が並んだ本棚",
		},
		gc.PropTypeBarrel: {
			Name:             "樽",
			Type:             gc.PropTypeBarrel,
			SpriteNumber:     22, // Barrel sprite
			BlocksMovement:   true,
			BlocksVisibility: false,
			SpawnWeight:      8,
			Description:      "何か入っているかもしれない樽",
		},
		gc.PropTypeCrate: {
			Name:             "木箱",
			Type:             gc.PropTypeCrate,
			SpriteNumber:     23, // Crate sprite
			BlocksMovement:   true,
			BlocksVisibility: false,
			SpawnWeight:      12,
			Description:      "頑丈そうな木箱",
		},
		gc.PropTypeBed: {
			Name:             "寝台",
			Type:             gc.PropTypeBed,
			SpriteNumber:     24, // Bed sprite
			BlocksMovement:   true,
			BlocksVisibility: false,
			SpawnWeight:      8,
			Description:      "簡素だが清潔な寝台",
		},
	}
}

// PropManager は置物を管理するシステム
type PropManager struct {
	configs map[gc.PropType]PropConfig
}

// NewPropManager は新しいPropManagerを作成する
func NewPropManager() *PropManager {
	return &PropManager{
		configs: DefaultPropConfigs(),
	}
}

// GetConfig は指定されたタイプの設定を取得する
func (pm *PropManager) GetConfig(propType gc.PropType) (PropConfig, bool) {
	config, exists := pm.configs[propType]
	return config, exists
}

// AddConfig は新しい置物設定を追加する
func (pm *PropManager) AddConfig(config PropConfig) {
	pm.configs[config.Type] = config
}

// GetAllConfigs は全ての設定を取得する
func (pm *PropManager) GetAllConfigs() map[gc.PropType]PropConfig {
	return pm.configs
}
