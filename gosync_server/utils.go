package main

import (
	"crypto/md5"
	"encoding/hex"
	//"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func md5File(filePath string) (string, error) {
	var rst string
	file, err := os.Open(filePath)
	if err != nil {
		return rst, err
	}
	defer file.Close()

	h := md5.New()
	if _, err = io.Copy(h, file); err != nil {
		return rst, err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func writeFile(file, content string) bool {
	if !fileExist(file) {
		_, err := os.Create(file)
		if err != nil {
			log.Fatal(err)
		}
	}

	ioutil.WriteFile(file, []byte(content), 0777)
	return true
}

func fileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
