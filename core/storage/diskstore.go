package storage

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// DiskStore is an on-disk filestore implementation that is particularly useful
// for development
type DiskStore struct {
	rootDir string
}

// NewDiskStore returns a new Filestore that is on-disk
func NewDiskStore(root string) (*DiskStore, error) {
	ds := &DiskStore{
		rootDir: root,
	}
	err := ds.ensureDirectoryStructure(root)
	if err != nil {
		return nil, err
	}
	return ds, nil
}

// Get abstract method to get JSON content from a Filestore
func (ds *DiskStore) Get(t StoredDataType, filename string) (string, error) {
	path := ds.path(t, filename)
	_, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// List is method to list available files
func (ds *DiskStore) List(t StoredDataType) ([]string, error) {
	ioutil.ReadDir(ds.rootDir)
	return []string{""}, nil
}

// Put is method to Put a file to disk
func (ds *DiskStore) Put(t StoredDataType, filename string, json string) error {
	path := ds.path(t, filename)
	return ioutil.WriteFile(path, []byte(json), os.ModePerm)
}

func (ds *DiskStore) path(t StoredDataType, filename string) string {
	return fmt.Sprintf("%s/%s/%s", ds.rootDir, t, filename)
}

// ensureDirectoryStructure ensures that paths exist for each
// of the expected file outputs
func (ds *DiskStore) ensureDirectoryStructure(rootDir string) error {
	fileInfo, err := os.Stat(rootDir)
	if err != nil {
		return err
	}
	if !fileInfo.IsDir() {
		return errors.New("Directory does not exist")
	}
	for _, t := range storedDataTypes {
		curDir := fmt.Sprintf("%s/%s", ds.rootDir, t)
		fileInfo, err := os.Stat(curDir)
		if err != nil {
			os.Mkdir(curDir, 0755)
		} else if !fileInfo.IsDir() {
			return errors.New("file is in place of directory")
		}
	}
	return nil
}
