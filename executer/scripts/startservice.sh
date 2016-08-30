#!/bin/bash

lsof -i :$1 | grep LISTEN | cut -f 2 -d ' ' | xargs kill
netstat -lp | grep $1 | awk '{print $6}' | grep / | cut -f 1 -d "/" | xargs kill
nohup kcptun/server_linux_amd64 -l :$1 -t 127.0.0.1:$1 --crypt salsa20 --key $2 --mtu 1200 --nocomp --mode fast2 --dscp 46 &
nohup ss-server -p $1 -k $2 -m chacha20 -A --fast-open > /dev/null 2>&1 &