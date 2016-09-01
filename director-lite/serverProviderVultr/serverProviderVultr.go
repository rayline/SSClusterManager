package serverProviderVultr

//this is the server provider driver for vultr

//in this version it will be the simplest version that follows everything from the config file and the user will really need to edit the config file carefully and manually

import "SSClusterManager/director-lite/server"
import "SSClusterManager/director-lite/servers"
import "SSClusterManager/director-lite/util"
import "net/http"
import "net/url"
import "log"
import "io/ioutil"
import "encoding/json"
import "net"
import "time"
import "bytes"
import "os"
import "github.com/dynport/gossh"

var configFile = "vultr_config.json"
var configs map[string]interface{}
var apikey string

func init() {
	configs = util.LoadConfig(configFile)
	apikey = configs["apikey"].(string)
	serverProvider := &VultrServers{}
	servers.RegisterProvider(serverProvider)

	respData := requestToVultr("/v1/server/list?tag="+configs["creation"].(map[string]interface{})["tag"].(string), "GET", nil)
	var serverList = map[string]interface{}{}
	json.Unmarshal(respData, &serverList)
	for index, data := range serverList {
		serverInfo := data.(map[string]interface{})
		s := &VultrServer{
			Subid:    index,
			IPv4:     net.ParseIP(serverInfo["main_ip"].(string)),
			IPv6:     net.ParseIP(serverInfo["v6_main_ip"].(string)),
			Password: serverInfo["default_password"].(string),
		}
		servers.AddServer(s)
	}
}

type VultrServers struct{}

type VultrServer struct {
	Subid    string
	IPv4     net.IP
	IPv6     net.IP
	Password string
}

func (v *VultrServer) Addr() net.IP {
	return v.IPv4
}

func (v *VultrServer) AddrV6() net.IP {
	return v.IPv6
}

func (v *VultrServer) RootPassword() string {
	return v.Password
}

func (v *VultrServer) Reboot() {
	form := url.Values{
		"SUBID": []string{v.Subid},
	}
	requestToVultr("/v1/server/reboot", "POST", []byte(form.Encode()))
}

func (v *VultrServer) Destroy() {
	form := url.Values{
		"SUBID": []string{v.Subid},
	}
	requestToVultr("/v1/server/destroy", "POST", []byte(form.Encode()))
}

func (v *VultrServer) Do(uri string, data []byte) {
	var method string
	if data == nil {
		method = "GET"
	} else {
		method = "POST"
	}
	req, _ := http.NewRequest(method, "http://"+v.IPv4.String()+uri, bytes.NewBuffer(data))
	req.Header.Add("key", util.Configs["SecurityKey"].(string))
	client := &http.Client{}
	for _, err := client.Do(req); err != nil; _, err = client.Do(req) {
		log.Println("Error accessing ", v.IPv4.String(), " ", err)
		time.Sleep(time.Second)
	}
}

func (v *VultrServer) Status() int {
	//log.Println("Checking Server status")
	req, err := http.NewRequest("GET", "https://api.vultr.com/v1/server/list?SUBID="+v.Subid, nil)
	req.Header.Add("API-Key", apikey)
	client := &http.Client{}
	var resp *http.Response
	for resp, err = client.Do(req); err != nil || resp.StatusCode == 503; resp, err = client.Do(req) {
		if err != nil {
			log.Println("Error access Vultr for new server : ", err)
		} else {
			log.Println("Accessing too fase...")
		}
		time.Sleep(time.Second)
	}
	if resp.StatusCode != 200 {
		return server.StateDestroyed
	}
	respData, err := ioutil.ReadAll(resp.Body)
	var serverList = map[string]interface{}{}
	json.Unmarshal(respData, &serverList)
	//log.Println("Got status ", serverList["server_state"])
	if serverList["server_state"] == "ok" {
		return server.StateOK
	} else {
		return server.StateNotReady
	}
}

func (v *VultrServers) Create() server.Server {
	var s server.Server
	for goodServerCreated := false; goodServerCreated == false; {
		log.Println("Creating Vultr Server...")
		creationConfig := configs["creation"].(map[string]interface{})
		var form = url.Values{}
		for index, value := range creationConfig {
			form.Add(index, value.(string))
		}
		respData := requestToVultr("/v1/server/create", "POST", []byte(form.Encode()))
		var respFrom = map[string]string{}
		json.Unmarshal(respData, &respFrom)
		subid := respFrom["SUBID"]
		var serverList = map[string]interface{}{}
		var server *VultrServer
		log.Println("Waiting for basic setup of server with SUBID ", subid)
		for p, ok := serverList["default_password"].(string); p == "" || ok == false; p, ok = serverList["default_password"].(string) {
			time.Sleep(time.Second * 2)
			respData = requestToVultr("/v1/server/list?SUBID="+subid, "GET", nil)
			json.Unmarshal(respData, &serverList)
		}
		log.Println("Retrieved basic information of server on : ", serverList["main_ip"].(string))
		server = &VultrServer{
			Subid:    subid,
			Password: serverList["default_password"].(string),
			IPv4:     net.ParseIP(serverList["main_ip"].(string)),
			IPv6:     net.ParseIP(serverList["v6_main_ip"].(string)),
		}

		//We want a server that can be connected anyway
		log.Println("Waiting for server to respond so we can know it is accessible from here")

		sshClient := gossh.New(server.Addr().String(), "root")
		defer sshClient.Close()
		sshClient.SetPassword(server.RootPassword())

		//We will wait five minutes until it respond or we will just destroy it and get a new one
		startPoint := time.Now()
		for _, err := sshClient.Execute("ls"); err != nil; _, err = sshClient.Execute("ls") {
			time.Sleep(time.Second * 2)
			if time.Now().Sub(startPoint) > time.Minute*5 {
				log.Println("Server on ", server.Addr().String(), " seem never going to respond, destroying...")
				server.Destroy()
				continue
			}
		}
		goodServerCreated = true
		s = server
	}
	return s
}

func requestToVultr(uri string, method string, data []byte) []byte {
	req, err := http.NewRequest(method, "https://api.vultr.com"+uri, bytes.NewBuffer(data))
	if err != nil {
		log.Fatalln("Failed to create a HTTP request :", err)
	}
	req.Header.Add("API-Key", apikey)
	client := &http.Client{}

	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	var resp *http.Response
	for resp, err = client.Do(req); err != nil || resp.StatusCode == 503; resp, err = client.Do(req) {
		if err != nil {
			log.Println("Error access Vultr for new server : ", err)
		} else {
			log.Println("Accessing too fast...")
		}
		time.Sleep(time.Second)
	}
	if resp.StatusCode != 200 {
		if resp.StatusCode == 412 {
			respData, _ := ioutil.ReadAll(resp.Body)
			os.Stderr.Write(respData)
			os.Stderr.Write([]byte{'\n'})
		}
		log.Fatalln("Server responded with code ", resp.StatusCode, " , which may be caused a change of API or Service tempororily unavailbale")
	}
	respData, err := ioutil.ReadAll(resp.Body)
	return respData
}
