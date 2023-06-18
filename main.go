package main

import (
	"github.com/exler/nurli/internal/cmd"
	_ "modernc.org/sqlite"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
