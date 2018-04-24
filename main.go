package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	BASE_URL   string      = "https://www.lifeofpix.com/gallery/"
	USER_AGENT string      = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:50.0) Gecko/20100101 Firefox/50.0"
	BASE_DIR   string      = "gallery"
	DIR_MODE   os.FileMode = 0777
)

var wg sync.WaitGroup

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

	if os.MkdirAll(BASE_DIR, DIR_MODE) != nil {
		fmt.Println("mkdir gallery failed")
		return
	}

	for i := 0; i < 1; i++ {
		dir := BASE_DIR + "/" + title[i]
		os.MkdirAll(dir, DIR_MODE)
		doc, err = goquery.NewDocument(links[i])
		if err != nil {
			fmt.Println(err)
			return
		}

		parseImgUrl(doc, dir)

		link, ok := doc.Find("div.pagination").Find("a").Eq(0).Attr("href")
		for ok {
			fmt.Println("next page: ", link)
			doc, err = goquery.NewDocument(link)
			if err != nil {
				fmt.Println(err)
				break
			}
			parseImgUrl(doc, dir)

			link, ok = doc.Find("div.pagination").Find("a").Eq(1).Attr("href")
			if !ok {
				break
			}
		}
	}

	wg.Wait()
}

func parseImgUrl(doc *goquery.Document, dir string) {
	doc.Find("div.node").Each(func(i int, s *goquery.Selection) {
		//		if imgUrl, ok := s.Find("img").Attr("src"); ok {
		//			if !strings.Contains(imgUrl, "adobe") {
		//				if pos := strings.LastIndex(imgUrl, "-"); pos != -1 {
		//					pos2 := strings.LastIndex(imgUrl, ".")
		//					oriImgUrl := imgUrl[:pos] + imgUrl[pos2:]
		//					fmt.Println(oriImgUrl)
		//					//wg.Add(1)
		//					//go downloadImg(oriImgUrl, dir)
		//				}
		//			}
		s.Find("div.actions").Find("a").Each(func(j int, ss *goquery.Selection) {
			if imgUrl, ok := ss.Attr("download"); ok {
				fmt.Println(imgUrl)
				wg.Add(1)
				go downloadImg(imgUrl, dir)
			}
		})
	})
}

func downloadImg(imgUrl, dir string) {
	defer wg.Done()
	url, err := url.Parse(imgUrl)
	if err != nil {
		log.Println("parse url failed: ", imgUrl, err)
		return
	}
	filename := dir + "/" + url.Path[strings.LastIndex(url.Path, "/")+1:]
	if isExist(filename) {
		return
	}

	log.Println("downloading: " + imgUrl)
	resp, err := http.Get(imgUrl)
	if err != nil {
		log.Println("get img url failed: ", imgUrl, err)
		return
	}
	defer resp.Body.Close()

	respdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read img data failed: ", err)
		return
	}
	image, err := os.Create(filename)
	if err != nil {
		log.Println("create img file failed: ", err)
		return
	}

	defer image.Close()
	_, err = image.Write(respdata)
	if err != nil {
		log.Println("write img file failed: ", filename)
		os.Remove(filename)
		return
	}
}

func isExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
