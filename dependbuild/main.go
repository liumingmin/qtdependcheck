package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "args must more than 2, clientDIr, qtDir")
		return
	}
	clientDir := os.Args[1]
	qtDir := os.Args[2]

	outputPath := filepath.Join(GetCurrPath(), "result.txt")

	clientMap := make(map[string]string)
	buildDirAllFiles(clientDir, clientMap)

	qtMap := make(map[string]string)
	buildDirAllFiles(filepath.Join(qtDir, "qml"), qtMap)
	buildDirAllFiles(filepath.Join(qtDir, "plugins"), qtMap)
	buildDirAllFiles(filepath.Join(qtDir, "bin"), qtMap)

	if len(clientMap) == 0 || len(qtMap) == 0 {
		fmt.Fprintln(os.Stderr, "dir no files")
		return
	}

	outputStrs := make([]string, 0)
	for key, value := range clientMap {
		if _, ok := qtMap[key]; ok {
			outputStrs = append(outputStrs, key+"|"+value)
		}
	}
	sort.Strings(outputStrs)

	os.MkdirAll(filepath.Dir(outputPath), 0666)
	saveMd5File, _ := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer saveMd5File.Close()

	saveMd5File.WriteString(strings.Join(outputStrs, "\n"))

	fmt.Fprintln(os.Stdout, "depend check file is generated: "+outputPath)
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

func GetCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	return filepath.Dir(path)
}
