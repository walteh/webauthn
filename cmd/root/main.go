package root

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/walteh/snake"
	myversion "github.com/walteh/webauthn/version"
)

type Root struct {
	Quiet   bool
	Debug   bool
	Version bool
}

var _ snake.Snakeable = (*Root)(nil)

func (me *Root) BuildCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buildrc",
		Short: "buildrc is a tool to help with building releases",
	}

	cmd.PersistentFlags().BoolVarP(&me.Quiet, "quiet", "q", false, "Do not print any output")
	cmd.PersistentFlags().BoolVarP(&me.Debug, "debug", "d", false, "Print debug output")
	cmd.PersistentFlags().BoolVarP(&me.Version, "version", "v", false, "Print version and exit")

	cmd.SetOutput(os.Stdout)

	return cmd
}

func (me *Root) ParseArguments(ctx context.Context, cmd *cobra.Command, args []string) error {

	var level zerolog.Level
	if me.Debug {
		level = zerolog.TraceLevel
	} else if me.Quiet {
		level = zerolog.NoLevel
	} else {
		level = zerolog.InfoLevel
	}

	ctx = zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Caller().Logger().Level(level).WithContext(ctx)

	if me.Version {
		cmd.Printf("%s %s %s\n", myversion.Package, myversion.Version, myversion.Revision)
		os.Exit(0)
	}

	root := afero.NewOsFs()

	ctx = snake.Bind(ctx, (*afero.Fs)(nil), root)

	cmd.SetContext(ctx)

	return nil
}
