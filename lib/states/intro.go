// ゲームの導入テキストを表示するステート
package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	ec "github.com/kijimaD/sokotwo/lib/engine/components"
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	"github.com/kijimaD/sokotwo/lib/engine/states"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	"github.com/kijimaD/sokotwo/lib/resources"
	"github.com/kijimaD/sokotwo/lib/utils/msg"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type IntroState struct {
	queue msg.Queue
	cycle int
	bg    *ebiten.Image
}

var introText = `
[image source="bg_urban1"]
<導入>[p]
大陸に散らばる遺跡...[p]
それは古代の失われた文明が眠る場所。[l]
その中では古代文明の財宝と多くの怪物たちが待ち構えている。[p]
<昔>[p]
かつて、勇者たちは剣や槍という粗末な武器を持って遺跡に挑み...[p]
[image source="bg_crystal1"]
[wait time="400"]
文字通り命がけの冒険の代償にわずかばかりの貴金属を持ち帰った。[p]
<変遷>[p]
やがて時代が進むと、剣は銃に変わり大砲を乗せた乗り物が現れた。[p]
遺跡から見つかる古代の財宝も失われた科学技術の品であり利用できることがわかってくると遺跡の価値はさらに上がりその支配権をめぐって、しばしば争いも起こるようになった。[p]
<テーマ>[p]
重武装の乗り物を駆って遺跡に挑み怪物たちやライバルたちと戦い古代の遺品を集めてくるプロ。[p]
彼らの乗り物を、人々は「バトルディッガー」あるいはモグラとそれに乗って遺跡を冒険する者達を「モグラ乗り」と呼んだ。[p]
これは、そんなモグラ乗りの物語である。`

// State interface ================

func (st *IntroState) OnPause(world w.World) {}

func (st *IntroState) OnResume(world w.World) {}

func (st *IntroState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	loader.AddEntities(world, prefabs.Intro)
	l := msg.NewLexer(introText)
	p := msg.NewParser(l)
	program := p.ParseProgram()
	e := msg.Evaluator{}
	e.Eval(program)
	st.queue = msg.NewQueue(e.Events)
}

func (st *IntroState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

func (st *IntroState) Update(world w.World) states.Transition {
	// アニメーションに便利なので、グローバルにあっていいかもしれない
	var queueResult msg.QueueState

	if v, ok := st.queue.Head().(*msg.ChangeBg); ok {
		world.Manager.Join(world.Components.Engine.SpriteRender, world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
			sprite := world.Components.Engine.SpriteRender.Get(entity).(*ec.SpriteRender)
			new := (*world.Resources.SpriteSheets)[v.Source]
			sprite.SpriteSheet = &new
		}))
	}
	if st.cycle%2 == 0 {
		queueResult = st.queue.RunHead()
		st.cycle = 0
	}
	st.cycle++

	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter):
		queueResult = st.queue.Pop()
	case inpututil.IsKeyJustPressed(ebiten.KeyEscape):
		// debug
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}}
	}

	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)
		text.Text = st.queue.Display()
	}))

	switch queueResult {
	case msg.QueueStateFinish:
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}}
	}
	return states.Transition{}
}

func (st *IntroState) Draw(world w.World, screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	if st.bg != nil {
		screen.DrawImage(st.bg, opts)
	}
}
