package main

//the selected server provider should be referred here by main to register itself in package servers in its init function

import _ "SSClusterManager/director-lite/util"
import _ "SSClusterManager/director-lite/serverProviderVultr"
import "SSClusterManager/director-lite/servers"
import "SSClusterManager/director-lite/server"
import "io/ioutil"
import "os"
import "flag"

var count int

func init() {
	flag.IntVar(&count, "c", 1, "the count of servers required")
}

func main() {
	flag.Parse()

	ch := make(chan int, count)

	for i := 1; i <= count; i++ {
		go func() {
			s := servers.Provider.Create()
			servers.AddServer(s)
			server.Setup(s)
			ch <- 1
		}()
	}

	for i := 1; i <= count; i++ {
		<-ch
	}

	ioutil.WriteFile("servers.json", servers.JSON(), os.ModePerm)

}
