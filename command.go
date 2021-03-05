package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func plate() error {
	objects := make(map[string]fs.FileInfo)
	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err == nil {
			if !info.IsDir() && !mustIgnore(path) {
				objects[path] = info
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	changed, deleted, unstaged, err := dbCheckObjects(objects)
	if err != nil {
		return err
	}
	if len(changed)+len(deleted) > 0 {
		fmt.Println("indigested ingredients:")
		for _, ch := range changed {
			fmt.Printf("  %s\n", ch)
		}
		for _, d := range deleted {
			fmt.Printf("  %s (deleted)\n", d)
		}
	}
	if len(unstaged) > 0 {
		fmt.Println("raw ingredients:")
		for _, u := range unstaged {
			fmt.Printf("  %s\n", u)
		}
	}
	return nil
}

func mustIgnore(path string) bool {
	return strings.Index(path, ".stomach") == 0 ||
		strings.Index(path, ".git") == 0 ||
		strings.Index(path, "vendor") == 0
}
