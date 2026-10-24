// Command localcloud is the entry point of the localcloud CLI.
package main

import (
	"os"

	"github.com/alzasr/localcloud/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
