package worldhelper

import (
	"errors"
	"log"
	"math/rand/v2"

	"github.com/kijimaD/ruins/lib/mathutil"
	ecs "github.com/x-hgg-x/goecs/v2"
	"github.com/yourbasic/bit"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
)

var errReachEdge = errors.New("reach edge error")

// Party はグルーピングする単位。味方あるいは敵がある
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

// NewParty はmemberは仲間入れ替えなどをしないと減ったりしない
// 派閥を指定して取得する
// 最初にセットされるインデックスは生存しているエンティティである
// みんな生きていない場合は想定していない。エラーを返す
func NewParty(world w.World, factionType gc.FactionType) (Party, error) {
	var q *bit.Set
	switch factionType {
	case gc.FactionAlly:
		q = world.Manager.Join(
			world.Components.FactionAlly,
			world.Components.InParty,
			world.Components.Pools,
			world.Components.Attributes,
		)
	case gc.FactionEnemy:
		q = world.Manager.Join(
			world.Components.FactionEnemy,
			world.Components.Pools,
			world.Components.Attributes,
			world.Components.CommandTable,
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
		if member.HasComponent(world.Components.Dead) {
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

// NewByEntity はentityから派閥を判定して、partyを初期化する
func NewByEntity(world w.World, entity ecs.Entity) (Party, error) {
	var party Party
	var err error

	switch {
	case entity.HasComponent(world.Components.FactionAlly):
		party, err = NewParty(world, gc.FactionAlly)
		if err != nil {
			return party, err
		}
	case entity.HasComponent(world.Components.FactionEnemy):
		party, err = NewParty(world, gc.FactionEnemy)
		if err != nil {
			return party, err
		}
	default:
		return party, errors.New("味方でも敵でもないエンティティが指定された")
	}

	return party, nil
}

// Value は選択中のentityを返す
func (p *Party) Value() *ecs.Entity {
	return p.lives[p.cur]
}

// LivesLen は生存エンティティの数を返す
func (p *Party) LivesLen() int {
	count := 0
	for _, l := range p.lives {
		if l != nil {
			count++
		}
	}

	return count
}

// Next はcurを進める
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

// Prev はcurを戻す
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

// GetNext はcurを進めずに取得だけする
func (p *Party) GetNext() (ecs.Entity, error) {
	cur := p.cur
	for {
		memo := cur
		cur = mathutil.Min(cur+1, len(p.members)-1)
		if memo == cur {
			// 末端に到達してcurが変化しなかった
			return 0, errReachEdge
		}
		if p.lives[cur] == nil {
			continue
		}

		break
	}

	return *p.lives[cur], nil
}

// GetPrev はcurを戻さずに取得だけする
func (p *Party) GetPrev() (ecs.Entity, error) {
	cur := p.cur
	for {
		memo := cur
		cur = mathutil.Max(cur-1, 0)
		if memo == cur {
			// 末端に到達してcurが変化しなかった
			return 0, errReachEdge
		}
		if p.lives[cur] == nil {
			continue
		}

		break
	}

	return *p.lives[cur], nil
}

// GetRandom は生存エンティティからランダムに選択する
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
	p.cur = mathutil.Min(p.cur+1, len(p.members)-1)
	if memo == p.cur {
		// 末端に到達してcurが変化しなかった
		return errReachEdge
	}

	return nil
}

func (p *Party) prev() error {
	memo := p.cur
	p.cur = mathutil.Max(p.cur-1, 0)
	if memo == p.cur {
		// 末端に到達してcurが変化しなかった
		return errReachEdge
	}

	return nil
}
