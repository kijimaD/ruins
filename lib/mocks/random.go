package mocks

import (
	"math/rand"
)

// 乱数生成のインターフェース
type RandomGenerator interface {
	Float64() float64
	Intn(n int) int
	Shuffle(n int, swap func(i, j int))
	Seed(seed int64)
}

// テスト用のモック乱数生成器
type MockRandomGenerator struct {
	values    []float64
	index     int
	intValues []int
	intIndex  int
}

// 新しいモック乱数生成器を作成する
func NewMockRandomGenerator() *MockRandomGenerator {
	return &MockRandomGenerator{
		values:    []float64{},
		intValues: []int{},
	}
}

// 予め設定された値を順番に返す
func (m *MockRandomGenerator) Float64() float64 {
	if len(m.values) == 0 {
		return 0.5 // デフォルト値
	}
	value := m.values[m.index%len(m.values)]
	m.index++
	return value
}

// 予め設定された整数値を順番に返す
func (m *MockRandomGenerator) Intn(n int) int {
	if len(m.intValues) == 0 {
		return n / 2 // デフォルト値
	}
	value := m.intValues[m.intIndex%len(m.intValues)]
	m.intIndex++
	return value % n
}

// 何もしない（テストでは順序を固定したい場合が多い）
func (m *MockRandomGenerator) Shuffle(n int, swap func(i, j int)) {
	// 何もしない、または予め設定されたパターンでシャッフル
}

// シードは無視する
func (m *MockRandomGenerator) Seed(seed int64) {
	// 何もしない
}

// テスト用のヘルパーメソッド

// Float64で返す値を設定する
func (m *MockRandomGenerator) SetFloat64Values(values ...float64) {
	m.values = values
	m.index = 0
}

// Intnで返す値を設定する
func (m *MockRandomGenerator) SetIntValues(values ...int) {
	m.intValues = values
	m.intIndex = 0
}

// 次に返すFloat64値を追加する
func (m *MockRandomGenerator) AddFloat64Value(value float64) {
	m.values = append(m.values, value)
}

// 次に返すInt値を追加する
func (m *MockRandomGenerator) AddIntValue(value int) {
	m.intValues = append(m.intValues, value)
}

// インデックスをリセットする
func (m *MockRandomGenerator) Reset() {
	m.index = 0
	m.intIndex = 0
}

// 実際のrand.Randを使用する乱数生成器
type RealRandomGenerator struct {
	rng *rand.Rand
}

// 新しい実際の乱数生成器を作成する
func NewRealRandomGenerator(seed int64) *RealRandomGenerator {
	return &RealRandomGenerator{
		rng: rand.New(rand.NewSource(seed)),
	}
}

func (r *RealRandomGenerator) Float64() float64 {
	return r.rng.Float64()
}

func (r *RealRandomGenerator) Intn(n int) int {
	return r.rng.Intn(n)
}

func (r *RealRandomGenerator) Shuffle(n int, swap func(i, j int)) {
	r.rng.Shuffle(n, swap)
}

func (r *RealRandomGenerator) Seed(seed int64) {
	r.rng.Seed(seed)
}
