package main

import (
	"fmt"
	"gocker/internal/commands"
	"os"
)


func main(){
	if len(os.Args)<2{
		fmt.Println("Usage:gocker <commands> [args...]")
		return
	}
	cmd := os.Args[1]
	args:= os.Args[2:]
	switch cmd{
	case "run":
		// add commands
		commands.RunContainer(args)
	case "ps":
		// add commands
		commands.ListContainer()
	case "logs":
		// add commands 
		commands.ShowLogs(args)
	case "stop":
		commands.StopContainer(args)
	case "rm":
		commands.RemoveContainer(args)	
	case "exec":
		commands.ExecContainer(args)
	
	default:
		fmt.Println("Unknown command",cmd)
	}
}