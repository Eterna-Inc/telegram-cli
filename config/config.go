package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var (
	homeDir, _    = os.UserHomeDir()
	PlainPath     = filepath.Join(homeDir, ".telegram_config")
	EncryptedPath = filepath.Join(homeDir, ".telegram_config.enc")
)

// MustLoadAny: düz veya şifreli yapılandırmayı yükler; yoksa varsayılan döner.
func MustLoadAny() Config {
	if c, err := LoadPlain(); err == nil {
		return c
	}
	if c, err := LoadEncryptedFromEnvOrPrompt(); err == nil {
		return c
	}
	// default boş config
	return Config{
		ActiveProfile: "default",
		Profiles:      map[string]Profile{"default": {}},
		Threads:       map[string]Threads{},
		Encrypted:     false,
	}
}

// düz dosyadan yükle
func LoadPlain() (Config, error) {
	var c Config
	b, err := os.ReadFile(PlainPath)
	if err != nil {
		return c, err
	}
	if err := json.Unmarshal(b, &c); err != nil {
		return c, err
	}
	return c, nil
}

// düz dosyaya kaydet
func SavePlain(c Config) error {
	b, _ := json.MarshalIndent(c, "", "  ")
	return os.WriteFile(PlainPath, b, 0600)
}

// genel kaydetme (Encrypted=true ise şifreli, aksi halde düz)
func Save(c Config) error {
	if c.Encrypted {
		key := os.Getenv("TELEGRAM_KEY")
		if key == "" {
			return fmt.Errorf("config marked Encrypted but TELEGRAM_KEY is empty")
		}
		return SaveEncrypted(c, key)
	}
	return SavePlain(c)
}

// aktif profili çözer
func (c Config) ResolveProfile(name string) Profile {
	if name == "" {
		name = c.ActiveProfile
	}
	p := c.Profiles[name]
	p.Name = name
	return p
}

// thread kaydet
func SaveThread(c Config, profile, thread string, id int) {
	if c.Threads == nil {
		c.ThreadSaneInit()
	}
	if c.Threads[profile] == nil {
		c.Threads[profile] = Threads{}
	}
	c.Threads[profile][thread] = id
	_ = Save(c)
}

func (c *Config) ThreadSaneInit() {
	if c.Threads == nil {
		c.Threads = map[string]Threads{}
	}
}
