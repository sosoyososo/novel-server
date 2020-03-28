package Html

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"../Config"
	"../Encoding"
)

var (
	host string
)

func Start(config *Config.Config,
	finish func(),
	pageStart func(string),
	pageFinish func(string, map[string]interface{}, error)) func() {

	p := NewParser(*config)

	u, err := url.Parse(config.EntryURL)
	if nil != err {
		panic(err)
	}
	host = u.Host
	spider := SpiderWorker{}
	spider.Run(config.EntryURL,
		10,
		func(url string) bool {
			if isInsiteURL(url) == false {
				return false
			}
			if isPicURL(url) {
				return false
			}
			if isUrlSkipedByConfig(url, config) {
				return false
			}
			return true
		}, func(worker *HtmlWorker) {
			//config request
			worker.CookieStrig = config.CookiesString
			encoding := strings.ToLower(config.Encoding)
			encoder := Encoding.Encoders[encoding]
			worker.Encoder = encoder
		}, func(url string, work *HtmlWorker) {
			//handle result
			p.ParseDocument(work)
		}, func(url string) string {
			//url in href tag as parameter here ,
			// u should convert and make sure result is legal for http
			return fullURL(url)
		}, finish)

	//FIXME: return action can stop this progress
	return func() {}
}
func isInsiteURL(URL string) bool {
	if strings.HasPrefix(URL, "/") ||
		strings.HasPrefix(URL, "../") ||
		strings.HasPrefix(URL, "http://"+host) ||
		strings.HasPrefix(URL, host) {
		return true
	}
	return false
}

func isUrlSkipedByConfig(url string, config *Config.Config) bool {
	skiped := false
	j := 0
	for j < len(config.SkipUrls) {
		re := config.SkipUrls[j]
		validID := regexp.MustCompile(re)
		if validID.MatchString(url) {
			skiped = true
			break
		}
		j++
	}
	return skiped
}

func isPicURL(URL string) bool {
	if strings.HasSuffix(URL, ".jpg") ||
		strings.HasSuffix(URL, ".png") {
		return true
	}
	return false
}

func fullURL(url string) string {
	if strings.HasPrefix(url, "www.") {
		return url
	}
	if strings.HasPrefix(url, "http://") {
		return url
	}
	if strings.HasPrefix(url, "/") {
		return "http://" + host + url
	}
	if strings.HasPrefix(url, "./") {
		url = strings.TrimLeft(url, ".")
		return "http://" + host + url
	}
	fmt.Printf("无法处理的url: %s\n", url)
	return url
}
