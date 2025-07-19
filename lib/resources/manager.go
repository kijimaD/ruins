package resources

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/kijimaD/ruins/assets"
	"github.com/kijimaD/ruins/lib/components"
	er "github.com/kijimaD/ruins/lib/engine/resources"
	"github.com/kijimaD/ruins/lib/raw"
)

// ResourceManager はすべてのリソースの読み込みを統括するインターフェース
type ResourceManager interface {
	// フォント関連
	LoadFonts() (map[string]er.Font, error)
	// スプライトシート関連
	LoadSpriteSheets() (map[string]components.SpriteSheet, error)
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
	RawsPath         string
}

// ResourceCache は読み込み済みのリソースをキャッシュする
type ResourceCache struct {
	Fonts        map[string]er.Font
	SpriteSheets map[string]components.SpriteSheet
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
func (rm *DefaultResourceManager) LoadSpriteSheets() (map[string]components.SpriteSheet, error) {
	// キャッシュがあれば返す
	if rm.cache.SpriteSheets != nil {
		return rm.cache.SpriteSheets, nil
	}

	type spriteSheetMetadata struct {
		SpriteSheets map[string]components.SpriteSheet `toml:"sprite_sheet"`
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

// LoadRaws はRawデータを読み込む
func (rm *DefaultResourceManager) LoadRaws() (*raw.Master, error) {
	// キャッシュがあれば返す
	if rm.cache.RawMaster != nil {
		return rm.cache.RawMaster, nil
	}

	rawMaster, err := raw.LoadFromFile(rm.config.RawsPath)
	if err != nil {
		return nil, err
	}
	rm.cache.RawMaster = &rawMaster

	return &rawMaster, nil
}

// LoadAll はすべてのリソースを一括で読み込む
func (rm *DefaultResourceManager) LoadAll(_ []string, _ []string) error {
	// フォントの読み込み
	if _, err := rm.LoadFonts(); err != nil {
		return fmt.Errorf("フォントの読み込みに失敗: %w", err)
	}

	// スプライトシートの読み込み
	if _, err := rm.LoadSpriteSheets(); err != nil {
		return fmt.Errorf("スプライトシートの読み込みに失敗: %w", err)
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
