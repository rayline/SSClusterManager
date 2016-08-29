package UserManager

//the package UserManager is responsible for maintaining users, and start the service for each user when executer starts or the user added

import "SSClusterManager/executer/UserManager/DB"
import "SSClusterManager/executer/UserManager/UserType"
import "SSClusterManager/executer/ServiceManager"

type User UserType.User

func init() {
	userList := DB.GetAll()
	for _, value := range userList {
		ServiceManager.StartService(value.Port, value.Password)
	}
}

func WriteUser(u User) {
	previousU, exist := DB.Get(name)
	if exist {
		ServiceManager.StopService(previousU.Port)
	}
	DB.Add(UserType.User(u))
	ServiceManager.StartService(u.Port, u.Password)
}

func DelUser(name string) {
	u, exist := DB.Get(name)
	if !exist {
		return
	}
	ServiceManager.StopService(u.Port)
	DB.Del(name)
}
