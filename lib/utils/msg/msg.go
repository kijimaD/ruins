package msg

import (
	"fmt"
	"log"
)

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
	// 行を正しく返した
	SubmitLineOK = submitResult("OK_LINE")
	// 終端に到達した
	SubmitMsgFinish  = submitResult("MSG_FINISH")
	SubmitPageFinish = submitResult("PAGE_FINISH")
	SubmitLineFinish = submitResult("LINE_FINISH")
	SubmitFail       = submitResult("FAIL")
	// 予期せぬエラー
	SubmitUnexpectFail = submitResult("UNEXPECT_FAIL")
)

type submitStrategy string

const (
	// ページ送り待ちする
	submitPageNext = submitStrategy("PAGE")
	// 改行待ちする
	submitLineNext = submitStrategy("LINE")
	submitStop     = submitStrategy("STOP")
)

func (m *Msg) Buf(opt submitStrategy) string {
	str, result := m.submit(opt)

	switch result {
	case SubmitMsgFinish:
		fmt.Println("終了。要状態遷移")
	case SubmitPageFinish:
		m.buf = ""
	case SubmitLineFinish:
		m.buf += "\n"
	case SubmitLineOK:
	default:
		log.Fatalf("予期しないエラー: %s", result)
	}
	m.buf = m.buf + str
	return m.buf
}

func (m *Msg) submit(opt submitStrategy) (string, submitResult) {
	if m.pos > len(m.pages)-1 {
		return "", SubmitMsgFinish
	}

	str, result := m.pages[m.pos].submit(opt)
	switch result {
	case SubmitLineOK:
		return str, result
	case SubmitPageFinish:
		if opt == submitPageNext {
			m.pos++
			str, result = m.submit(opt)
			return str, result
		} else {
			return "", SubmitPageFinish
		}
	case SubmitLineFinish:
		return "", SubmitLineFinish
	}
	return "", SubmitUnexpectFail
}

type page struct {
	lines []line
	pos   int // linesのpos
}

func (p *page) submit(opt submitStrategy) (string, submitResult) {
	if p.pos > len(p.lines)-1 {
		return "", SubmitPageFinish
	}

	str, result := p.lines[p.pos].submit()
	switch result {
	case SubmitLineOK:
		return str, result
	case SubmitLineFinish:
		if opt == submitLineNext || opt == submitPageNext {
			p.pos++
			str, result = p.submit(opt)
			return str, result
		} else {
			return "", SubmitLineFinish
		}
	}
	return "", SubmitUnexpectFail
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
