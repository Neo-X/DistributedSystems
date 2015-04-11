package serversInterfacing


import (
	"../../govec"
	"fmt"
	"log"
//	"math/rand"
//	"net"
	"net/rpc"
//	"os"
	"strconv"
	"time"
	"strings"
	"../header"
	"encoding/json"
)

// args in get(args)
type GetArgs struct {
	Key    string // key to look up
	VStamp []byte // vstamp(nil)
}

// args in put(args)
type PutArgs struct {
	Key    string // key to associate value with
	Val    string // value
	VStamp []byte // vstamp(nil)
}

// args in testset(args)
type TestSetArgs struct {
	Key     string // key to test
	TestVal string // value to test against actual value
	NewVal  string // value to use if testval equals to actual value
	VStamp  []byte // vstamp(nil)
}

// Reply from service for all three API calls above.
type ValReply struct {
	Val    string // value; depends on the call
	VStamp []byte // vstamp(nil)
}

// Value in the key-val store.
type MapVal struct {
	value  string       // the underlying value representation
	logger *govec.GoLog // GoVector instance for the *key* that this value is mapped to
}

// Map implementing the key-value store.
var kvmap map[string]*MapVal
var port string // port I am running on

// Reserved value in the service that is used to indicate that the key
// is unavailable: used in return values to clients and internally.
const unavail string = "unavailable"




//*************************************
var mypriority, masterpriority, masterpriorityShouldbe int64
var counter int64
var masterCounter int64
var myid string
var logger (*govec.GoLog)
var c = make(chan int64)
var keyCount int64

func nodemain(ip_port string, logfile string){
//	logger = govec.Initialize("Log"+ os.Args[3], os.Args[3])
	logger = govec.Initialize("Log"+ logfile, logfile)
	client, err := rpc.Dial("tcp", ip_port )
	if err != nil {
		log.Fatal("dailing:",err)
	}
	counter = 1
	masterCounter =0
	masterpriorityShouldbe = 0
	//test if there is priority key entry in the table
	var rply_getPriority ValReply
	var arg GetArgs
	arg.Key = "priority"
	arg.VStamp = logger.PrepareSend("get:priority", nil)
	err = client.Call("KeyValService.Get",&arg, &rply_getPriority)
	if err != nil {
		log.Fatal("KeyValService.Get:", err.Error())
	}
	fmt.Println("Received priority: ", rply_getPriority.Val )
		//act as not-first node
		// (1) update the priority key
		mypriority,err = strconv.ParseInt(rply_getPriority.Val,10,32)
		priority_val := strconv.FormatInt(mypriority+1, 10)
		var args_setPriority PutArgs
		args_setPriority.Key = "priority"
		args_setPriority.Val = priority_val
		args_setPriority.VStamp = logger.PrepareSend("put:"+ priority_val, nil)
		var reply_setPriority ValReply
		err = client.Call("KeyValService.Put",&args_setPriority, &reply_setPriority)
		if err != nil {
			log.Fatal("KeyValService.Put:", err.Error())
		}
		/*
		if reply_setPriority.Val == ""{
			fmt.Println("updated priority key: success")
		}
		*/
		regularNode(client)
}


func regularNode(client * rpc.Client) {
	keyCount =0
	needElection := false
	tickChannel_updatetable := time.NewTicker(250 * time.Millisecond).C
	tickChannel_pollforActive := time.NewTicker(1000 * time.Millisecond).C
	for {
		select {
			case <- c:
					keyCount++
					fmt.Printf("KeyCount %d\n", keyCount)

			case <- tickChannel_updatetable:
				// (2) set up its own entry
				counter++
				//ownEntry_val := strconv.FormatInt(mypriority,10) + string(";") + strconv.FormatInt(counter, 10)  // add timestamp
				var d header.AgentDB
				d.Client = header.MyClientName
				d.Agent.Name = header.MyAgent.Name
				d.Agent.Location = header.MyAgent.Location
				b,err := json.Marshal(d)
					if err!= nil {
						fmt.Println("Problem marshalling struct")
						fmt.Println(err)
					}
				ownEntry_val := string(b) + string(";;;") + strconv.FormatInt(counter, 10)  // add timestamp

				var args2 PutArgs
				args2.Key = myid
				args2.Val = ownEntry_val
				args2.VStamp = logger.PrepareSend("put:"+ ownEntry_val, nil)
				var reply_ownEntry ValReply
				err = client.Call("KeyValService.Put",&args2, &reply_ownEntry)
				if err != nil {
					log.Fatal("KeyValService.Put:", err.Error())
				}
				/*
				if reply_ownEntry.Val == ""{
					fmt.Println("set up its own entry: success")
				}
				*/
				// (3) get and set entry in nodes
				var reply_getNodes ValReply
				var arg_getNodes GetArgs
				arg_getNodes.Key = "nodes"
				arg_getNodes.VStamp = logger.PrepareSend("get:nodes", nil)
				err = client.Call("KeyValService.Get",&arg_getNodes, &reply_getNodes)
				if err != nil {
					log.Fatal("KeyValService.Get:", err.Error())
				}
				//fmt.Println("Received nodes: ", reply_getNodes.Val);
				//check if entry is already there in nodes
				var nodes string
				nodes = ""
				found := false
				msgparts := strings.Split(reply_getNodes.Val,";;")	
				for i := range msgparts {
					keyVal := strings.Split(msgparts[i],";")	
					if keyVal[0] == myid  {
						found = true
						val,err := strconv.ParseInt(keyVal[1],10,64)
						if err != nil{
							log.Fatal("couldn't parse key val:", err.Error())
						}
						if val != mypriority {
							if nodes != "" {
								nodes = nodes + string(";;")
							}
							fmt.Printf("->KeyCount %d\n", keyCount)
							nodes = nodes + keyVal[0] +string(";")+ strconv.FormatInt(mypriority,10) + string(";") + strconv.FormatInt(keyCount,10)
							
						} else {
							if nodes != "" {
								nodes = nodes + string(";;")
							}
							//nodes = nodes + msgparts[i]
							nodes = nodes + keyVal[0] +string(";")+ strconv.FormatInt(mypriority,10) + string(";") + strconv.FormatInt(keyCount,10)
						}
					}else {
						if nodes != "" {
							nodes = nodes + string(";;")
						}
						nodes = nodes + msgparts[i]
					}
				}
				var setNodes_val string	
				if found == false {
					// (4) put its node in nodes
				fmt.Printf("-->KeyCount %d\n", keyCount)
					setNodes_val = reply_getNodes.Val + string(";;") + myid + string(";") + strconv.FormatInt(mypriority,10) + string(";")+ strconv.FormatInt(keyCount,10)
				} else {
					setNodes_val = nodes
				}
//				fmt.Println("setting nodes value to ", setNodes_val)
				var reply_setNodes ValReply
				var arg_setNodes PutArgs
				arg_setNodes.Key = "nodes"
				arg_setNodes.Val = setNodes_val
				arg_setNodes.VStamp = logger.PrepareSend("Put:" + setNodes_val ,nil)
				err = client.Call("KeyValService.Put",&arg_setNodes, &reply_setNodes)
				if err != nil {
					log.Fatal("KeyValService.Put:", err.Error())
				}
				/*
				if reply_setNodes.Val == "" {
					fmt.Println("updated nodes key: success")
				}
				*/
			case <- tickChannel_pollforActive:
				// (4) wait for (t+2) sec and get list of active nodes
				var reply_getActive ValReply
				var arg_getActive GetArgs
				arg_getActive.Key = "active"
				arg_getActive.VStamp = logger.PrepareSend("get:active nodes", nil)
				err := client.Call("KeyValService.Get",&arg_getActive, &reply_getActive)
				if err != nil {
					log.Fatal("KeyValService.Get:", err.Error())
				}
				fmt.Println("active Nodes:", reply_getActive.Val)

				msgpart := strings.Split(reply_getActive.Val, ";;;")
				activeNodes := strings.Split(msgpart[0], ";;")
				t := time.Now().Local()
				currentTime := t.Format("20060102150405")
				for i := range activeNodes {
					if i == 0 { 
						continue 
					}else{
						part := strings.Split(activeNodes[i], ";")
						port := part[0]
						if header.OnlineNodes[port] == ""{ // New Node came up
							if port != header.ServiceIP_Port {
								sendState(port)
							}
						}
						header.OnlineNodes[port] = currentTime //t.Format("20060102150405")
						updateClientAgentMap(client, port)
					}
				}	
								
				// print all online nodes
				fmt.Println("header.OnlineNodes: ")
				for k := range header.OnlineNodes {
					if currentTime != header.OnlineNodes[k] { 
						header.OnlineNodes[k] = ""
						fmt.Println("Offline" + "\t"+ k + "\t" + header.OnlineNodes[k])
					}else {
						fmt.Println("Online" + "\t" + k + "\t" + header.OnlineNodes[k] +"\t" + header.IpToAgentDB[k].Client + "\t" +header.IpToAgentDB[k].Agent.Name )
					}
				}

				if masterCounter == 0 {
					msgpart := strings.Split(reply_getActive.Val, ";;;")
					masterdetail := strings.Split(msgpart[1], ";")
					masterCounter,err = strconv.ParseInt(masterdetail[2],10,64)
					masterpriority, err = strconv.ParseInt(masterdetail[1],10,64)
					if err != nil {
						log.Fatal("couldnot parse mastercounter:", err.Error())
					}
				}else {
					msgpart := strings.Split(reply_getActive.Val, ";;;")
					masterdetail := strings.Split(msgpart[1], ";")
					masterpriority, err = strconv.ParseInt(masterdetail[1],10,64)
					mCounter,err := strconv.ParseInt(masterdetail[2],10,64)
					if err != nil {
						log.Fatal("couldnot parse mastercounter:", err.Error())
					}
					if mCounter <= masterCounter {
						// need elcetion 
						needElection = true
						fmt.Println("Need Election")
						break
					}
					masterCounter = mCounter
				}

		}
		if needElection == true {
			break
		}
	}
}

func updateClientAgentMap(client *rpc.Client, port string) {
			var reply_getDetail ValReply
			var arg_getDetail GetArgs
			arg_getDetail.Key = port
			arg_getDetail.VStamp = logger.PrepareSend("get:port details", nil)
			err := client.Call("KeyValService.Get",&arg_getDetail, &reply_getDetail)
				if err != nil {
					log.Fatal("KeyValService.Get:", err.Error())
				}
			if reply_getDetail.Val != "" {
				reply := reply_getDetail.Val
				var buf []byte = make([]byte,1500)
				ownVal := strings.Split(reply, ";;;")
				buf = []byte(ownVal[0])	
				//copy(buf[:],ownVal)
				var d header.AgentDB
				err = json.Unmarshal(buf[0:], &d)
					if err != nil {
						fmt.Println("Error Unmarshalling message")
						fmt.Println(err)
					}
				header.IpToAgentDB[port] = d
				
				header.ClientAgentMap[d.Client] = d.Agent.Name
//				header.Nodes[d.Client] = 
			}
}

func sendState(port string) {
	fmt.Println("Send State to Newly added Node: "+ port)
}



// Main server loop.
func Main(ip_port, ip_port_frontend, logfile string) {
	// parse args
	/*
	usage := fmt.Sprintf("Usage: %s <ip:port> <id> <logfileName>\n", os.Args[0])
	if len(os.Args) != 4 {
		fmt.Printf(usage)
		os.Exit(1)
	}
	*/
	//parts:= strings.Split(ip_port,":")
	//port = parts[1]
	
	/// use port as my id
	myid = ip_port


 	nodemain(ip_port_frontend  , logfile)

	// setup key-value store and register service
//	kvmap = make(map[string]*MapVal)
//	kvservice := new(BackEnd)
//	rpc.Register(kvservice)
//	l, e := net.Listen("tcp", ip_port)
//	if e != nil {
//		log.Fatal("listen error:", e)
//	}
//	for {
//		conn, _ := l.Accept()
//		go rpc.ServeConn(conn)
//	}
}




