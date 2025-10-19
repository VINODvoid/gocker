package commands

import (
	"fmt"
	"gocker/internal/container"
)



func RunContainer(args[] string){

	if len(args) == 0{
		fmt.Println("Usage: gocker run <command>")
		return
	}
	
	container.Run(args)
}
func ListContainer(){
	container.List()
}
func ShowLogs(args[] string){
	if len(args) == 0{
		fmt.Println("Usage: gocker logs <command_id>")
		return
	}
	container.Logs(args[0])

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
