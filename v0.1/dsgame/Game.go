/***
*	File Name: Game.go
*	Description: This file encapsulates the function of the GAME
*				 
	Effectively keeps track of state changes in the game and provides game 
	functionality to server and client
	
*	Date: 	16 March 2015
*	Author: Glen & Ravjot
*/

package dsgame


// exported variables must start with a capital letter....
const UpdateLocationAction string = "UpdateLocation"
const JoinAction string = "Join"
const FireAction string = "Fire"
const DestroyAction string = "Destroy"
const AcceptJointAction string = "AcceptJoin"

/*
	Kind of hacky way to handle any kind of message
*/
type Message struct {
	Action string
    Client string
    Agent string
    TimeStamp int64
    Location [3]float64
    Target string
}

type Agent struct {
	Name string
	Location [3]float64
}
	
