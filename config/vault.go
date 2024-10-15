package config

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/vault-client-go"
)

type SecretPath struct {
	Path      string
	MountPath string
}

var SECRET_PATHS = []string{
	"db",
	"duitku",
	"jwt",
	"mailtrap",
	"minio",
	"oauth",
	"redis",
	"google",
}

func LoadVaultSecrets(mountPath string) (map[string]interface{}, error) {
	ctx := context.Background()

	VAULT_URL := os.Getenv("VAULT_URL")
	VAULT_TOKEN := os.Getenv("VAULT_TOKEN")

	client, err := vault.New(
		vault.WithAddress(VAULT_URL),
	)
	if err != nil {
		return nil, fmt.Errorf("[config.secret.LoadVaultSecrets()]: failed to create vault client: %w", err)
	}

	err = client.SetToken(VAULT_TOKEN)
	if err != nil {
		return nil, fmt.Errorf("[config.secret.LoadVaultSecrets()]: failed to set token: %w", err)
	}

	secrets := make(map[string]interface{})

	for _, path := range SECRET_PATHS {
		secret, err := client.Secrets.KvV2Read(ctx, path, vault.WithMountPath(mountPath))
		if err != nil {
			return nil, fmt.Errorf("[config.secret.LoadVaultSecrets()]: failed to read secret: %w", err)
		}

		for k, v := range secret.Data.Data {
			secrets[k] = v
		}
	}

	return secrets, nil
}
