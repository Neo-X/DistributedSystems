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
	"time"
	"net"
	"net/rpc"
	"os"
	"../dsgame"
//	"../s3dm"
	// "../fixed"
	"encoding/json"
	// "sync"
	"strconv"
	"./header"
	"./localClientInterfacing"
	// "./serversInterfacing"
	activityserver "../activityserver"
	// "os"
	"../govec"
	"strings"
)


var ServerID string
/***
*	Function Name: 	main()
*	Desc:			The main function for server
*	Pre-cond:		
*	Post-cond:		Call the service functions
*/
func main(){
	
	
	timePtr := flag.Int64("time", 0, "The initial time")
	slavesFilePtr := flag.String("slavesfile", "slavesfile.txt", "The filename for the slaves file.")
	logFilePtr := flag.String("logfile", "server", "The log file for this node.")
	idPtr := flag.Uint64("ID", 0, "The id for this node")
	ipandportPtr := flag.String("kvservice", "127.0.0.1:9999", "The ip address and port for the kvservice.")
	clientLinkPtr := flag.String("clientAddress", "127.0.0.1:10000" , "The ip address clients should use to connect to this service")
	
	flag.Parse()
	
	
	
    fmt.Println("time:", *timePtr)
    fmt.Println("slavesFile:", *slavesFilePtr)
    fmt.Println("logFile:", *logFilePtr)
    fmt.Println("ID:", *idPtr)
    fmt.Println("kvservice ip:port:", *ipandportPtr)
    fmt.Println("clientLink ip:port:", *clientLinkPtr)
    
    *logFilePtr = *logFilePtr + strconv.FormatUint(*idPtr, 10)
    ServerID = *logFilePtr
    logger := govec.Initialize(*logFilePtr, *logFilePtr)
    // var err error
	server, err := rpc.Dial("tcp", *ipandportPtr)
	if err != nil {
		fmt.Println("Dial failed:", err.Error())
		os.Exit(1)
	}
	activityserver.Server = server
	activityserver.Logger = logger
	
    go activityserver.Member(*ipandportPtr, *timePtr, *logFilePtr)
    activityserver.Put(*logFilePtr, *clientLinkPtr)
    go CheckForNewNodes()	
	
	// go serversInterfacing.Main(os.Args[1], header.KvService, os.Args[4])
	
	
    
  header.Nodes = make(map[string]*net.UDPConn)
  header.ClientAgentMap = make(map[string]string)
  header.AgentDB = make(map[string]dsgame.Agent)
  header.OnlineNodes = make(map[string]string)
  header.IpToAgentDB = make(map[string]header.AgentDB_)
    
  udpAddress, err := net.ResolveUDPAddr("udp4",*clientLinkPtr)

  if err != nil {
		fmt.Println("error resolving UDP address on ", *clientLinkPtr)
		fmt.Println(err)
		return
  }

  conn ,err := net.ListenUDP("udp",udpAddress)

	if err != nil {
		fmt.Println("error listening on UDP port ", *clientLinkPtr)
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
    	// fmt.Println("got message from ", address, " with n = ", n)

      if n > 0 {
      	// fmt.Println("from address", address, "got message:", string(buf[0:n]))
        ////// Everything should be good now
        handleMessage(conn , address, buf[0:n], *idPtr)
     		printState()       	
      }

      /* conn, err := net.DialUDP("udp", nil, address)
         if err != nil {
           	fmt.Println("Error connecting to UDP client")
		        fmt.Println(err)
         }*/
/*
		n, err :=	conn.WriteToUDP([]byte("Thank you for your message"), address)

		if err != nil {
		  	fmt.Println("WriteUDP Message", n)
		    fmt.Println(err)
		} 
*/

		}
	}


}

/*
	This is the function used for the node to poll for updated nodes. 
	If a new node is found it is added to the header.Nodes.
	If a node in header.Nodes does not appear in the active node, remove it from header.Nodes.
*/
func CheckForNewNodes() {
	
	for {
		
		activeMembers := activityserver.GetMembers()
		// fmt.Println("active members", activeMembers)
		if (activeMembers != "" ) { // THere should be other active members
			members := strings.Split(activeMembers,  ",")
			// add any new members
			for i := 0; i < len(members); i++ {
				/*
				if (members[i] == activityserver.ActivityServerKey) { // skip activity server
					continue
				}*/
				if _, ok := header.Nodes[members[i]]; ok {
				    // do nothing if it already exsists
				} else {
					// need to create the node and add it to header.Nodes
					address := activityserver.Get(members[i])
					if address == "" {
				        fmt.Println("Error getting address for node" , members[i])
				        // fmt.Println(err)
				        return
				    }
					
					server, err := net.ResolveUDPAddr("udp",address)
					conn, err := net.DialUDP("udp", nil, server)
					if err != nil {
				        fmt.Println("Error connecting to " , address)
				        fmt.Println(err)
				        return
				    }
					header.Nodes[members[i]] = conn 
					
					
				}
					
			}
			
			// Remove stale members
			for key, value := range header.Nodes {
				fmt.Println("node:", key, "ip:", value.RemoteAddr().String() )	
				// check that this key is in active members
				if ( strings.Contains(activeMembers, key) ) { // Not very fast...
					continue
				} else {
					// remove from nodes
					removeNode(key)
					
				}
				
		 	}
			
		}
		time.Sleep(1*time.Second)
		// printState()	
	}
	
	
}


/***
*	Function Name: 	removeNode()
*	Desc:			This function removes all of the data pertaining to the specified node
*	Pre-cond:		key to the node
*	Post-cond:		node should be removed from Nodes and all data for the agent related to that node should be removed
*/
func removeNode(key string) {
	fmt.Println("***** Removing inactive node from system:", key, "*******")
	// remove node
	delete(header.Nodes, key)
	
	// remove agent data
	agentName := header.ClientAgentMap[key]
	delete(header.AgentDB, agentName)	
	
}

/***
*	Function Name: 	serviceJoinReq()
*	Desc:			The function provide service to a client requesting to join the game
*	Pre-cond:		takes connection argument ....
*	Post-cond:		Client should be registard in the game
*/
func handleMessage(conn *net.UDPConn, clientAddr *net.UDPAddr, buf []byte, id uint64){
	var msg dsgame.Message
	err := json.Unmarshal(buf, &msg)
	if err != nil {
		fmt.Println("Error handling message")
		fmt.Println(err)
		return
	}
	
	if msg.Action == dsgame.JoinAction {
		localClientInterfacing.ServiceJoinReq(conn, clientAddr, msg, id)
	}else if msg.Action == dsgame.UpdateLocationAction {
		if ( localClientInterfacing.ServiceUpdateLocationReq(conn, msg) ) {
			// broadcast succesful move to other nodes if Agent is this nodes agent
			if ( msg.Agent == header.MyAgent.Name) {
				BroadcastClientLocationUpdate(msg)
			}
		}
	}else if msg.Action == dsgame.FireAction {
		localClientInterfacing.ServiceFireReq(conn, msg)
	} else if msg.Action == dsgame.DestroyAction {
		localClientInterfacing.HandleDestroyReq(msg)
  }else {
		fmt.Println("Message not understood: ")
  }
	
}

func BroadcastClientLocationUpdate(msg dsgame.Message) {
	
	for key, conn := range header.Nodes {
		if (key != ServerID) { // don't send message to self
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
		
		}	
 	}
	
}


/*
	Print the current state of the game for this server
*/
func printState() {
	
	fmt.Println("Game state")
	/*
	for key, value := range header.AgentsDB {
		fmt.Println("client:", key, "time:", value.TimeStamp, " agent Location:", value.Location)	
 	}
 	*/
	for key, value := range header.AgentDB {
		fmt.Println("agent:", key, " agent Location:", value.Location)	
 	}
	for key, value := range header.Nodes {
		fmt.Println("node:", key, "ip:", value.RemoteAddr().String() )	
 	}
}
