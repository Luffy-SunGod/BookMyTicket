package main

import (
	"fmt"
	"github.com/Luffy-SunGod/MyShowSeat/cmd/app"
)

const webPort = "8095"

const pgConnectionString = "host=localhost port=5432 user=rayanc dbname=tickets sslmode=disable"

func main() {
	fmt.Println("Main file started!!!")
	service := app.Service{}
	service.Start()
}
