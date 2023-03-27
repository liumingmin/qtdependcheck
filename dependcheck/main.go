package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "args must more than 2, clientDir, md5file")
		os.Exit(-1)
		return
	}
	clientDir := os.Args[1]
	md5file := os.Args[2]

	clientMap := make(map[string]string)
	buildDirAllFiles(clientDir, clientMap)

	md5Map := make(map[string]string)
	loadMd5File(md5file, md5Map)

	if len(clientMap) == 0 || len(md5Map) == 0 {
		fmt.Fprintln(os.Stderr, "dir no files  or md5file is empty")
		os.Exit(-1)
		return
	}

	for key, value := range clientMap {
		if filemd5, ok := md5Map[key]; ok {
			if filemd5 != value {
				fmt.Fprintln(os.Stderr, key+" md5 is not equal "+value+":"+filemd5)
				os.Exit(-2)
			}
		}
	}

	fmt.Fprintln(os.Stdout, "depend file is equal")
}

func buildDirAllFiles(dir string, filesMap map[string]string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path) //path是相对路径
		if err != nil {
			return nil
		}
		defer file.Close()
		value, _ := checksumMd5(file)

		fullPath, _ := filepath.Abs(path)
		shortPath, _ := filepath.Rel(dir, fullPath)
		shortPath = strings.ReplaceAll(shortPath, "\\", "/")

		filesMap[shortPath] = value
		return nil
	})
}

func checksumMd5(fileReader io.Reader) (string, error) {
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, fileReader); err != nil {
		return "", err
	}
	md5Val := md5hash.Sum(nil)
	hexVal := hex.EncodeToString(md5Val)
	return hexVal, nil
}

func loadMd5File(md5filepath string, fileMd5Map map[string]string) {
	bs, err := ioutil.ReadFile(md5filepath)
	if err != nil {
		return
	}

	lines := strings.Split(string(bs), "\n")
	for _, line := range lines {
		fileAndMd5 := strings.Split(line, "|")
		if len(fileAndMd5) < 2 {
			continue
		}

		filePath := fileAndMd5[0]
		md5 := fileAndMd5[1]

		fileMd5Map[filePath] = md5
	}
}
