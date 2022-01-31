package secretary

import (
	"os"
)

func PlainLoadSecret(name string) (Environment, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return ParseEnvironmentFile(name, string(data))
}
