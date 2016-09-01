#!/bin/bash

.source ~./bashrc
nohup ~/go/src/SSClusterManager/executer/executer > init.log 2> initerr.log &