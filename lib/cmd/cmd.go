package cmd

import (
	"fmt"

	"github.com/kijimaD/ruins/lib/consts"
	"github.com/urfave/cli/v2"
)

const splash = `───────────────────────────────────────────────────────
██████  ██    ██ ██ ███    ██ ███████
██   ██ ██    ██ ██ ████   ██ ██
██████  ██    ██ ██ ██ ██  ██ ███████
██   ██ ██    ██ ██ ██  ██ ██      ██
██   ██  ██████  ██ ██   ████ ███████
───────────────────────────────────────────────────────
`

// NewMainApp は新しいメインアプリケーションを作成する
func NewMainApp() *cli.App {
	app := cli.NewApp()
	app.Name = "ruins"
	app.Usage = "ruins [subcommand] [args]"
	app.Description = "This is RPG!"
	app.DefaultCommand = CmdPlay.Name
	app.Version = consts.AppVersion
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		CmdPlay,
		CmdScreenshot,
	}
	cli.AppHelpTemplate = fmt.Sprintf(`%s
%s
`, splash, cli.AppHelpTemplate)

	return app
}

// RunMainApp はメインアプリケーションを実行する
func RunMainApp(app *cli.App, args ...string) error {
	err := app.Run(args)
	if err != nil {
		return fmt.Errorf("コマンド実行が失敗した: %w", err)
	}

	return nil
}
