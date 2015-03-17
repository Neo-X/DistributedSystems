/***
*	File Name: client.go
*	Description: The client provides the two functionalities:
*				 - request move
*				 - request fire
*	Date: 	8 March 2015
*	Author: Glen & Ravjot
*/
package main
import(
	"fmt"
	"net"
	"time"
)


/***
*	Function Name: 	main()
*	Desc:			The main function for client
*	Pre-cond:		
*	Post-cond:		Call the request functions
*/

func main(){
	service := "127.0.0.1:10000"
	server, err := net.ResolveUDPAddr("udp",service)
	conn, err := net.DialUDP("udp", nil, server)

	if err != nil {
        fmt.Println("Error connecting to " , service)
        fmt.Println(err)
        return
    }
	for n := int64(0); n >= 0; n++ {
		n, err :=	conn.Write([]byte("Hello"))
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
		time.Sleep(1 * time.Second)
	}

	conn.Close()
}


/***
*	Function Name: 	move()
*	Desc:			The function requests move
*	Pre-cond:		takes connection argument and the new location
*	Post-cond:		retuen success or return failure
*/
func move(){
}


/***
*	Function Name: 	fire()
*	Desc:			The function requests fire
*	Pre-cond:		takes connection argument and name of client who is fired
*	Post-cond:		return success or returns failure
*/
func fire(){
}
