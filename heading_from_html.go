package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(html string) string {
	r := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return fmt.Sprint(err)
	}

	h1Heading := doc.Find("h1").Text()
	if h1Heading == "" {
		h2Heading := doc.Find("h2").Text()
		if h2Heading == "" {
			fmt.Println("h1 and h2 headers are not found")
			return ""
		}
		return h2Heading
	}
	return h1Heading
}
func getFirstParagraphFromHTML(html string) string {
	r := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return fmt.Sprint(err)
	}
	main := doc.Find("main")
	if main.Text() == "" {
		fmt.Println("main body is not found")
		return ""
	}
	paragraphs := main.Find("p")
	if paragraphs.Text() == "" {
		fmt.Println("main paragraph is not found")
		return ""
	}
	return paragraphs.First().Text()
}
