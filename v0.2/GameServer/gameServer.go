/***
*	File Name: gamServer.go
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
//	"../s3dm"
	// "../fixed"
	"encoding/json"
	// "sync"
	"strconv"
	"os"
	"./kvmanager"
)

var clientOffset int // to generate new client name and agent name for new joinee
const kv_IP_Port string = "127.0.0.1:12345"

/***
*	Function Name: 	main()
*	Desc:			The main function for server
*	Pre-cond:		
*	Post-cond:		Call the service functions
*/
func main(){
	
	if len(os.Args) != 2 {
		fmt.Printf("Syntax: %s  <GameServerIP:Port> \n", os.Args[0])
		os.Exit(0)
	}	
	
	go kvmanager.Main(kv_IP_Port, "manager")	

	clientOffset = 0
	
	ServiceIP_Port := os.Args[1]
	clientLink := flag.String("clientAddress", ServiceIP_Port, "The ip address clients should use to connect to this service")
	logFilePtr := flag.String("logfile", "GameServer", "The log file for the GameServer.")
	flag.Parse()
	
	fmt.Println("clientLink:", *clientLink)
  fmt.Println("logFile:", *logFilePtr)
    
    
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
         }
      n, err :=	conn.WriteToUDP([]byte("Thank you for your message"), address)

		 	if err != nil {
		  	fmt.Println("WriteUDP Message", n)
		    fmt.Println(err)
		  } 
			*/
		}
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
		ServiceJoinReq(conn, clientAddr, msg)
  }else {
		fmt.Println("Message not understood: ")
  }
	
}


/***
*	Function Name: 	serviceJoinReq()
*	Desc:			The function provide service to a client requesting to join the game
*	Pre-cond:		takes connection argument ....
*	Post-cond:		Client should be registard in the game
*/
func ServiceJoinReq(conn *net.UDPConn, clientAddr *net.UDPAddr, msg dsgame.Message){
	
	clientName, agentName := getNextClientID()

	// Prepare response for the request
	msg.Action = dsgame.AcceptJointAction
	msg.Client = clientName
	msg.Agent = agentName
	msg.TimeStamp = 0
	// pick random location
	msg.Location = dsgame.GetRandomLocation()
	
	
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

	fmt.Println(string(buf))
    
	// If everything was sent update the DB
	var _agent dsgame.Agents
	_agent.Location = msg.Location
  _agent.TimeStamp = 0
	
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


