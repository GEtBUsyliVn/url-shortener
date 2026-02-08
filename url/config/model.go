package config

type Config struct {
	DataBase DataBaseConfig `yaml:"database"`
	Grpc     GrpcConfig     `yaml:"grpc"`
}

type DataBaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DataBase string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type GrpcConfig struct {
	Port string `yaml:"port"`
}
