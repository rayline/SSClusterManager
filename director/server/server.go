package servers

//the package create, destroy, setup and directs the servers
//every server will have a gorountine to process the directions

import "net"

const StateDestroyed = "D"
const StateNotReady = "B"
const StateOK = "O"

type Server interface {
	Addr() net.IP
	AddrV6() net.IP
	RootPassword() string
	Status() string
	Do(uri string, data []byte) //There may be bugs here in future if extremely large amount of server and users appear, but not for now
	Reboot()
	Destroy()
}

func Available(s *Server) bool {
	// availability check
	//TODO: more checking items
	if s.Status() != StateOK {
		return false
	}
	return true
}
