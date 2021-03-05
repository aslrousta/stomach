package main

import (
	"flag"
	"os"
)

func main() {
	flag.Parse()
	if err := dbOpen(); err != nil {
		panic(err)
	}
	defer dbClose()
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "plate":
		if err := plate(); err != nil {
			panic(err)
		}
	}
}
