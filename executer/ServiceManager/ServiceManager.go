package ServiceManager

import "os/exec"
import "SSClusterManager/executer/util"
import "log"
import "strconv"

//The package is responsible for starting and stopping the service with the help of external script
//The main identifier of service of different user is only the port number

//guide on scripting
//[startScript] [port] [password]
//[stopScript] [port]
//The scripts should make sure it can run multiple instances of the service
//The script name must be defined in config file

//scripts should be stored in folder "scripts"

//At this version the program will not try to handle any error from the scripts

func StartService(port uint16, password string) {
	log.Print("Starting service on ", port, "...")
	scriptNameInterface, exist := util.Configs["startScript"]
	scriptName := ""
	if exist == false {
		scriptName = "startscript.sh"
	} else {
		scriptName = scriptNameInterface.(string)
	}
	cmd := exec.Command("scripts/"+scriptName, strconv.FormatUint(uint64(port), 10), password)
	err := cmd.Run()
	if err != nil {
		log.Println("Failed:", err)
	}
	log.Println("Successs")
}

func StopService(port uint16) {
	log.Print("Stopping service on ", port, "...")
	scriptNameInterface, exist := util.Configs["stopScript"]
	scriptName := ""
	if exist == false {
		scriptName = "stopscript.sh"
	} else {
		scriptName = scriptNameInterface.(string)
	}
	cmd := exec.Command("scripts/"+scriptName, strconv.FormatUint(uint64(port), 10))
	err := cmd.Run()
	if err != nil {
		log.Println("Failed:", err)
	}
	log.Println("Successs")
}
