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

const FirstQuadrant int64 = 1
const SecondQuadrant int64 = 2
const ThirdQuadrant int64 = 3
const FourthQuadrant int64 = 4

/*
	Kind of hacky way to handle any kind of message
*/
type Message struct {
	Action string
    Client string
    Agent string
    TimeStamp int64
    Location [3]float64
    Target FireTarget
}

type Agent struct {
	Name string
	Location [3]float64
}

type Agents struct {
	//Agent string
	TimeStamp int64
	Location [3]float64
}

type FireTarget struct { // Equation of line in 2D y = mx + c
	M float64 // slope of the line
	C float64 // constant 
	Q float64 // to define the direction of the fire
								// for an agent at (x1,y1), we place the coordinate system (as shown below) 
								// on the point and define the direction according to the quadrant division
								//						II	| I
								//					-------------
								//						III	| IV
								//
}
	
/*
func calculateLineParameters(x1,y1,x2,y2 int64) FireTarget{
	m := (y2- y1)/(x2 - x1)
	
}
*/
