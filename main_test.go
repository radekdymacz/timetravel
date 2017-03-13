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
	// let's run timetravel
	_, err = Timetravel(os.TempDir())
	if err != nil {
		t.Fatalf("%s", err)

	}
	info, err := file.Stat()
	if err != nil {
		t.Errorf("%s", err)
	}
	//check if timestamp stil in the future
	if info.ModTime().After(time.Now()) {
		t.Errorf("Timestamp in the future %s", info.ModTime())
	}

	//clean up
	defer os.Remove(file.Name())
	defer file.Close()
}
