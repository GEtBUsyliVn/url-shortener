package config

type Config struct {
	GRPC     GrpcConfig     `yaml:"grpc"`
	Database DataBaseConfig `yaml:"database"`
}

type GrpcConfig struct {
	Addr string `yaml:"addr"`
}

type DataBaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DataBase string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
