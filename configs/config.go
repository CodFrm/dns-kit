package configs

import (
	"context"
	"github.com/codfrm/cago/configs"
)

func EncryptKey() string {
	return configs.Default().String(context.Background(), "encrypt_key")
}
