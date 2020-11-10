package main

import (
	"fmt"
	"os"

	"github.com/murpg/gsp/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(-1)
	}
}
