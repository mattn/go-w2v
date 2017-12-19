package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/text/encoding/japanese"

	"github.com/PuerkitoBio/goquery"
)

func topicList() ([]string, error) {
	resp, err := http.Get("http://komachi.yomiuri.co.jp/?g=04&o=2")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(japanese.ShiftJIS.NewDecoder().Reader(resp.Body))
	if err != nil {
		return nil, err
	}
	uris := []string{}
	doc.Find(".topicslist .hd a").Each(func(n int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok {
			return
		}
		uris = append(uris, "http://komachi.yomiuri.co.jp"+href)
	})
	return uris, nil
}

func saveTopic(filename string, uri string) error {
	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(japanese.ShiftJIS.NewDecoder().Reader(resp.Body))
	if err != nil {
		return err
	}
	text := doc.Find("#topiccontent p").First().Text()
	err = ioutil.WriteFile(filename, []byte(text), 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	uris, err := topicList()
	if err != nil {
		log.Fatal(err)
	}

	for n, uri := range uris {
		log.Println(uri)
		err = saveTopic(fmt.Sprintf("data/data%03d.txt", n+1), uri)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(2 * time.Second)
	}
}
