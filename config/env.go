package config

import "github.com/kelseyhightower/envconfig"

type Env struct {
	S3Bucket     string `envconfig:"S3_BUCKET" required:"true"`
	S3Prefix     string `envconfig:"S3_PREFIX" default:""`
	AwsAccessKey string `envconfig:"AWS_ACCESS_KEY_ID" required:"true"`
	AwsSecretKey string `envconfig:"AWS_SECRET_ACCESS_KEY" required:"true"`
	InputFile    string `envconfig:"INPUT_FILE" required:"true"`
}

func Get() Env {
	env := Env{}
	envconfig.MustProcess("", &env)
	return env
}
