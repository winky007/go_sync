package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
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
