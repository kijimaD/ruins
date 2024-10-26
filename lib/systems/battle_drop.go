package systems

import (
	"math/rand"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/worldhelper/material"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 戦闘終了後に経験値や素材を獲得する
// 通貨も取得するか?
// 表示用に獲得した素材を返す
func BattleDropSystem(world w.World) []string {
	rawMaster := world.Resources.RawMaster.(raw.RawMaster)
	gameComponents := world.Components.Game.(*gc.Components)

	cands := []string{}
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.FactionEnemy,
		gameComponents.Attributes,
		gameComponents.DropTable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := gameComponents.Name.Get(entity).(*gc.Name)
		dt := rawMaster.GetDropTable(name.Name)
		for i := 0; i < 3; i++ {
			cands = append(cands, dt.SelectByWeight())
		}
	}))

	rand.Shuffle(len(cands), func(i, j int) { cands[i], cands[j] = cands[j], cands[i] })
	resultCands := cands[0:3]
	for _, cand := range resultCands {
		material.PlusAmount(cand, 1, world)
	}

	return resultCands
}
