package mapplanner

import (
	"math/rand/v2"
)

// RandomSource はシード値を管理して再現性のあるランダム生成を提供する
type RandomSource struct {
	// Seed はランダム生成のシード値
	Seed uint64
	// rng はシード値から作成したランダムソース
	rng *rand.Rand
}

// NewRandomSource は新しいRandomSourceを作成する
func NewRandomSource(seed uint64) *RandomSource {
	source := rand.NewPCG(seed, seed+1)
	return &RandomSource{
		Seed: seed,
		rng:  rand.New(source),
	}
}

// Intn は0からn-1までのランダムな整数を返す
func (r *RandomSource) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	return r.rng.IntN(n)
}

// Float64 は0.0から1.0までのランダムな浮動小数点数を返す
func (r *RandomSource) Float64() float64 {
	return r.rng.Float64()
}

// Shuffle はスライスの要素をシャッフルする
func (r *RandomSource) Shuffle(n int, swap func(i, j int)) {
	r.rng.Shuffle(n, swap)
}
