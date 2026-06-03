package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	argCount := len(os.Args)
	if argCount < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if argCount > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	if argCount == 2 {
		fmt.Println("starting crawl of: ", os.Args[1])
	}

	crawledPages := make(map[string]int)
	crawlPage(os.Args[1], os.Args[1], crawledPages)

	fmt.Println(crawledPages)
}

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "GoCrawler/1.0")

	c := &http.Client{}

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server error code: %v", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("resp header is not text/html")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	parsedBase, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Println("error parsing rawBaseURL: ", err)
		return
	}

	parsedCurrent, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println("error parsing rawCurrentURL: ", err)
		return
	}

	if parsedBase.Hostname() != parsedCurrent.Hostname() {
		fmt.Println("different domain between base and current url")
		return
	}

	normCurrent, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println("error normalizing rawCurrentURL: ", err)
		return
	}

	val, exists := pages[normCurrent]
	if exists {
		fmt.Println("page is already crawled")
		pages[normCurrent] = val + 1
		return
	} else {
		pages[normCurrent] = 1
	}

	fmt.Println("getting page html from current url: ", normCurrent)
	pageHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println("error getting html from normCurrent: ", err)
		return
	}
	fmt.Println("page html retrieval successful")

	pageURLs, err := getURLsFromHTML(pageHTML, parsedCurrent)
	if err != nil {
		fmt.Println("error getting urls from page: ", err)
		return
	}

	for _, pageURL := range pageURLs {
		crawlPage(rawBaseURL, pageURL, pages)
	}
}
