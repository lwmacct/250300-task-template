package main

import (
	"context"
	"fmt"
	"os"

	"github.com/lwmacct/251128-workspace/internal/command/client"
)

func main() {
	if err := client.Command.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
