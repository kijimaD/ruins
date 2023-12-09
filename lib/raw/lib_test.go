package raw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	str := `
[[item]]
name = "リペア"

[[item]]
name = "回復薬"
`
	raw := Load(str)

	expect := RawMaster{
		Raws: Raws{
			Items: []Item{
				Item{Name: "リペア"},
				Item{Name: "回復薬"},
			},
		},
		ItemIndex: map[string]int{
			"リペア": 0,
			"回復薬": 1,
		},
	}
	assert.Equal(t, expect, raw)
}

func TestGenerateItem(t *testing.T) {
	str := `
[[item]]
name = "リペア"
`
	raw := Load(str)
	entity := raw.GenerateItem("リペア")
	assert.NotNil(t, entity.Components.Name)
	assert.NotNil(t, entity.Components.Item)
}
