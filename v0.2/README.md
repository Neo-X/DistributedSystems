# ARM Game With Distributed States

##Description
The system is a prototype of distributed system that supports an Asynchronous Real-time Multiplayer Game (ARM Game). 

## System Architecture
![Alt text](https://github.com/Neo-X/DistributedSystems/blob/master/FinalReport/images/client-distributed-server-model-Activity.png "System Architecture")


##Running the System
The system has many different component which are required to run separately.

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



##Errors/Bugs:
There can be a bug or race condition between the ServiceLocationUpdate and Respawn agent. Due to the number of threads the location could be updated after it is set un Respawn by locationUpdate and the respawn is lost.
