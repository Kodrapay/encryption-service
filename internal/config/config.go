package config

import "os"

type Config struct {
	ServiceName   string
	Port          string
	EncryptionKey string
}

func Load(serviceName, defaultPort string) Config {
	return Config{
		ServiceName:   serviceName,
		Port:          getEnv("PORT", defaultPort),
		EncryptionKey: getEnv("ENCRYPTION_KEY", "default-encryption-key-change-me-32bytes"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
