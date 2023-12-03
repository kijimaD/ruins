package msg

// こんにちは。\n -- 改行
// 今日は晴れです。[p] -- 改ページクリック待ち
// ところで[l] -- 行末クリック待ち

type QueueState string

const (
	QueueStateNone   = QueueState("NONE")
	QueueStateFinish = QueueState("FINISH")
)

type Queue struct {
	events []event
	// 現在の表示文字列
	buf string
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

// キューの先端にあるイベントを実行する
func (q *Queue) RunHead() {
	if !q.active {
		return
	}
	if len(q.events) == 0 {
		return
	}
	q.events[0].PreHook()
	q.events[0].Run(q)

	return
}

// キューの先端を消して先に進める
func (q *Queue) Pop() QueueState {
	if len(q.events) == 0 {
		return QueueStateFinish
	}
	q.events = append(q.events[:0], q.events[1:]...)
	q.activate()
	return QueueStateNone
}

func (q *Queue) Display() string {
	return q.buf
}

func (q *Queue) SetEvents(es []event) {
	q.events = es
}

func (q *Queue) activate() {
	q.active = true
}

func (q *Queue) deactivate() {
	q.active = false
}

type event interface {
	PreHook()
	Run(*Queue)
}

// ================

// メッセージ表示
type msgEmit struct {
	body []rune
	pos  int
	// 自動改行カウント
	nlCount int
}

func (e *msgEmit) PreHook() {
	return
}

// 1つ位置を進めて1文字得る
func (e *msgEmit) Run(q *Queue) {
	const width = 14

	q.buf += string(e.body[e.pos])
	// 意図的に挿入された改行がある場合はリセット
	if string(e.body[e.pos]) == "\n" {
		e.nlCount = 0
	}
	if e.nlCount%width == width-1 {
		q.buf += "\n"
	}

	e.pos++
	e.nlCount++

	if e.pos > len(e.body)-1 {
		q.deactivate()
	}
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
	q.deactivate()
	q.Pop()
	return
}
