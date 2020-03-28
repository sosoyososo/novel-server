package main

import (
	"./NovelSpider"
	"./service"
	_ "./service/novel"
)

func main() {
	/**
	 * 1. load conf
	 * 2. start web service
	 * 3. start repeat task
	 * 4. repeat find new summary spider
	 * 5. repeat load and check catelog page
	 */
	// NovelSpider.StartSummarySpider()
	go NovelSpider.StartCatelogSpider()
	// NovelSpider.UpdateCatelogNovelID()
	service.CreateService()
}
