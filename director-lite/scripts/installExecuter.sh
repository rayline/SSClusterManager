#!/bin/bash

# The file should be distributed alone since the script will get the latest executer source code from GitHub

cd ~

#installing ShadowSocks-libev

killall ss-server
killall ss-local
rm -rf ./shadowsocks-libev
git clone https://github.com/shadowsocks/shadowsocks-libev.git

cd shadowsocks-libev
apt -y install --no-install-recommends build-essential autoconf libtool libssl-dev gawk debhelper dh-systemd init-system-helpers pkg-config asciidoc xmlto apg
dpkg-buildpackage -b -us -uc -i
cd ..
sudo dpkg -i shadowsocks-libev*.deb

#installing Executer of SSClusterManager

killall executer
rm -rf ./SSClusterManager
git clone https://github.com/rayline/SSClusterManager.git

chmod +x SSClusterManager/executer/scripts/*

SSClusterManager/executer/scripts/golanginstall.sh --64
source ~/.bashrc

cp -rf SSClusterManager $GOPATH/src 
cd $GOPATH/src/SSClusterManager/executer 

go build

cp scripts/startexecuter.sh /etc/init.d 

#specially add kcptun installation
# NOTE: not compile installation so the script may need to be updated when kcptun updates, and no good solution before its upload to apt source
mkdir kcptun
cd kcptun
wget https://github.com/xtaci/kcptun/releases/download/v20160830/kcptun-linux-amd64-20160830.tar.gz
tar -zxvf kcptun-linux-amd64-20160830.tar.gz
cd ..

scripts/startexecuter.sh
