package storage

import (
	"context"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestDiskGetNoFileErrors(t *testing.T) {
	ctx := context.Background()
	dir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	ds, err := NewDiskStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ds.Get(ctx, Games, "testfile.json")
	if err == nil {
		t.Fatal(err)
	}
}

func TestDiskStorePutAndGet(t *testing.T) {
	ctx := context.Background()
	dir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	ds, err := NewDiskStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	contents := "{\"some\": \"value\"}"
	ds.Put(nil, Games, "testfile.json", contents)

	result, err := ds.Get(ctx, Games, "testfile.json")
	if err != nil {
		t.Fatal(err)
	}
	if result != contents {
		t.Errorf("Stored and retrieve result not equal to expected '%v' '%v'", result, contents)
	}
}

func TestEnsurePathExistsCreatesDir(t *testing.T) {
	tmp := os.TempDir()
	root := path.Join(tmp, "otrego_data_test"+strconv.Itoa(rand.Int()))
	defer os.RemoveAll(root)

	err := os.Mkdir(root, DefaultDirPerms)
	if err != nil {
		t.Fatal(err, " Could not make a directory in temp dir")
	}
	_, err = NewDiskStore(root)

	if err != nil {
		t.Fatal(err, "Could not make the expected directories for otrego")
	}
	var files []string
	if err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	for _, dirName := range storedDataTypes {
		found := false
		for _, file := range files {
			if strings.Contains(file, string(dirName)) {
				found = true
			}
		}
		if !found {
			t.Fatalf("could not find %s", dirName)
		}
	}
}

func TestEnsurePathExistsErrors(t *testing.T) {
	tmp := os.TempDir()
	root := path.Join(tmp, "nonexistent_"+strconv.Itoa(rand.Int()))

	d := &DiskStore{rootDir: root}

	if err := d.makeGenDirs(); err == nil {
		t.Fatal("got nil error, there should be an error when if dir doesn't exist")
	}
}
