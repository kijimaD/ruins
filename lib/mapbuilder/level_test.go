package mapbuilder

import "testing"

func TestGetSpriteNumberForWallType(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		wallType WallType
		expected int
	}{
		{WallTypeTop, 10},
		{WallTypeBottom, 11},
		{WallTypeLeft, 12},
		{WallTypeRight, 13},
		{WallTypeTopLeft, 14},
		{WallTypeTopRight, 15},
		{WallTypeBottomLeft, 16},
		{WallTypeBottomRight, 17},
		{WallTypeGeneric, 1},
	}

	for _, tc := range testCases {
		actual := getSpriteNumberForWallType(tc.wallType)
		if actual != tc.expected {
			t.Errorf("壁タイプ %s のスプライト番号が間違っています。期待値: %d, 実際: %d",
				tc.wallType.String(), tc.expected, actual)
		}
	}
}
