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

import(
	"fmt"
	"../s3dm"
	"../fixed"
	"math/rand"
)

// exported variables must start with a capital letter....
const UpdateLocationAction string = "UpdateLocation"
const JoinAction string = "Join"
const FireAction string = "Fire"
const DestroyAction string = "Destroy"
const AcceptJointAction string = "AcceptJoin"

const FireDistance float64 = 7.0
const AgentRadius float64 = 0.5
const GameLowerBound float64 = -10.0
const GameUpperBound float64 = 10.0 

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
    Location s3dm.V3
    Target s3dm.V3
}

type Agent struct {
	Name string
	Location s3dm.V3
}

type Agents struct {
	//Agent string
	TimeStamp int64
	Location s3dm.V3 
}

/*
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
}*/
	
/*
func calculateLineParameters(x1,y1,x2,y2 int64) FireTarget{
	m := (y2- y1)/(x2 - x1)
	
}
*/

func RayHitsAgent(agentLoc, rayOrigin, rayDir s3dm.V3) bool {
	var _ray s3dm.Ray
	
	_ray.Origin = s3dm.Position{fixed.New(rayOrigin.X), fixed.New(rayOrigin.Y), fixed.New(rayOrigin.Z)}
	_ray.Dir = rayDir.Unit()
	var _sphere *s3dm.Sphere
	_sphere = s3dm.NewSphere(s3dm.Position{fixed.New(agentLoc.X), fixed.New(agentLoc.Y), fixed.New(agentLoc.Z)}, AgentRadius)
	
	_hit, _pos, _dir := _sphere.Intersect(&_ray)
	
	fmt.Println("Intersection at", _pos.V3() , " in direction ", _dir)
	return _hit
	
} 

var _rand *rand.Rand = rand.New(rand.NewSource(99))
func GetRandomLocation() s3dm.V3 {
	// _rand := rand.New(rand.NewSource(99))
	x := (_rand.Float64() * (GameUpperBound - GameLowerBound)) + GameLowerBound
	y := (_rand.Float64() * (GameUpperBound - GameLowerBound)) + GameLowerBound
	z := (_rand.Float64() * (GameUpperBound - GameLowerBound)) + GameLowerBound 
	
	return s3dm.V3{x,y,z}
}