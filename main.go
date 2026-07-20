package main

import (
	"os"

	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/cmd"
	"github.com/sharanyashwant27-tech/GitPulse-CLI-GO/internal/utils"
)

func main() {
	defer utils.Sync()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
