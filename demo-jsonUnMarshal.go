package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type employee struct {
	ID           int
	EmployeeName string
	Tel          string
	Email        string
}

func main() {
	e := employee{}
	err := json.Unmarshal([]byte(`{"ID":101, "EmployeeName":"Pontakorn Changpat", "Tel":"0876667777", "Email":"pontakorn322@gmail.com"}`), &e)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(e.ID)
}
