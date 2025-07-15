package resources

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/kijimaD/ruins/assets"
	ec "github.com/kijimaD/ruins/lib/engine/components"
	er "github.com/kijimaD/ruins/lib/engine/resources"
	"github.com/kijimaD/ruins/lib/raw"
)

// ResourceManager はすべてのリソースの読み込みを統括するインターフェース
type ResourceManager interface {
	// フォント関連
	LoadFonts() (map[string]er.Font, error)
	// スプライトシート関連
	LoadSpriteSheets() (map[string]ec.SpriteSheet, error)
	// コントロール関連
	LoadControls(axes []string, actions []string) (er.Controls, er.InputHandler, error)
	// Raw(エンティティ定義)関連
	LoadRaws() (*raw.Master, error)
	// すべてのリソースを一括読み込み
	LoadAll(axes []string, actions []string) error
}

// DefaultResourceManager はResourceManagerのデフォルト実装
type DefaultResourceManager struct {
	// 設定ファイルのパス
	config ResourceConfig
	// キャッシュされたリソース
	cache *ResourceCache
}

// ResourceConfig はリソースファイルのパスを管理する設定
type ResourceConfig struct {
	FontsPath        string
	SpriteSheetsPath string
	ControlsPath     string
	RawsPath         string
}

// ResourceCache は読み込み済みのリソースをキャッシュする
type ResourceCache struct {
	Fonts        map[string]er.Font
	SpriteSheets map[string]ec.SpriteSheet
	Controls     *er.Controls
	InputHandler *er.InputHandler
	RawMaster    *raw.Master
}

// NewResourceManager は新しいResourceManagerを作成する
func NewResourceManager(config ResourceConfig) ResourceManager {
	return &DefaultResourceManager{
		config: config,
		cache:  &ResourceCache{},
	}
}

// NewDefaultResourceManager はデフォルトのパス設定でResourceManagerを作成する
func NewDefaultResourceManager() ResourceManager {
	return NewResourceManager(ResourceConfig{
		FontsPath:        "metadata/fonts/fonts.toml",
		SpriteSheetsPath: "metadata/spritesheets/spritesheets.toml",
		ControlsPath:     "config/controls.toml",
		RawsPath:         "metadata/entities/raw/raw.toml",
	})
}

// LoadFonts はフォントリソースを読み込む
func (rm *DefaultResourceManager) LoadFonts() (map[string]er.Font, error) {
	// キャッシュがあれば返す
	if rm.cache.Fonts != nil {
		return rm.cache.Fonts, nil
	}

	type fontMetadata struct {
		Fonts map[string]er.Font `toml:"font"`
	}

	var metadata fontMetadata
	bs, err := assets.FS.ReadFile(rm.config.FontsPath)
	if err != nil {
		return nil, fmt.Errorf("フォントファイルの読み込みに失敗: %w", err)
	}

	if _, err := toml.Decode(string(bs), &metadata); err != nil {
		return nil, fmt.Errorf("フォントメタデータのデコードに失敗: %w", err)
	}

	rm.cache.Fonts = metadata.Fonts
	return metadata.Fonts, nil
}

// LoadSpriteSheets はスプライトシートリソースを読み込む
func (rm *DefaultResourceManager) LoadSpriteSheets() (map[string]ec.SpriteSheet, error) {
	// キャッシュがあれば返す
	if rm.cache.SpriteSheets != nil {
		return rm.cache.SpriteSheets, nil
	}

	type spriteSheetMetadata struct {
		SpriteSheets map[string]ec.SpriteSheet `toml:"sprite_sheet"`
	}

	var metadata spriteSheetMetadata
	bs, err := assets.FS.ReadFile(rm.config.SpriteSheetsPath)
	if err != nil {
		return nil, fmt.Errorf("スプライトシートファイルの読み込みに失敗: %w", err)
	}

	if _, err := toml.Decode(string(bs), &metadata); err != nil {
		return nil, fmt.Errorf("スプライトシートメタデータのデコードに失敗: %w", err)
	}

	// 名前を設定
	for k, v := range metadata.SpriteSheets {
		v.Name = k
		metadata.SpriteSheets[k] = v
	}

	rm.cache.SpriteSheets = metadata.SpriteSheets
	return metadata.SpriteSheets, nil
}

// LoadControls はコントロール設定を読み込む
func (rm *DefaultResourceManager) LoadControls(axes []string, actions []string) (er.Controls, er.InputHandler, error) {
	// キャッシュがあれば返す
	if rm.cache.Controls != nil && rm.cache.InputHandler != nil {
		return *rm.cache.Controls, *rm.cache.InputHandler, nil
	}

	type controlsConfig struct {
		Controls er.Controls `toml:"controls"`
	}

	var config controlsConfig
	bs, err := assets.FS.ReadFile(rm.config.ControlsPath)
	if err != nil {
		return er.Controls{}, er.InputHandler{}, fmt.Errorf("コントロールファイルの読み込みに失敗: %w", err)
	}

	if _, err := toml.Decode(string(bs), &config); err != nil {
		return er.Controls{}, er.InputHandler{}, fmt.Errorf("コントロール設定のデコードに失敗: %w", err)
	}

	// InputHandlerの初期化
	inputHandler := er.InputHandler{
		Axes:    make(map[string]float64),
		Actions: make(map[string]bool),
	}

	// 軸の検証
	for _, axis := range axes {
		if _, ok := config.Controls.Axes[axis]; !ok {
			return er.Controls{}, er.InputHandler{}, fmt.Errorf("軸 '%s' の設定が見つかりません", axis)
		}
		inputHandler.Axes[axis] = 0
	}

	// アクションの検証
	for _, action := range actions {
		if _, ok := config.Controls.Actions[action]; !ok {
			return er.Controls{}, er.InputHandler{}, fmt.Errorf("アクション '%s' の設定が見つかりません", action)
		}
		inputHandler.Actions[action] = false
	}

	rm.cache.Controls = &config.Controls
	rm.cache.InputHandler = &inputHandler

	return config.Controls, inputHandler, nil
}

// LoadRaws はRawデータを読み込む
func (rm *DefaultResourceManager) LoadRaws() (*raw.Master, error) {
	// キャッシュがあれば返す
	if rm.cache.RawMaster != nil {
		return rm.cache.RawMaster, nil
	}

	rawMaster := raw.LoadFromFile(rm.config.RawsPath)
	rm.cache.RawMaster = &rawMaster

	return &rawMaster, nil
}

// LoadAll はすべてのリソースを一括で読み込む
func (rm *DefaultResourceManager) LoadAll(axes []string, actions []string) error {
	// フォントの読み込み
	if _, err := rm.LoadFonts(); err != nil {
		return fmt.Errorf("フォントの読み込みに失敗: %w", err)
	}

	// スプライトシートの読み込み
	if _, err := rm.LoadSpriteSheets(); err != nil {
		return fmt.Errorf("スプライトシートの読み込みに失敗: %w", err)
	}

	// コントロールの読み込み
	if _, _, err := rm.LoadControls(axes, actions); err != nil {
		return fmt.Errorf("コントロールの読み込みに失敗: %w", err)
	}

	// Rawデータの読み込み
	if _, err := rm.LoadRaws(); err != nil {
		return fmt.Errorf("rawデータの読み込みに失敗: %w", err)
	}

	return nil
}

// GetCache は現在のキャッシュを取得する（テスト用）
func (rm *DefaultResourceManager) GetCache() *ResourceCache {
	return rm.cache
}

// ClearCache はキャッシュをクリアする（テスト用）
func (rm *DefaultResourceManager) ClearCache() {
	rm.cache = &ResourceCache{}
}

// GetResourcePath はリソースのフルパスを取得する（ヘルパー関数）
func GetResourcePath(basePath, filename string) string {
	return filepath.Join(basePath, filename)
}
