package raw

import "math/rand/v2"

// DropTable はドロップテーブル
type DropTable struct {
	Name    string
	XpBase  float64
	Entries []DropTableEntry `toml:"entries"`
}

// DropTableEntry はドロップテーブルのエントリ
type DropTableEntry struct {
	Material string
	Weight   float64
}

// SelectByWeight は重みで選択する
func (ct DropTable) SelectByWeight() string {
	var totalWeight float64
	for _, entry := range ct.Entries {
		totalWeight += entry.Weight
	}
	randomValue := rand.Float64() * totalWeight

	// 累積ウェイトで判定
	var cumulativeWeight float64
	for _, entry := range ct.Entries {
		cumulativeWeight += entry.Weight
		if randomValue < cumulativeWeight {
			return entry.Material
		}
	}

	return ""
}
