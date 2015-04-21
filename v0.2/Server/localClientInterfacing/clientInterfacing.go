package localClientInterfacing

import(
	"fmt"
//	"flag"
	"net"
	"../../dsgame"
	"../../s3dm"
	// "../fixed"
	"encoding/json"
	// "sync"
	"strconv"
	"../header"
	"time"
	
)


/***
*	Function Name: 	ServiceUpdateLocationReq()
*	Desc:			The function provide service client move requests
*	Pre-cond:		takes connection argument and the new location
*	Post-cond:		list is updated with new location or return failure
Also return failure on this being an invalid update??
*/
func ServiceUpdateLocationReq(conn *net.UDPConn, msg dsgame.Message) bool {
	// update server agent database
	_timeNow := time.Now().UnixNano()
	// if header.Nodes[msg.Client] == conn { // && header.ClientAgentMap[msg.Client] == msg.Agent { // [Glen] Maybe use some other invarients
	if ( header.MyAgent.Name == msg.Agent) { // Need to validate location update
		// (NEW location - old location).length()*deltaTime < maxVelocity*deltaTime
		fmt.Println("Checking for invalid location")
		distance := msg.Location.Sub(header.AgentDB[msg.Agent].Location).Length()
		
		clientDeltaTime := _timeNow - header.AgentDB[msg.Agent].LastUpdateTime // Timestamps should be controlled by server?

		deltaTime := float64(clientDeltaTime)/1000000000.0
		fmt.Println("Checking for invalid location, msgTime", _timeNow, " agent time", header.AgentDB[msg.Agent].LastUpdateTime, " delta time in seconds", deltaTime)
		if ( (distance/deltaTime) > dsgame.GameMaxVelocity ) {
			fmt.Println("Invalid location submitted by client:", msg.Client, " location", msg.Location, " deltaTime", deltaTime )
			// need to send position override now.
			SendPositionOverrideforAgent()
			return false
		} 
	
	}
		var tmpObj dsgame.Agent // This covers the case when the agent has not been initialized yet, but I don't think we want this
		tmpObj.TimeStamp = msg.TimeStamp
		tmpObj.Location = s3dm.V3{msg.Location.X,msg.Location.Y,msg.Location.Z}
		tmpObj.LastUpdateTime = _timeNow
		header.AgentDB[msg.Agent] = tmpObj
		header.ClientAgentMap[msg.Client] = msg.Agent // This does not need to be set so often
		// fmt.Println("Location updated by:" + msg.Client)
		SendPositionforAgent(msg)
		return true
	// } else {
		// fmt.Println("something seems fishy with the agent updates")
	// }
	// return false
}


/***
*	Function Name: 	serviceFireReq()
*	Desc:			The function provide service to client's fire request
*	Pre-cond:		takes connection argument and name of client who is fired
*	Post-cond:		Destroy the client or returns failure
*/
func ServiceFireReq(conn *net.UDPConn, msg dsgame.Message){
	fmt.Println()
	fmt.Println("Firing projectile --->> ", msg.Target)
	fmt.Println("From agent  ", msg.Client)
	fmt.Println()
	pos := header.AgentDB[msg.Agent].Location
	// destroy the client if valid

	if header.MyAgent.Name == msg.Agent { // if it is from my client
		for key, val := range header.ClientAgentMap {
 			 fmt.Println("Client:", key, " agent :", val)
 			// fmt.Println("agent:", key, " agent Location:", value.Location)
		 value :=	header.AgentDB[val]
 	 	 if (val != msg.Agent ) { // Ignore intersections with self
				if (dsgame.RayHitsAgent(value.Location, pos, msg.Target)) {
	 	   		// fmt.Println("Ray hit agent", key)
					//send that client a destroy request
 		   		sendDestroyReq(key, msg) 
					break; // assuming it hits single agent only
 		   	} 	   	
 		  }
 		}

	} else { // if it is from other servers
		// verify if it hits me
		if (dsgame.RayHitsAgent(header.MyAgent.Location,msg.Location,msg.Target)) { // if it hits me, broadcast destroy msg
			BroadcastDestroyMeReq(msg)
		} else { //if doesn't, ignore msg
			fmt.Println("----->I updated my loc, just before getting hit")
		}
	}
	
}

func sendDestroyReq(expectedhit_ClientName string, msg dsgame.Message) {
/*	conn, err := net.ResolveUDPAddr("udp",header.Nodes[expectedhit_ClientName])
/	conn_server, err := net.DialUDP("udp", nil, header.Nodes[expectedhit_ClientName])
		if err != nil {
        fmt.Println("Error connecting to " , expectedhit_ClientName)
        fmt.Println(err)
        return
    }
*/
	//fmt.Println("Inside sendDestroyReq")
	conn_Server := header.Nodes[expectedhit_ClientName]
	b, err := json.Marshal(msg)
		if err != nil {
        fmt.Println("Problem marshalling struct")
        fmt.Println(err)
    } 
	
	_, err = conn_Server.Write(b)
//	var buf []byte = make([]byte, 1500) 
	fmt.Println("sendDestroyReq sent.")
/*
	_, _, err = conn_Server.ReadFromUDP(buf) 
    if err != nil {
        fmt.Println("ReadFromUDP")
        fmt.Println(err)
    }
*/
}

func BroadcastDestroyMeReq(msg dsgame.Message) {
	fmt.Println("I should broadcast destroy message")
	msg.Action = dsgame.DestroyAction
	msg.Client = header.MyClientName
	msg.Agent = header.MyAgent.Name
	msg.Location = dsgame.GetRandomDirection()
	for _, conn := range header.Nodes {
//-	for key, conn := range header.Nodes {
//-		if (key != header.MyClientName) { // don't send message to self
			// fmt.Println("node:", key, "ip:", value.RemoteAddr().String() )
		
			b, err := json.Marshal(msg)
			if err != nil {
		        fmt.Println("Problem marshalling struct")
		        fmt.Println(err)
		    } 
			
			_, err =	conn.Write(b)
		    if err != nil {
		        fmt.Println("WriteUDP")
		        fmt.Println(err)
		    } 
		
	//-	}	
 	}
}

/***
*	Function Name: 	handleDestroyReq()
*	Desc:			The function provide service to destroy agents that are shot
*	Pre-cond:		takes connection argument and name of agent that is destoroid
*	Post-cond:		Destroy the agent, send it a new random location
*/
func HandleDestroyReq(msg dsgame.Message) {
	var tmp dsgame.Agent
	tmp = header.AgentDB[msg.Agent]
	tmp.Location = msg.Location
	header.AgentDB[msg.Agent] = tmp
	SendPositionforAgent(msg)
	fmt.Println(msg.Client +" is Destroyed!!!!!!")
		
}



/***
*	Function Name: 	serviceJoinReq()
*	Desc:			The function provide service to a client requesting to join the game
*	Pre-cond:		takes connection argument ....
*	Post-cond:		Client should be registard in the game
*/
func ServiceJoinReq(conn *net.UDPConn, clientAddr *net.UDPAddr, msg dsgame.Message, id uint64){

/*
	server, err := net.ResolveUDPAddr("udp",header.CentralServerIP_Port)
	conn_centralServer, err := net.DialUDP("udp", nil, server)
		if err != nil {
        fmt.Println("Error connecting to " , server)
        fmt.Println(err)
        return
    }

	b, err := json.Marshal(msg)
		if err != nil {
        fmt.Println("Problem marshalling struct")
        fmt.Println(err)
    } 
	
	_, err = conn_centralServer.Write(b)
	var buf []byte = make([]byte, 1500) 

	n, _, err := conn_centralServer.ReadFromUDP(buf) 
    if err != nil {
        fmt.Println("ReadFromUDP")
        fmt.Println(err)
    }

	var m dsgame.Message
	err = json.Unmarshal(buf[0:n],&m)
    if err != nil {
        fmt.Println("Error unmarshalling message")
        fmt.Println(err)
    }
    */
	/*newly added */
	header.MyClientName = "node" + strconv.FormatUint(id, 10)
	header.MyAgent.Name = "agent"  + strconv.FormatUint(id, 10)
	header.MyAgent.Location = dsgame.GetRandomLocation()
	header.ClientLink = clientAddr
	header.Connection = conn
	
	var _msg dsgame.Message
	_msg.Action = dsgame.AcceptJointAction
	_msg.Client = header.MyClientName
	_msg.Agent = header.MyAgent.Name
	_msg.Location = header.MyAgent.Location
	
	 
	b, err := json.Marshal(_msg)
		if err != nil {
        fmt.Println("Problem marshalling struct")
        fmt.Println(err)
    } 

	_, err = conn.WriteToUDP(b,clientAddr)
	
	var tmpObj dsgame.Agent
	tmpObj.TimeStamp = msg.TimeStamp
	tmpObj.LastUpdateTime = time.Now().UnixNano()
	tmpObj.Location = _msg.Location
	tmpObj.Name = _msg.Agent
	header.AgentDB[_msg.Agent] = tmpObj
	header.ClientAgentMap[_msg.Client] = _msg.Agent 
	
	

	//fmt.Println(string(buf[0:n]))

	
}



/***
*	Function Name: 	getNextClientID()
*	Desc:			The function provide next client and agent name to a client requesting to join the game
*	Pre-cond:		takes connection argument ....
*	Post-cond:		new client and agent name should be returned
*/
func getNextClientID() (string,string){
	header.ClientOffset++
	return "client" + strconv.Itoa(header.ClientOffset), "agent" + strconv.Itoa(header.ClientOffset)
}

/***
*	Function Name: 	SendPositionOverrideforAgent()
*	Desc:			Sends a message back to the client overwriting the location of the agent
*	Pre-cond:		Don't think it needs any arguments.
*	Post-cond:		The message should be sent and hopefully the client will update the agent location
*/
func SendPositionOverrideforAgent() {
	var _msg dsgame.Message
	_msg.Action = dsgame.PositionOverrideAction
	_msg.Client = header.MyClientName
	_msg.Agent = header.MyAgent.Name
	_msg.Location = header.AgentDB[header.MyAgent.Name].Location
	
	 
	b, err := json.Marshal(_msg)
		if err != nil {
        fmt.Println("Problem marshalling struct")
        fmt.Println(err)
    } 

	_, err = header.Connection.WriteToUDP(b,header.ClientLink)
	if ( err != nil ) {
		fmt.Println("Problem sending Position override to client")
        fmt.Println(err)
	}
}

/***
*	Function Name: 	SendPositionforAgent()
*	Desc:			Sends a message back to the client for the location of the agent
*	Pre-cond:		Don't think it needs any arguments.
*	Post-cond:		The message should be sent and hopefully the client will update the agent location
*/
func SendPositionforAgent(_msg dsgame.Message) {
	 
	b, err := json.Marshal(_msg)
		if err != nil {
        fmt.Println("Problem marshalling struct")
        fmt.Println(err)
    } 

	if ( header.Connection != nil ) { // There is kind of a gray area for this first message, before the client has joined properly
		_, err = header.Connection.WriteToUDP(b,header.ClientLink)
		if ( err != nil ) {
			fmt.Println("Problem sending Position override to client")
	        fmt.Println(err)
		}	
	}
	
}


/*
func broadcast(msg dsgame.Message){
	for k,_ := range header.AgentsDB {
		
	}
}
*/
