package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
) 
func main(){
	switch os.Args[1]{
		case "server":
		runAsServer()
		default:
		if len(os.Args)<3 {
			panic("client name must be specified")
		}
	}
		runAsClient()
}

func handleConnection(conn net.Conn){

	for {
		request := ClientRequest{}
		response := ClientRequest{Name: "server"}
		decoder := gob.NewDecoder(conn)
		err := decoder.Decode(&request)
		if err!=nil{
			fmt.Println("Got error while decoding ", err)
			break
		}
		if request.Content == ".quit"{
			response.Content = "Bye!!!"
			conn.Write(response.Serialize())
			break
		} else {
			response.Content = "Received"
			conn.Write(response.Serialize())
		}
	}
	conn.Close()
}


func runAsServer(){
	fmt.Println("Running as a server")
	l, err := net.Listen("tcp",":8080")
	if err!=nil{
		fmt.Println(err)
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
	defer conn.Close()
	if err!=nil{
		fmt.Println(err)
		os.Exit(1)
	}
	sysReader := bufio.NewReader(os.Stdin)
	connWriter := bufio.NewWriter(conn)
	request := ClientRequest{
		Name : client_name,
	}

	rcvMessage := ClientRequest{}

	for{
		fmt.Printf(">> ")
		bct, err := sysReader.ReadBytes('\n')
		bct = bct[:len(bct)-1]
		request.Content = string(bct)
		if err!=nil{
			fmt.Println(err)
		}
		connWriter.Write(request.Serialize())
		err = connWriter.Flush()
		if err!=nil{
			fmt.Println(err)
		}
		decoder := gob.NewDecoder(conn)
		err = decoder.Decode(&rcvMessage)
		if err!=nil{
			fmt.Println(err)
		}
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
