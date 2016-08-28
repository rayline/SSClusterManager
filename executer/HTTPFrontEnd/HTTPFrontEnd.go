package HTTPFrontEnd

import "SSClusterManager/executer/util"
import "net/http"
import "log"
import "strconv"

func Start() {
	log.Print("Starting HTTP Service...")
	defer log.Println("Success")
	port, exist := util.Configs["HTTPPort"].(uint16)
	if exist == false {
		port = 41600
	}
	log.Print("On ", port, "...")
	err := http.ListenAndServe(":"+strconv.FormatUint(uint64(port), 10), securityMux)
	if err != nil {
		log.Fatalln("Failed to listen and serve ", err)
	}
}
