package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/pathwar/go-pathwar/pkg/api"
)

func main() {
	cli := api.NewAPIPathwar(os.Getenv("PATHWAR_TOKEN"), os.Getenv("PATHWAR_DEBUG"))
	_, err := cli.PostRequest("organization-coupons", map[string]string{"organization": "283e4414-34cc-472c-91b5-4b7ed0cf8d92", "coupon": "poeut"})
	if err != nil {
		logrus.Fatal(err)
	}
}
