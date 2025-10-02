// Package props は置物（Prop）の設定とマネージャーを提供する
// 置物の種類、属性、スプライト情報などを管理し、
// ゲーム世界への置物配置をサポートする
package props

// PropConfig は置物の設定情報を保持する
type PropConfig struct {
	Name             string // 置物の名前（識別用キー）
	SpriteKey        string // スプライトキー
	BlocksMovement   bool   // 移動を阻害するか
	BlocksVisibility bool   // 視界を遮るか
	SpawnWeight      int    // ランダム配置での重み
	Description      string // 説明文
}

// DefaultPropConfigs はデフォルトの置物設定を返す
func DefaultPropConfigs() map[string]PropConfig {
	return map[string]PropConfig{
		"table": {
			Name:             "テーブル",
			SpriteKey:        "prop_table",
			BlocksMovement:   true,
			BlocksVisibility: false,
			SpawnWeight:      10,
			Description:      "頑丈な木製のテーブル",
		},
		"chair": {
			Name:             "椅子",
			SpriteKey:        "prop_chair",
			BlocksMovement:   true,
			BlocksVisibility: false,
			SpawnWeight:      15,
			Description:      "座り心地の良さそうな椅子",
		},
		"bookshelf": {
			Name:             "本棚",
			SpriteKey:        "prop_bookshelf",
			BlocksMovement:   true,
			BlocksVisibility: true,
			SpawnWeight:      5,
			Description:      "たくさんの本が並んだ本棚",
		},
		"barrel": {
			Name:             "樽",
			SpriteKey:        "prop_barrel",
			BlocksMovement:   true,
			BlocksVisibility: false,
			SpawnWeight:      8,
			Description:      "何か入っているかもしれない樽",
		},
		"crate": {
			Name:             "木箱",
			SpriteKey:        "prop_crate",
			BlocksMovement:   true,
			BlocksVisibility: false,
			SpawnWeight:      12,
			Description:      "頑丈そうな木箱",
		},
		"bed": {
			Name:             "寝台",
			SpriteKey:        "prop_bed",
			BlocksMovement:   true,
			BlocksVisibility: false,
			SpawnWeight:      8,
			Description:      "簡素だが清潔な寝台",
		},
	}
}

// PropManager は置物を管理するシステム
type PropManager struct {
	configs map[string]PropConfig
}

// NewPropManager は新しいPropManagerを作成する
func NewPropManager() *PropManager {
	return &PropManager{
		configs: DefaultPropConfigs(),
	}
}

// GetConfig は指定された名前の設定を取得する
func (pm *PropManager) GetConfig(name string) (PropConfig, bool) {
	config, exists := pm.configs[name]
	return config, exists
}

// AddConfig は新しい置物設定を追加する
func (pm *PropManager) AddConfig(name string, config PropConfig) {
	pm.configs[name] = config
}

// GetAllConfigs は全ての設定を取得する
func (pm *PropManager) GetAllConfigs() map[string]PropConfig {
	return pm.configs
}
