# Shadowsocks Cluster Manager

Shadowsocks Cluster Manager is an application intending to create and maintain Shadowsocks servers automically. The project will be mainly based on Golang.  

## Roles

+ Director, the controlling server of the cluster, which is responsible for creating, maintaining and destroying other servers which run the actual service. It is also the interface for the user to learn about the information of the service, such as server address, and manage users if there are multiple users. (NOT IMPLEMENTED)
+ Execture, the servers that runs the service. It should run both the actual service we want and the service to recieve order from director.
+ Director-Lite, the lightweighted CLI tool that controlls executers, for single user use on personal computer

## Modules

+ director service(NOT IMPLEMENTED)
	- server manager (create and send order to executers)
	- maintainer (ensures there are a right number of executer running service and sensor their status)
	- user interface (HTTP server to recieve order from user and provide service information)
+ executer service
	- service controler (runs and stops the service)
	- director interface (recieve order from director)
	- user manager (manages all users sent from director)
+ director-lite
	- server manager
	- CLI interface (reloading the program itself is the only purpose, everything is configured through the config file)

## Behaviour description

The director should run on a server. Configured with HTTP interface, it creats server from server provider, on which director system is automatically installed and configured. Users of the service should also be configured through HTTP, whose information include port and password. At this version the authentication of the user interface is only a pre-shared key in the HTTP header. There should be TLS for the interface if possible. 
The executer runs the Shadowsocks instances, with user information (name, port, password) sent from the director. It should run an alternable script to start Shadowsocks. The implemention of Shadowsocks should not be limited, as long as it works well when multiple instances run well on different ports. The config information from director should be stored locally in case of a restart.
The director-lite only loads the config file, checkout the servers and setup new servers or destroy redundant servers depend on the situation.

## Conventions

+ Shadowsocks will use chacha20 as cipher
+ One-time authentication should be enabled
+ default port of executer is 41600
+ default port of director api is 8080
+ default port of director interface is 80
+ UDP relay will be on
+ IMPORTANT: Shadowsocks should ONLY take the EVEN port number, while KCPTUN will run on the port number plus 1
+ KCPTUN uses crypt salsa20, and the same key as Shadowsocks

## System requirement

Ubuntu 16 x64