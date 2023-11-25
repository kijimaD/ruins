package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/sokotwo/lib/resources"
	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"
)

type IntroState struct {
	progress int
}

// 全体のテキスト数から位置を計算してるっぽいので、改行してても文章が長いと左端に寄る
var introText = []string{
	"大陸に\n散らばる遺跡...",
	"それは古代の\n失われた文明が眠る場所",
	"その中では古代文明の財宝と\n多くの怪物たちが待ち構えている。",
	// "かつて、勇者たちは剣や槍という\n粗末な武器を持って遺跡に挑み...",
	// "文字通り命がけの冒険の代償に\nわずかばかりの貴金属を持ち帰った。",

	// "やがて時代が進むと、剣は銃に変わり\n大砲を乗せた乗り物が現れた。",
	// "遺跡から見つかる古代の財宝も失われた\n科学技術の品であり利用できることがわかってくると",
	// "遺跡の価値はさらに上がりその支配権をめぐって、\nしばしば争いも起こるようになった。",
	// "重武装の乗り物を駆って遺跡に挑み怪物たちや\nライバルたちと戦い古代の遺品を集めてくるプロ。",
	// "彼らの乗り物を、人々は「バトルディッガー」\nあるいはモグラとそれに乗って遺跡を冒険する者達「モグラ乗り」と呼んだ。",
	// "これは、そんなモグラ乗りの物語である。",
}

// State interface ================

func (st *IntroState) OnPause(world w.World) {}

func (st *IntroState) OnResume(world w.World) {}

func (st *IntroState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	loader.AddEntities(world, prefabs.Intro)
}

func (st *IntroState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

func (st *IntroState) Update(world w.World) states.Transition {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter):
		st.nextPage()
	case inpututil.IsKeyJustPressed(ebiten.KeyBackspace):
		st.prevPage()
	}

	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)
		text.Text = introText[st.progress]
	}))
	return states.Transition{}
}

// utils ================

func (st *IntroState) nextPage() {
	if st.progress < len(introText)-1 {
		st.progress += 1
	}
}

func (st *IntroState) prevPage() {
	if st.progress > 0 {
		st.progress -= 1
	}
}
