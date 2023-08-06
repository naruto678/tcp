package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"bytes"
	"os"
	"time"
) 
func main(){
	switch os.Args[1]{
		case "server":
		runAsServer()
		default:
		if len(os.Args)<3 {
			panic("client name must be specified")
		}
		runAsClient()
	}
}

func handleConnection(conn net.Conn){

	decoder := gob.NewDecoder(conn)

	for {
		request := ClientRequest{}
		err := decoder.Decode(&request)
		if err!=nil{
			fmt.Println("Got error while decoding ", err)
			break
		}
		fmt.Printf("Received %s\n", request)

	}
	(conn).Close()
}


func runAsServer(){
	fmt.Println("Running as a server")
	l, err := net.Listen("tcp",":8080")
	if err!=nil{
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err!=nil{
			fmt.Println(err)
		}
		go handleConnection(conn)
	}

}
func runAsClient(){
	client_name := os.Args[2]
	fmt.Println("Running as a client ", client_name)
	conn, err := net.Dial("tcp", ":8080")
	if err!=nil{
		fmt.Println(err)
	}
	sysReader := bufio.NewReader(os.Stdin)
	connWriter := bufio.NewWriter(conn)

	time.Sleep(1*time.Second)
	request := ClientRequest{
		Name : client_name,
	}
	for{
		fmt.Printf(">> ")
		bct, err := sysReader.ReadBytes('\n')
		request.Content = string(bct)
		if err!=nil{
			fmt.Println(err)
		}
		connWriter.Write(request.Serialize())
		err = connWriter.Flush()
		fmt.Println(err)
	}
}

type ClientRequest struct{
	Name string
	Content string
}



func (cr ClientRequest) String() string{
	return fmt.Sprintf("ClientRequest(name = %s, content = %s)", cr.Name, cr.Content)
}

func (cr ClientRequest) Serialize() []byte {
	result := bytes.Buffer{}
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(cr)
	if err!=nil{
		fmt.Println("Error during serializing ", err)
	}
	return result.Bytes()
}


func Deserialize(request []byte ) ClientRequest{
	fmt.Print("Deserializing ", string(request))
	cr := ClientRequest{}
	dec := gob.NewDecoder(bytes.NewBuffer(request))
	err := dec.Decode(&cr)
	if err!=nil{
		fmt.Println(err)
	}
	return cr

}
