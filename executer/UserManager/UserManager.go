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

func AddUser(u User) {
	DB.Add(UserType.User(u))
	ServiceManager.StartService(u.Port, u.Password)
}

func DelUser(name string) {
	ServiceManager.StopService(DB.Get(name).Port)
	DB.Del(name)
}
