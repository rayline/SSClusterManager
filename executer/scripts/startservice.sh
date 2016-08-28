#!/bin/bash

nohup ss-server -p $1 -k $2 -m aes-256-cfb -u -A --fast-open > /dev/null 2>&1 &