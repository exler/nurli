package main

import (
	"github.com/exler/nurli/internal/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
