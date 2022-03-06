package main

import (
	"encoding/json"
	"fmt"
)

type employee struct {
	ID           int
	EmployeeName string
	Tel          string
	Email        string
}

func main() {

	data, _ := json.Marshal(&employee{101, "Pontakorn Changpat", "0876667777", "pontakorn322@gmail.com"})
	fmt.Println(string(data))

}
