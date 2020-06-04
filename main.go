package main

import (
	// "./service"
	"./NovelSpider"
	_ "./service/novel"
)

func main() {
	// service.CreateService()
	NovelSpider.StartSummarySpider()
}
