package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Name          string `yaml:"name"`
		BaseURL       string `yaml:"base_url"`
		WebAppBaseURL string `yaml:"web_app_base_url"`
	} `yaml:"app"`

	Http struct {
		Port string `yaml:"port"`
	} `yaml:"http"`

	Vault struct {
		MountPath string `yaml:"mount_path"`
	} `yaml:"vault"`

	DB struct {
		Name            string `yaml:"name"`
		Host            string `yaml:"host"`
		Port            string `yaml:"port"`
		User            string `yaml:"user"`
		Password        string `yaml:"password"`
		TLS             bool   `yaml:"tls"`
		MultiStatements string `yaml:"multi_statements"`
	} `yaml:"db"`

	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		Secret   string `yaml:"secret"`
	}

	Session struct {
		MaxAge int `yaml:"max_age"`
	} `yaml:"session"`

	Log struct {
		Dir string `yaml:"dir"`
	} `yaml:"log"`

	Jwt struct {
		Secret     string `yaml:"secret"`
		Timeout    int    `yaml:"timeout"`
		MaxRefresh int    `yaml:"max_refresh"`
	} `yaml:"jwt"`

	Email struct {
		EmailSender string `yaml:"email_sender"`
		// TemplateUUID_TournamentRegistrationSuccess string `yaml:"template_uuid_tournament_registration_success"`
		// TemplateUUID_CoachingRegistrationSuccess   string `yaml:"template_uuid_coaching_registration_success"`
	}

	MailTrap struct {
		ApiURL string `yaml:"api_url"`
		Token  string `yaml:"token"`
	} `yaml:"mailtrap"`

	Minio struct {
		URL       string `yaml:"url"`
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
		Bucket    string `yaml:"bucket"`
	} `yaml:"minio"`

	OAuth struct {
		GoogleClientID    string `yaml:"google_client_id"`
		GoogleSecret      string `yaml:"google_secret"`
		GoogleCallbackURL string `yaml:"google_callback_url"`

		FacebookAppID       string `yaml:"facebook_app_id"`
		FacebookSecret      string `yaml:"facebook_secret"`
		FacebookCallbackURL string `yaml:"facebook_callback_url"`
		FacebookConfigID    string `yaml:"facebook_config_id"`
	} `yaml:"oauth"`

	Duitku struct {
		ApiURL       string `yaml:"api_url"`
		ApiKey       string `yaml:"api_key"`
		MerchantCode string `yaml:"merchant_code"`
	} `yaml:"duitku"`

	Google struct {
		ServiceAccountCreds string `yaml:"service_account_creds"`
		RecaptchaSiteKey    string `yaml:"recaptcha_site_key"`
		RecaptchaSecretKey  string `yaml:"recaptcha_secret_key"`
	} `yaml:"google"`
}

func getConfigPath() string {
	if os.Getenv("ENV") == "" || os.Getenv("ENV") == "local" {
		return "config/config.local.yml"
	}

	if os.Getenv("ENV") == "dev" || os.Getenv("ENV") == "development" {
		return "config/config.dev.yml"
	}

	if os.Getenv("ENV") == "prod" || os.Getenv("ENV") == "production" {
		return "config/config.prod.yml"
	}

	if os.Getenv("ENV") == "test" {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			return ""
		}
		dirname := filepath.Dir(filename)
		return dirname + "/config.test.yml"
	}

	return "config/config.local.yml"
}

func NewConfig() (*Config, error) {
	configPath := getConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile(): failed to read config.yml: %w", err)
	}

	d := yaml.NewDecoder(strings.NewReader(string(data)))
	d.KnownFields(true)

	c := &Config{}
	err = d.Decode(c)
	if err != nil {
		return nil, fmt.Errorf("yaml.Unmarshal(): failed to unmarshal config.yml: %w", err)
	}

	// this is used when doing migrations
	if c.DB.MultiStatements == "" {
		c.DB.MultiStatements = "false"
	}

	err = loadSecrets(c)
	if err != nil {
		return nil, fmt.Errorf("NewConfig(): failed to load secrets: %w", err)
	}

	return c, nil
}

func loadSecrets(c *Config) error {
	vaultSecrets, err := LoadVaultSecrets(c.Vault.MountPath)
	if err != nil {
		return fmt.Errorf("[config.config.loadSecrets()]: failed to load vault secrets: %w", err)
	}

	configValue := reflect.ValueOf(c).Elem()

	missingEnvVars := []string{}

	for i := 0; i < configValue.NumField(); i++ {
		fieldName := configValue.Type().Field(i).Name
		fieldValue := configValue.Field(i)

		if fieldValue.Kind() == reflect.String {
			if strings.HasPrefix(fieldValue.String(), "$") {
				envVarName := strings.ToUpper(fieldName)
				envVarValue := os.Getenv(envVarName)

				if envVarValue != "" {
					fieldValue.SetString(envVarValue)
				}
			}

			continue
		}

		if fieldValue.Kind() == reflect.Struct {
			// parentFieldName := strings.ToUpper(fieldName)
			parentFieldValue := configValue.Field(i)

			for j := 0; j < parentFieldValue.NumField(); j++ {
				// fieldName := parentFieldValue.Type().Field(j).Tag.Get("yaml")
				fieldValue := parentFieldValue.Field(j)

				if strings.HasPrefix(fieldValue.String(), "$") {
					envVarName := strings.ToUpper(strings.TrimPrefix(fieldValue.String(), "$"))

					val, ok := vaultSecrets[envVarName]
					if !ok {
						missingEnvVars = append(missingEnvVars, envVarName)
						continue
					}

					valStr, ok := val.(string)
					if !ok {
						missingEnvVars = append(missingEnvVars, envVarName)
						continue
					}

					fieldValue.SetString(valStr)
				}
			}
		}
	}

	if len(missingEnvVars) > 0 {
		return fmt.Errorf("loadSecrets(): missing environment variables: %v", "\n"+strings.Join(missingEnvVars, "\n"))
	}

	return nil
}
