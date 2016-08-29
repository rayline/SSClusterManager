package HTTPFrontEnd

//the file is where controllers, those who handle HTTP requests, are designated.
//Since we need to do a security check for all requests, controllers should be set under the new Mux

import "net/http"
import "SSClusterManager/executer/util"
import "log"
import "SSClusterManager/executer/UserManager"
import "strconv"

var controllerMux *http.ServeMux
var securityMux *http.ServeMux

func init() {
	controllerMux = http.NewServeMux()
	setupControllers()
	securityMux = http.NewServeMux()
	securityMux.HandleFunc("/", securityController)
}

func securityController(w http.ResponseWriter, r *http.Request) {
	//at this step we use a pre-shared key in requests from Directors
	//if the key was not configured, the program should have a fatal Error
	SecurityKey, exist := util.Configs["SecurityKey"].(string)
	if exist == false {
		log.Fatalln("SecurityKey not set or wrongly set. The key must be set and pre shared among director and executers for security concern")
	}
	if r.Header.Get("key") == SecurityKey {
		controllerMux.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func setupControllers() {
	controllerMux.HandleFunc("/writeuser", WriteUserController)
	controllerMux.HandleFunc("/deluser", DelUserController)
}

// /write POST adds a user to the executer
// name :the user name
// password :the user password
// port :the port the user is supposed to use
func WriteUserController(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	portStr := r.FormValue("port")
	port64, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	port := uint16(port64)
	UserManager.WriteUser(UserManager.User{
		Port:     port,
		Password: password,
	})
	w.WriteHeader(http.StatusOK)
}

// /deluser POST deletes a user from the executer
// port :the port of service
func DelUserController(w http.ResponseWriter, r *http.Request) {
	portStr := r.FormValue("port")
	port64, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	port := uint16(port64)
	UserManager.DelUser(port)
	w.WriteHeader(http.StatusOK)
}
