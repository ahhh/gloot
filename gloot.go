package gloot

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/alexmullins/zip"
)

// Keyz global string array to track files
var Keyz []string
var ignoreNames = []string{}
var ignoreContent = []string{}
var includeNames = []string{}
var includeContent = []string{}

// Searcher is a public function designed to be called as a package
func Searcher(pathToDir string, igNames, igContent, inNames, inContent []string) []string {
	ignoreNames = igNames
	ignoreContent = igContent
	includeNames = inNames
	includeContent = inContent
	// Start recursive search
	searchForFiles(pathToDir)
	if Keyz != nil {
		return Keyz
	}
	return nil
}

// searchForFiles is a private function that recurses through directories, running our searchFileForCriteria function on every file
func searchForFiles(pathToDir string) {
	files, err := ioutil.ReadDir(pathToDir)
	if err != nil {
		return
	}
	for _, file := range files {
		if stringLooper(file.Name(), ignoreNames) {
		} else {
			if file.IsDir() {
				dirName := file.Name() + "/"
				fullPath := strings.Join([]string{pathToDir, dirName}, "")
				searchForFiles(fullPath)
			} else {
				if searchFileForCriteria(pathToDir, file.Name()) {
					fullPath := strings.Join([]string{pathToDir, file.Name()}, "")
					Keyz = append(Keyz, fullPath)
				}
			}
		}
	}
}

func searchFileForCriteria(pathToDir, fileName string) bool {
	fullPath := strings.Join([]string{pathToDir, fileName}, "")
	if stringLooper(fullPath, includeNames) {
		return true
	}
	fileData, _ := ioutil.ReadFile(fullPath)
	fileLines := strings.Split(string(fileData), "\n")
	for _, line := range fileLines {
		if stringLooper(line, ignoreContent) {
			return false
		}
		if stringLooper(line, includeContent) {
			return true
		}
	}
	return false
}

// A function to loop over our string slices and match any of our globally defined content
func stringLooper(target string, list []string) bool {
	for _, loot := range list {
		if strings.Contains(target, loot) {
			return true
		}
	}
	return false
}

// ZipFiles compresses one or many files into a single zip archive file
func ZipFiles(filename string, files []string, encryptPassword string) error {
	newfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newfile.Close()
	zipWriter := zip.NewWriter(newfile)
	defer zipWriter.Close()
	encryptedWriter, err := zipWriter.Encrypt("", encryptPassword)
	if err != nil {
		return err
	}
	encryptedZipWriter := zip.NewWriter(encryptedWriter)
	for _, file := range files {
		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Method = zip.Deflate
		writer, err := encryptedZipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, zipfile)
		if err != nil {
			return err
		}
	}
	return nil
}
