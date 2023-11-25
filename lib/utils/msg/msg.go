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

type MsgBuilder struct {
	pages []page
	pos   int
}

type page struct {
	lines []line
	pos   int
}

// lineごとにクリック待ちが発生する。改行は入るかもしれない
// 表示はstr1つ1つ
type line struct {
	str string
	pos int
}

type chars []string

func New(raw string) *MsgBuilder {
	return &MsgBuilder{}
}
