package config

type Config struct {
	Grpc GrpcConfig `yaml:"grpc"`
}

type GrpcConfig struct {
	ShortenerAddr string `yaml:"shortener_addr"`
	CacheAddr     string `yaml:"cache_addr"`
}
