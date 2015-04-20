/***
*	File Name: main.go
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
//	"encoding/json"
	"../dsgame"
//"../fixed"
	"../s3dm"
	"./dispatcher"	
//	"./listener"
	"./header"
	"flag"
)





/***
*	Function Name: 	main()
*	Desc:			The main function for client
*	Pre-cond:		
*	Post-cond:		Call the request functions

This is going to need to have two threads. One it is action on and another 
in the background listening for notofications from the server.
*/

func main(){
	//service := header.LocalServerIP_Port
	
	
	dsgame.GameMessageDeltaTime, _ = time.ParseDuration("300ms")
	/*
	if err != nil {
        fmt.Println("Error parsing time duration ", dsgame.GameDeltaTime)
        fmt.Println(err)
        return
    } 
	*/
	
	timePtr := flag.Int64("time", 0, "The initial time")
	logFilePtr := flag.String("logfile", "client0", "The log file for this node.")
	serverLinkPtr := flag.String("serverAddress", "127.0.0.1:10000" , "The ip address clients should use to connect to this service")
	
	flag.Parse()
	
    fmt.Println("time:", *timePtr)
    fmt.Println("logFile:", *logFilePtr)
    fmt.Println("clientLink ip:port:", *serverLinkPtr)

	server, err := net.ResolveUDPAddr("udp",*serverLinkPtr)
	conn, err := net.DialUDP("udp", nil, server)
	if err != nil {
        fmt.Println("Error connecting to " , *serverLinkPtr)
        fmt.Println(err)
        return
    }

	dispatcher.Join(conn)	

	fmt.Println("Game Started...")

	go dispatcher.UpdateServerFrame(conn)
	go RunAgent(conn)

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
				header.MyAgent.Location.X = x
				header.MyAgent.Location.Y = y
				header.MyAgent.Location.Z = z
				

		} else if option == 2 {
				var target s3dm.V3
				// add target value
				fmt.Println("Enter new m,c,quadrant:")
				var m,c,q float64
				fmt.Scanf("%f",&m)
				fmt.Scanf("%f",&c)
				fmt.Scanf("%f",&q)
				target.X = m
				target.Y = c
				target.Z = q
				dispatcher.Fire(conn,target)
		}	else {
				fmt.Println("Please Retry.")
		}
		// time.Sleep(1 * time.Second)
	}

	// To update a location of an agent
	// 1. update location in agent.Location
	// 2. call sendUpdateLocation()

	conn.Close()
}



func RunAgent(conn *net.UDPConn) {
	
	fireTime := dsgame.GetNextFireTime()
	for {

				
		header.MyAgent.Location = header.MyAgent.Location.Add(header.MyAgent.Direction.Muls(dsgame.GameDeltaTime))
		fmt.Println("Agent: moving", header.MyAgent.Location)
		if (header.MyAgent.Location.X > dsgame.GameUpperBound) || (header.MyAgent.Location.X < dsgame.GameLowerBound) {
			header.MyAgent.Direction.X = header.MyAgent.Direction.X*-1.0 
		}
		if (header.MyAgent.Location.Z > dsgame.GameUpperBound) || (header.MyAgent.Location.Z < dsgame.GameLowerBound) {
			header.MyAgent.Direction.Z = header.MyAgent.Direction.Z*-1.0 
		}
		if (header.MyAgent.Location.Y > dsgame.GameUpperBound) || (header.MyAgent.Location.Y < dsgame.GameLowerBound) {
			header.MyAgent.Direction.Y = header.MyAgent.Direction.Y*-1.0 
		}


		// Ramdom fire requests
		if (fireTime == 0) {
			direction := dsgame.GetRandomDirection()
			fmt.Println("Agent: shooting", direction) 
			dispatcher.Fire(conn,direction)
			fireTime = dsgame.GetNextFireTime() + 1
		}
		fireTime = fireTime - 1;
		time.Sleep(dsgame.GameMessageDeltaTime)
	}
	
}