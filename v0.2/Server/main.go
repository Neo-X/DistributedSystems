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
//	"../s3dm"
	// "../fixed"
	"encoding/json"
	// "sync"
//	"strconv"
	"./header"
	"./localClientInterfacing"
	"./serversInterfacing"
	"os"
)


/***
*	Function Name: 	main()
*	Desc:			The main function for server
*	Pre-cond:		
*	Post-cond:		Call the service functions
*/
func main(){
	
	if len(os.Args) != 5 {
		fmt.Printf("Syntax: %s <ServiceIP:Port> <LocalClientIP:Port> <GameServerIP:Port> <logfile>\n", os.Args[0])
		os.Exit(0)
	}	
	
	go serversInterfacing.Main(os.Args[1], header.KvService, os.Args[4])
	
	header.ClientOffset = 0
	header.ServiceIP_Port = os.Args[1]
	clientLink := flag.String("clientAddress", header.ServiceIP_Port, "The ip address clients should use to connect to this service")
	logFilePtr := flag.String("logfile", "GameServer", "The log file for the GameServer.")
	flag.Parse()
	
	fmt.Println("clientLink:", *clientLink)
  fmt.Println("logFile:", *logFilePtr)
    
  header.Nodes = make(map[string]*net.UDPConn)
  header.ClientAgentMap = make(map[string]string)
  header.AgentsDB = make(map[string]dsgame.Agents)
  header.OnlineNodes = make(map[string]string)
    
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
     		printState()       	
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
		localClientInterfacing.ServiceJoinReq(conn, clientAddr, msg)
	}else if msg.Action == dsgame.UpdateLocationAction {
		localClientInterfacing.ServiceUpdateLocationReq(conn, msg)
	}else if msg.Action == dsgame.FireAction {
		localClientInterfacing.ServiceFireReq(conn, msg)
  }else {
		fmt.Println("Message not understood: ")
  }
	
}
/*
	Print the current state of the game for this server
*/
func printState() {
	
	fmt.Println("Server game state")
	for key, value := range header.AgentsDB {
		fmt.Println("agent:", key, "time:", value.TimeStamp, " agent Location:", value.Location)	
 	}
}
