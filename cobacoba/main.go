package main

import (
	"fmt"
)

type Student struct{
	NIM string
	FirstName string
	LastName string
	Class string
	Age int
}

func main(){

	
	map1 := map[string]string{
		"key1" : "value1",
		"key2" : "value2",
		"key3" : "value2",
	}

	map1["key1"] = "value4"

	if value, found := map1["key1"]; found {
		fmt.Println(value)
	}

}