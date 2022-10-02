package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cespare/xxhash"
)

var (
	reportOutputType = "list" // In the future it will be (or won't be) a flag item [list|json|table]
)

// ProcessedFile consists of all the data about
// file that have been processed we need:
// Full path to file
// hash sums (sha256, xxhash64)
// Size in bytes
// Creation time
// Change time
type ProcessedFile struct {
	Path        string
	Sha256sum   string
	Xxhash64sum string
	Size        uint64
	Created     time.Time
	Changed     time.Time
}

func (pf ProcessedFile) String() string {
	return fmt.Sprintf("Path:\t\t%s\nxxhash64:\t%s\nsha256:\t\t%s\nSize:\t\t%d\n\n", pf.Path, pf.Xxhash64sum, pf.Sha256sum, pf.Size)
}

type Report struct {
	Files []ProcessedFile
}

// String returns string that contains full report
// for the Report struct
func (r Report) String() string {
	var sRep string

	switch strings.ToLower(reportOutputType) {
	case "list":
		for _, pf := range r.Files {
			sRep += pf.String()
		}
	case "json":

	case "table":
	default:

	}

	return sRep
}

func main() {
	var userPath string

	if len(os.Args) < 2 {
		userPath = "./"
	} else {
		userPath = os.Args[1]
	}

	// Start test
	f, d, err := dirContent(userPath)
	checkErr(err)

	rep := new(Report)
	loopFiles(userPath, f, d, rep)
	fmt.Println(*rep)
	return
	// End test

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

// dirContent returns the list of files in ([]os.DirEntry)
// and a slice of strings that contains names of all directories
// in the given path.
// If the given path is a file, then error will be returned
func dirContent(dn string) (files []os.DirEntry, directories []string, err error) {
	fe, err := os.ReadDir(dn)
	if err != nil {
		return files, directories, err
	}

	for _, v := range fe {
		if v.IsDir() {
			directories = append(directories, v.Name())
		} else {
			files = append(files, v)
		}
	}

	return
}

func loopFiles(dirPath string, files []os.DirEntry, folders []string, rep *Report) {

	for _, de := range files {
		fc, err := os.ReadFile(filepath.Join(dirPath, de.Name()))
		checkErr(err)
		pf := ProcessedFile{
			Path:        filepath.Join(dirPath, de.Name()),
			Sha256sum:   sha256Sum(fc),
			Xxhash64sum: xxhash64Sum(fc),
			Size:        uint64(len(fc)),
		}
		(*rep).Files = append((*rep).Files, pf)
	}

	for _, fn := range folders {
		fullDirPath := filepath.Join(dirPath, fn)
		fls, fldrs, err := dirContent(fullDirPath)
		checkErr(err)
		loopFiles(fullDirPath, fls, fldrs, rep)
	}
}

func loopThroughFiles(dirPath string, des []os.DirEntry) {
	var wg sync.WaitGroup
	for _, de := range des {
		wg.Add(1)
		if de.IsDir() { // recursive call
			dn := filepath.Join(dirPath, de.Name())
			loopThroughFiles(dn, filesInDir(dn))
			wg.Done()
			continue
		}

		go func(de os.DirEntry) { // base case, check file's hash
			// Get file content
			fc, err := os.ReadFile(filepath.Join(dirPath, de.Name()))
			checkErr(err)
			startAt := time.Now()

			fmt.Printf("Path:\t\t%s\nName:\t\t%s\nxxhash64:\t%s\nsha256:\t\t%s\nSize:\t\t%d\nTime:\t\t%d\n", dirPath, de.Name(), xxhash64Sum(fc), sha256Sum(fc), len(fc), time.Since(startAt))
			fmt.Println("--------")
			wg.Done()
		}(de)
		wg.Wait() // It probably removes the effect from concurrency
		// I have to manage amount of bytes that are being reading at the current time and wait
		// if it exeeds certain value
	}
}
