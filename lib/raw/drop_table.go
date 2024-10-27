package raw

import "math/rand/v2"

type DropTable struct {
	Name    string
	XpBase  float64
	Entries []DropTableEntry `toml:"entries"`
}

type DropTableEntry struct {
	Material string
	Weight   float64
}

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
