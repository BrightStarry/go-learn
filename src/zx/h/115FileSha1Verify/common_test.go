package main

import (
	"testing"
	"fmt"
	"os"
	"crypto/sha1"
	"io"
	"encoding/hex"
)

func Test1(t *testing.T) {
	file, err := os.Open(`E:\m3u8Downloader\新建文本文档.txt`)
	if err != nil {
		return
	}
	defer file.Close()
	h := sha1.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return
	}
	fmt.Println(hex.EncodeToString(h.Sum(nil)))

}
