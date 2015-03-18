/***
*	File Name: client.go
*	Description: The client provides the two functionalities:
*				 - request move
*				 - request fire
*	Date: 	8 March 2015
*	Author: Glen & Ravjot
*/
package main
import(
	"fmt"
	"net"
	"time"
	"encoding/json"
	"../dsgame"
)



/*
	global varaibles for client
*/

var client string // string to idetify this client
var agent dsgame.Agent // string to identify the agent for this client
// var conn *net.UDPConn // Connection to the GameServer
var simulationTime int64


/***
*	Function Name: 	main()
*	Desc:			The main function for client
*	Pre-cond:		
*	Post-cond:		Call the request functions
*/

func main(){
	service := "127.0.0.1:10000"
	server, err := net.ResolveUDPAddr("udp",service)
	conn, err := net.DialUDP("udp", nil, server)

	if err != nil {
        fmt.Println("Error connecting to " , service)
        fmt.Println(err)
        return
    }

	join(conn)	
	
	for n := int64(0); n >= 0; n++ {
		m := dsgame.Message{dsgame.UpdateLocationAction, client, agent.Name, simulationTime, agent.Location, ""}
		b, err := json.Marshal(m)
		if err != nil {
	        fmt.Println("Problem marshalling struct")
	        fmt.Println(err)
	    } 
		
		n, err :=	conn.Write(b)
	    if err != nil {
	        fmt.Println("WriteUDP")
	        fmt.Println(err)
	    } 
	
		var buf []byte = make([]byte, 1500) 
		n, _, err = conn.ReadFromUDP(buf) 
		// _, _, err = conn.ReadFromUDP(buf) 
	    if err != nil {
	        fmt.Println("ReadFromUDP")
	        fmt.Println(err)
	    } 
		// fmt.Println("Wrote ", n, "bytes")
		fmt.Println(string(buf[0:n]))
		time.Sleep(1 * time.Second)
		agent.Location[1] = agent.Location[1] + 1.0
	}

	conn.Close()
}


/***
*	Function Name: 	move()
*	Desc:			The function requests move
*	Pre-cond:		takes connection argument and the new location
*	Post-cond:		retuen success or return failure
*/
func move(){
}


/***
*	Function Name: 	fire()
*	Desc:			The function requests fire
*	Pre-cond:		takes connection argument and name of client who is fired
*	Post-cond:		return success or returns failure
*/
func fire(){
}

/***
*	Function Name: 	join()
*	Desc:			The function requests to join the game
*/
func join( conn *net.UDPConn ){
	
	m := dsgame.Message{dsgame.JoinAction, "", "", 0, [3]float64{0.0,0.0,0.0}, ""}
	b, err := json.Marshal(m)
	if err != nil {
        fmt.Println("Problem marshalling struct")
        fmt.Println(err)
    } 
	
	n, err := conn.Write(b)
	var buf []byte = make([]byte, 1500) 
	n, _, err = conn.ReadFromUDP(buf) 
	// _, _, err = conn.ReadFromUDP(buf) 
    if err != nil {
        fmt.Println("ReadFromUDP")
        fmt.Println(err)
    } 
	// fmt.Println("Wrote ", n, "bytes")
	fmt.Println(string(buf[0:n]))
	err = json.Unmarshal(buf[0:n], &m)
	if err != nil {
		fmt.Println("Error unmarshalling message")
        fmt.Println(err)
	}
	client = m.Client
	agent = dsgame.Agent{m.Agent, m.Location}
	simulationTime = m.TimeStamp
	
	fmt.Println("clientName: " ,  client)
	fmt.Println("agent: " ,  agent.Name, " location: ", agent.Location)
	fmt.Println("simulationTime: ", simulationTime)
		
	
}