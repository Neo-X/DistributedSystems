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
var ClientAgentMap map[string]string // A map from client to the agent controlled by that client
var Nodes map[string]*net.UDPConn // A map of the node name to the client connection to that node

var AgentsDB map[string]dsgame.Agents // Database of all agent details [client -> <details>]
var ClientOffset int // to generate new client name and agent name for new joinee

var ServiceIP_Port string = "127.0.0.1:5000"
var CentralServerIP_Port string = "127.0.0.1:10000"
const KvService string = "127.0.0.1:12345"

var OnlineNodes map[string]string // contains all online nodes
