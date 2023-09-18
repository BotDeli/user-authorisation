package config

type Config struct {
	Grpc *GRPCConfig `yaml:"grpc" env-required:"true"`
}

type GRPCConfig struct {
	Network string `yaml:"network" env-default:"tcp"`
	Address string `yaml:"address" env-required:"true"`
}
