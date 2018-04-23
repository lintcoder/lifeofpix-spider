package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	BASE_URL   string = "https://www.lifeofpix.com/gallery/"
	USER_AGENT string = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:50.0) Gecko/20100101 Firefox/50.0"
)

func main() {
	doc, err := goquery.NewDocument(BASE_URL)
	if err != nil {
		fmt.Println(err)
		return
	}

	var title []string
	var links []string
	var count []int
	doc.Find("div.col-xs-10").Each(func(i int, s *goquery.Selection) {
		sub := s.Find("a")
		if t := sub.Find("div.title"); t.Nodes != nil {
			title = append(title, t.Text())
			l, _ := sub.Attr("href")
			links = append(links, l)
			num, _ := strconv.Atoi(strings.Fields(sub.Find("div.count").Text())[0])
			count = append(count, num)
		}
	})

	doc, err = goquery.NewDocument(links[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	doc.Find("div.node").Each(func(i int, s *goquery.Selection) {
		imgurl, _ := s.Find("img").Attr("src")
		fmt.Println(imgurl)
	})
}
