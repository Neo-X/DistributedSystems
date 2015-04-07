package kvmanager


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

// Local store to keep track of keys
var localmap map[string]string

const unavail string = "unavailable"


//************************************ 
var mypriority, masterpriority, masterpriorityShouldbe int64
var counter int64
var masterCounter int64
var myid string
var logger (*govec.GoLog)
func nodemain(ip_port string, logfile string){
	logger = govec.Initialize("Log"+ logfile, logfile)
	client, err := rpc.Dial("tcp", ip_port )
	if err != nil {
		log.Fatal("dailing:",err)
	}
	counter = 1
	masterCounter =0
	masterpriorityShouldbe = 0
	//test if there is priority key entry in the table
		//implies that it is the first node, so set up the table

		// (1) set up priority key
		priority_val := "2"
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
		if reply_setPriority.Val == "" {
			fmt.Println("Setting priority: success"  )
		}
		*/
		// (2) set up its own entry
		mypriority = 1
		ownEntry_val := strconv.FormatInt(mypriority, 10)+ string(";;;")+ strconv.FormatInt(counter,10) // add timestamp
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
		if reply_ownEntry.Val == "" {
			fmt.Println("Setting own entry: success"  )
		}
		*/

		// (3) set up nodes key
		setNode_val := myid + string(";1")
		var args_setNodes PutArgs
		args_setNodes.Key = "nodes"
		args_setNodes.Val = setNode_val
		args_setNodes.VStamp = logger.PrepareSend("put:"+ setNode_val, nil)
		var reply_setNodes ValReply
		err = client.Call("KeyValService.Put",&args_setNodes, &reply_setNodes)
		if err != nil {
			log.Fatal("KeyValService.Put:", err.Error())
		}
		/*
		if reply_setNodes.Val == "" {
			fmt.Println("Setting up value in nodes: success"  )
		}
		*/
		domainMaster(client)
}

func domainMaster(client *rpc.Client){
	// (4) Accept the role of leader, that is start updating the active key		

	// create a local table for regular node database
	var idtopreviousCounter map[string]int64
	idtopreviousCounter = make(map[string]int64)	

	tickChannel := time.NewTicker(time.Millisecond * 500).C
	for {
		select {
			case <- tickChannel:
				counter++
				// (4.1) get and set entry in nodes
				var reply_getNodes ValReply
				var arg_getNodes GetArgs
				arg_getNodes.Key = "nodes"
				arg_getNodes.VStamp = logger.PrepareSend("get:nodes", nil)
				err := client.Call("KeyValService.Get",&arg_getNodes, &reply_getNodes)
				if err != nil {
					log.Fatal("KeyValService.Get:", err.Error())
				}
				fmt.Println("Received nodes: ", reply_getNodes.Val)
				//check if entry is already there in nodes
				//found := false
				var active string
				active = myid + string(";") + strconv.FormatInt(mypriority,10)
				msgparts := strings.Split(reply_getNodes.Val,";;")	
				for i := range msgparts {
					keyVal := strings.Split(msgparts[i],";")	
					// check if the node is still alive or not
					var reply_getNodes ValReply
					var arg_getNodes GetArgs
					arg_getNodes.Key = keyVal[0]
					arg_getNodes.VStamp = logger.PrepareSend("get: "+ keyVal[0], nil)
					err := client.Call("KeyValService.Get",&arg_getNodes, &reply_getNodes)
					if err != nil {
						log.Fatal("KeyValService.Get:", err.Error())
					}
					strcounter := strings.Split(reply_getNodes.Val,";;;")
					counter,err := strconv.ParseInt(strcounter[1],10,64)
					if err != nil {
						log.Fatal("couldn't parse from string to int:", err.Error())
					}
					if idtopreviousCounter[keyVal[0]] == 0 {
						idtopreviousCounter[keyVal[0]] = counter
					} else {
						if idtopreviousCounter[keyVal[0]] < counter {
							// it alive, append it
							if active != "" {
								active = active + string(";;")
							}
							active = active + keyVal[0] + string(";") + keyVal[1] + string(";") + keyVal[2]
							idtopreviousCounter[keyVal[0]] = counter
							fmt.Println(strcounter[0])
						} else {
							idtopreviousCounter[keyVal[0]] = 0
						}
					}
					// append alive in active
				}
				fmt.Println("active processes: ", active)
				// put active key in local map
				localmap["active"] = active	

				if active != "" {
					 //update the active
					setActive_val := active + string(";;;") + myid +string(";") +strconv.FormatInt(mypriority,10) +(";")+strconv.FormatInt(counter,10) //actually it should be active
					var arg_setActive PutArgs
					arg_setActive.Key = "active"
					arg_setActive.Val = setActive_val
					arg_setActive.VStamp = logger.PrepareSend("Put:" + setActive_val ,nil)
					var reply_setActive ValReply
					err = client.Call("KeyValService.Put",&arg_setActive, &reply_setActive)
					if err != nil {
						log.Fatal("KeyValService.Put:", err.Error())
					}
/*					// update nodes
					arg_setActive.Key = "nodes"
					arg_setActive.VStamp = logger.PrepareSend("Put:" + setActive_val ,nil)
					err = client.Call("KeyValService.Put",&arg_setActive, &reply_setActive)
					if err != nil {
						log.Fatal("KeyValService.Put:", err.Error())
					}*/
				}
		}
	}	
}


//************************************

/**
*	Function Name: main
*	Desc : start function
*/

func Main(ip_port_kvInterface, logfileName string) {
	// parse args
	/*
	usage := fmt.Sprintf("Usage: %s <kv-interfacing ip:port> <client-interfacing ip:port> <logfile> \n", os.Args[0])
	if len(os.Args) != 4 {
		fmt.Printf(usage)
		os.Exit(1)
	}

	ip_port_kvInterface := os.Args[1]

	ip_port_clientInterface := os.Args[2]
	
	logfileName = os.Args[3]
	*/
	myid = "1"

	localmap = make(map[string]string)
	//go frontEndInterface(ip_port_kvInterface)

	//go kvInterface(ip_port_clientInterface)
	//time.Sleep(1 * time.Second)	
	nodemain(ip_port_kvInterface, logfileName)

}
