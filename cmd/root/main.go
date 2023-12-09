package root

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/walteh/snake"
	myversion "github.com/walteh/webauthn/version"
)

func Build(ctx context.Context) (*cobra.Command, error) {

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "webauthn",
		Short: "Webauthn server",
		Long:  "Webauthn server",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	rootCmd.AddCommand(myversion.Build(ctx))
	rootCmd.AddCommand(snake.Build(ctx))

	return rootCmd, nil
}
