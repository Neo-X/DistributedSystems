package header

import(
//	"fmt"
//	"flag"
	"net"
	"../../dsgame"
//	"../../s3dm"
	// "../fixed"
//	"encoding/json"
	// "sync"
//	"strconv"
)

var MyClientName string // string to idetify this client
var MyAgent dsgame.Agent // string to identify the agent for this client

var ClientAgentMap map[string]string // A map from node to the agent controlled by that node
var Nodes map[string]*net.UDPConn // A map of the node name to the client connection to that node
var ClientLink *net.UDPAddr // Comminucation link to the client for this node
var Connection *net.UDPConn // Comminucation connection used by this server to send messages

var AgentDB map[string]dsgame.Agent // Database of all agent details [agent -> <details>]
var ClientOffset int // to generate new client name and agent name for new joinee

var ServiceIP_Port string = "127.0.0.1:5000"
var CentralServerIP_Port string = "127.0.0.1:10000"
const KvService string = "127.0.0.1:12345"


/* Intermidiate state databases */

var OnlineNodes map[string]string // contains all online nodes

// I am not sure how these two items will be used [Glen]
type AgentDB_ struct {
	Client string
	Agent dsgame.Agent
}

var IpToAgentDB map[string]AgentDB_  // [IP:Port -- >  AgentDB]
