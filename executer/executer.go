package main

import "log"
import _ "SSClusterManager/executer/util"
import "SSClusterManager/executer/HTTPFrontEnd"

func main() {
	log.Println("Executer started")
	HTTPFrontEnd.Start()
}
