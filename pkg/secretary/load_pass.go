package secretary

import (
	"context"
	"os"

	"github.com/gopasspw/gopass/pkg/gopass/api"
)

func PassSetStoreDir(dir string) {
	os.Setenv("PASSWORD_STORE_DIR", dir)
}

func PassSetConfig(conf string) {
	os.Setenv("GOPASS_CONFIG", conf)
}

func PassLoadSecret(name string) (Environment, error) {
	pstore, err := api.New(context.Background())
	if err != nil {
		return nil, err
	}
	secret, err := pstore.Get(context.Background(), name, "latest")
	if err != nil {
		return nil, err
	}
	data := string(secret.Bytes())
	return ParseEnvironmentFile(name, data)
}
