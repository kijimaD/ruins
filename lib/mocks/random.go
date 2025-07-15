package mocks

import (
	"math/rand"
)

// RandomGenerator は乱数生成のインターフェース
type RandomGenerator interface {
	Float64() float64
	Intn(n int) int
	Shuffle(n int, swap func(i, j int))
	Seed(seed int64)
}

// MockRandomGenerator はテスト用のモック乱数生成器
type MockRandomGenerator struct {
	values    []float64
	index     int
	intValues []int
	intIndex  int
}

// NewMockRandomGenerator は新しいモック乱数生成器を作成する
func NewMockRandomGenerator() *MockRandomGenerator {
	return &MockRandomGenerator{
		values:    []float64{},
		intValues: []int{},
	}
}

// Float64 は予め設定された値を順番に返す
func (m *MockRandomGenerator) Float64() float64 {
	if len(m.values) == 0 {
		return 0.5 // デフォルト値
	}
	value := m.values[m.index%len(m.values)]
	m.index++
	return value
}

// Intn は予め設定された整数値を順番に返す
func (m *MockRandomGenerator) Intn(n int) int {
	if len(m.intValues) == 0 {
		return n / 2 // デフォルト値
	}
	value := m.intValues[m.intIndex%len(m.intValues)]
	m.intIndex++
	return value % n
}

// Shuffle は何もしない（テストでは順序を固定したい場合が多い）
func (m *MockRandomGenerator) Shuffle(_ int, _ func(i, j int)) {
	// 何もしない、または予め設定されたパターンでシャッフル
}

// Seed はシードを無視する
func (m *MockRandomGenerator) Seed(_ int64) {
	// 何もしない
}

// テスト用のヘルパーメソッド

// SetFloat64Values はFloat64で返す値を設定する
func (m *MockRandomGenerator) SetFloat64Values(values ...float64) {
	m.values = values
	m.index = 0
}

// SetIntValues はIntnで返す値を設定する
func (m *MockRandomGenerator) SetIntValues(values ...int) {
	m.intValues = values
	m.intIndex = 0
}

// AddFloat64Value は次に返すFloat64値を追加する
func (m *MockRandomGenerator) AddFloat64Value(value float64) {
	m.values = append(m.values, value)
}

// AddIntValue は次に返すInt値を追加する
func (m *MockRandomGenerator) AddIntValue(value int) {
	m.intValues = append(m.intValues, value)
}

// Reset はインデックスをリセットする
func (m *MockRandomGenerator) Reset() {
	m.index = 0
	m.intIndex = 0
}

// RealRandomGenerator は実際のrand.Randを使用する乱数生成器
type RealRandomGenerator struct {
	rng *rand.Rand
}

// NewRealRandomGenerator は新しい実際の乱数生成器を作成する
func NewRealRandomGenerator(seed int64) *RealRandomGenerator {
	return &RealRandomGenerator{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// Float64 は0.0と1.0の間の擬似乱数を返す
func (r *RealRandomGenerator) Float64() float64 {
	return r.rng.Float64()
}

// Intn は0からn-1までの擬似乱数を返す
func (r *RealRandomGenerator) Intn(n int) int {
	return r.rng.Intn(n)
}

// Shuffle は疑似乱数を使って要素をシャッフルする
func (r *RealRandomGenerator) Shuffle(n int, swap func(i, j int)) {
	r.rng.Shuffle(n, swap)
}

// Seed は乱数生成器のシードを設定する
func (r *RealRandomGenerator) Seed(seed int64) {
	r.rng.Seed(seed)
}
