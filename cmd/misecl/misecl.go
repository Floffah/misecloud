package main

import (
	"github.com/floffah/misecloud/internal/app/misecl"
)

func main() {
	misecl.InitCli()

	err := misecl.RootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
