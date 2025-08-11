package config

import (
	"fmt"
	"log"
	"sync"
)

var (
	instance *Config
	once     sync.Once
)

// Get はアプリケーション設定のシングルトンインスタンスを返す
func Get() *Config {
	once.Do(func() {
		var err error
		instance, err = Load()
		if err != nil {
			log.Fatalf("設定の読み込みに失敗しました: %v", err)
		}

		if err := instance.Validate(); err != nil {
			log.Fatalf("設定の検証に失敗しました: %v", err)
		}
	})
	return instance
}

// Reset は設定インスタンスをリセットする（主にテスト用）
func Reset() {
	once = sync.Once{}
	instance = nil
}

// MustGet は設定を取得し、エラーがあればパニックする
func MustGet() *Config {
	return Get()
}

// String は設定の文字列表現を返す（デバッグ用）
func (c *Config) String() string {
	return fmt.Sprintf(`Config{
	Profile: %s,
	WindowWidth: %d, WindowHeight: %d, Fullscreen: %t,
	Debug: %t, LogLevel: %s, DebugPProf: %t, PProfPort: %d,
	StartingState: %s,
	TargetFPS: %d,
	ProfileMemory: %t, ProfileCPU: %t, ProfileMutex: %t, ProfileTrace: %t,
	ProfilePath: %s
}`,
		c.Profile,
		c.WindowWidth, c.WindowHeight, c.Fullscreen,
		c.Debug, c.LogLevel, c.DebugPProf, c.PProfPort,
		c.StartingState,
		c.TargetFPS,
		c.ProfileMemory, c.ProfileCPU, c.ProfileMutex, c.ProfileTrace,
		c.ProfilePath)
}
