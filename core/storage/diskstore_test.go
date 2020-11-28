package storage

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiskGetNoFileErrors(t *testing.T) {
	dir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	ds, _ := NewDiskStore(dir)
	_, err = ds.Get(SGFInput, "testfile.json")
	if err == nil {
		log.Fatal(err, "An error should have been generated for a file that didn't exist")
	}

}

func TestDiskStorePutAndGet(t *testing.T) {
	dir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	ds, _ := NewDiskStore(dir)
	contents := "{\"some\": \"value\"}"
	ds.Put(SGFInput, "testfile.json", contents)

	result, err := ds.Get(SGFInput, "testfile.json")
	if err != nil {
		log.Fatal(err)
	}
	if result != contents {
		t.Errorf("Stored and retrieve result not equal to expected '%v' '%v'", result, contents)
	}
}

func TestEnsurePathExistsCreatesDir(t *testing.T) {
	tmp := os.TempDir()
	root := fmt.Sprintf("%s%s%d", tmp, "otrego_data_test", rand.Int())
	defer os.RemoveAll(root)

	err := os.Mkdir(root, 0755)
	if err != nil {
		log.Fatal(err, " Could not make a directory in temp dir")
	}
	ds, err := NewDiskStore(root)
	if err != nil {
		log.Fatal(err)
	}

	err = ds.ensureDirectoryStructure(root)
	if err != nil {
		log.Fatal(err, "Could not make the expected directories for otrego")

	}
	var files []string
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, dirName := range storedDataTypes {
		found := false
		for _, file := range files {
			if strings.Contains(file, string(dirName)) {
				found = true
			}
		}
		if !found {
			panic(fmt.Sprintf("could not find %s", dirName))
		}
	}
}

func TestEnsurePathExistsErrors(t *testing.T) {
	tmp := os.TempDir()
	root := fmt.Sprintf("%s%s%d", tmp, "otrego_data_test", rand.Int())

	_, err := NewDiskStore(root)
	if err == nil {
		log.Fatal("there should be an error when if dir doesn't exist")
	}
}
