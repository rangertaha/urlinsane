package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	// Existing JSON data
	jsonData := []byte(`{"name": "John", "age": 30}`)

	// Unmarshal the JSON data into a struct
	var person Person
	err := json.Unmarshal(jsonData, &person)
	if err != nil {
		panic(err)
	}

	// Extend the data with a new field
	person.City = "New York"

	// Marshal the extended data back to JSON
	newJsonData, err := json.Marshal(person)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(newJsonData)) // Output: {"name": "John", "age": 30, "City": "New York"}
}
