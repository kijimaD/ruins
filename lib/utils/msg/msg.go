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

func (m *Msg) Buf() string {
	str, _ := m.submit()
	// submit statusがlineのとき、\nを入れる
	//                pageのとき、bufを消す
	//                msgのとき、終了ステータスを返す
	// 送りオプションがついてるときだけ、次に進めるようにする
	m.buf = m.buf + str
	return m.buf
}

func (m *Msg) submit() (string, bool) {
	if m.pos > len(m.pages)-1 {
		return "", false
	}

	str, ok := m.pages[m.pos].submit()
	if !ok {
		m.pos++
		str, ok := m.submit()
		return str, ok
	}
	return str, true
}

type page struct {
	lines []line
	pos   int // linesのpos
}

func (p *page) submit() (string, bool) {
	if p.pos > len(p.lines)-1 {
		return "", false
	}
	str, ok := p.lines[p.pos].submit()
	if !ok {
		p.pos++
		str, ok := p.submit()
		return str, ok
	}
	return str, true
}

// lineごとにクリック待ちが発生する。改行は入るかもしれない
// 表示はstr1つ1つ
type line struct {
	str string
	pos int // strのpos
}

func (l *line) submit() (string, bool) {
	// fmt.Println(l.pos, len(l.str)-1)
	if l.pos > len(l.str)-1 {
		return "", false
	}
	str := l.str[l.pos]
	l.pos++
	return string(str), true
}
