package commands

import (
	"fmt"
	"gocker/internal/container"
	"strings"
)

func RunContainer(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: gocontainer run [--env=KEY=VAL] <command> [args...]")
		return
	}

	envVars := []string{}
	cmdArgs := []string{}

	// Separate env vars from command
	for _, a := range args {
		if strings.HasPrefix(a, "--env=") {
			envVars = append(envVars, strings.TrimPrefix(a, "--env="))
		} else {
			cmdArgs = append(cmdArgs, a)
		}
	}

	if len(cmdArgs) == 0 {
		fmt.Println("No command specified")
		return
	}

	container.RunWithEnv(cmdArgs, envVars)
}

func ListContainer() {
	container.List()
}
func StopContainer(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: gocker stop <container_id>")
		return
	}
	container.Stop(args[0])
}

func RemoveContainer(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: gocker rm <container_id>")
		return
	}
	container.Remove(args[0])
}

func ShowLogs(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: gocker logs <container_id> [-f]")
		return
	}
	if len(args) > 1 && args[1] == "-f" {
		container.LogsFollow(args[0])
	} else {
		container.Logs(args[0])
	}
}
func ExecContainer(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: gocker exec <container_id> <command> [args...]")
		return
	}

	containerID := args[0]
	cmdArgs := args[1:]

	container.Exec(containerID, cmdArgs)
}
