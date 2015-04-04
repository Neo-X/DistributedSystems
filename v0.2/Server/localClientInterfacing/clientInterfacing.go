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
)


/***
*	Function Name: 	serviceMoveReq()
*	Desc:			The function provide service to client's move request
*	Pre-cond:		takes connection argument and the new location
*	Post-cond:		list is updated with new location or return failure
*/
func ServiceUpdateLocationReq(conn *net.UDPConn, msg dsgame.Message){
	// update server agent database
	if header.Nodes[msg.Client] == conn && header.ClientAgentMap[msg.Client] == msg.Agent {
		var tmpObj dsgame.Agents
		tmpObj.TimeStamp = msg.TimeStamp
		tmpObj.Location = s3dm.V3{msg.Location.X,msg.Location.Y,msg.Location.Z}
		header.AgentsDB[msg.Client] = tmpObj
		fmt.Println("Location updated by:" + msg.Client)
	} else {
		// if something seems fishy
	}
}


/***
*	Function Name: 	serviceFireReq()
*	Desc:			The function provide service to client's fire request
*	Pre-cond:		takes connection argument and name of client who is fired
*	Post-cond:		Destroy the client or returns failure
*/
func ServiceFireReq(conn *net.UDPConn, msg dsgame.Message){
	fmt.Println("Firing projectile  ", msg.Target)
	fmt.Println("From agent  ", msg.Client)
	pos := header.AgentsDB[msg.Client].Location
	// destroy the client if valid
	
	for key, value := range header.AgentsDB {
 		// fmt.Println("agent:", key, " agent Location:", value.Location)
 	  if (key != msg.Client ) { // Ignore intersections with self
			if (dsgame.RayHitsAgent(value.Location, pos, msg.Target)) {
 	   		// fmt.Println("Ray hit agent", key)
 	   		handleDestroyReq(conn, value) 
 	   	} 	   	
 	  }
 	}

}

/***
*	Function Name: 	handleDestroyReq()
*	Desc:			The function provide service to destroy agents that are shot
*	Pre-cond:		takes connection argument and name of agent that is destoroid
*	Post-cond:		Destroy the agent, send it a new random location
*/
func handleDestroyReq(conn *net.UDPConn, agent dsgame.Agents) {
	
	
}



/***
*	Function Name: 	serviceJoinReq()
*	Desc:			The function provide service to a client requesting to join the game
*	Pre-cond:		takes connection argument ....
*	Post-cond:		Client should be registard in the game
*/
func ServiceJoinReq(conn *net.UDPConn, clientAddr *net.UDPAddr, msg dsgame.Message){
	
	clientName, agentName := getNextClientID()

	// Add to server agent database
	header.ClientAgentMap[clientName] = agentName
	header.Nodes[clientName] = conn

	var tmpObj dsgame.Agents
	//	tmpObj.Agent = agentName
	tmpObj.TimeStamp = 0
	
	tmpObj.Location = s3dm.V3{1.0,3.0,-2.0}
	header.AgentsDB[clientName] = tmpObj

	// Prepare response for the request
	msg.Action = dsgame.AcceptJointAction
	msg.Client = clientName
	msg.Agent = agentName
	msg.TimeStamp = 0
	// pick random location
	msg.Location = dsgame.GetRandomLocation()
	
	//msg.Target = ""
	
	buf, err := json.Marshal(msg)
	if err != nil {
        fmt.Println("Problem Marshaling Joint Req message")
        fmt.Println(err)
    } 
	
	n, err :=	conn.WriteToUDP(buf, clientAddr)
	if err != nil {
		fmt.Println("WriteUDP Message", n)
		fmt.Println(err)
	}
    
	// If everything was sent update the DB
	var _agent dsgame.Agents
	_agent.Location = msg.Location
  _agent.TimeStamp = 0
  header.AgentsDB[msg.Client] = _agent 
	
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


/*
func broadcast(msg dsgame.Message){
	for k,_ := range header.AgentsDB {
		
	}
}
*/
