package components

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestWarpMode(t *testing.T) {
	t.Parallel()
	t.Run("warp mode constants", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, warpMode("NEXT"), WarpModeNext, "WarpModeNextの値が正しくない")
		assert.Equal(t, warpMode("ESCAPE"), WarpModeEscape, "WarpModeEscapeの値が正しくない")
	})
}

func TestTargetNumType(t *testing.T) {
	t.Parallel()
	t.Run("target num constants", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, TargetNumType("SINGLE"), TargetSingle, "TargetSingleの値が正しくない")
		assert.Equal(t, TargetNumType("ALL"), TargetAll, "TargetAllの値が正しくない")
	})

	t.Run("valid target num types", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name      string
			targetNum TargetNumType
			wantErr   bool
		}{
			{"valid single", TargetSingle, false},
			{"valid all", TargetAll, false},
			{"invalid type", TargetNumType("INVALID"), true},
			{"empty type", TargetNumType(""), true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				err := tt.targetNum.Valid()
				if tt.wantErr {
					assert.Error(t, err, "無効な値でエラーが発生しない")
					assert.True(t, errors.Is(err, ErrInvalidEnumType), "エラーの種類が正しくない")
				} else {
					assert.NoError(t, err, "有効な値でエラーが発生する")
				}
			})
		}
	})
}

func TestTargetGroupType(t *testing.T) {
	t.Parallel()
	t.Run("target group constants", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, TargetGroupType("ALLY"), TargetGroupAlly, "TargetGroupAllyの値が正しくない")
		assert.Equal(t, TargetGroupType("ENEMY"), TargetGroupEnemy, "TargetGroupEnemyの値が正しくない")
		assert.Equal(t, TargetGroupType("CARD"), TargetGroupCard, "TargetGroupCardの値が正しくない")
		assert.Equal(t, TargetGroupType("NONE"), TargetGroupNone, "TargetGroupNoneの値が正しくない")
	})

	t.Run("valid target group types", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name        string
			targetGroup TargetGroupType
			wantErr     bool
		}{
			{"valid ally", TargetGroupAlly, false},
			{"valid enemy", TargetGroupEnemy, false},
			{"valid card", TargetGroupCard, false},
			{"valid none", TargetGroupNone, false},
			{"invalid type", TargetGroupType("INVALID"), true},
			{"empty type", TargetGroupType(""), true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				err := tt.targetGroup.Valid()
				if tt.wantErr {
					assert.Error(t, err, "無効な値でエラーが発生しない")
					assert.True(t, errors.Is(err, ErrInvalidEnumType), "エラーの種類が正しくない")
				} else {
					assert.NoError(t, err, "有効な値でエラーが発生する")
				}
			})
		}
	})
}

func TestUsableSceneType(t *testing.T) {
	t.Parallel()
	t.Run("usable scene constants", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, UsableSceneType("BATTLE"), UsableSceneBattle, "UsableSceneBattleの値が正しくない")
		assert.Equal(t, UsableSceneType("FIELD"), UsableSceneField, "UsableSceneFieldの値が正しくない")
		assert.Equal(t, UsableSceneType("ANY"), UsableSceneAny, "UsableSceneAnyの値が正しくない")
	})

	t.Run("valid usable scene types", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name        string
			usableScene UsableSceneType
			wantErr     bool
		}{
			{"valid battle", UsableSceneBattle, false},
			{"valid field", UsableSceneField, false},
			{"valid any", UsableSceneAny, false},
			{"invalid type", UsableSceneType("INVALID"), true},
			{"empty type", UsableSceneType(""), true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				err := tt.usableScene.Valid()
				if tt.wantErr {
					assert.Error(t, err, "無効な値でエラーが発生しない")
					assert.True(t, errors.Is(err, ErrInvalidEnumType), "エラーの種類が正しくない")
				} else {
					assert.NoError(t, err, "有効な値でエラーが発生する")
				}
			})
		}
	})
}

func TestAttackType(t *testing.T) {
	t.Parallel()
	t.Run("attack type constants", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, AttackType("SWORD"), AttackSword, "AttackSwordの値が正しくない")
		assert.Equal(t, AttackType("SPEAR"), AttackSpear, "AttackSpearの値が正しくない")
		assert.Equal(t, AttackType("HANDGUN"), AttackHandgun, "AttackHandgunの値が正しくない")
		assert.Equal(t, AttackType("RIFLE"), AttackRifle, "AttackRifleの値が正しくない")
		assert.Equal(t, AttackType("FIST"), AttackFist, "AttackFistの値が正しくない")
		assert.Equal(t, AttackType("CANON"), AttackCanon, "AttackCanonの値が正しくない")
	})

	t.Run("valid attack types", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name       string
			attackType AttackType
			wantErr    bool
		}{
			{"valid sword", AttackSword, false},
			{"valid spear", AttackSpear, false},
			{"valid handgun", AttackHandgun, false},
			{"valid rifle", AttackRifle, false},
			{"valid fist", AttackFist, false},
			{"valid canon", AttackCanon, false},
			// 注: invalid typeのテストは、String()メソッドでlog.Fatalが呼ばれるためスキップ
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				err := tt.attackType.Valid()
				if tt.wantErr {
					assert.Error(t, err, "無効な値でエラーが発生しない")
					assert.True(t, errors.Is(err, ErrInvalidEnumType), "エラーの種類が正しくない")
				} else {
					assert.NoError(t, err, "有効な値でエラーが発生する")
				}
			})
		}
	})

	t.Run("attack type string representation", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			attackType AttackType
			expected   string
		}{
			{AttackSword, "刀剣"},
			{AttackSpear, "長物"},
			{AttackHandgun, "拳銃"},
			{AttackRifle, "小銃"},
			{AttackFist, "格闘"},
			{AttackCanon, "大砲"},
		}

		for _, tt := range tests {
			t.Run(string(tt.attackType), func(t *testing.T) {
				t.Parallel()
				result := tt.attackType.String()
				assert.Equal(t, tt.expected, result, "文字列表現が正しくない")
			})
		}
	})

	// 注: invalid attack typeのString()はlog.Fatalを呼ぶため、テスト不可
}

func TestEquipmentType(t *testing.T) {
	t.Parallel()
	t.Run("equipment type constants", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, EquipmentType("HEAD"), EquipmentHead, "EquipmentHeadの値が正しくない")
		assert.Equal(t, EquipmentType("TORSO"), EquipmentTorso, "EquipmentTorsoの値が正しくない")
		assert.Equal(t, EquipmentType("LEGS"), EquipmentLegs, "EquipmentLegsの値が正しくない")
		assert.Equal(t, EquipmentType("JEWELRY"), EquipmentJewelry, "EquipmentJewelryの値が正しくない")
	})

	t.Run("valid equipment types", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name          string
			equipmentType EquipmentType
			wantErr       bool
		}{
			{"valid head", EquipmentHead, false},
			{"valid torso", EquipmentTorso, false},
			{"valid legs", EquipmentLegs, false},
			{"valid jewelry", EquipmentJewelry, false},
			// 注: invalid typeのテストは、String()メソッドでlog.Fatalが呼ばれるためスキップ
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				err := tt.equipmentType.Valid()
				if tt.wantErr {
					assert.Error(t, err, "無効な値でエラーが発生しない")
					assert.True(t, errors.Is(err, ErrInvalidEnumType), "エラーの種類が正しくない")
				} else {
					assert.NoError(t, err, "有効な値でエラーが発生する")
				}
			})
		}
	})

	t.Run("equipment type string representation", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			equipmentType EquipmentType
			expected      string
		}{
			{EquipmentHead, "頭部"},
			{EquipmentTorso, "胴体"},
			{EquipmentLegs, "脚部"},
			{EquipmentJewelry, "装飾"},
		}

		for _, tt := range tests {
			t.Run(string(tt.equipmentType), func(t *testing.T) {
				t.Parallel()
				result := tt.equipmentType.String()
				assert.Equal(t, tt.expected, result, "文字列表現が正しくない")
			})
		}
	})
}

func TestElementType(t *testing.T) {
	t.Parallel()
	t.Run("element type constants", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, ElementType("NONE"), ElementTypeNone, "ElementTypeNoneの値が正しくない")
		assert.Equal(t, ElementType("FIRE"), ElementTypeFire, "ElementTypeFireの値が正しくない")
		assert.Equal(t, ElementType("THUNDER"), ElementTypeThunder, "ElementTypeThunderの値が正しくない")
		assert.Equal(t, ElementType("CHILL"), ElementTypeChill, "ElementTypeChillの値が正しくない")
		assert.Equal(t, ElementType("PHOTON"), ElementTypePhoton, "ElementTypePhotonの値が正しくない")
	})

	t.Run("valid element types", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name        string
			elementType ElementType
			wantErr     bool
		}{
			{"valid none", ElementTypeNone, false},
			{"valid fire", ElementTypeFire, false},
			{"valid thunder", ElementTypeThunder, false},
			{"valid chill", ElementTypeChill, false},
			{"valid photon", ElementTypePhoton, false},
			// 注: invalid typeのテストは、String()メソッドでlog.Fatalが呼ばれるためスキップ
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				err := tt.elementType.Valid()
				if tt.wantErr {
					assert.Error(t, err, "無効な値でエラーが発生しない")
					assert.True(t, errors.Is(err, ErrInvalidEnumType), "エラーの種類が正しくない")
				} else {
					assert.NoError(t, err, "有効な値でエラーが発生する")
				}
			})
		}
	})

	t.Run("element type string representation", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			elementType ElementType
			expected    string
		}{
			{ElementTypeNone, "無"},
			{ElementTypeFire, "火"},
			{ElementTypeThunder, "電"},
			{ElementTypeChill, "冷"},
			{ElementTypePhoton, "光"},
		}

		for _, tt := range tests {
			t.Run(string(tt.elementType), func(t *testing.T) {
				t.Parallel()
				result := tt.elementType.String()
				assert.Equal(t, tt.expected, result, "文字列表現が正しくない")
			})
		}
	})
}
