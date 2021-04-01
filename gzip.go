package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"log"
)

func gUnzipData(data []byte) (resData []byte, err error) {
	b := bytes.NewBuffer(data)

	var r io.Reader
	r, err = gzip.NewReader(b)
	if err != nil {
		return
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return
	}

	resData = resB.Bytes()

	return
}

func gUnzipDataNew(data []byte) []byte {

	szTempData, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		log.Println("Payload not gzip")
		return data
	}
	szTemp2, err := gUnzipData(szTempData)
	if err != nil {
		log.Println("Payload not gzip")
		return data
	}
	return szTemp2
}
