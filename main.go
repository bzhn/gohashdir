package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/cespare/xxhash"
)

func main() {
	var userPath string

	if len(os.Args) < 2 {
		userPath = "./"
	} else {
		userPath = os.Args[1]
	}

	if de, err := os.Stat(userPath); !de.IsDir() {
		checkErr(err)
		fc, err := os.ReadFile(userPath)
		checkErr(err)
		startAt := time.Now()

		fmt.Printf("Path:\t\t%s\nName:\t\t%s\nxxhash64:\t%s\nsha256:\t\t%s\nSize:\t\t%d\nTime:\t\t%d\n", userPath, de.Name(), xxhash64Sum(fc), sha256Sum(fc), len(fc), time.Since(startAt))
		return
	}
	fe := filesInDir(userPath)

	loopThroughFiles(userPath, fe)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func xxhash64Sum(b []byte) string {
	n := xxhash.Sum64(b)
	return fmt.Sprintf("%016s", strconv.FormatUint(n, 16))
}

func sha256Sum(b []byte) string {
	hasher := sha256.New()
	hasher.Write(b)

	n := hasher.Sum(nil)
	return fmt.Sprintf("%x", n)
}

func filesInDir(dn string) []os.DirEntry {
	fe, err := os.ReadDir(dn)
	checkErr(err)

	return fe
}

func loopThroughFiles(dirPath string, des []os.DirEntry) {
	var wg sync.WaitGroup
	for _, de := range des {
		wg.Add(1)
		if de.IsDir() { // recursive call
			dn := path.Join(dirPath, de.Name())
			loopThroughFiles(dn, filesInDir(dn))
			wg.Done()
			continue
		}

		go func(de os.DirEntry) { // base case, check file's hash
			// Get file content
			fc, err := os.ReadFile(path.Join(dirPath, de.Name()))
			checkErr(err)
			startAt := time.Now()

			fmt.Printf("Path:\t\t%s\nName:\t\t%s\nxxhash64:\t%s\nsha256:\t\t%s\nSize:\t\t%d\nTime:\t\t%d\n", dirPath, de.Name(), xxhash64Sum(fc), sha256Sum(fc), len(fc), time.Since(startAt))
			fmt.Println("--------")
			wg.Done()
		}(de)
		wg.Wait() // It probably removes the effect from concurrency
		// I have to manage amount of bytes that are being reading in the current time and wait
		// if it exeeds certain value
	}
}
