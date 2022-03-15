package main

import (
	"bufio"
	"errors"
	"fmt"
	"hash/crc32"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var crc32q *crc32.Table
var filesMap map[uint32][]string
var isRemoveDuplicates bool

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("usage: check-duplicate-files <root directory path> [remove]")
	} else {
		srcDirPath := os.Args[1]
		crc32q = crc32.MakeTable(0xD5828281)
		filesMap = make(map[uint32][]string)
		isRemoveDuplicates = false

		if len(os.Args) > 2 && os.Args[2] == "remove" {
			isRemoveDuplicates = true
		}

		filepath.WalkDir(srcDirPath, walk)
		duplicates := outputMap()

		if isRemoveDuplicates {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("\n\nAre you sure you want to move %d duplicates to the directory \"REMOVED\" [Y/N] ", duplicates)
			text, _ := reader.ReadString('\n')
			if strings.Compare(text, "Y\n") == 0 {
				moveDuplicates(srcDirPath)
			}
		}
	}
}

func walk(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() {
		fileBuffer, err := ioutil.ReadFile(s)
		if err != nil {
			log.Fatalf("Failed to open file: %s", err)
		}
		checkSum := crc32.Checksum(fileBuffer, crc32q)
		paths, exist := filesMap[checkSum]
		if exist {
			paths = append(paths, s)
			filesMap[checkSum] = paths
		} else {
			newPaths := []string{s}
			filesMap[checkSum] = newPaths
		}
	}
	return nil
}

func outputMap() int {
	duplicates := 0

	for checkSum, paths := range filesMap {
		if len(paths) > 1 { // Output only paths with duplicates
			paths, err := bubbleSortPaths(paths)
			if err != nil {
				log.Println(err) // Non fatal error
			} else {
				filesMap[checkSum] = paths // Update map, so paths are already sorted for remove function
				for index, path := range paths {
					fmt.Printf("%08x;%d;%s\n", checkSum, index+1, path)
					if index > 0 {
						duplicates++
					}
				}
			}
		}
	}
	return duplicates
}

func convertPathToTime(path string) (time.Time, error) {
	layout := "2006-01-02_15.04.05"
	filename := filepath.Base(path)

	if len(filename) < 19 {
		return time.Time{}, errors.New("Convertion failed, path too short " + filename)
	}

	t, err := time.Parse(layout, filename[0:19])
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func bubbleSortPaths(paths []string) ([]string, error) {
	for i := 0; i < len(paths)-1; i++ {
		for j := 0; j < len(paths)-i-1; j++ {
			timeJ, errJ := convertPathToTime(paths[j])
			if errJ != nil {
				return nil, errJ
			}
			timeJ1, errJ1 := convertPathToTime(paths[j+1])
			if errJ1 != nil {
				return nil, errJ1
			}
			if timeJ.After(timeJ1) {
				paths[j], paths[j+1] = paths[j+1], paths[j]
			}
		}
	}
	return paths, nil
}

func moveDuplicates(srcDirPath string) {
	removedDirectory := filepath.Join(srcDirPath, "REMOVED")

	if _, err := os.Stat(removedDirectory); errors.Is(err, os.ErrNotExist) { // Folder doesn't exist
		err := os.Mkdir(removedDirectory, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, paths := range filesMap {
		if len(paths) > 1 { // Move only paths with duplicates
			for index, path := range paths {
				if index > 0 {
					filename := filepath.Base(path)
					err := os.Rename(path, filepath.Join(removedDirectory, filename))
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}
}
