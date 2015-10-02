package pathwar

import (
	"os"

	"github.com/pathwar/go-pathwar/pkg/api"
)

func Start(argv []string) (int, error) {
	cli := api.NewAPIPathwar(os.Getenv("PATHWAR_TOKEN"), os.Getenv("PATHWAR_DEBUG"))
	_, err := cli.PostResquest("organization-coupons", map[string]string{"organization": "283e4414-34cc-472c-91b5-4b7ed0cf8d92", "coupon": "poeut"})
	if err != nil {
		return 1, err
	}
	return 0, nil
}
