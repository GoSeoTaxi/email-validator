package config

import (
	"flag"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	GRPCPort       string
	DNSHosts       []string
	RedisHost      string
	RedisPort      string
	RedisDB        int
	RedisMaxMemory string
}

func LoadConfig() *Config {

	grpcPort := flag.String("grpc-port", "50051", "gRPC Port")
	dnsHosts := flag.String("dns-hosts", "1.1.1.1,1.0.0.1", "DNS servers")
	redisHost := flag.String("redis-host", "localhost", "Redis host")
	redisPort := flag.String("redis-port", "6379", "Redis port")
	redisDB := flag.Int("redis-db", 0, "Redis DB")
	redisMaxMemory := flag.String("redis-maxmemory", "100mb", "Maximum Redis memory size (for example, 100mb - approximately 1M records")

	flag.Parse()

	viper.AutomaticEnv()

	return &Config{
		GRPCPort:       getStringValue("GRPC_PORT", *grpcPort),
		DNSHosts:       strings.Split(getStringValue("DNS_HOSTS", *dnsHosts), ","),
		RedisHost:      getStringValue("REDIS_HOST", *redisHost),
		RedisPort:      getStringValue("REDIS_PORT", *redisPort),
		RedisDB:        getIntValue("REDIS_DB", *redisDB),
		RedisMaxMemory: getStringValue("REDIS_MAXMEMORY", *redisMaxMemory),
	}
}

func getStringValue(envKey string, defaultValue string) string {
	if v := viper.GetString(envKey); v != "" {
		return v
	}
	return defaultValue
}

func getIntValue(envKey string, defaultValue int) int {
	if v := viper.GetInt(envKey); v != 0 {
		return v
	}
	return defaultValue
}
