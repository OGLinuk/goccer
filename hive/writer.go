package hive

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

type URLFile struct {
	file *os.File
	urls map[string]struct{}
}

type URLWriter struct {
	path     string
	urlFiles map[string]*URLFile
}

func NewURLFile(fPath string) *URLFile {
	file, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Printf("os.OpenFile (NewUrlFile) err: %s", err.Error())
	}
	return &URLFile{
		file: file,
		urls: make(map[string]struct{}),
	}
}

func NewURLWriter(dPath string) *URLWriter {
	err := os.MkdirAll(dPath, 0777)
	if err != nil {
		log.Printf("os.MkdirAll (NewUrlWriter) err: %s", err.Error())
	}
	return &URLWriter{
		path:     dPath,
		urlFiles: make(map[string]*URLFile),
	}
}

func (uw *URLWriter) Write(URL string) error {
	u, err := url.Parse(URL)
	if err != nil {
		return err
	}

	base := u.Hostname()
	if base == "" || base == " " {
		base = "error"
	}
	fileDir := filepath.Join(uw.path, base)

	if _, ok := uw.urlFiles[base]; !ok {
		uw.urlFiles[base] = NewURLFile(fileDir)
	}

	uf := uw.urlFiles[base]

	decoded, err := url.QueryUnescape(URL)
	if err != nil {
		return err
	}

	if _, exists := uf.urls[decoded]; !exists {
		uf.file.WriteString(fmt.Sprintf("%s\n", decoded))
		uf.urls[decoded] = struct{}{}
	}

	uf.file.Close()

	return nil
}

func (q *Queen) Aggregate() {
	for k := range q.rw.urlFiles {
		fileDir := filepath.Join(q.rw.path, k)

		fd, err := os.Open(fileDir)
		if err != nil {
			log.Printf("os.Open (Aggregate) err: %s", err.Error())
		}

		scanner := bufio.NewScanner(fd)
		for scanner.Scan() {
			scanned := scanner.Text()
			q.aw.file.WriteString(fmt.Sprintf("%s\n", scanned))
		}
		fd.Close()
	}
}