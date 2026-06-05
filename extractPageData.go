package main

import (
	"fmt"
	"net/url"
)

type PageData struct {
	URL            string   `json:"url"`
	Heading        string   `json:"heading"`
	FirstParagraph string   `json:"first_paragraph"`
	OutgoingLinks  []string `json:"outgoing_links"`
	ImageURLs      []string `json:"image_urls"`
}

func extractPageData(html, pageURL string) PageData {
	normalized, err := normalizeURL(pageURL)
	if err != nil {
		fmt.Println("error normalizing url: ", err)
	}
	fmt.Println(normalized)
	u, err := url.Parse(pageURL)
	if err != nil {
		fmt.Println("error parsing url: ", err)
	}

	heading := getHeadingFromHTML(html)

	firstP := getFirstParagraphFromHTML(html)

	links, err := getURLsFromHTML(html, u)
	if err != nil {
		fmt.Println("error getting links from html: ", err)
	}

	images, err := getImagesFromHTML(html, u)
	if err != nil {
		fmt.Println("error getting images from html: ", err)
	}

	return PageData{
		URL:            pageURL,
		Heading:        heading,
		FirstParagraph: firstP,
		OutgoingLinks:  links,
		ImageURLs:      images,
	}
}
