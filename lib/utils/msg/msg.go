package msg

// 指定数で、自動で改行する
// 取り出す関数を実行するごとに1文字ずつ取り出す。これによってアニメーションっぽくする
// 下にはみ出るサイズの文字列が入れられるとエラーを出す
// 改行と改ページが意図的に行える。自動でも行える
// 改ページせずに入らなくなった場合は上にスクロールする

// こんにちは。\n -- 改行
// 今日は晴れです。[p] -- 改ページクリック待ち
// ところで[l] -- 行末クリック待ち

// 改ページは押したときに現在の表示文字をフラッシュして、上から表示にすること。

// 階層構造だとやりにくいし、ほかの構造に対応しにくい気がしてきた。テキスト表示だけじゃなくって、BGMなんかもあるんだぞ
// イベントキュー形式はどうだろうか?

type queueResult string

var (
	queueEmpty      = queueResult("EMPTY")
	queueProcessing = queueResult("PROCESSING")
	queueWait       = queueResult("WAIT")
)

type queue struct {
	events []event
	now    event
	buf    string
	active bool
}

func NewQueue() queue {
	return queue{
		active: true,
	}
}

func (q *queue) Exec() queueResult {
	for {
		result := q.exec()
		if result == queueWait || result == queueEmpty {
			break
		}
	}
	return queueWait
}

func (q *queue) exec() queueResult {
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

func (q *queue) Next() {
	q.events = append(q.events[:0], q.events[1:]...)
}

type event interface {
	PreHook()
	Run(*queue)
}

// ================

type msg struct {
	body []rune
	pos  int
}

func (e *msg) PreHook() {
	return
}

func (e *msg) Run(q *queue) {
	q.buf += string(e.body[e.pos])
	e.pos++
	if e.pos > len(e.body)-1 {
		q.active = false
	}
	return
}

// ================

type wait struct{}

func (e *wait) PreHook() {
	return
}

func (e *wait) Run(q *queue) {
	q.active = false
	return
}

// ================

type resume struct{}

func (e *resume) PreHook() {
	return
}

func (e *resume) Run(q *queue) {
	q.active = true
	return
}

// ================

type flush struct{}

func (e *flush) PreHook() {
	return
}

func (e *flush) Run(q *queue) {
	q.buf = ""
	return
}
