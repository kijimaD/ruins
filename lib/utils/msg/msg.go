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

func (q *queue) Exec() {
	if !q.active {
		return
	}
	q.events[0].PreHook()
	q.events[0].Run(q)
}

type event interface {
	PreHook() func()
	Run(*queue)
}

// ================

type msg struct {
	body []rune
	pos  int
	done bool
}

func (e *msg) PreHook() func() {
	return func() {}
}

func (e *msg) Run(q *queue) {
	if e.done {
		return
	}
	q.buf += string(e.body[e.pos])
	e.pos++
	if e.pos > len(e.body)-1 {
		e.done = true
	}
	return
}

// ================

type wait struct{}

func (e *wait) PreHook() func() {
	return func() {}
}

func (e *wait) Run(q queue) {
	q.active = false
	return
}

// ================

type resume struct{}

func (e *resume) PreHook() func() {
	return func() {}
}

func (e *resume) Run(q queue) {
	q.active = true
	return
}
