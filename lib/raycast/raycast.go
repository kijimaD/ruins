// https://ebitengine.org/en/examples/raycasting.html を参考にした
package raycast

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	padding = 20
)

var (
	baseImage   *ebiten.Image // 一番下にある黒背景
	bgImage     *ebiten.Image // 床を表現する
	shadowImage *ebiten.Image // 影を表現する
	visionImage *ebiten.Image // 視界を表現する黒背景
	blackImage  *ebiten.Image // 影生成時の、マスクのベースとして使う黒画像

	vertices   []ebiten.Vertex
	visionNgon = 20
)

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

func (o Object) points() [][2]float64 {
	// すべてのセグメントの終点を取得
	var points [][2]float64
	for _, wall := range o.walls {
		points = append(points, [2]float64{wall.X2, wall.Y2})
	}
	// パスが閉じてない場合、最初のポイントを足すことでパスを閉じる
	p := [2]float64{o.walls[0].X1, o.walls[0].Y1}
	if p[0] != points[len(points)-1][0] && p[1] != points[len(points)-1][1] {
		points = append(points, [2]float64{o.walls[0].X1, o.walls[0].Y1})
	}
	return points
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

// rayCasting returns a slice of line originating from point cx, cy and intersecting with objects
// cx, cyは視点
func rayCasting(cx, cy float64, objects []Object) []line {
	const rayLength = 1000 // something large enough to reach all objects

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

				// rayの視点から最も近い点(=距離が最も小さい点)を求める
				min := math.Inf(1)
				minI := -1
				for i, p := range points {
					d2 := (cx-p[0])*(cx-p[0]) + (cy-p[1])*(cy-p[1])
					if d2 < min {
						min = d2
						minI = i
					}
				}
				rays = append(rays, line{cx, cy, points[minI][0], points[minI][1]})
			}
		}
	}

	// Sort rays based on angle, otherwise light triangles will not come out right
	sort.Slice(rays, func(i int, j int) bool {
		return rays[i].angle() < rays[j].angle()
	})
	return rays
}

func rayVertices(x1, y1, x2, y2, x3, y3 float64) []ebiten.Vertex {
	return []ebiten.Vertex{
		{DstX: float32(x1), DstY: float32(y1), ColorR: 0, ColorG: 0, ColorB: 0, ColorA: 0},
		{DstX: float32(x2), DstY: float32(y2), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x3), DstY: float32(y3), ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
}

func visionVertices(num int, x int, y int, r int) []ebiten.Vertex {
	vs := []ebiten.Vertex{}
	for i := 0; i < num; i++ {
		rate := float64(i) / float64(num)
		cr := 0.0
		cg := 0.0
		cb := 0.0
		vs = append(vs, ebiten.Vertex{
			DstX:   float32(float64(r)*math.Cos(2*math.Pi*rate)) + float32(x),
			DstY:   float32(float64(r)*math.Sin(2*math.Pi*rate)) + float32(y),
			SrcX:   0,
			SrcY:   0,
			ColorR: float32(cr),
			ColorG: float32(cg),
			ColorB: float32(cb),
			ColorA: 1,
		})
	}

	vs = append(vs, ebiten.Vertex{
		DstX:   float32(x),
		DstY:   float32(y),
		SrcX:   0,
		SrcY:   0,
		ColorR: 0,
		ColorG: 0,
		ColorB: 0,
		ColorA: 0,
	})

	return vs
}

type Game struct {
	showRays     bool
	Px, Py       int
	Objects      []Object
	ScreenWidth  int
	ScreenHeight int
}

func (g *Game) Prepare() {
	baseImage = ebiten.NewImage(g.ScreenWidth, g.ScreenHeight)
	shadowImage = ebiten.NewImage(g.ScreenWidth, g.ScreenHeight)
	visionImage = ebiten.NewImage(g.ScreenWidth, g.ScreenHeight)
	blackImage = ebiten.NewImage(g.ScreenWidth, g.ScreenHeight)

	img, _, err := image.Decode(bytes.NewReader(images.Tile_png))
	if err != nil {
		log.Fatal(err)
	}
	bgImage = ebiten.NewImageFromImage(img)

	baseImage.Fill(color.Black)
	blackImage.Fill(color.Black)

	// Add outer walls
	g.Objects = append(g.Objects, Object{rect(padding, padding, float64(g.ScreenWidth)-2*padding, float64(g.ScreenHeight)-2*padding)})

	// Rectangles
	g.Objects = append(g.Objects, Object{rect(45, 50, 70, 20)})
	g.Objects = append(g.Objects, Object{rect(150, 50, 30, 60)})

	g.Objects = append(g.Objects, Object{rect(95, 90, 70, 20)})
	g.Objects = append(g.Objects, Object{rect(120, 150, 30, 60)})

	g.Objects = append(g.Objects, Object{rect(200, 210, 5, 5)})
	g.Objects = append(g.Objects, Object{rect(220, 210, 5, 5)})

	g.Objects = append(g.Objects, Object{rect(100, 250, 500, 20)})
	g.Objects = append(g.Objects, Object{rect(100, 400, 500, 20)})
}

func (g *Game) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.showRays = !g.showRays
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.Px += 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.Py += 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.Px -= 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.Py -= 2
	}

	// +1/-1 is to stop player before it reaches the border
	if g.Px >= g.ScreenWidth-padding {
		g.Px = g.ScreenWidth - padding - 1
	}

	if g.Px <= padding {
		g.Px = padding + 1
	}

	if g.Py >= g.ScreenHeight-padding {
		g.Py = g.ScreenHeight - padding - 1
	}

	if g.Py <= padding {
		g.Py = padding + 1
	}

	return
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Reset the shadowImage
	shadowImage.Fill(color.Black)
	visionImage.Fill(color.Black)
	// position, blockviewをクエリ
	rays := rayCasting(float64(g.Px), float64(g.Py), g.Objects)

	// 全面が黒の画像から、三角形の部分をブレンドで引いて、影になっている部分だけ黒で残す
	{
		opt := &ebiten.DrawTrianglesOptions{}
		opt.Address = ebiten.AddressRepeat
		opt.Blend = ebiten.BlendSourceOut
		for i, line := range rays {
			nextLine := rays[(i+1)%len(rays)]

			// Draw triangle of area between rays
			// vertices: 頂点
			v := rayVertices(float64(g.Px), float64(g.Py), nextLine.X2, nextLine.Y2, line.X2, line.Y2)
			shadowImage.DrawTriangles(v, []uint16{0, 1, 2}, blackImage, opt)
		}
	}

	// Draw background
	screen.DrawImage(baseImage, nil)
	{
		tileWidth, tileHeight := bgImage.Size()
		// 背景画像を敷き詰める
		for i := 0; i < g.ScreenWidth; i += tileWidth {
			for j := 0; j < g.ScreenHeight; j += tileHeight {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i), float64(j))
				screen.DrawImage(bgImage, op)
			}
		}
	}

	// Draw walls
	for _, obj := range g.Objects {
		for _, w := range obj.walls {
			vector.StrokeLine(screen, float32(w.X1), float32(w.Y1), float32(w.X2), float32(w.Y2), 20, color.RGBA{0, 0, 0, 100}, true)
		}
	}

	if g.showRays {
		// Draw rays
		for _, r := range rays {
			vector.StrokeLine(screen, float32(r.X1), float32(r.Y1), float32(r.X2), float32(r.Y2), 1, color.RGBA{255, 255, 0, 150}, true)
		}
	}

	// 視界以外をグラデーションを入れながら塗りつぶし
	{
		vs := visionVertices(visionNgon, g.Px, g.Py, 200)
		opt := &ebiten.DrawTrianglesOptions{}
		opt.Blend = ebiten.BlendSourceIn
		indices := []uint16{}
		for i := 0; i < visionNgon; i++ {
			indices = append(indices, uint16(i), uint16(i+1)%uint16(visionNgon), uint16(visionNgon))
		}
		visionImage.DrawTriangles(vs, indices, blackImage, opt)
	}

	// 光源の中心付近を明るくする
	{
		vs := visionVertices(visionNgon, g.Px, g.Py, 20)
		opt := &ebiten.DrawTrianglesOptions{}
		opt.Blend = ebiten.BlendClear
		indices := []uint16{}
		for i := 0; i < visionNgon; i++ {
			indices = append(indices, uint16(i), uint16(i+1)%uint16(visionNgon), uint16(visionNgon))
		}
		visionImage.DrawTriangles(vs, indices, blackImage, opt)
	}

	// 建物影
	{
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(1)
		screen.DrawImage(shadowImage, op)
	}
	// 視界影
	{
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(1)
		screen.DrawImage(visionImage, op)
	}

	// Draw player as a rect
	vector.DrawFilledRect(screen, float32(g.Px)-16, float32(g.Py)-16, 32, 32, color.Black, true)
	vector.DrawFilledRect(screen, float32(g.Px)-1, float32(g.Py)-1, 2, 2, color.RGBA{255, 100, 100, 255}, true)

	if g.showRays {
		ebitenutil.DebugPrintAt(screen, "R: hide rays", g.ScreenWidth-padding*8, 0)
	} else {
		ebitenutil.DebugPrintAt(screen, "R: show rays", g.ScreenWidth-padding*8, 0)
	}
}

func rect(x, y, w, h float64) []line {
	return []line{
		{x, y, x, y + h},
		{x, y + h, x + w, y + h},
		{x + w, y + h, x + w, y},
		{x + w, y, x, y},
	}
}
