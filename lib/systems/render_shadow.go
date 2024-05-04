package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	w "github.com/kijimaD/ruins/lib/engine/world"
)

var (
	shadowImage = ebiten.NewImage(1000, 1000) // 影生成時の、マスクのベースとして使う黒画像
)

func RenderShadowSystem(world w.World, screen *ebiten.Image) {
	// gameComponents := world.Components.Game.(*gc.Components)

	// var pos *gc.Position
	// world.Manager.Join(
	// 	gameComponents.Position,
	// 	gameComponents.Player,
	// ).Visit(ecs.Visit(func(entity ecs.Entity) {
	// 	pos = gameComponents.Position.Get(entity).(*gc.Position)
	// }))

	// shadowImage.Fill(color.Black)
	// rays := rayCasting(float64(pos.X), float64(pos.Y), g.Objects)

	// // 全面が黒の画像から、三角形の部分をブレンドで引いて、影になっている部分だけ黒で残す
	// {
	// 	opt := &ebiten.DrawTrianglesOptions{}
	// 	opt.Address = ebiten.AddressRepeat
	// 	opt.Blend = ebiten.BlendSourceOut
	// 	for i, line := range rays {
	// 		nextLine := rays[(i+1)%len(rays)]

	// 		// Draw triangle of area between rays
	// 		// vertices: 頂点
	// 		v := rayVertices(float64(pos.X), float64(pos.Y), nextLine.X2, nextLine.Y2, line.X2, line.Y2)
	// 		shadowImage.DrawTriangles(v, []uint16{0, 1, 2}, blackImage, opt)
	// 	}
	// }
}

func rayVertices(x1, y1, x2, y2, x3, y3 float64) []ebiten.Vertex {
	return []ebiten.Vertex{
		{DstX: float32(x1), DstY: float32(y1), ColorR: 0, ColorG: 0, ColorB: 0, ColorA: 0},
		{DstX: float32(x2), DstY: float32(y2), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x3), DstY: float32(y3), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
}

// func rayCasting(cx, cy float64, objects []Object) []line {
// 	const rayLength = 1000 // something large enough to reach all objects

// 	var rays []line
// 	for _, obj := range objects {
// 		// Cast two rays per point
// 		for _, p := range obj.points() {
// 			l := line{cx, cy, p[0], p[1]}
// 			angle := l.angle()

// 			// 微妙に角度をつけて影を自然に見せる(直線にならないようにする)
// 			for _, offset := range []float64{-0.005, 0.005} {
// 				points := [][2]float64{}
// 				ray := newRay(cx, cy, rayLength, angle+offset)

// 				// Unpack all objects
// 				for _, o := range objects {
// 					for _, wall := range o.walls {
// 						if px, py, ok := intersection(ray, wall); ok {
// 							points = append(points, [2]float64{px, py})
// 						}
// 					}
// 				}

// 				// rayの視点から最も近い点(=距離が最も小さい点)を求める
// 				min := math.Inf(1)
// 				minI := -1
// 				for i, p := range points {
// 					d2 := (cx-p[0])*(cx-p[0]) + (cy-p[1])*(cy-p[1])
// 					if d2 < min {
// 						min = d2
// 						minI = i
// 					}
// 				}
// 				rays = append(rays, line{cx, cy, points[minI][0], points[minI][1]})
// 			}
// 		}
// 	}

// 	// Sort rays based on angle, otherwise light triangles will not come out right
// 	sort.Slice(rays, func(i int, j int) bool {
// 		return rays[i].angle() < rays[j].angle()
// 	})
// 	return rays
// }
