package util

import "time"
import "os"
import "log"

var LogFileName string
var Configs map[string]interface{}

const LogDirecotry = "log/"

func init() {
	LogFileName = time.Now().Format(time.RFC3339) + ".log"
	w, err := os.OpenFile(LogDirecotry+LogFileName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalln("Failed to create log file", err, "...Exiting")
	}
	log.SetOutput(w)
	Configs = LoadConfig("config.json") // the config file is loaded from the same directory with the binary
}
