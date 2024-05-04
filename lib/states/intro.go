// ゲームの導入テキストを表示するステート
package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	ec "github.com/kijimaD/ruins/lib/engine/components"
	"github.com/kijimaD/ruins/lib/engine/loader"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/utils/msg"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type IntroState struct {
	queue msg.Queue
	cycle int
	bg    *ebiten.Image
}

func (st IntroState) String() string {
	return "Intro"
}

var introText = `
[image source="bg_urban1"]

人類は古代遺跡とともにあった。[p]
先史時代から粗末な装備で怪物と財宝に満ちた遺跡に挑み、[p]

[image source="bg_crystal1"]
[wait time="500"]

得られたささいな品から生活を発展させてきた。[p]
国家の成立と工業の進展とともに、[l]
遺跡から得られる技術が利用できることがわかってくると、[p]
しばしば支配権をめぐって戦争が行われるようになった。[p]`

// State interface ================

func (st *IntroState) OnPause(world w.World) {}

func (st *IntroState) OnResume(world w.World) {}

func (st *IntroState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	loader.AddEntities(world, prefabs.Intro)
	st.queue = msg.NewQueueFromText(introText)
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
