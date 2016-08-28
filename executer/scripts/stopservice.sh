#!/bin/bash

lsof -i :$1 | grep LISTEN | cut -f 2 -d ' ' | xargs kill