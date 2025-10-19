package container

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Container struct {
	ID string
	command string
	CreatedAt time.Time
}

var containers []Container

func Run(args[]string){
	cmd:= exec.Command(args[0],args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	containerID := fmt.Sprintf("c-%d",time.Now().UnixNano())
	logFile := fmt.Sprintf("logs/%s",containerID)
	f,_:= os.Create(logFile)
	defer f.Close()
	cmd.Stdout = f
	cmd.Stderr = f


	fmt.Println("Running container: ",containerID)
	if err := cmd.Start(); err != nil{
		fmt.Println("Error: ",err)
		return
	}
	containers = append(containers, Container{
		ID: containerID,
		command: args[0],
		CreatedAt: time.Now(),
	})


	go func ()  {
		cmd.Wait()
		fmt.Println("Container exited: ",containerID)
	}()

}


func List(){
	fmt.Println("Running Containers: ")
	for _,c:= range containers{
		fmt.Printf("%s %s %v\n",c.ID,c.command,c.CreatedAt.Format(time.RFC3339))
	}
}

func Logs(containerID string){
	filePath := fmt.Sprintf("Logs /%s.log",containerID)
	data,err := os.ReadFile(filePath)
	if err != nil{
		fmt.Println("No logs found for: ",containerID)
		return
	}
	fmt.Println(string(data))
}