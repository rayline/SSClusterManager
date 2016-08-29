#!/bin/bash

lsof -i :$1 | grep LISTEN | cut -f 2 -d ' ' | xargs kill
nohup ss-server -p $1 -k $2 -m aes-256-cfb -u -A --fast-open > /dev/null 2>&1 &