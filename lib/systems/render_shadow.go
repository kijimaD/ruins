package systems

import (
	"image/color"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	ec "github.com/kijimaD/ruins/lib/engine/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var (
	shadowImage = ebiten.NewImage(1000, 1000) // 影生成時の、マスクのベースとして使う黒画像
)

func RenderShadowSystem(world w.World, screen *ebiten.Image) {
	gameComponents := world.Components.Game.(*gc.Components)

	var pos *gc.Position
	world.Manager.Join(
		gameComponents.Position,
		gameComponents.Player,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pos = gameComponents.Position.Get(entity).(*gc.Position)
	}))

	shadowImage.Fill(color.Black)
	rays := rayCasting(float64(pos.X), float64(pos.Y), world)

	// 全面が黒の画像から、三角形の部分をブレンドで引いて、影になっている部分だけ黒で残す
	{
		opt := &ebiten.DrawTrianglesOptions{}
		opt.Address = ebiten.AddressRepeat
		opt.Blend = ebiten.BlendSourceOut
		for i, line := range rays {
			nextLine := rays[(i+1)%len(rays)]

			// Draw triangle of area between rays
			// vertices: 頂点
			v := rayVertices(float64(pos.X), float64(pos.Y), nextLine.X2, nextLine.Y2, line.X2, line.Y2)
			shadowImage.DrawTriangles(v, []uint16{0, 1, 2}, blackImage, opt)
		}
	}

	// Draw rays
	// for _, r := range rays {
	// 	vector.StrokeLine(screen, float32(r.X1), float32(r.Y1), float32(r.X2), float32(r.Y2), 1, color.RGBA{255, 255, 0, 150}, true)
	// }

	{
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(1)
		screen.DrawImage(shadowImage, op)
	}
}

func rayVertices(x1, y1, x2, y2, x3, y3 float64) []ebiten.Vertex {
	return []ebiten.Vertex{
		{DstX: float32(x1), DstY: float32(y1), ColorR: 0, ColorG: 0, ColorB: 0, ColorA: 0},
		{DstX: float32(x2), DstY: float32(y2), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x3), DstY: float32(y3), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
}

func rayCasting(cx, cy float64, world w.World) []line {
	const rayLength = 10000 // something large enough to reach all objects

	objects := []Object{}
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Position,
		gameComponents.SpriteRender,
		gameComponents.BlockView,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		if !entity.HasComponent(gameComponents.Player) {
			pos := gameComponents.Position.Get(entity).(*gc.Position)
			spriteRender := gameComponents.SpriteRender.Get(entity).(*ec.SpriteRender)
			sprite := spriteRender.SpriteSheet.Sprites[spriteRender.SpriteNumber]

			x := float64(pos.X - sprite.Width/2)
			y := float64(pos.Y - sprite.Height/2)
			w := float64(sprite.Width)
			h := float64(sprite.Height)
			objects = append(objects, Object{rect(x, y, w, h)})
		}
	}))

	// 外周の壁。rayが必ずどこかに当たるようにしないといけない
	{
		screenWidth := float64(world.Resources.ScreenDimensions.Width)
		screenHeight := float64(world.Resources.ScreenDimensions.Height)
		padding := float64(20)
		objects = append(objects, Object{rect(padding, padding, float64(screenWidth-2*padding), float64(screenHeight-2*padding))})
	}

	var rays []line
	for _, obj := range objects {
		// Cast two rays per point
		for _, p := range obj.points() {
			l := line{cx, cy, p[0], p[1]}
			angle := l.angle()

			// 微妙に角度をつけて影を自然に見せる(直線にならないようにする)
			for _, offset := range []float64{-0.005, 0.005} {
				points := [][2]float64{}
				ray := newRay(cx, cy, rayLength, angle+offset)

				// Unpack all objects
				for _, o := range objects {
					for _, wall := range o.walls {
						if px, py, ok := intersection(ray, wall); ok {
							points = append(points, [2]float64{px, py})
						}
					}
				}

				// 視点から最も近い交点までの線分を rays スライスに追加する
				min := math.Inf(1)
				minIdx := -1
				for i, p := range points {
					d2 := (cx-p[0])*(cx-p[0]) + (cy-p[1])*(cy-p[1])
					if d2 < min {
						min = d2
						minIdx = i
					}
				}
				rays = append(rays, line{cx, cy, points[minIdx][0], points[minIdx][1]})
			}
		}
	}

	// Sort rays based on angle, otherwise light triangles will not come out right
	sort.Slice(rays, func(i int, j int) bool {
		return rays[i].angle() < rays[j].angle()
	})
	return rays
}

// 始点と終点
type line struct {
	X1, Y1, X2, Y2 float64
}

func (l *line) angle() float64 {
	return math.Atan2(l.Y2-l.Y1, l.X2-l.X1)
}

type Object struct {
	walls []line
}

// 基本的に4点
func (o Object) points() [][2]float64 {
	// すべてのセグメントの終点を取得
	var points [][2]float64
	for _, wall := range o.walls {
		points = append(points, [2]float64{wall.X2, wall.Y2})
	}
	// パスが閉じてない場合、最初のポイントを足すことでパスを閉じる
	// p := [2]float64{o.walls[0].X1, o.walls[0].Y1}
	// if p[0] != points[len(points)-1][0] && p[1] != points[len(points)-1][1] {
	// 	points = append(points, [2]float64{o.walls[0].X1, o.walls[0].Y1})
	// }

	return points
}

func rect(x, y, w, h float64) []line {
	return []line{
		{x, y, x, y + h},
		{x, y + h, x + w, y + h},
		{x + w, y + h, x + w, y},
		{x + w, y, x, y},
	}
}

// 角度を変えて新しい線を生成する
func newRay(x, y, length, angle float64) line {
	return line{
		X1: x,
		Y1: y,
		X2: x + length*math.Cos(angle),
		Y2: y + length*math.Sin(angle),
	}
}

// intersection calculates the intersection of given two lines.
func intersection(l1, l2 line) (float64, float64, bool) {
	// https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection#Given_two_points_on_each_line
	denom := (l1.X1-l1.X2)*(l2.Y1-l2.Y2) - (l1.Y1-l1.Y2)*(l2.X1-l2.X2)
	tNum := (l1.X1-l2.X1)*(l2.Y1-l2.Y2) - (l1.Y1-l2.Y1)*(l2.X1-l2.X2)
	uNum := -((l1.X1-l1.X2)*(l1.Y1-l2.Y1) - (l1.Y1-l1.Y2)*(l1.X1-l2.X1))

	if denom == 0 {
		return 0, 0, false
	}

	t := tNum / denom
	if t > 1 || t < 0 {
		return 0, 0, false
	}

	u := uNum / denom
	if u > 1 || u < 0 {
		return 0, 0, false
	}

	x := l1.X1 + t*(l1.X2-l1.X1)
	y := l1.Y1 + t*(l1.Y2-l1.Y1)
	return x, y, true
}
