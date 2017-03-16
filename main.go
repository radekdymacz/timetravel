// Copied and extended from https://blog.golang.org/pipelines/bounded.go

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// walkFiles starts a goroutine to walk the directory tree at root and send the
// path of each regular file on the string channel.  It sends the result of the
// walk on the error channel.  If done is closed, walkFiles abandons its work.
func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)
	go func() { // HL
		// Close the paths channel after Walk returns.
		defer close(paths) // HL
		// No select needed for this send, since errc is buffered.
		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error { // HL
			if err != nil {
				log.Printf("%s", err)
				//continue after error
				return nil
			}
			if !info.Mode().IsRegular() {
				//log.Printf("%s ", info.Mode().IsRegular)
				return nil
			}
			// Return anly files with  modified time in the future
			if !info.ModTime().After(time.Now()) {
				//log.Printf("Timestamp OK %s", path)
				return nil
			}

			log.Printf("Found file with wrong timestamp  %s %s", path, info.ModTime())
			select {
			case paths <- path: // HL
			case <-done: // HL
				return errors.New("walk canceled")
			}
			return nil
		})
	}()
	return paths, errc
}

// A result is the product of reading and summing a file using MD5.
type result struct {
	path string
	err  error
}

// modTime reads pathas names form paths and fixes metadata if it in the furure
func modTime(done <-chan struct{}, paths <-chan string, c chan<- result) {
	for path := range paths { // HLpaths
		log.Printf(" Adjust timestamp for %s", path)
		err := os.Chtimes(path, time.Now(), time.Now())
		select {
		case c <- result{path, err}:
		case <-done:
			return
		}
	}
}

// Timetravel reads all the files in the file tree rooted at root and start digesters/time travel machines to fix metadata.
func Timetravel(root string) (map[string]string, error) {
	//set log stdout

	log.SetOutput(os.Stdout)
	// Timetravel closes the done channel when it returns; it may do so before
	// receiving all the values from c and errc.
	done := make(chan struct{})
	defer close(done)

	paths, errc := walkFiles(done, root)

	// Start a fixed number of goroutines to read and digest files.
	c := make(chan result) // HLc
	var wg sync.WaitGroup
	const numDigesters = 20
	wg.Add(numDigesters)
	for i := 0; i < numDigesters; i++ {
		go func() {
			modTime(done, paths, c) // HLc
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(c) // HLc
	}()
	// End of pipeline.

	m := make(map[string]string)
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.path
	}
	// Check whether the Walk failed.
	if err := <-errc; err != nil { // HLerrc
		return nil, err
	}
	return m, nil
}

func main() {
	// Fix timestamps for all files under the specified directory,
	// then print the results sorted by path name.
	m, err := Timetravel(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	var paths []string
	for path := range m {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x  %s\n", m[path], path)
	}
}
