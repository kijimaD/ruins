package systems

import (
	"math"
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/stretchr/testify/assert"
)

func TestVisionVertices(t *testing.T) {
	t.Parallel()
	t.Run("create vision vertices with basic parameters", func(t *testing.T) {
		t.Parallel()
		num := 4
		x := gc.Pixel(100)
		y := gc.Pixel(200)
		r := gc.Pixel(50)

		vertices := visionVertices(num, x, y, r)

		// 実際の頂点数は要求されたnum+1（中心点が追加される）
		assert.Len(t, vertices, num+1, "頂点の数が正しくない")

		// 円周上の頂点（最初のnum個）の基本的な性質を確認
		for i := 0; i < num; i++ {
			vertex := vertices[i]
			assert.Equal(t, float32(1), vertex.ColorA, "頂点%dのアルファ値が正しくない", i)
			assert.Equal(t, float32(0), vertex.SrcX, "頂点%dのSrcXが正しくない", i)
			assert.Equal(t, float32(0), vertex.SrcY, "頂点%dのSrcYが正しくない", i)
			assert.Equal(t, float32(0), vertex.ColorR, "頂点%dのColorRが正しくない", i)
			assert.Equal(t, float32(0), vertex.ColorG, "頂点%dのColorGが正しくない", i)
			assert.Equal(t, float32(0), vertex.ColorB, "頂点%dのColorBが正しくない", i)
		}

		// 中心点（最後の頂点）の確認
		centerVertex := vertices[num]
		assert.Equal(t, float32(x), centerVertex.DstX, "中心点のX座標が正しくない")
		assert.Equal(t, float32(y), centerVertex.DstY, "中心点のY座標が正しくない")
		assert.Equal(t, float32(0), centerVertex.ColorA, "中心点のアルファ値が正しくない")
	})

	t.Run("verify circular positioning", func(t *testing.T) {
		t.Parallel()
		num := 4
		x := gc.Pixel(0)
		y := gc.Pixel(0)
		r := gc.Pixel(100)

		vertices := visionVertices(num, x, y, r)

		// 4つの頂点で円を描く場合の座標を確認
		// 0度: (100, 0), 90度: (0, 100), 180度: (-100, 0), 270度: (0, -100)
		tolerance := 0.001

		// 第1頂点 (0度)
		assert.InDelta(t, 100.0, vertices[0].DstX, tolerance, "第1頂点のX座標が正しくない")
		assert.InDelta(t, 0.0, vertices[0].DstY, tolerance, "第1頂点のY座標が正しくない")

		// 第2頂点 (90度)
		assert.InDelta(t, 0.0, vertices[1].DstX, tolerance, "第2頂点のX座標が正しくない")
		assert.InDelta(t, 100.0, vertices[1].DstY, tolerance, "第2頂点のY座標が正しくない")

		// 第3頂点 (180度)
		assert.InDelta(t, -100.0, vertices[2].DstX, tolerance, "第3頂点のX座標が正しくない")
		assert.InDelta(t, 0.0, vertices[2].DstY, tolerance, "第3頂点のY座標が正しくない")

		// 第4頂点 (270度)
		assert.InDelta(t, 0.0, vertices[3].DstX, tolerance, "第4頂点のX座標が正しくない")
		assert.InDelta(t, -100.0, vertices[3].DstY, tolerance, "第4頂点のY座標が正しくない")
	})

	t.Run("verify distance from center", func(t *testing.T) {
		t.Parallel()
		num := 8
		x := gc.Pixel(50)
		y := gc.Pixel(75)
		r := gc.Pixel(30)

		vertices := visionVertices(num, x, y, r)

		// 円周上の頂点（最初のnum個）が中心からr距離にあることを確認
		for i := 0; i < num; i++ {
			vertex := vertices[i]
			dx := float64(vertex.DstX - float32(x))
			dy := float64(vertex.DstY - float32(y))
			distance := math.Sqrt(dx*dx + dy*dy)
			assert.InDelta(t, float64(r), distance, 0.001, "頂点%dが中心からの距離が正しくない", i)
		}

		// 中心点は距離0であることを確認
		centerVertex := vertices[num]
		dx := float64(centerVertex.DstX - float32(x))
		dy := float64(centerVertex.DstY - float32(y))
		distance := math.Sqrt(dx*dx + dy*dy)
		assert.InDelta(t, 0.0, distance, 0.001, "中心点が中心からの距離が正しくない")
	})

	t.Run("empty vertices for zero count", func(t *testing.T) {
		t.Parallel()
		num := 0
		x := gc.Pixel(100)
		y := gc.Pixel(200)
		r := gc.Pixel(50)

		vertices := visionVertices(num, x, y, r)

		// 0個の円周点でも中心点は追加されるため、1個の頂点が返される
		assert.Len(t, vertices, 1, "0個の頂点要求で中心点のみのスライスが返されていない")

		// 中心点の確認
		centerVertex := vertices[0]
		assert.Equal(t, float32(x), centerVertex.DstX, "中心点のX座標が正しくない")
		assert.Equal(t, float32(y), centerVertex.DstY, "中心点のY座標が正しくない")
		assert.Equal(t, float32(0), centerVertex.ColorA, "中心点のアルファ値が正しくない")
	})

	t.Run("single vertex", func(t *testing.T) {
		t.Parallel()
		num := 1
		x := gc.Pixel(10)
		y := gc.Pixel(20)
		r := gc.Pixel(5)

		vertices := visionVertices(num, x, y, r)

		// 1個の円周点+1個の中心点=2個の頂点
		assert.Len(t, vertices, 2, "1個の頂点要求で2個のスライスが返されていない")

		// 単一頂点は0度の位置 (x+r, y)
		assert.InDelta(t, 15.0, vertices[0].DstX, 0.001, "単一頂点のX座標が正しくない")
		assert.InDelta(t, 20.0, vertices[0].DstY, 0.001, "単一頂点のY座標が正しくない")

		// 中心点の確認
		centerVertex := vertices[1]
		assert.Equal(t, float32(x), centerVertex.DstX, "中心点のX座標が正しくない")
		assert.Equal(t, float32(y), centerVertex.DstY, "中心点のY座標が正しくない")
		assert.Equal(t, float32(0), centerVertex.ColorA, "中心点のアルファ値が正しくない")
	})
}
