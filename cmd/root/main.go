package root

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
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

	return rootCmd, nil
}
