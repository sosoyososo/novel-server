package main

import (
	"./service"
	_ "./service/novel"
)

func main() {
	service.CreateService()
}
