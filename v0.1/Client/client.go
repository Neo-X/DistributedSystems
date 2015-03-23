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
//	"time"
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

	fmt.Println("Game Started...")
	// simulating temporary sendUpdateLocation()	
	for {
		fmt.Println("Select Option:\n 1. Move\n 2. Fire")
		//option,_ := reader.ReadString()
		var option int
		fmt.Scanf("%d",&option)
		if option == 1 {
				fmt.Println("Enter new x,y,z:")
				var x,y,z float64
				fmt.Scanf("%f",&x)
				fmt.Scanf("%f",&y)
				fmt.Scanf("%f",&z)
				//x,_ := reader.ReadString('\n')
				//y,_ := reader.ReadString('\n')
				//z,_ := reader.ReadString('\n')
				
				// update location in local database
				//agent.Location[0],_ = strconv.ParseFloat(x,64)
				//agent.Location[1],_ = strconv.ParseFloat(y,64)
				//agent.Location[2],_ = strconv.ParseFloat(z,64)
				agent.Location[0] = x
				agent.Location[1] = y
				agent.Location[2] = z
				
				sendUpdateLocation(conn)

		} else if option == 2 {
				var target dsgame.FireTarget
				// add target value
				fmt.Println("Enter new m,c,quadrant:")
				var m,c,q float64
				fmt.Scanf("%f",&m)
				fmt.Scanf("%f",&c)
				fmt.Scanf("%f",&q)
				//m,_ := reader.ReadString('\n')
				//c,_ := reader.ReadString('\n')
				//q,_ := reader.ReadString('\n')
				
				//target.M,_ = strconv.ParseFloat(m,64)
				//target.C,_ = strconv.ParseFloat(c,64)
				//target.Q,_ = strconv.ParseFloat(q,64)
				target.M = m
				target.C = c
				target.Q = q
				fire(conn,target)
		}	else {
				fmt.Println("Please Retry.")
		}
		//time.Sleep(1 * time.Second)
	}

	// To update a location of an agent
	// 1. update location in agent.Location
	// 2. call sendUpdateLocation()

	conn.Close()
}


/***
*	Function Name: 	sendUpdateLocation()
*	Desc:			The function requests updateLocation
*	Pre-cond:		takes connection argument
*	Post-cond:		retuen success or return failure
*/
func sendUpdateLocation(conn *net.UDPConn){
		m := dsgame.Message{dsgame.UpdateLocationAction, client, agent.Name, simulationTime, agent.Location,dsgame.FireTarget{0.0,0.0,0.0}}
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
}


/***
*	Function Name: 	fire()
*	Desc:			The function requests fire
*	Pre-cond:		takes connection argument and name of client who is fired
*	Post-cond:		return success or returns failure
*/
func fire(conn * net.UDPConn, target dsgame.FireTarget){
		m := dsgame.Message{dsgame.FireAction, client, agent.Name, simulationTime, agent.Location, target}
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
}

/***
*	Function Name: 	join()
*	Desc:			The function requests to join the game
*/
func join( conn *net.UDPConn ){
	
	m := dsgame.Message{dsgame.JoinAction, "", "", 0, [3]float64{0.0,0.0,0.0}, dsgame.FireTarget{0.0,0.0,0.0}}
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
