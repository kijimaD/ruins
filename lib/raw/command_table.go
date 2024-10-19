package raw

import (
	"math/rand/v2"
)

type CommandTable struct {
	Name    string
	Entries []CommandTableEntry `toml:"entries"`
}

type CommandTableEntry struct {
	Card   string
	Weight float64
}

func (ct CommandTable) SelectByWeight() string {
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
			return entry.Card
		}
	}

	return ""
}
