package goccer

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// HTTPCrawler for HTTP URLs
type HTTPCrawler struct {
	seed string
	wg   *sync.WaitGroup
}

// NewHTTPCrawler constructor
func NewHTTPCrawler(s string) HTTPCrawler {
	return HTTPCrawler{
		seed: s,
		wg:   &sync.WaitGroup{},
	}
}

// Crawl c.seed and extract all URLs
func (c HTTPCrawler) Crawl() ([]string, error) {
	var collected []string

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Second * 7,
	}

	resp, err := client.Get(c.seed)
	if err != nil {
		log.Printf("crawlers::Crawl::client.Get(%s)::ERROR: %s", c.seed, err.Error())
	}
	defer resp.Body.Close()

	if resp == nil {
		log.Printf("crawlers::Crawl::resp::NIL")
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		for _, URL := range c.extract(resp, c.seed) {
			collected = append(collected, URL)
		}
	} else {
		err = fmt.Errorf("crawlers::Crawl::resp.StatusCode(%d): %s", resp.StatusCode, c.seed)
	}

	return collected, nil
}

// extract links from resp
func (c HTTPCrawler) extract(resp *http.Response, seed string) []string {
	if resp == nil {
		return nil
	}
	links := collectLinks(resp.Body)
	rebuiltLinks := []string{}

	for _, link := range links {
		url := rebuildURL(link, seed)
		if url != "" {
			rebuiltLinks = append(rebuiltLinks, url)
		}
	}

	return rebuiltLinks
}

// collectLinks from httpBody
func collectLinks(httpBody io.Reader) []string {
	links := make(map[string]struct{})
	col := []string{}
	page := html.NewTokenizer(httpBody)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			l := []string{}
			for k := range links {
				l = append(l, k)
			}
			return l
		}
		token := page.Token()
		if tokenType == html.StartTagToken && token.DataAtom.String() == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					tl := trimHash(attr.Val)
					col = append(col, tl)
					for _, link := range col {
						if _, exists := links[link]; !exists {
							links[link] = struct{}{}
						}
					}
				}
			}
		}
	}
}

// trimHash of (l)ink
func trimHash(l string) string {
	if strings.Contains(l, "#") {
		for n, str := range l {
			if strconv.QuoteRune(str) == "'#'" {
				return l[:n]
			}
		}
	}
	return l
}

// rebuildURL using href and base
func rebuildURL(href, base string) string {
	url, err := url.Parse(href)
	if err != nil {
		return ""
	}

	baseURL, err := url.Parse(base)
	if err != nil {
		return ""
	}

	return baseURL.ResolveReference(url).String()
}
