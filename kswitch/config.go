package kswitch

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func defaultConfig() Config {
	return Config{
		Log: LogConfig{
			Level:     slog.LevelInfo,
			Format:    "json",
			AddSource: false,
		},
		Server: ServerConfig{
			ListenAddr: ":8080",
		},
	}
}

func LoadConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	cfg := defaultConfig()
	if err = toml.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to decode config file: %w", err)
	}
	return cfg, nil
}

type Config struct {
	Log     LogConfig          `toml:"log"`
	Server  ServerConfig       `toml:"server"`
	Configs []KillSwitchConfig `toml:"configs"`
}

func (c Config) String() string {
	return fmt.Sprintf("\n log: %v\n server: %v\n configs: %v",
		c.Log,
		c.Server,
		c.Configs,
	)
}

type LogConfig struct {
	Level     slog.Level `toml:"level"`
	Format    string     `toml:"format"`
	AddSource bool       `toml:"add_source"`
}

func (l LogConfig) String() string {
	return fmt.Sprintf("\n  level: %s\n  format: %s\n  add_source: %t",
		l.Level.String(),
		l.Format,
		l.AddSource,
	)
}

type ServerConfig struct {
	ListenAddr string `toml:"listen_addr"`
}

func (s ServerConfig) String() string {
	return fmt.Sprintf("\n  listen_addr: %s",
		s.ListenAddr,
	)
}

type KillSwitchConfig struct {
	Name     string `toml:"name"`
	Endpoint string `toml:"endpoint"`

	Address  string `toml:"address"`
	Insecure bool   `toml:"insecure"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Relay    int    `toml:"relay"`

	Labels Labels `toml:"labels"`
}

func (k KillSwitchConfig) String() string {
	return fmt.Sprintf("\n   name: %s\n   endpoint: %s\n   address: %s\n   insecure: %t\n   username: %s\n   password: %s\n   relay: %d\n   labels: %v",
		k.Name,
		k.Endpoint,
		k.Address,
		k.Insecure,
		k.Username,
		k.Password,
		k.Relay,
		k.Labels,
	)
}

type Labels map[string]string

func (l Labels) Contains(labels Labels) bool {
	for k, v := range labels {
		if l[k] != v {
			return false
		}
	}
	return true
}

func (l Labels) String() string {
	var s string
	for k, v := range l {
		s += fmt.Sprintf("\n    %s: %s", k, v)
	}
	return s
}
