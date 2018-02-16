package main

import (
	"os"
	"io"
	"github.com/satori/go.uuid"
)

func uploadFile(file io.Reader) (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	err = os.MkdirAll("data/", 0777)
	if err != nil {
		return "", err
	}
	f, err := os.OpenFile("data/"+uid.String(), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	io.Copy(f, file)

	return uid.String(), err
}
