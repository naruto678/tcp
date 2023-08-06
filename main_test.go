package main

import (
	"testing"
	"fmt"
)


func TestSerialize(t *testing.T){
	cr := ClientRequest{
		Name : "first-client",
		Content : "this is from first client",
	}
	result := cr.Serialize()

	request := Deserialize(result)
	if request.Content == cr.Content && request.Name==cr.Name{
		fmt.Println("test passed")
	} else{
		t.Fatal("test failed ")
	}
}
