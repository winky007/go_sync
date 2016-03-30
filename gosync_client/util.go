package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"math"
	"os"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func fileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func md5String(str string) (string, error) {
	h := md5.New()
	io.WriteString(h, str)
	return hex.EncodeToString(h.Sum(nil)), nil
}

func getFileContent(filePath string) (string, error) {
	var content string

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var fileSize int64 = fileInfo.Size()
	//fmt.Println(fileSize)
	if fileSize < 1 {
		return "", nil
	}

	const fileChunk = 5 * (1 << 20) //5MB
	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))
	//fmt.Println("total num:", totalPartsNum)
	for i := uint64(1); i <= totalPartsNum; i++ {
		partSize := 0
		if i == totalPartsNum {
			partSize = int(float64(fileSize) - float64((totalPartsNum-1)*fileChunk))
		} else {
			partSize = int(fileChunk)
		}
		partBuffer := make([]byte, partSize)
		file.Read(partBuffer)
		content = content + string(partBuffer)
	}

	return content, nil
}
