// Version 1.0
//
// A simple key-value store that supports three API calls over rpc:
// - get(key)
// - put(key,val)
// - testset(key,testval,newval)
//
// Usage: go run kvservicemainl.go ip:port key-fail-prob
// - ip:port : the ip and port on which the service will listen for connections
// - key-fail-prob : probability in range [0,1] of the key becoming
//   unavailable during one of the above operations (permanent key unavailability)
//
// TODOs:
// - needs some serious refactoring
// - simulate netw. partitioning failures
//
// Dependencies:
// - GoVector: https://github.com/arcaneiceman/GoVector

package main

import (
	"../govec"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"
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

// Reserved value in the service that is used to indicate that the key
// is unavailable: used in return values to clients and internally.
const unavail string = "unavailable"

type KeyValService int

// Lookup a key, and if it's used for the first time, then initialize its value.
func lookupKey(key string) *MapVal {
	// lookup key in store
	val := kvmap[key]
	if val == nil {
		// key used for the first time: create and initialize a MapVal instance to associate with a key
		val = &MapVal{
			value:  "",
			logger: govec.Initialize("key-"+key, "key-"+key),
		}
		kvmap[key] = val
	}
	return val
}

// The probability with which a key operation triggers permanent key unavailability.
var failProb float64

// Check whether a key should fail with independent fail probability.
func checkFail(val *MapVal) bool {
	if val.value == unavail {
		return true
	}
	if rand.Float64() < failProb {
		val.value = unavail // permanent unavailability
		return true
	}
	return false
}

// GET
func (kvs *KeyValService) Get(args *GetArgs, reply *ValReply) error {
	val := lookupKey(args.Key)
	val.logger.UnpackReceive("get(k:"+args.Key+")", args.VStamp)

	if checkFail(val) {
		reply.Val = unavail
		reply.VStamp = val.logger.PrepareSend("get-re:"+unavail, nil)
		return nil
	}

	reply.Val = val.value // execute the get
	reply.VStamp = val.logger.PrepareSend("get-re:"+val.value, nil)
	return nil
}

// PUT
func (kvs *KeyValService) Put(args *PutArgs, reply *ValReply) error {
	val := lookupKey(args.Key)
	val.logger.UnpackReceive("put(k:"+args.Key+",v:"+args.Val+")", args.VStamp)

	if checkFail(val) {
		reply.Val = unavail
		reply.VStamp = val.logger.PrepareSend("put-re:"+unavail, nil)
		return nil
	}

	val.value = args.Val // execute the put
	reply.Val = ""
	reply.VStamp = val.logger.PrepareSend("put-re", nil)
	return nil
}

// TESTSET
func (kvs *KeyValService) TestSet(args *TestSetArgs, reply *ValReply) error {
	val := lookupKey(args.Key)
	val.logger.UnpackReceive("testset(k:"+args.Key+",tv:"+args.TestVal+",nv:"+args.NewVal+")", args.VStamp)

	if checkFail(val) {
		reply.Val = unavail
		reply.VStamp = val.logger.PrepareSend("testset-re:"+unavail, nil)
		return nil
	}

	// execute the testset
	if val.value == args.TestVal {
		val.value = args.NewVal
	}

	reply.Val = val.value
	reply.VStamp = val.logger.PrepareSend("testset-re:"+val.value, nil)
	return nil
}

// Main server loop.
func main() {
	// parse args
	usage := fmt.Sprintf("Usage: %s ip:port key-fail-prob\n", os.Args[0])
	if len(os.Args) != 3 {
		fmt.Printf(usage)
		os.Exit(1)
	}

	ip_port := os.Args[1]
	arg, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if arg < 0 || arg > 1 {
		fmt.Printf(usage)
		fmt.Printf("\tkey-fail-prob arg must be in range [0,1]\n")
		os.Exit(1)
	}
	failProb = arg

	// setup randomization
	rand.Seed(time.Now().UnixNano())

	// setup key-value store and register service
	kvmap = make(map[string]*MapVal)
	kvservice := new(KeyValService)
	rpc.Register(kvservice)
	l, e := net.Listen("tcp", ip_port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	for {
		conn, _ := l.Accept()
		go rpc.ServeConn(conn)
	}
}

