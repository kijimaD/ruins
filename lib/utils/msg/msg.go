package msg

// こんにちは。\n -- 改行
// 今日は晴れです。[p] -- 改ページクリック待ち
// ところで[l] -- 行末クリック待ち

type queueResult string

var (
	queueEmpty      = queueResult("EMPTY")
	queueProcessing = queueResult("PROCESSING")
	queueWait       = queueResult("WAIT")
)

type Queue struct {
	events []event
	buf    string
	// trueの場合キューを処理する
	active bool
}

func NewQueue(events []event) Queue {
	q := Queue{
		active: true,
		events: events,
	}
	return q
}

func (q *Queue) Exec() queueResult {
	if !q.active {
		return queueWait
	}
	if len(q.events) == 0 {
		return queueEmpty
	}
	q.events[0].PreHook()
	q.events[0].Run(q)

	return queueProcessing
}

// キューの先端を消して先に進める
func (q *Queue) Next() queueResult {
	q.events = append(q.events[:0], q.events[1:]...)
	q.active = true
	// for {
	// 	result := q.Exec()
	// 	if result == queueWait || result == queueEmpty {
	// 		break
	// 	}
	// }
	return queueWait
}

func (q *Queue) Display() string {
	return q.buf
}

func (q *Queue) SetEvents(es []event) {
	q.events = es
}

type event interface {
	PreHook()
	Run(*Queue)
}

// ================

// メッセージ表示
type msg struct {
	body []rune
	pos  int
}

func (e *msg) PreHook() {
	return
}

func (e *msg) Run(q *Queue) {
	q.buf += string(e.body[e.pos])
	e.pos++
	if e.pos > len(e.body)-1 {
		q.active = false
	}
	return
}

// ================

// キューを待ち状態にする
type wait struct{}

func (e *wait) PreHook() {
	return
}

func (e *wait) Run(q *Queue) {
	q.active = false
	return
}

// ================

// キューを実行可能状態にする
type resume struct{}

func (e *resume) PreHook() {
	return
}

func (e *resume) Run(q *Queue) {
	q.active = true
	return
}

// ================

// ページをフラッシュする
type flush struct{}

func (e *flush) PreHook() {
	return
}

func (e *flush) Run(q *Queue) {
	q.buf = ""
	q.active = false
	return
}
