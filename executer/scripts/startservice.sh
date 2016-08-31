#!/bin/bash

lsof -i :$1 | grep LISTEN | cut -f 2 -d ' ' | xargs kill
nohup ss-server -p $1 -k $2 -m chacha20 -A -u --fast-open > /dev/null 2>&1 &
./kcp-server.sh $(($1 + 1)) $1 $2