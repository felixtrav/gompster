// Package config loads server configuration from flags and environment variables.
// Precedence (highest to lowest): CLI flags > environment variables > defaults.
package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

// Config holds all tunable server parameters.
type Config struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	LogRequests  bool
}

// Load parses environment variables first, then overlays CLI flags.
func Load() *Config {
	c := &Config{
		Host:         envStr("HOST", "0.0.0.0"),
		Port:         envStr("PORT", "80"),
		ReadTimeout:  envDuration("READ_TIMEOUT", 30*time.Second),
		WriteTimeout: envDuration("WRITE_TIMEOUT", 30*time.Second),
		IdleTimeout:  envDuration("IDLE_TIMEOUT", 120*time.Second),
		LogRequests:  envBool("LOG_REQUESTS", true),
	}

	flag.StringVar(&c.Host, "host", c.Host, "listener host address (env: HOST)")
	flag.StringVar(&c.Port, "port", c.Port, "listener port (env: PORT)")
	flag.DurationVar(&c.ReadTimeout, "read-timeout", c.ReadTimeout, "HTTP read timeout (env: READ_TIMEOUT)")
	flag.DurationVar(&c.WriteTimeout, "write-timeout", c.WriteTimeout, "HTTP write timeout (env: WRITE_TIMEOUT)")
	flag.DurationVar(&c.IdleTimeout, "idle-timeout", c.IdleTimeout, "HTTP idle timeout (env: IDLE_TIMEOUT)")
	flag.BoolVar(&c.LogRequests, "log", c.LogRequests, "log incoming requests (env: LOG_REQUESTS)")
	flag.Parse()

	return c
}

func envStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func envDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}
