package main

import (
	"io/ioutil"
	"io"
	"log"
	"path"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"strconv"
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

// DownloadFileIntoSubFolder downloads file into specific subfolder from the url
func (storage *Storage) DownloadFileIntoSubFolder(url string, name string, subfolder string) (err error) {
	folder := path.Join(storage.config.Folder, subfolder)
	return downloadFileIntoSubFolder(url, name, folder)
}

// DownloadFileIntoFolder downloads file from the url
func (storage *Storage) DownloadFileIntoFolder(url string, name string) (err error) {
	folder := storage.config.Folder
	return downloadFileIntoSubFolder(url, name, folder)
}


func downloadFileIntoSubFolder(url string, name string, folder string) (err error) {
	// Download File
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
	defer resp.Body.Close()
	

	// Find the allowed name
	fullPath := path.Join(folder, name)
	_, err = os.Stat(fullPath)

	index := 0

	for {
		if (!os.IsNotExist(err)){
			break
		}

		ext := filepath.Ext(name)
		basename := strings.TrimSuffix(name, ext)
		index++
		name = basename + "_" + strconv.Itoa(index) + ext

		fullPath = path.Join(folder, name)
		_, err = os.Stat(fullPath)
	}

	// Create file
    out, err := os.Create(fullPath)
    if err != nil {
        return err
    }
    defer out.Close()

	// Write file
    _, err = io.Copy(out, resp.Body)
    return err
}

