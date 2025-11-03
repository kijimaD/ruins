package raw

import "math/rand/v2"

// ItemTable はアイテム出現テーブル
type ItemTable struct {
	Name    string
	Entries []ItemTableEntry `toml:"entries"`
}

// ItemTableEntry はアイテムテーブルのエントリ
type ItemTableEntry struct {
	ItemName string
	Weight   float64
}

// SelectByWeight は重みで選択する
func (it ItemTable) SelectByWeight(rng *rand.Rand) string {
	var totalWeight float64
	for _, entry := range it.Entries {
		totalWeight += entry.Weight
	}

	if totalWeight == 0 {
		return ""
	}

	randomValue := rng.Float64() * totalWeight

	// 累積ウェイトで判定
	var cumulativeWeight float64
	for _, entry := range it.Entries {
		cumulativeWeight += entry.Weight
		if randomValue < cumulativeWeight {
			return entry.ItemName
		}
	}

	return ""
}
