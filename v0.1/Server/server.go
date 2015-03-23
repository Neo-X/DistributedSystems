/***
*	File Name: server.go
*	Description: The server provides the two functionalities:
*				 - service move request
*				 - service fire request
*	Date: 	8 March 2015
*	Author: Glen & Ravjot
*/
package main

import(
	"fmt"
	"flag"
	"net"
	"../dsgame"
	"encoding/json"
	// "sync"
	"strconv"
)

var clientAgentMap map[string]string // A map from client to the agent controlled by that client
var nodes map[string]*net.UDPConn // A map of the node name to the client connection to that node

var agentsDB map[string]dsgame.Agents // Database of all agent details [client -> <details>]
var clientOffset int // to generate new client name and agent name for new joinee


/***
*	Function Name: 	main()
*	Desc:			The main function for server
*	Pre-cond:		
*	Post-cond:		Call the service functions
*/
func main(){
	
	clientOffset = 0
	
	clientLink := flag.String("clientAddress", "127.0.0.1:10000", "The ip address clients should use to connect to this service")
	logFilePtr := flag.String("logfile", "GameServer", "The log file for the GameServer.")
	flag.Parse()
	
	fmt.Println("clientLink:", *clientLink)
    fmt.Println("logFile:", *logFilePtr)
    
    nodes = make(map[string]*net.UDPConn)
    clientAgentMap = make(map[string]string)
    agentsDB = make(map[string]dsgame.Agents)
    
    udpAddress, err := net.ResolveUDPAddr("udp4",*clientLink)

    if err != nil {
       fmt.Println("error resolving UDP address on ", *clientLink)
       fmt.Println(err)
       return
   }

   conn ,err := net.ListenUDP("udp",udpAddress)

   if err != nil {
        fmt.Println("error listening on UDP port ", *clientLink)
        fmt.Println(err)
        return
   }
        
    defer conn.Close()
	var buf []byte = make([]byte, 1500)   
    for n := int64(0); n >= 0; n++ {
	 	n,address, err := conn.ReadFromUDP(buf)
	 	if err != nil {
            fmt.Println("error reading data from connection")
            fmt.Println(err)
            return
        }
	 	
        if address != nil {
            fmt.Println("got message from ", address, " with n = ", n)
            if n > 0 {
            	fmt.Println("from address", address, "got message:", string(buf[0:n]), n)
            	////// Everything should be good now
            	handleMessage(conn , address, buf[0:n])
            	
            }
            /* conn, err := net.DialUDP("udp", nil, address)
            if err != nil {
            	fmt.Println("Error connecting to UDP client")
		        fmt.Println(err)
            }*/
            n, err :=	conn.WriteToUDP([]byte("Thank you for your message"), address)
		    if err != nil {
		        fmt.Println("WriteUDP Message", n)
		        fmt.Println(err)
		    } 
        }
	 }
   
	
}


/***
*	Function Name: 	serviceMoveReq()
*	Desc:			The function provide service to client's move request
*	Pre-cond:		takes connection argument and the new location
*	Post-cond:		list is updated with new location or return failure
*/
func serviceUpdateLocationReq(conn *net.UDPConn, msg dsgame.Message){
	// update server agent database
	if nodes[msg.Client] == conn && clientAgentMap[msg.Client] == msg.Agent {
		var tmpObj dsgame.Agents
		tmpObj.TimeStamp = msg.TimeStamp
		tmpObj.Location = [3]float64{msg.Location[0],msg.Location[1],msg.Location[2]}
		agentsDB[msg.Client] = tmpObj
		fmt.Printf("Location updated by:" + msg.Client)
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
func serviceFireReq(conn *net.UDPConn, msg dsgame.Message){
	fmt.Printf("Target Received M: %f  C: %f  Q: %f\n", msg.Target.M , msg.Target.C, msg.Target.Q)
	// destroy the client if valid
}

/***
*	Function Name: 	serviceJoinReq()
*	Desc:			The function provide service to a client requesting to join the game
*	Pre-cond:		takes connection argument ....
*	Post-cond:		Client should be registard in the game
*/
func serviceJoinReq(conn *net.UDPConn, clientAddr *net.UDPAddr, msg dsgame.Message){
	
	clientName, agentName := getNextClientID()

	// Add to server agent database
	clientAgentMap[clientName] = agentName
	nodes[clientName] = conn

	var tmpObj dsgame.Agents
	//	tmpObj.Agent = agentName
	tmpObj.TimeStamp = 0
	tmpObj.Location = [3]float64{1.0,2.0,3.0}
	agentsDB[clientName] = tmpObj

	// Prepare response for the request
	msg.Action = dsgame.AcceptJointAction
	msg.Client = clientName
	msg.Agent = agentName
	msg.TimeStamp = 0
	msg.Location = [3]float64{1.0,2.0,3.0}
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
	
}

/***
*	Function Name: 	serviceJoinReq()
*	Desc:			The function provide service to a client requesting to join the game
*	Pre-cond:		takes connection argument ....
*	Post-cond:		Client should be registard in the game
*/
func handleMessage(conn *net.UDPConn, clientAddr *net.UDPAddr, buf []byte){
	var msg dsgame.Message
	err := json.Unmarshal(buf, &msg)
	if err != nil {
		fmt.Println("Error handling message")
		fmt.Println(err)
		return
	}
	
	if msg.Action == dsgame.JoinAction {
		serviceJoinReq(conn, clientAddr, msg)
	} else if msg.Action == dsgame.UpdateLocationAction {
        serviceUpdateLocationReq(conn, msg)
	} else if msg.Action == dsgame.FireAction {
				serviceFireReq(conn, msg)
  } else {
     fmt.Println("Message not understood: ")
  }
	
}


/***
*	Function Name: 	getNextClientID()
*	Desc:			The function provide next client and agent name to a client requesting to join the game
*	Pre-cond:		takes connection argument ....
*	Post-cond:		new client and agent name should be returned
*/
func getNextClientID() (string,string){
	clientOffset++
	return "client" + strconv.Itoa(clientOffset), "agent" + strconv.Itoa(clientOffset)
}


