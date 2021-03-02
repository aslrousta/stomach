package main

import (
	"fmt"
	"os"

	"github.com/aslrousta/stomach/myers"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("use: stomach [before] [after]")
		os.Exit(1)
	}

	before, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	after, err := os.ReadFile(os.Args[2])
	if err != nil {
		panic(err)
	}
	edits := myers.Diff(string(before), string(after))
	for _, edit := range edits {
		var sign byte
		if edit.Type == myers.Delete {
			sign = '-'
		} else {
			sign = '+'
		}
		fmt.Printf("%c %3d %s\n", sign, edit.Index, edit.Line)
	}
}
