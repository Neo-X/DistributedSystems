package header

import(
	"../../dsgame"
	)


/*
	global varaibles for client
*/

var MyClientName string // string to idetify this client
var MyAgent dsgame.Agent // string to identify the agent for this client
// var conn *net.UDPConn // Connection to the GameServer
var SimulationTime int64

var ServiceIP_Port string = "127.0.0.1:4000"
var LocalServerIP_Port string = "127.0.0.1:5000"
