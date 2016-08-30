+ /register POST
	- name
	- password
+ /login POST sets cookie if the name and password matches
	- name
	- password
+ /{username}/port GET returns the port of the service that this user use 
+ /{username}/changepassword POST
	- password
+ /servers GET returns server info in JSON, including server IP and root password, must be admin
+ /servers/add GET/POST order to add a server, must be admin
	- count (default : 1)
+ /servers/del GET/POST removes a server, must be admin
	- count (default : 1)
+ /servicelist GET returns Shadowsocks service info, in the format like most SS client config file
+ /servicelist/v6 GET returns Shadowsocks service info, in the format like most SS client config file, specially only includes the IPv6 address
+ /servicelist/v4 GET returns Shadowsocks service info, in the format like most SS client config file, specially only includes the IPv4 address