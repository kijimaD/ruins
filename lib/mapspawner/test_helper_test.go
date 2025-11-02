package mapspawner

import (
	"github.com/kijimaD/ruins/lib/raw"
)

// createMapspawnerTestRawMaster はテスト用のRawMasterを作成する
// raw.tomlファイルから読み込む
func createMapspawnerTestRawMaster() *raw.Master {
	rawMaster, err := raw.LoadFromFile("metadata/entities/raw/raw.toml")
	if err != nil {
		panic(err)
	}
	return &rawMaster
}
