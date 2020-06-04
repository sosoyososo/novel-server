package NovelSpider

import (
	"fmt"
	"testing"
)

func TestSpider(t *testing.T) {
}

func TestSpiderConf(t *testing.T) {
	for _, confRef := range confMap {
		conf := *confRef
		fmt.Println(conf.IsSummaryPage("/book/18006/"))
	}
}
