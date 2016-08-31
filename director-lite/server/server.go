package server

//the package create, destroy, setup and directs the servers
//every server will have a gorountine to process the directions

import "net"
import "net/http"
import "net/url"
import "github.com/dynport/gossh"
import "io/ioutil"
import "log"
import "bytes"
import "time"
import "strings"

const StateDestroyed = -1
const StateNotReady = 0
const StateOK = 1

var installScript string

func init() {
	//cache the executer install script here
	installScriptBytes, err := ioutil.ReadFile("scripts/installExecuter.sh")
	if err != nil {
		log.Fatalln("Error Reading executer install script", err)
	}
	b := bytes.NewBuffer(installScriptBytes)
	installScript = b.String()

	//all \r must removed
	replacer := strings.NewReplacer("\r", "")
	installScript = replacer.Replace(installScript)

}

type Server interface {
	Addr() net.IP
	AddrV6() net.IP
	RootPassword() string
	Status() int
	Do(uri string, data []byte) //There may be bugs here in future if extremely large amount of server and users appear, but not for now
	Reboot()
	Destroy()
}

func Available(s Server) bool {
	// availability check
	//TODO: more checking items
	if s.Status() != StateOK {
		return false
	}
	if _, err := http.Get("http://" + s.Addr().String() + "/"); err != nil {
		return false
	}
	return true
}

func Setup(s Server) {
	for state := s.Status(); state <= StateNotReady; {
		if state == StateDestroyed {
			return
		}
		time.Sleep(time.Second)
	}
	for !Available(s) {
		log.Println("Using SSH to setup server")
		sshClient := gossh.New(s.Addr().String(), "root")
		defer sshClient.Close()
		sshClient.SetPassword(s.RootPassword())
		o, err := sshClient.Execute(installScript)
		if err != nil {
			log.Println("Error executing ssh ", err, "With output :stderr: ", o.Stderr(), " , stdout: ", o.Stdout())
		}
	}
	log.Println("done setup on server on ", s.Addr().String())
}

func WriteUser(s Server, password string, port string) {
	s.Do("/writeuser", []byte(url.Values{
		"password": []string{password},
		"port":     []string{port},
	}.Encode()))
}
