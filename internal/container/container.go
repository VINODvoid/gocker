package container

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Container struct {
	ID      string
	command string
	PID     int
	Env     []string
	CreatedAt time.Time
}

var containers []Container

func init() {
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
		command:   strings.Join(args," "),
		PID: cmd.Process.Pid,
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

func LogsFollow(containerID string) {
	filePath := fmt.Sprintf("logs/%s.log", containerID)
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("No logs found for:", containerID)
		return
	}
	defer f.Close()

	fmt.Printf("Following logs for %s...\n", containerID)
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		fmt.Print(line)
	}
}

func RunWithEnv(args []string, envVars []string) {
	loadState()

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = append(os.Environ(), envVars...) // Add env vars

	// Make sure we capture stdout/stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	containerID := fmt.Sprintf("c-%d", time.Now().UnixNano())
	logFile := fmt.Sprintf("logs/%s.log", containerID)
	f, err := os.Create(logFile)
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}
	defer f.Close()


	mw := io.MultiWriter(os.Stdout, f)
	cmd.Stdout = mw
	cmd.Stderr = mw

	fmt.Println("Running container:", containerID)
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	containers = append(containers, Container{
		ID:        containerID,
		command:   strings.Join(args, " "),
		PID:       cmd.Process.Pid,
		Env:       envVars,
		CreatedAt: time.Now(),
	})
	saveState()

	go func() {
		cmd.Wait()
		fmt.Println("Container exited:", containerID)
	}()
}



func Exec(containerID string, args []string) {
	for _, c := range containers {
		if c.ID == containerID {
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Dir = fmt.Sprintf("workspace/%s", c.ID)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Env = append(os.Environ(), c.Env...)
			if err := cmd.Run(); err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
	}
	fmt.Println("Container not found:", containerID)
}
