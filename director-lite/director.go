package main

import _ "SSClusterManager/director-lite/util"
import _ "SSClusterManager/director-lite/serverProviderVultr"
import "SSClusterManager/director-lite/servers"
import "log"
import "io/ioutil"
import "os"
import "time"

func main() {
	log.Println("Director Started")
	servers.CheckOutServerCount()

	ioutil.WriteFile("servers.json", servers.JSON(), os.ModePerm)

	for {
		time.Sleep(time.Hour)
	}
}
