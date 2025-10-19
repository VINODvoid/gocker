package container

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Container struct {
	ID        string
	command   string
	CreatedAt time.Time
}

var containers []Container

func init(){
	loadState()
}
func Run(args []string) {
	loadState()

	cmd := exec.Command(args[0], args[1:]...)
	containerID := fmt.Sprintf("c-%d", time.Now().UnixNano())

	workDir := fmt.Sprintf("workspace/%s", containerID)
	os.MkdirAll(workDir, 0755)
	cmd.Dir = workDir

	logFile := fmt.Sprintf("logs/%s.log", containerID)
	f, _ := os.Create(logFile)
	defer f.Close()
	cmd.Stdout = f
	cmd.Stderr = f

	fmt.Println("Running container: ", containerID)
	if err := cmd.Start(); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	containers = append(containers, Container{
		ID:        containerID,
		command:   args[0],
		CreatedAt: time.Now(),
	})
	saveState()

	go func() {
		cmd.Wait()
		fmt.Println("Container exited: ", containerID)
	}()

}

func List() {
	fmt.Println("Running Containers: ")
	for _, c := range containers {
		fmt.Printf("%s %s %v\n", c.ID, c.command, c.CreatedAt.Format(time.RFC3339))
	}
}

func Logs(containerID string) {
	filePath := fmt.Sprintf("Logs /%s.log", containerID)
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("No logs found for: ", containerID)
		return
	}
	fmt.Println(string(data))
}

func Stop(containerID string) {
	fmt.Println("Stopping container:", containerID)
	for i, c := range containers {
		if c.ID == containerID {
			fmt.Println("Container stopped:", c.ID)
			containers = append(containers[:i], containers[i+1:]...)
			saveState()
			return
		}
	}
	fmt.Println("Container not found.")
}

func Remove(containerID string) {
	fmt.Println("Removing container:", containerID)
	logFile := fmt.Sprintf("logs/%s.log", containerID)
	os.Remove(logFile)

	for i, c := range containers {
		if c.ID == containerID {
			containers = append(containers[:i], containers[i+1:]...)
			saveState()
			return
		}
	}
	fmt.Println("Container not found.")
}
