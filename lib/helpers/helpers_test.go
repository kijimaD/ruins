package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPtr(t *testing.T) {
	t.Run("int value", func(t *testing.T) {
		value := 42
		ptr := GetPtr(value)
		assert.NotNil(t, ptr, "ポインタがnilであってはいけない")
		assert.Equal(t, value, *ptr, "ポインタの値が元の値と一致しない")
		// GetPtrは値のコピーを作成するため、元の変数のアドレスとは異なる
		// この動作が意図されている
	})

	t.Run("string value", func(t *testing.T) {
		value := "hello"
		ptr := GetPtr(value)
		assert.NotNil(t, ptr, "ポインタがnilであってはいけない")
		assert.Equal(t, value, *ptr, "ポインタの値が元の値と一致しない")
	})

	t.Run("float64 value", func(t *testing.T) {
		value := 3.14
		ptr := GetPtr(value)
		assert.NotNil(t, ptr, "ポインタがnilであってはいけない")
		assert.Equal(t, value, *ptr, "ポインタの値が元の値と一致しない")
	})

	t.Run("struct value", func(t *testing.T) {
		type TestStruct struct {
			Field int
		}
		value := TestStruct{Field: 100}
		ptr := GetPtr(value)
		assert.NotNil(t, ptr, "ポインタがnilであってはいけない")
		assert.Equal(t, value, *ptr, "ポインタの値が元の値と一致しない")
	})
}
