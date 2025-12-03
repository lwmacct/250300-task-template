package main

import (
	"context"
	"fmt"
	"os"

	"github.com/lwmacct/251128-workspace/internal/command/api"
)

func main() {
	if err := api.Command.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
