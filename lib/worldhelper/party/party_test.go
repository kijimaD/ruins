package party

import (
	"testing"

	"github.com/kijimaD/ruins/lib/utils"
	"github.com/stretchr/testify/assert"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestNext(t *testing.T) {
	party := Party{
		members: []ecs.Entity{0, 1, 2, 3},
		lives: []*ecs.Entity{
			utils.GetPtr(ecs.Entity(0)),
			nil,
			utils.GetPtr(ecs.Entity(2)),
			utils.GetPtr(ecs.Entity(3)),
		},
		cur: 0,
	}
	{
		// nilを飛ばして0->2で取得できる
		err := party.Next()
		assert.NoError(t, err)
		assert.Equal(t, 2, int(*party.Value()))
		assert.Equal(t, 2, party.cur)
	}
	{
		// 2->3で取得できる
		err := party.Next()
		assert.NoError(t, err)
		assert.Equal(t, 3, int(*party.Value()))
		assert.Equal(t, 3, party.cur)
	}
	{
		// 末端に到達した
		err := party.Next()
		assert.Error(t, err)
	}
}

func TestPrev(t *testing.T) {
	party := Party{
		members: []ecs.Entity{0, 1, 2},
		lives: []*ecs.Entity{
			utils.GetPtr(ecs.Entity(0)),
			nil,
			utils.GetPtr(ecs.Entity(2)),
		},
		cur: 0,
	}
	{
		// 末端に到達した
		err := party.Prev()
		assert.Error(t, err)
	}
	{
		// 0->2で進める
		err := party.Next()
		assert.NoError(t, err)
		assert.Equal(t, 2, int(*party.Value()))
		assert.Equal(t, 2, party.cur)
	}
	{
		// 2->0で戻れる
		err := party.Prev()
		assert.NoError(t, err)
		assert.Equal(t, 0, int(*party.Value()))
		assert.Equal(t, 0, party.cur)
	}
}
