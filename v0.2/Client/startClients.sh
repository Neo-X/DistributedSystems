#!/bin/bash

# use ./startClientss <num of servers to start>

numClients=$1
i=0
port=10000
while [ $i -lt $numClients ]
do
	./client -serverAddress 127.0.0.1:$((port+i)) &
	i=$((i+1))
done


