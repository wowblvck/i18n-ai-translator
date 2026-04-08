package main

import (
	"os"

	"github.com/wowblvck/i18n-translator/internal/cli"
)

func main() {
	os.Exit(cli.Execute(os.Args[0], os.Args[1:]))
}
