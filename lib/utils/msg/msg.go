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

type Msg struct {
	pages []page
	pos   int // pageのpos
	buf   string
}

type submitResult string

const (
	SubmitOK         = submitResult("OK")
	SubmitPageOK     = submitResult("OK_Page")
	SubmitLineOK     = submitResult("OK_LINE")
	SubmitMsgFinish  = submitResult("MSG_FINISH")
	SubmitPageFinish = submitResult("PAGE_FINISH")
	SubmitLineFinish = submitResult("LINE_FINISH")
	SubmitFail       = submitResult("FAIL")
)

type submitOpt struct {
	pageNext bool
	lineNext bool
}

func (m *Msg) Buf(opt submitOpt) string {
	str, _ := m.submit(opt)
	// submit statusがlineのとき、\nを入れる
	//                pageのとき、bufを消す
	//                msgのとき、終了ステータスを返す
	// 送りオプションがついてるときだけ、次に進めるようにする
	m.buf = m.buf + str
	return m.buf
}

func (m *Msg) submit(opt submitOpt) (string, submitResult) {
	if m.pos > len(m.pages)-1 {
		return "", SubmitMsgFinish
	}

	str, result := m.pages[m.pos].submit(opt)
	if result == SubmitLineOK {
		return str, result
	}
	if result == SubmitPageFinish && opt.pageNext {
		m.pos++
		str, result = m.submit(opt)
		return str, result
	}
	return "", SubmitFail
}

type page struct {
	lines []line
	pos   int // linesのpos
}

func (p *page) submit(opt submitOpt) (string, submitResult) {
	if p.pos > len(p.lines)-1 {
		return "", SubmitPageFinish
	}

	str, result := p.lines[p.pos].submit()
	if result == SubmitLineOK {
		return str, result
	}
	if result == SubmitLineFinish && (opt.lineNext || opt.pageNext) {
		p.pos++
		str, result = p.submit(opt)
		return str, result
	}
	return "", SubmitFail
}

// lineごとにクリック待ちが発生する。改行は入るかもしれない
// 表示はstr1つ1つ
type line struct {
	str string
	pos int // strのpos
}

func (l *line) submit() (string, submitResult) {
	if l.pos > len(l.str)-1 {
		return "", SubmitLineFinish
	}
	str := l.str[l.pos]
	l.pos++
	return string(str), SubmitLineOK
}
