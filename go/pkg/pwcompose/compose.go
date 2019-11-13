package pwcompose

import (
	"fmt"

	"go.uber.org/zap"
)

func Prepare(path string, noPush bool, logger *zap.Logger) error {
	logger.Debug("prepare", zap.Bool("no-push", noPush), zap.String("path", path))
	// parse yaml
	// check yaml
	// create tmp yaml
	// exec docker-compose bundle --push-images
	// parse .dab image sha
	// create final yaml
	// cleanup
	// print yaml
	return fmt.Errorf("not implemented")
}

func Up(preparedCompose string, instanceKey string, logger *zap.Logger) error {
	logger.Debug(
		"up",
		zap.String("compose", preparedCompose),
		zap.String("instance-key", instanceKey),
	)
	// parse compose labels to get flavor info
	// generate instanceID with flavor+instanceKey
	// exec docker-compose up -d with params
	// print instanceID
	// later: print instance passphrases
	return fmt.Errorf("not implemented")
}

func Down(ids []string, logger *zap.Logger) error {
	logger.Debug("down", zap.Strings("ids", ids))
	// id can be an instance_id or a flavor_id
	return fmt.Errorf("not implemented")
}

func PS(depth int, logger *zap.Logger) error {
	logger.Debug("ps", zap.Int("depth", depth))
	return fmt.Errorf("not implemented")
}
