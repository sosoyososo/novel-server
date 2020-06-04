package main

import (
	"./NovelSpider"
	"./service"
	_ "./service/novel"
)

func main() {
	go NovelSpider.StartSummarySpider()
	service.CreateService()
}
