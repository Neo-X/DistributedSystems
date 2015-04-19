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
    time.Sleep(1 * time.Second)
  }
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
		//fmt.Println("Sending location update to: ", header.MyAgent.Location)
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
	
		var buf []byte = make([]byte, 1500) 
		//n, _, err = conn.ReadFromUDP(buf) 
		 _, _, err = conn.ReadFromUDP(buf) 
	    if err != nil {
	        fmt.Println("ReadFromUDP")
	        fmt.Println(err)
	    } 
		// fmt.Println("Wrote ", n, "bytes")
		//fmt.Println(string(buf[0:n]))
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
		
		n, err :=	conn.Write(b)
	    if err != nil {
	        fmt.Println("WriteUDP")
	        fmt.Println(err)
	    } 
	
		var buf []byte = make([]byte, 1500) 
		n, _, err = conn.ReadFromUDP(buf) 
		// _, _, err = conn.ReadFromUDP(buf) 
	    if err != nil {
	        fmt.Println("ReadFromUDP")
	        fmt.Println(err)
	    } 
		// fmt.Println("Wrote ", n, "bytes")
		fmt.Println(string(buf[0:n]))
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
	header.MyAgent = dsgame.Agent{m.Agent, m.Location}
	header.SimulationTime = m.TimeStamp
	
	fmt.Println("clientName: " ,  header.MyClientName)
	fmt.Println("agent: " ,  header.MyAgent.Name, " location: ", header.MyAgent.Location)
	fmt.Println("simulationTime: ", header.SimulationTime)
		
	
}
