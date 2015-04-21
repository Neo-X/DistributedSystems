package main

import (
	"fmt"
//	"net"
	"net/rpc"
	"flag"
//	"bufio"
	"os"
	"../govec"
// 	"./kvservice.GetArgs"
	"strconv"
	"time"
//	"math"
	"strings"
	kvservice "./kvservice"
)



var Server *rpc.Client
var Logger *govec.GoLog
/**
	Its not pretty but it does the job

*/
func main() {
	
	masterPtr := flag.Bool("m", false, "a bool")
	slavePtr := flag.Bool("s", false, "a bool")
	timePtr := flag.Int64("time", 0, "The initial time")
	deltaPtr := flag.Int64("d", 0, "delta time")
	slavesFilePtr := flag.String("slavesfile", "slavesfile.txt", "The filename for the slaves file.")
	logFilePtr := flag.String("logfile", "activityserver", "The log file for this node.")
	ipandportPtr := flag.String("address", "127.0.0.1:9999", "The op address and port for the node.")
	
	flag.Parse()
	_timeThreshold = (6 * time.Second).Nanoseconds()
	
	
	fmt.Println("master:", *masterPtr)
    fmt.Println("slave:", *slavePtr)
    fmt.Println("time:", *timePtr)
    fmt.Println("delta:", *deltaPtr)
    fmt.Println("slavesFile:", *slavesFilePtr)
    fmt.Println("logFile:", *logFilePtr)
    fmt.Println("ip:port:", *ipandportPtr)
    
    Logger = govec.Initialize(*logFilePtr, *logFilePtr)
    var err error
	Server, err = rpc.Dial("tcp", *ipandportPtr)
	if err != nil {
		fmt.Println("Dial failed:", err.Error())
		os.Exit(1)
	}
    
    Put(*logFilePtr, "127.0.0.1:9995")
    Member(*ipandportPtr, *timePtr, *logFilePtr)
    

}

func leader( _time int64, delta int64, address string, slavesFile string, logFile string ) {

    
}

func Get(key string) string {
	message := key// + strconv.FormatInt(init_time, 10)
	lMessage := (Logger).PrepareSend("getting value for key", []byte(message))
	var reply kvservice.ValReply	
	args := &kvservice.GetArgs{string(message), lMessage}
	
	err := Server.Call("KeyValService.Get", args, &reply)
	if err != nil {
		fmt.Println("Error running get on kvserver:", err.Error())
		os.Exit(1)
	}
	val := reply.Val
	// fmt.Printf("Key: %s, Vstamp: %s, val: %s\n", key, Logger.UnpackReceive("get(k:"+args.Key+")", reply.VStamp), val)
	return 	val
}

func Put(key string, value string) string {
	message := key// + strconv.FormatInt(init_time, 10)
	lMessage := (Logger).PrepareSend("putting value for key", []byte(message))
	var reply kvservice.ValReply	
	args := &kvservice.PutArgs{key, value, lMessage}
	
	err := Server.Call("KeyValService.Put", args, &reply)
	if err != nil {
		fmt.Println("Error running put on kvserver:", err.Error())
		os.Exit(1)
	}
	val := reply.Val
	// fmt.Printf("Key: %s, value: %s, Vstamp: %s, val: %s\n",key , reply.Val, Logger.UnpackReceive("get(k:"+args.Key+")", reply.VStamp), val)
	return 	val
}

func TestSet(key string, testVal string, setVal string) string {
	lMessage := (Logger).PrepareSend("testSet value for key", []byte(key))
	var reply kvservice.ValReply	
	args := &kvservice.TestSetArgs{key, testVal, setVal, lMessage}
	
	err := Server.Call("KeyValService.TestSet", args, &reply)
	if err != nil {
		fmt.Println("Error running testSet on kvserver:", err.Error())
		os.Exit(1)
	}
	val := reply.Val
	// fmt.Printf("Key: %s, value: %s, Vstamp: %s, val: %s\n",key , reply.Val, Logger.UnpackReceive("get(k:"+args.Key+")", reply.VStamp), val)
	return 	val
}

var leaderKey = "leader"
var membersKey = "members"
var newMembersKey = "memberUpdateQueue"
const ActivityServerKey = "activityserver"
var memDelimiter = ","
var keyValDelimiter = ":"
var members int64
// members = 0;
var _timeThreshold int64
var _leader_string string

func getLeader() string {
	return Get(leaderKey)
}

func setLeader(leader string, oldLeader string) string {
	return TestSet(leaderKey, oldLeader, leader)
}


func postMembers() { // post list of actime members
	// All members with timestamp < threshold
	_timeNow := time.Now().UnixNano()
	
	var activeMembers string
	for key, value := range _members {
    	fmt.Println("Key:", key, "Value:", value)
    	if ( (_timeNow - value) < _timeThreshold ) { // Last update less than threshold
    		activeMembers = activeMembers + key + memDelimiter 
    	}
	}
	Put(membersKey, activeMembers[:len(activeMembers)-1])
}

func GetMembers() string {
	return Get(membersKey)
}

func getNewMemberQueue() string { // gets and clears member queue
	newMembers := Get(newMembersKey)
	TestSet(newMembersKey, newMembers, "")
	return newMembers
}

func checkLeaderActive(leader string) bool {
	pair := strings.Split(getLeader(),  keyValDelimiter)
	_timeNow := time.Now().UnixNano()
	/*
	_timeStamp , err := strconv.ParseInt(pair[1], 10, 64)
	if (err != nil ) {
		fmt.Println("Error parsing int: ", err.Error()) 
	}*/
	
	if ( leader != _leader_string ) { // timestamp has been updated
		_members[pair[0]] = _timeNow
		_leader_string = leader
	}
	
	if ( ( (_timeNow - _members[pair[0]]) > _timeThreshold ) ) { // last timestamp greater then x distnace ago
		return false
	}		
	
	return true

}

func isLeader(id string) bool {
	_leader := getLeader()
	return strings.Contains(_leader, id)
}

func updateMember(memberID string) string {
	var val string
	val = " "
	var appendMembers string
	for ( val != appendMembers ) {
		newMembers := Get(newMembersKey)
		if (newMembers == "" ){
			appendMembers = memberID + keyValDelimiter + strconv.FormatInt(time.Now().UnixNano(), 10)
		} else {
			appendMembers = newMembers + memDelimiter + memberID + keyValDelimiter + strconv.FormatInt(time.Now().UnixNano(), 10)
		}
		val = TestSet(newMembersKey,newMembers, appendMembers)
	} 
	
	return val
}

func GetIDNumber(memberID string) string {
	newMembers := Get(newMembersKey)
	var appendMembers string
	if (newMembers == "" ){
		appendMembers = memberID + keyValDelimiter + strconv.FormatInt(time.Now().UnixNano(), 10)
	} else {
		appendMembers = newMembers + memDelimiter + memberID + keyValDelimiter + strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	val := TestSet(newMembersKey,newMembers, appendMembers) 
	
	return val
}

var _members map[string] int64 // stores members last timestamps

func updateMemberTimeMap(queue string) {
	pairs := strings.Split(queue,  ",")
	fmt.Println("Member Queue: " + queue)
	for i := 0; i < len(pairs); i++ {
		parts := strings.Split(pairs[i], keyValDelimiter)
		_member := parts[0]
		// timeStamp := parts[1]
		// fmt.Println("id: " + _member + ", timeStamp: " + timeStamp)
		/*_timeStamp , err := strconv.ParseInt(timeStamp, 10, 64)
		if (err != nil ) {
			fmt.Println("Error parsing int: ", err.Error()) 
		}*/
		_members[_member] = time.Now().UnixNano()
		// _members[_member] = _timeStamp
	}
}

func Member(address string, _time int64, logFile string) {
	// init_time := time.Now().UnixNano()
	id := logFile
	_members = make(map[string] int64)
	
	/**
	val := get("hello")
	strconv.FormatInt(init_time, 10)
	
	fmt.Println("Val: ", val)
	val = put("hello", "howdy do")
	
	fmt.Println("put Val: ", val)
	val = get("hello")
	
	fmt.Println("Val: ", val)
	val = testSet("hello", "howdy do", "hi")
	fmt.Println("Val: ", val)
	*/
	for i := 0; i >= 0; i++ {
		// always update timestamp for member
		updateMember(id)		
		leader := getLeader()
		if (leader == "")  { // no leader yet
			// Set self to leader
			_timeNow := strconv.FormatInt(time.Now().UnixNano(), 10)
			setLeader(id+keyValDelimiter+_timeNow, leader)			
		} else if ( isLeader(id) )  { // I am the leader
    		// Check to see if members have timed out and empty member Queue
    		_timeNow := strconv.FormatInt(time.Now().UnixNano(), 10)
			setLeader(id+keyValDelimiter+_timeNow, leader)
    		membersList := getNewMemberQueue()
    		updateMemberTimeMap(membersList)
			    		
    		// post list of members
    		postMembers()
    		time.Sleep(1 * time.Second)
    	} else { // I am member
    		// update newMembers with timestamp
    		// updateMember(id)
    		if ( !checkLeaderActive(leader) ) {
    			_timeNow := strconv.FormatInt(time.Now().UnixNano(), 10)
    			setLeader(id+keyValDelimiter+_timeNow, leader)
    		}
    		time.Sleep(2 * time.Second)
    	}
    	
		fmt.Println("\n Leader: ", leader)
		fmt.Println("active Members: ", GetMembers())
		fmt.Println("MembersUpdateQueue: ", Get(newMembersKey))
	}
	
	Server.Close()


}
