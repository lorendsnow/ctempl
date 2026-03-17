package main

import (
	"fmt"
	"os"

	"github.com/lorendsnow/ctempl/cmd/cproject"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(topLevelHelp)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "c":
		cmd := cproject.NewCProject(os.Args[2:])
		if err := cmd.Run(); err != nil {
			fmt.Println("error occurred while trying to set up C project:", err.Error())
			os.Exit(1)
		}
	case "cxx":
		fmt.Println("to be implemented...")
	default:
		fmt.Println("Valid commands are 'c' and 'cxx'")
		os.Exit(1)
	}
}
