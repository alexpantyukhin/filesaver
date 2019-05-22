package main

import (
	"io/ioutil"
	"log"
	"path"
)

// Config for storage
type Config struct {
	Folder string
}

// Storage sontains methods for store files.
type Storage struct {
	config Config
}

// GetInnerFolders returns the slice of the inner folders
func (storage *Storage) GetInnerFolders() []string {
	files, err := ioutil.ReadDir(storage.config.Folder)

	if err != nil {
		log.Fatal(err)
	}

	res := make([]string, len(files))

	for i, f := range files {
		res[i] = f.Name()
	}

	return res
}

// PutFileIntoFolder puts file into folder
func (storage *Storage) PutFileIntoFolder(content []byte, name string) (err error) {
	return storage.putFileIntoSubFolder(content, name, storage.config.Folder)
}

// PutFileIntoSubFolder puts file into subfolder
func (storage *Storage) PutFileIntoSubFolder(content []byte, name string, subfolder string) (err error) {
	folder := path.Join(storage.config.Folder, subfolder)
	return storage.putFileIntoSubFolder(content, name, folder)
}

func (storage *Storage) putFileIntoSubFolder(content []byte, name string, folder string) (err error) {
	err = ioutil.WriteFile(folder, content, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return err
}
