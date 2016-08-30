package main

import "SSClusterManager/director/util"
import "log"
import "net/http"

func main() {
	log.Println("Director Started")

	//start the file server so we can have a friendlier user interface
	log.Fatalln(http.ListenAndServe(":"+util.Configs["staticFileServer"].(string), http.FileServer(http.Dir(util.Configs["staticFilePath"].(string)))))

	//start the API server

}
