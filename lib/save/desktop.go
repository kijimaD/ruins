//go:build !js || !wasm

package save

import (
	"fmt"
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

// saveFileExistsWasm はデスクトップ環境では使用できない
func (sm *SerializationManager) saveFileExistsWasm(_ string) bool {
	return false
}

// getSaveFileTimestampWasm はデスクトップ環境では使用できない
func (sm *SerializationManager) getSaveFileTimestampWasm(_ string) (time.Time, error) {
	return time.Time{}, fmt.Errorf("WASM functions not available in desktop environment")
}