package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/walteh/snake"
	"github.com/walteh/webauthn/cmd/root"
)

func init() {

	zerolog.TimeFieldFormat = time.RFC3339Nano

}

func main() {

	ctx := context.Background()

	rootCmd := snake.NewRootCommand(ctx, &root.Root{})

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		_, err = fmt.Fprintf(os.Stderr, "[%s] (error) %+v\n", rootCmd.Name(), err)
		if err != nil {
			panic(err)
		}
	}

}
