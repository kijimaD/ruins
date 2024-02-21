package raw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	str := `
[[item]]
name = "リペア"
description = "半分程度回復する"

[[item]]
name = "回復薬"
description = "半分程度回復する"
`
	raw := Load(str)

	expect := RawMaster{
		Raws: Raws{
			Items: []Item{
				Item{Name: "リペア", Description: "半分程度回復する"},
				Item{Name: "回復薬", Description: "半分程度回復する"},
			},
		},
		ItemIndex: map[string]int{
			"リペア": 0,
			"回復薬": 1,
		},
		MemberIndex:   map[string]int{},
		MaterialIndex: map[string]int{},
		RecipeIndex:   map[string]int{},
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
	assert.NotNil(t, entity.Name)
	assert.NotNil(t, entity.Item)
	assert.NotNil(t, entity.Description)
}
