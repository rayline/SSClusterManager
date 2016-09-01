package servers

import "SSClusterManager/director-lite/server"
import "SSClusterManager/director-lite/util"
import "sync"
import "strconv"
import "encoding/json"
import "log"

//the package manages servers and intergrates all drivers for server provider
//all drivers must register themselves during the initializing routine, otherwise the driver may lose chance competing for request
//also drivers should better retrieve the created servers
//NOTE: at this version only one provider will be allowed and there is no information sent to the provider when creating server

type ServerProvider interface {
	Create() server.Server
}

var provider ServerProvider
var mutex sync.Mutex
var Servers = map[server.Server]bool{}

func AddServer(s server.Server) {
	mutex.Lock()
	defer mutex.Unlock()
	Servers[s] = true
	if !server.Available(s) {
		log.Println("Setting up server on ", s.Addr().String())
		go func() {
			server.Setup(s)
			server.WriteUser(s, util.Configs["password"].(string), util.Configs["port"].(string))
		}()
	} else {
		log.Println("Server at ", s.Addr().String(), " good for service")
		server.WriteUser(s, util.Configs["password"].(string), util.Configs["port"].(string))
	}
}

func RegisterProvider(serverProvider ServerProvider) {
	mutex.Lock()
	defer mutex.Unlock()
	provider = serverProvider
}

func CheckOutServerCount() {
	requiredServerCount := int(util.Configs["serverCnt"].(float64))
	for len(Servers) < requiredServerCount {
		log.Println("Creating a new server...")
		AddServer(provider.Create())
	}
	for len(Servers) > requiredServerCount {
		log.Println("Destroying needless server")
		for s, _ := range Servers {
			s.Destroy()
			delete(Servers, s)
			break
		}
	}
}

type ssserver struct {
	Server      string `json:"server"`
	Server_port uint16 `json:"server_port"`
	Password    string `json:"password"`
	Method      string `json:"method"`
	Remarks     string `json:"remarks"`
	Auth        bool   `json:"auth"`
}

var port uint16 = 0
var password = ""

func newSSserver(IP string) ssserver {
	if port == 0 {
		portStr := util.Configs["port"].(string)
		port64, _ := strconv.ParseUint(portStr, 10, 64)
		port = uint16(port64)
		password = util.Configs["password"].(string)
	}

	return ssserver{
		Server:      IP,
		Server_port: port,
		Password:    password,
		Method:      "chacha20",
		Remarks:     "",
		Auth:        true,
	}
}

func JSON() []byte {
	mutex.Lock()
	defer mutex.Unlock()
	ssservers := []ssserver{}
	for s, _ := range Servers {
		ssservers = append(ssservers, newSSserver(s.Addr().String()))
		ssservers = append(ssservers, newSSserver(s.AddrV6().String()))
	}
	data, err := json.MarshalIndent(ssservers, "", "	")
	str = by
	if err != nil {
		log.Fatalln("Error encoding SS server list to JSON ", err)
	}
	return data
}
