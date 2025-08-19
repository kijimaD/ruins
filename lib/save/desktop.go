//go:build !js || !wasm

package save

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// saveToLocalStorage はデスクトップ環境では使用できない
func (sm *SerializationManager) saveToLocalStorage(_ string, _ []byte) error {
	return fmt.Errorf("localStorage is not available in desktop environment")
}

// loadFromLocalStorage はデスクトップ環境では使用できない
func (sm *SerializationManager) loadFromLocalStorage(_ string) ([]byte, error) {
	return nil, fmt.Errorf("localStorage is not available in desktop environment")
}

// saveFileExistsImpl はデスクトップ環境でセーブファイルが存在するかチェックする
func (sm *SerializationManager) saveFileExistsImpl(slotName string) bool {
	fileName := filepath.Join(sm.saveDirectory, slotName+".json")
	_, err := os.Stat(fileName)
	return err == nil
}

// getSaveFileTimestampImpl はデスクトップ環境でセーブファイルのタイムスタンプを取得する
func (sm *SerializationManager) getSaveFileTimestampImpl(slotName string) (time.Time, error) {
	fileName := filepath.Join(sm.saveDirectory, slotName+".json")
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get file info: %w", err)
	}

	return fileInfo.ModTime(), nil
}
