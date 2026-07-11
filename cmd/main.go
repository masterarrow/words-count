package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: wc <filepath>\nOr wc <filepath1> <filepath2> ...")
	}

	appStart := time.Now()

	files := os.Args[1:]

	const maxConcurrency = 5
	jobs := make(chan string, maxConcurrency)
	results := make(chan Result, maxConcurrency)
	var wg sync.WaitGroup

	for w := 0; w < maxConcurrency; w++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for path := range jobs {
				results <- processFile(path)
			}
		}()
	}

	go func() {
		for _, path := range files {
			jobs <- path
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if res.Err != nil {
			fmt.Printf("%s: %v [took: %v]\n", res.FileName, res.Err, res.Duration)
			continue
		}

		fmt.Printf("%s: words %d, lines %d [took: %v]\n", res.FileName, res.Words, res.Lines, res.Duration)
	}

	fmt.Printf("\nTotal execution time: %v\n", time.Since(appStart))
}

type Result struct {
	FileName string
	Words    int
	Lines    int
	Err      error
	Duration time.Duration
}

func processFile(path string) Result {
	start := time.Now()

	file, err := os.Open(path)
	if err != nil {
		return Result{FileName: path, Err: err, Duration: time.Since(start)}
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return Result{FileName: path, Err: err, Duration: time.Since(start)}
	}
	if info.IsDir() {
		return Result{FileName: path, Err: fmt.Errorf("is a directory, skipped"), Duration: time.Since(start)}
	}

	counter := NewCounter(file)
	if err := counter.Execute(); err != nil {
		return Result{FileName: path, Err: err, Duration: time.Since(start)}
	}

	return Result{
		FileName: path,
		Words:    counter.Words,
		Lines:    counter.Lines,
		Duration: time.Since(start),
	}
}

type Counter struct {
	Words  int
	Lines  int
	reader io.Reader
	inWord bool
}

func NewCounter(r io.Reader) *Counter {
	return &Counter{
		reader: r,
	}
}

func (c *Counter) Execute() error {
	data := make([]byte, 32*1024)

	for {
		n, err := c.reader.Read(data)
		c.count(data[:n])

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	return nil
}

func (c *Counter) count(chunk []byte) {
	for _, b := range chunk {
		if b == '\n' {
			c.Lines++
		}

		isSpace := b == ' ' || b == '\n' || b == '\t' || b == '\r'

		if !isSpace && !c.inWord {
			c.inWord = true
			c.Words++
		} else if isSpace && c.inWord {
			c.inWord = false
		}
	}
}
