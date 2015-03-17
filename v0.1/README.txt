README.txt

The UDP message format in JSON


{"action":"<action>", "agent":"<agent>", "client":"<cient>", "location","<location>"}

Depending on the action the data at the end will be different

# location update sends the current location of the agent (not the client)
locationUpdate()
{"action":"locationUpdate", "agent":"0", "client":"0", "location", [x, y, z]}

# Fire commands are interesting
# I have read that sending a target is not "safe". Instead the fire command should 
#just be a direction, then the server computes if position of client firing in direction hit anything.
fire()
{"action":"fire", "agent":"0", "client":"0", "target", "2"}


# Destroy
# originates from the server and is used to terminate an agent...
# MUST BE ACHKNOLEDGED BE EACH CLIENT? OR JUST THE CLIENT CONTROLLING THE DESTROIED AGENT?
destory()
{"action":"destory", "agent":"9"}
This is good for now

#################
# An action that is separate from the game (sort of)
# sent from a new client to the server
{"action":"join"}

# Responce to new client joining game
# Clients id and agent id
{"action":"join-accept", "client":0, "agent":0, "location":[x,y,z]}


