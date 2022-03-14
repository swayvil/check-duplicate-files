package main

import (
	"fmt"
	"hash/crc32"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var crc32q *crc32.Table
var filesMap map[uint32][]string

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("usage: check-duplicate-files <root directory path>")
	} else {
		srcDirPath := os.Args[1]
		crc32q = crc32.MakeTable(0xD5828281)
		filesMap = make(map[uint32][]string)
		filepath.WalkDir(srcDirPath, walk)
	}
	outputMap()
}

func walk(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() {
		fileBuffer, err := ioutil.ReadFile(s)
		if err != nil {
			log.Fatalf("Failed opening file: %s", err)
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

func outputMap() {
	for checkSum, paths := range filesMap {
		for index, path := range paths {
			fmt.Printf("%08x;%d;%s\n", checkSum, index, path)
		}
	}
}
