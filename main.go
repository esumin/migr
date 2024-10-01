package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mig/pkg/migrator"
)

func TraverseAndModifyFiles(root string, handlers migrator.MigrationHandlers) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" {
			fmt.Printf("Processing file: %s ...", path)
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer func() {
				_ = file.Close()
			}()

			scanner := bufio.NewScanner(file)

			result := strings.Builder{}
			saveFile := false

			for scanner.Scan() {
				line := scanner.Text()
				for _, handler := range handlers {
					modified := handler(line)
					if modified != "" {
						saveFile = true
						line = modified
						break
					}
				}
				result.WriteString(line + "\n")
			}

			if err := scanner.Err(); err != nil {
				return err
			}

			if !saveFile {
				fmt.Printf(" no changes\n")
			} else {
				fmt.Printf(" changed\n")
				// Now let's save builder content back to file
				err := os.WriteFile(path, []byte(result.String()), 0644)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %v: %v\n", root, err)
	}
}

func main() {
	matchers, err := migrator.GetMigratorHandlers(migrator.V1)
	if err != nil {
		fmt.Println(err)
		return
	}

	TraverseAndModifyFiles(
		"/Users/eugen.sumin/git/kanister/pkg/blockstorage",
		matchers,
	)
}
