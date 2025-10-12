// Package testutil はテスト用のユーティリティ関数を提供する
package testutil

import (
	"sync"
	"testing"

	"github.com/kijimaD/ruins/lib/maingame"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/stretchr/testify/require"
)

// initWorldMu はInitWorldの並行実行を防ぐためのmutex
var initWorldMu sync.Mutex

// InitTestWorld はテスト用にWorldを初期化する
// 並行実行時の競合を防ぐためmutexでガードする
func InitTestWorld(t *testing.T) w.World {
	t.Helper()
	initWorldMu.Lock()
	defer initWorldMu.Unlock()

	world, err := maingame.InitWorld(960, 720)
	require.NoError(t, err)
	return world
}
