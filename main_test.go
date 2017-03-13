package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestTimetravel(t *testing.T) {
	log.Printf("creating test file in %s", os.TempDir())
	file, err := ioutil.TempFile(os.TempDir(), "timetravel-test")

	if err != nil {
		t.Fatalf("%s", err)
	}
	// Chtimes to have timestamp in future
	fd := time.Now().Add(time.Hour * 24 * 3000)
	log.Printf("Change timestamps 3 days ins future for %s", file.Name())
	err = os.Chtimes(file.Name(), fd, fd)
	if err != nil {
		t.Fatalf("%s", err)
	}
	_, err = Timetravel(os.TempDir())
	if err != nil {
		t.Fatalf("%s", err)

	}
	info, err := file.Stat()

	log.Printf("%s", info.ModTime())
	defer file.Close()
	// var paths []string
	// for path := range m {
	// 	paths = append(paths, path)
	// }
	// sort.Strings(paths)
	// for _, path := range paths {
	// 	log.Printf("%s \n", path)
	// }
	defer os.Remove(file.Name())
}
