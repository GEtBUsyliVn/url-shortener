package config

import "time"

type GrpcConfig struct {
	Addr string `yaml:"addr"`
}

type RedisConfig struct {
	Address  string        `yaml:"addr"`
	Password string        `yaml:"password"`
	DB       int           `yaml:"db"`
	CacheTTl time.Duration `yaml:"redis_cache_ttl"`
}

type MemoryStorageConfig struct {
	CacheTTl time.Duration `yaml:"memory_cache_ttl"`
}

type WorkerConfig struct {
	Interval time.Duration `yaml:"cleanup_interval"`
}

type Config struct {
	Grpc          GrpcConfig          `yaml:"grpc"`
	Redis         RedisConfig         `yaml:"redis"`
	MemoryStorage MemoryStorageConfig `yaml:"memory_storage"`
	Worker        WorkerConfig        `yaml:"worker"`
}
