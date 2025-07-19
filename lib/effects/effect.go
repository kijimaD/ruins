package effects

// 型エイリアスは不要（同じパッケージ内なので直接参照）

// Logger はエフェクト処理のログ出力インターフェース
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// defaultLogger はデフォルトのログ実装（何も出力しない）
type defaultLogger struct{}

func (d defaultLogger) Info(msg string, args ...interface{})  {}
func (d defaultLogger) Error(msg string, args ...interface{}) {}
func (d defaultLogger) Debug(msg string, args ...interface{}) {}
