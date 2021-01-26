package storage

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// DiskStore is an on-disk filestore implementation that is particularly useful
// for development
type DiskStore struct {
	rootDir string
}

// DefaultDirPerms is the Default directory permissions expected by DiskStore
const DefaultDirPerms = 0755

// DefaultFilePerms is the Default file permissions expected by DiskStore
const DefaultFilePerms = 0644

// NewDiskStore returns a new Filestore that is on-disk
func NewDiskStore(root string) (*DiskStore, error) {
	ds := &DiskStore{
		rootDir: root,
	}
	err := ds.makeGenDirs()
	if err != nil {
		return nil, err
	}
	return ds, nil
}

// makeGenDirs ensures that paths exist for each of the expected file outputs
func (ds *DiskStore) makeGenDirs() error {
	fileInfo, err := os.Stat(ds.rootDir)
	if err != nil {
		return err
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("root directory %v does not exist", ds.rootDir)
	}
	for _, t := range storedDataTypes {
		curDir := path.Join(ds.rootDir, string(t))
		fileInfo, err := os.Stat(curDir)
		if err != nil {
			os.Mkdir(curDir, DefaultDirPerms)
		} else if !fileInfo.IsDir() {
			return fmt.Errorf("non-directory is in place of directory %v", curDir)
		}
	}
	return nil
}

// Get is a method to retrieve string content from a file on Filestore
func (ds *DiskStore) Get(ctx context.Context, t StoredDataType, filename string) (string, error) {
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
func (ds *DiskStore) List(ctx context.Context, t StoredDataType) ([]string, error) {
	files, err := ioutil.ReadDir(ds.rootDir)
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(files))
	for i := 0; i < len(files); i++ {
		ret[i] = files[i].Name()
	}

	return ret, nil
}

// Put is method to Put a file to disk
func (ds *DiskStore) Put(ctx context.Context, t StoredDataType, filename string, contents string) error {
	p := ds.path(t, filename)
	return ioutil.WriteFile(p, []byte(contents), DefaultFilePerms)
}

func (ds *DiskStore) path(t StoredDataType, filename string) string {
	return path.Join(ds.rootDir, string(t), filename)
}

var _ Filestore = &DiskStore{}
