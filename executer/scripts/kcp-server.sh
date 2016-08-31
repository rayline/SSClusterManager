#!/bin/bash

# usage: kcp-server.sh kcp-listen-port targer-port key

netstat -lp | grep $1 | awk '{print $6}' | grep / | cut -f 1 -d "/" | xargs kill
nohup kcptun/server_linux_amd64 -l :$1 -t 127.0.0.1:$2 --crypt salsa20 --key $3 --mtu 1200 --nocomp --mode fast2 --dscp 46 &
