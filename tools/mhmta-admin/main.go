package main

import (
	"fmt"
	"os"
)

const usage = `Usage: mhmta-admin [command] args...

Commands:
    add-user    add a user to the authentication registry
`

func main() {
	if len(os.Args) >= 2 {
		cmd := os.Args[1]
		switch cmd {
		case "add-user":
		default:
			fmt.Printf("Unrecognised command '%s'\n", cmd)
			os.Exit(1)
		}
	}
	fmt.Printf("%s", usage)
}
