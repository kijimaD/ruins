package worldhelper

import (
	"errors"
	"log"
	"math/rand/v2"

	"github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/utils"
	ecs "github.com/x-hgg-x/goecs/v2"
	"github.com/yourbasic/bit"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
)

var reachEdgeError = errors.New("reach edge error")

// グルーピングする単位。味方あるいは敵がある
type Party struct {
	// メンバー一覧
	// entityの番号順に並んでいるという前提で書いている
	// 後々並びが変わることもあるだろうから、そのときに対応する
	members []ecs.Entity
	// 死んでいる場合はnilが入る。membersと同じ長さになる
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
		log.Fatalf("invalid case: %v", factionType)
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

// entityから派閥を判定して、partyを初期化する
func NewByEntity(world w.World, entity ecs.Entity) (Party, error) {
	var party Party
	var err error

	gameComponents := world.Components.Game.(*gc.Components)
	switch {
	case entity.HasComponent(gameComponents.FactionAlly):
		party, err = NewParty(world, gc.FactionAlly)
		if err != nil {
			return party, err
		}
	case entity.HasComponent(gameComponents.FactionEnemy):
		party, err = NewParty(world, gc.FactionEnemy)
		if err != nil {
			return party, err
		}
	default:
		return party, errors.New("味方でも敵でもないエンティティが指定された")
	}

	return party, nil
}

// 選択中のentityを返す
func (p *Party) Value() *ecs.Entity {
	return p.lives[p.cur]
}

// 生存エンティティの数を返す
func (p *Party) LivesLen() int {
	count := 0
	for _, l := range p.lives {
		if l != nil {
			count += 1
		}
	}

	return count
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

// curを進めずに取得だけする
func (p *Party) GetNext() (ecs.Entity, error) {
	cur := p.cur
	for {
		memo := cur
		cur = utils.Min(cur+1, len(p.members)-1)
		if memo == cur {
			// 末端に到達してcurが変化しなかった
			return 0, reachEdgeError
		}
		if p.lives[cur] == nil {
			continue
		}

		break
	}

	return *p.lives[cur], nil
}

// curを戻さずに取得だけする
func (p *Party) GetPrev() (ecs.Entity, error) {
	cur := p.cur
	for {
		memo := cur
		cur = utils.Max(cur-1, 0)
		if memo == cur {
			// 末端に到達してcurが変化しなかった
			return 0, reachEdgeError
		}
		if p.lives[cur] == nil {
			continue
		}

		break
	}

	return *p.lives[cur], nil
}

// 生存エンティティからランダムに選択する
func (p *Party) GetRandom() (ecs.Entity, error) {
	lives := []ecs.Entity{}
	for _, live := range p.lives {
		lives = append(lives, *live)
	}
	if len(lives) == 0 {
		return 0, errors.New("生存エンティティが存在しない")
	}
	idx := rand.IntN(len(lives) - 1)

	return lives[idx], nil
}

func (p *Party) next() error {
	memo := p.cur
	p.cur = utils.Min(p.cur+1, len(p.members)-1)
	if memo == p.cur {
		// 末端に到達してcurが変化しなかった
		return reachEdgeError
	}

	return nil
}

func (p *Party) prev() error {
	memo := p.cur
	p.cur = utils.Max(p.cur-1, 0)
	if memo == p.cur {
		// 末端に到達してcurが変化しなかった
		return reachEdgeError
	}

	return nil
}
