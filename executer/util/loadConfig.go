package util

//this file loads config from the config file

import "os"
import "io/ioutil"
import "log"
import "encoding/json"

func LoadConfig(filename string) map[string]interface{} {
	log.Print("Loading config file", filename, " ...")
	defer log.Println("Success")
	file, err := ioutil.ReadFile(filename)
	if err == os.ErrNotExist {
		return map[string]interface{}{}
	}
	if err != nil {
		log.Fatalln("Error reading config", err)
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(file, &result)
	if err != nil {
		log.Fatalln("Error decoding config", err)
	}
	return result
}
