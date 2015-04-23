# ARM Game With Distributed States

##Description
The system is a prototype of distributed system that supports an Asynchronous Real-time Multiplayer Game (ARM Game). 

##Running the System

###Order for running the system

- Start the kvservice: ./kvservicemain 127.0.0.1:9999 0
- Start the activityserver: ./activityServer -address 127.0.0.1:9999
- Start some servers: ./server -clientAddress 127.0.0.1:10000 -logfile node -ID 0
- Start some servers: ./server -clientAddress 127.0.0.1:10001 -logfile node -ID 1
- Start the clients for the servers: ./client -serverAddress 127.0.0.1:10000
- Start the clients for the servers: ./client -serverAddress 127.0.0.1:10001

###Run Large Number of Nodes

- Start the kvservice: ./kvservicemain 127.0.0.1:9999 0
- Start the activityserver: ./activityServer -address 127.0.0.1:9999
- Start some servers: ./startServers.sh 30
- Start some clients: ./startClients.sh 30

You will probably need to execute killall server and killall client to end the processes.


##Errors/Bugs:
There are still few open issues.
