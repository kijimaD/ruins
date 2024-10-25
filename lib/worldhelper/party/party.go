package party

import (
	"errors"
	"log"

	"github.com/kijimaD/ruins/lib/components"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/utils/mathutil"
	ecs "github.com/x-hgg-x/goecs/v2"
	"github.com/yourbasic/bit"
)

// グルーピングする単位
// 味方あるいは敵がある
type Party struct {
	// メンバー一覧
	members []ecs.Entity
	// 死んでいる場合はnilが入る
	lives []*ecs.Entity
	// 現在のインデックス。0始まり
	cur int
}

// memberは仲間入れ替えなどをしないと減ったりしない
// 派閥を指定して取得する
// 最初にセットされるインデックスは生存しているエンティティである
// みんな生きていない場合は想定していない。エラーを返す
func NewParty(world w.World, factionType gc.FactionType) (Party, error) {
	gameComponents := world.Components.Game.(*gc.Components)

	var q *bit.Set
	switch factionType {
	case components.FactionAlly:
		q = world.Manager.Join(
			gameComponents.FactionAlly,
			gameComponents.Pools,
			gameComponents.Attributes,
		)
	case components.FactionEnemy:
		q = world.Manager.Join(
			gameComponents.FactionEnemy,
			gameComponents.Pools,
			gameComponents.Attributes,
			gameComponents.CommandTable,
		)
	default:
		log.Fatalf("invalid case: %v")
	}
	members := []ecs.Entity{}
	q.Visit(ecs.Visit(func(entity ecs.Entity) {
		members = append(members, entity)
	}))

	lives := []*ecs.Entity{}
	for _, member := range members {
		pools := gameComponents.Pools.Get(member).(*gc.Pools)
		if pools.HP.Current == 0 {
			lives = append(lives, nil)
		} else {
			lives = append(lives, &member)
		}
	}

	party := Party{
		members: members,
		lives:   lives,
		cur:     0,
	}
	if party.lives[party.cur] == nil {
		err := party.Next()
		if err != nil {
			return Party{}, errors.New("生存Entityが存在しない")
		}
	}

	return party, nil
}

// entityを返す
func (p *Party) Value() *ecs.Entity {
	return p.lives[p.cur]
}

var reachEdgeError = errors.New("reach edge error")

func (p *Party) next() error {
	memo := p.cur
	p.cur = mathutil.Min(p.cur+1, len(p.members)-1)
	if memo == p.cur {
		// 末端に到達してcurが変化しなかった
		return reachEdgeError
	}

	return nil
}

// curを進める
func (p *Party) Next() error {
	for {
		err := p.next()
		if err != nil {
			// 末端に到達した
			return err
		}
		if p.lives[p.cur] == nil {
			continue
		}

		return nil
	}
}

func (p *Party) prev() error {
	memo := p.cur
	p.cur = mathutil.Max(p.cur-1, 0)
	if memo == p.cur {
		// 末端に到達してcurが変化しなかった
		return reachEdgeError
	}

	return nil
}

// curを戻す
func (p *Party) Prev() error {
	for {
		err := p.prev()
		if err != nil {
			// 末端に到達した
			return err
		}
		if p.lives[p.cur] == nil {
			continue
		}

		return nil
	}
}
