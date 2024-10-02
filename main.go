package main

import (
	"fmt"
	"os"

	"mig/pkg/migrator"
	"mig/pkg/traverser"
)

func main() {
	// Check if a path is provided in the command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: myapp <path>")
		return
	}

	// Get the path from the command line arguments
	path := os.Args[1]

	matchers, err := migrator.GetMigratorHandlers(migrator.V2)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Pass the command line path to the TraverseAndModifyFiles function
	traverser.TraverseAndModifyFiles(
		path,
		matchers,
	)
}
