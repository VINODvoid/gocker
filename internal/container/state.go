package container

import (
	"encoding/json"
	"os"
)

const stateFile = "container_state.json"

func saveState(){
	data,_:= json.MarshalIndent(containers,""," ")
	os.WriteFile(stateFile,data,0644)
}
func loadState(){
	data,err:= os.ReadFile(stateFile)
	if err!=nil{
		return
	}
	json.Unmarshal(data,&containers)
}

