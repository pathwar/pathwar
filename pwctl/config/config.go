package config // import "pathwar.land/pwctl/config"

import "encoding/json"

type Config struct {
	Passphrases []string
}

func (c Config) String() string {
	out, _ := json.Marshal(c)
	return string(out)
}
