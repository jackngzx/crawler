package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func newConfig(rawBaseURL string, maxConcurrency int, maxPages int) (*config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse baseURL: %s", err)
	}
	return &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}, nil
}

func main() {
	argCount := len(os.Args)
	if argCount < 4 {
		fmt.Println("missing arguments: usage ./crawler URL maxConcurrency maxPages")
		os.Exit(1)
	}
	if argCount > 4 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	fmt.Println("starting crawl of: ", os.Args[1])

	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println(err)
		return
	}

	c, err := newConfig(os.Args[1], maxConcurrency, maxPages)
	if err != nil {
		fmt.Print(err)
		return
	}

	c.wg.Add(1)
	go c.crawlPage(os.Args[1])
	c.wg.Wait()

	for normalizedURL := range c.pages {
		fmt.Printf("urls crawled: %s\n", normalizedURL)
	}
	fmt.Println("finished crawling pages")
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

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	if cfg.pagesLen() >= cfg.maxPages {
		return
	}

	parsedBase, err := url.Parse(cfg.baseURL.String())
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
		return
	}

	normCurrent, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println("error normalizing rawCurrentURL: ", err)
		return
	}

	addPage := cfg.addPageVisit(normCurrent)
	if !addPage {
		return
	}

	pageHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println("error getting html from normCurrent: ", err)
		return
	}

	pageData := extractPageData(pageHTML, rawCurrentURL)
	cfg.setPageData(normCurrent, pageData)

	for _, pageURL := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(pageURL)
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	val, exists := cfg.pages[normalizedURL]
	if exists {
		return false
	} else {
		cfg.pages[normalizedURL] = val
		return true
	}
}

func (cfg *config) setPageData(normalizedURL string, data PageData) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.pages[normalizedURL] = data
}

func (cfg *config) pagesLen() int {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	return len(cfg.pages)
}
