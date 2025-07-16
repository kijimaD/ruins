package msg

import "time"

// QueueState はキューの状態を表す型
type QueueState string

const (
	// QueueStateNone はキューが非アクティブ状態
	QueueStateNone = QueueState("NONE")
	// QueueStateFinish はキューが完了状態
	QueueStateFinish = QueueState("FINISH")
)

// Queue はイベントキューを管理する構造体
type Queue struct {
	events []Event
	// 現在の表示文字列
	buf string
	// trueの場合キューを処理する
	active bool
}

// NewQueue は新しいQueueを作成する
func NewQueue(events []Event) Queue {
	q := Queue{
		active: true,
		events: events,
	}
	return q
}

// NewQueueFromText はテキストからQueueを作成する
func NewQueueFromText(text string) Queue {
	l := NewLexer(text)
	p := NewParser(l)
	program := p.ParseProgram()
	e := NewEvaluator(program)

	return NewQueue(e.Events)
}

// RunHead はキューの先端にあるイベントを実行する
func (q *Queue) RunHead() QueueState {
	if !q.active {
		return QueueStateNone
	}
	if len(q.events) == 0 {
		return QueueStateFinish
	}
	q.events[0].PreHook()
	q.events[0].Run(q)

	return QueueStateNone
}

// Head はキューの先端イベントを返す
func (q *Queue) Head() Event {
	if len(q.events) == 0 {
		return &notImplement{}
	}
	return q.events[0]
}

// Pop はキューの先端を消して先に進める
func (q *Queue) Pop() QueueState {
	if len(q.events) == 0 {
		return QueueStateFinish
	}
	q.events = append(q.events[:0], q.events[1:]...)
	q.activate()
	return QueueStateNone
}

// Display は現在の表示文字列を返す
func (q *Queue) Display() string {
	return q.buf
}

// SetEvents はイベントを設定する
func (q *Queue) SetEvents(es []Event) {
	q.events = es
}

func (q *Queue) activate() {
	q.active = true
}

func (q *Queue) deactivate() {
	q.active = false
}

// Event はイベントのインターフェース
type Event interface {
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

func (e *msgEmit) PreHook() {}

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
}

// ================

// ページをフラッシュする
type flush struct{}

func (e *flush) PreHook() {}

func (e *flush) Run(q *Queue) {
	q.buf = ""
	q.deactivate()
	q.Pop()
}

// ================

// ChangeBg は背景変更イベント
type ChangeBg struct {
	Source string
}

// PreHook はイベント実行前の処理
func (c *ChangeBg) PreHook() {}

// Run はイベントを実行する
func (c *ChangeBg) Run(q *Queue) {
	q.Pop()
}

// ================

// 行末クリック待ち
type lineEndWait struct{}

func (l *lineEndWait) PreHook() {}

func (l *lineEndWait) Run(q *Queue) {
	q.buf = q.buf + "\n"
	q.deactivate()
	q.Pop()
}

// ================

// 未実装
type notImplement struct{}

func (l *notImplement) PreHook() {}

func (l *notImplement) Run(q *Queue) {
	q.buf = ""
	q.deactivate()
	q.Pop()
}

// ================
type wait struct {
	duration time.Duration
}

func (w *wait) PreHook() {}

func (w *wait) Run(q *Queue) {
	time.Sleep(w.duration)
	q.buf = ""
	q.deactivate()
	q.Pop()
}
