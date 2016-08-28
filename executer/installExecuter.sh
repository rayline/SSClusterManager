#!/bin/bash

# The file should be distributed alone since the script will get the latest executer source code from GitHub

cd ~

#installing ShadowSocks-libev

git clone https://github.com/shadowsocks/shadowsocks-libev.git

cd shadowsocks-libev
apt -y install --no-install-recommends build-essential autoconf libtool libssl-dev gawk debhelper dh-systemd init-system-helpers pkg-config asciidoc xmlto apg
dpkg-buildpackage -b -us -uc -i
cd ..
sudo dpkg -i shadowsocks-libev*.deb

#installing Executer of SSClusterManager

git clone https://github.com/rayline/SSClusterManager.git

chmod +x SSClusterManager/executer/scripts/*

SSClusterManager/executer/scripts/golanginstall.sh --64

cp -r SSClusterManager $GOPATH/src 
cd $GOPATH/src/SSClusterManager/executer 

go build

cp scripts/startexecuter.sh /etc/init.d 

scripts/startexecuter.sh
