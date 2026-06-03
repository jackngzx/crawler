package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	listURLs := []string{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return listURLs, err
	}
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			crawled, err := url.Parse(href)
			if err != nil {
				fmt.Println(err)
			}
			finalURL := baseURL.ResolveReference(crawled)
			listURLs = append(listURLs, finalURL.String())
		}
	})
	return listURLs, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	listIMGs := []string{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return listIMGs, err
	}
	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if ok {
			crawled, err := url.Parse(src)
			if err != nil {
				fmt.Println(err)
			}
			finalURL := baseURL.ResolveReference(crawled)
			listIMGs = append(listIMGs, finalURL.String())
		}
	})
	return listIMGs, nil
}
