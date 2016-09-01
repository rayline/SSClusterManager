package main

import (
	"github.com/dynport/gossh"
	"log"
)

// returns a function of type gossh.Writer func(...interface{})
// MakeLogger just adds a prefix (DEBUG, INFO, ERROR)
func MakeLogger(prefix string) gossh.Writer {
	return func(args ...interface{}) {
		log.Println((append([]interface{}{prefix}, args...))...)
	}
}

func main() {
	client := gossh.New("45.63.120.170", "root")
	// my default agent authentication is used. use
	client.SetPassword("Nx)9X{rT(AP=]GMC")
	// for password authentication
	client.DebugWriter = MakeLogger("DEBUG")
	client.InfoWriter = MakeLogger("INFO ")
	client.ErrorWriter = MakeLogger("ERROR")

	defer client.Close()

	rsp, e := client.Execute(`
		source ~/.bashrc;
		go --version`)
	if e != nil {
		client.ErrorWriter(e.Error())
		client.ErrorWriter("STDOUT: " + rsp.Stdout())
		client.ErrorWriter("STDERR: " + rsp.Stderr())
	}
}
