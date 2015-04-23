# ARM Game With Distributed States

##Description
The system is a prototype of distributed system that supports an Asynchronous Real-time Multiplayer Game (ARM Game). 

##Running the System

###Order for running the system

1) Start the kvservice: ./kvservicemain 127.0.0.1:9999 0
2) Start the activityserver: ./activityServer -address 127.0.0.1:9999
3) Start some servers: ./server -clientAddress 127.0.0.1:10000 -logfile node -ID 0
4) Start some servers: ./server -clientAddress 127.0.0.1:10001 -logfile node -ID 1
5) Start the clients for the servers: ./client -serverAddress 127.0.0.1:10000
6) Start the clients for the servers: ./client -serverAddress 127.0.0.1:10001

###Run Large Number of Nodes

1) Start the kvservice: ./kvservicemain 127.0.0.1:9999 0
2) Start the activityserver: ./activityServer -address 127.0.0.1:9999
3) Start some servers: ./startServers.sh 30
4) Start some clients: ./startClients.sh 30

You will probably need to execute killall server and killall client to end the processes.


##Errors/Bugs:
There are still few open issues.
