package dispatcher

import(
	"fmt"
	"net"
	"time"
	"encoding/json"
	"../../dsgame"
//"../../fixed"
	"../../s3dm"
	"../header"
)

/***
* Func Name: updateServerFrame()
* Desc:
* Pre-cond:
* Post-cond:
**/
  
func UpdateServerFrame(conn *net.UDPConn){ 
  for {
    sendUpdateLocation(conn)
    // time.Sleep(0.2 * time.Second)
    time.Sleep(dsgame.GameMessageDeltaTime)
    // header.PrintState()
  }
}

/***
*	Function Name: 	listenForMessages
*	Desc:			This function listens on the connection for incoming messages
*	Pre-cond:		takes connection argument
*	Post-cond:		Runs until the end of the game
*/
func ListenForMessages(conn *net.UDPConn){ 
  for {
    var buf []byte = make([]byte, 1500) 
	n, _, err := conn.ReadFromUDP(buf) 
	// _, _, err = conn.ReadFromUDP(buf) 
    if err != nil {
        fmt.Println("ReadFromUDP")
        fmt.Println(err)
    } 
	// fmt.Println("Wrote ", n, "bytes")
	// fmt.Println(string(buf[0:n]))
	var msg dsgame.Message
	err = json.Unmarshal(buf[0:n], &msg)
	if err != nil {
		fmt.Println("Error handling message")
		fmt.Println(err)
		return
	}
	processMessage(msg)
  }
}

/***
*	Function Name: 	processMessage
*	Desc:			This function will process a message
*	Pre-cond:		takes a msg to be processed
*	Post-cond:		Depends on the message but the message should be handled properly.
*/
func processMessage(msg dsgame.Message) {

	if ( msg.Action == dsgame.PositionOverrideAction ) {
		// override current agent location
		fmt.Println("Overriding agent location:", msg.Location)
		header.MyAgent.Location = msg.Location
		tmp := header.AgentDB[header.MyAgent.Name]
		tmp.Location = msg.Location
		header.AgentDB[header.MyAgent.Name] = tmp
	} else if ( msg.Action == dsgame.UpdateLocationAction ) {
		ServiceUpdateLocationReq(msg)
	} else if (msg.Action == dsgame.DestroyAction) {
		tmp := header.AgentDB[msg.Agent]
		tmp.Location = msg.Location
		header.AgentDB[msg.Agent] = tmp
		fmt.Println(msg.Agent, " is destroyed !!!!!!!")
	}

}


/***
*	Function Name: 	ServiceUpdateLocationReq()
*	Desc:			The function provide service client move requests
*	Pre-cond:		takes connection argument and the new location
*	Post-cond:		list is updated with new location or return failure
Also return failure on this being an invalid update??
*/
func ServiceUpdateLocationReq(msg dsgame.Message) bool {
	// update server agent database
	_timeNow := time.Now().UnixNano()

	
	var tmpObj dsgame.Agent // This covers the case when the agent has not been initialized yet, but I don't think we want this
	tmpObj.TimeStamp = msg.TimeStamp
	tmpObj.Location = s3dm.V3{msg.Location.X,msg.Location.Y,msg.Location.Z}
	tmpObj.LastUpdateTime = _timeNow
	header.AgentDB[msg.Agent] = tmpObj
	// fmt.Println("Location updated by:" + msg.Client)
	// SendPositionforAgent(msg)
	return true
}



/***
*	Function Name: 	sendUpdateLocation()
*	Desc:			The function requests updateLocation
*	Pre-cond:		takes connection argument
*	Post-cond:		retuen success or return failure
*/
func sendUpdateLocation(conn *net.UDPConn) {
		var _dir s3dm.V3
		//_ray_ray.Origin = s3dm.Position{fixed.New(0.0),fixed.New(0.0),fixed.New(0.0)}
		// fmt.Println("Sending location update to: ", header.MyAgent.Location)
		m := dsgame.Message{dsgame.UpdateLocationAction, header.MyClientName, header.MyAgent.Name, header.SimulationTime, header.MyAgent.Location,_dir}
		b, err := json.Marshal(m)
		if err != nil {
	        fmt.Println("Problem marshalling struct")
	        fmt.Println(err)
	    } 
		
		_, err =	conn.Write(b)
	    if err != nil {
	        fmt.Println("WriteUDP")
	        fmt.Println(err)
	    } 
	
/*
		var buf []byte = make([]byte, 1500) 
		//n, _, err = conn.ReadFromUDP(buf) 
		 _, _, err = conn.ReadFromUDP(buf) 
	    if err != nil {
	        fmt.Println("ReadFromUDP")
	        fmt.Println(err)
	    } 
		// fmt.Println("Wrote ", n, "bytes")
		//fmt.Println(string(buf[0:n]))
*/
}


/***
*	Function Name: 	fire()
*	Desc:			The function requests fire
*	Pre-cond:		takes connection argument and name of client who is fired
*	Post-cond:		return success or returns failure
*/
func Fire(conn * net.UDPConn, _dir s3dm.V3){
		m := dsgame.Message{dsgame.FireAction, header.MyClientName, header.MyAgent.Name, header.SimulationTime, header.MyAgent.Location, _dir}
		b, err := json.Marshal(m)
		if err != nil {
	        fmt.Println("Problem marshalling struct")
	        fmt.Println(err)
	    } 
		
		_, err =	conn.Write(b)
	    if err != nil {
	        fmt.Println("WriteUDP")
	        fmt.Println(err)
	    } 
	/*
		var buf []byte = make([]byte, 1500) 
		n, _, err = conn.ReadFromUDP(buf) 
		// _, _, err = conn.ReadFromUDP(buf) 
	    if err != nil {
	        fmt.Println("ReadFromUDP")
	        fmt.Println(err)
	    } 
		// fmt.Println("Wrote ", n, "bytes")
		fmt.Println(string(buf[0:n]))
*/
}

/***
*	Function Name: 	join()
*	Desc:			The function requests to join the game
*/
func Join( conn *net.UDPConn ){
	
	var m dsgame.Message
	m.Action = dsgame.JoinAction //, "", "", 0, [3]float64{0.0,0.0,0.0}, dsgame.FireTarget{0.0,0.0,0.0}}
	b, err := json.Marshal(m)
	if err != nil {
        fmt.Println("Problem marshalling struct")
        fmt.Println(err)
    } 
	
	n, err := conn.Write(b)
	var buf []byte = make([]byte, 1500) 
	n, _, err = conn.ReadFromUDP(buf) 
	// _, _, err = conn.ReadFromUDP(buf) 
    if err != nil {
        fmt.Println("ClientReadFromUDP")
        fmt.Println(err)
    } 
	// fmt.Println("Wrote ", n, "bytes")
	fmt.Println(string(buf[0:n]))
	err = json.Unmarshal(buf[0:n], &m)
	if err != nil {
		fmt.Println("Error unmarshalling message")
        fmt.Println(err)
	}
	header.MyClientName = m.Client
	header.MyAgent = dsgame.Agent{m.Agent, m.Location, 0, time.Now().UnixNano(), dsgame.GetRandomDirection()}
	header.SimulationTime = m.TimeStamp
	header.AgentDB[header.MyAgent.Name] = header.MyAgent
	
	fmt.Println("clientName: " ,  header.MyClientName)
	fmt.Println("agent: " ,  header.MyAgent.Name, " location: ", header.MyAgent.Location)
	fmt.Println("simulationTime: ", header.SimulationTime)
		
	
}
