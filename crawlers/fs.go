package crawlers

import (
	"io/ioutil"
	"log"
)

// FSCrawler for filesystems
type FSCrawler struct {
	path string
}

// NewFSCrawler constructor
func NewFSCrawler(p string) FSCrawler {
	return FSCrawler{
		path: p,
	}
}

// Crawl c.path
func (c FSCrawler) Crawl() []string {
	var paths []string

	infos, err := ioutil.ReadDir(c.path)
	if err != nil {
		log.Printf("crawlers::fs.go::Crawl::ioutil.ReadDir(%s)::ERROR: %s", c.path, err.Error())
		return nil
	}

	for _, info := range infos {
		log.Printf("Found: %s", info.Name())
		paths = append(paths, info.Name())
	}

	return paths
}
