package main

import "github.com/kelseyhightower/envconfig"

type Config struct {
	ServerPort int `required:"true" split_words:"true"`

	DBDatabase string `required:"true" split_words:"true"`

	PrimaryDBHost     string `required:"true" split_words:"true"`
	PrimaryDBPort     string `required:"true" split_words:"true"`
	PrimaryDBPassword string `required:"true" split_words:"true"`
	PrimaryDBUsername string `required:"true" split_words:"true"`

	ReplicaDBHost     string `required:"true" split_words:"true"`
	ReplicaDBPort     string `required:"true" split_words:"true"`
	ReplicaDBPassword string `required:"true" split_words:"true"`
	ReplicaDBUsername string `required:"true" split_words:"true"`
}

func LoadConfig() (*Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
