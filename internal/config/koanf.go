package config

import (
	"log"
	"os"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

var K = koanf.New(".")

func Load() {
	configFile := "./config.yaml"

	if _, err := os.Stat(configFile); err == nil {
		if err := K.Load(file.Provider(configFile), yaml.Parser()); err != nil {
			log.Fatalf("❌ Error loading YAML config: %v", err)
		}
	} else {
		log.Printf("⚠️ config.yaml not found, skipping file loading")
	}

	if err := K.Load(env.Provider("APP_", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(s, "APP_")), "_", ".")
	}), nil); err != nil {
		log.Fatalf("❌ Error loading ENV config: %v", err)
	}

	log.Println("✅ Config loaded successfully")
}

func Get() *koanf.Koanf {
	return K
}
