package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/pathwar/go-pathwar/pkg/pathwar"
)

func main() {
	ec, err := pathwar.Start(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
	os.Exit(ec)
}
