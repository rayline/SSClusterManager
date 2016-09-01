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
	installScriptBytes, err := ioutil.ReadFile("scripts/installExecuterRemote.sh")
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
	log.Println("Wating for server to be active...")
	for state := s.Status(); state <= StateNotReady; state = s.Status() {
		if state == StateDestroyed {
			return
		}
		time.Sleep(time.Second * 2)
	}
	sshClient := gossh.New(s.Addr().String(), "root")
	defer sshClient.Close()
	sshClient.SetPassword(s.RootPassword())
	log.Println("Waiting For remote SSH to be ready...")
	for _, err := sshClient.Execute("ls"); err != nil; _, err = sshClient.Execute("ls") {
		time.Sleep(time.Second * 2)
	}
	log.Println("Remote ready for configuration, Using SSH to setup server")
	o, err := sshClient.Execute(installScript)
	log.Println(o.Stderr())
	log.Println(o.Stdout())
	if err != nil {
		log.Println("Error executing ssh ", err, "With output :stderr: ", o.Stderr(), " , stdout: ", o.Stdout())
	}
	log.Println("Waiting for installation of executer to complete...")
	for !Available(s) {
		time.Sleep(time.Second * 2)
	}
	log.Println("done setup on server on ", s.Addr().String())
}

func WriteUser(s Server, password string, port string) {
	s.Do("/writeuser", []byte(url.Values{
		"password": []string{password},
		"port":     []string{port},
	}.Encode()))
}
