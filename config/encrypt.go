package config

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
)

func Handle(args []string) {
	// Alt komutlar:
	//   config show
	//   config --token XXX [--profile name]
	//   config --chatid YYY [--profile name]
	//   config use <profile>
	//   config encrypt
	//   config decrypt
	// Not: Bayrakları basitçe elle parse edelim (küçük ve bağımsız kalsın)
	token := ""
	chatid := ""
	profile := "default"

	// kaba bayrak ayrıştırma
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--token":
			if i+1 < len(args) {
				token = args[i+1]
				i++
			}
		case "--chatid":
			if i+1 < len(args) {
				chatid = args[i+1]
				i++
			}
		case "--profile":
			if i+1 < len(args) {
				profile = args[i+1]
				i++
			}
		case "show":
			cfg := MustLoadAny()
			b, _ := json.MarshalIndent(cfg, "", "  ")
			fmt.Println(string(b))
			return
		case "use":
			if i+1 >= len(args) {
				fmt.Println("usage: telegram config use <profile>")
				return
			}
			cfg := MustLoadAny()
			name := args[i+1]
			if cfg.Profiles == nil {
				cfg.Profiles = map[string]Profile{}
			}
			if _, ok := cfg.Profiles[name]; !ok {
				cfg.Profiles[name] = Profile{}
			}
			cfg.ActiveProfile = name
			if err := Save(cfg); err != nil {
				fmt.Println("save error:", err)
				return
			}
			fmt.Println("✅ active_profile =", name)
			return
		case "encrypt":
			cfg := MustLoadAny()
			pw := os.Getenv("TELEGRAM_KEY")
			if pw == "" {
				fmt.Print("🔐 Enter encryption passphrase: ")
				pw = readLine()
				fmt.Println()
			}
			cfg.Encrypted = true
			if err := SaveEncrypted(cfg, pw); err != nil {
				fmt.Println("encrypt error:", err)
				return
			}
			fmt.Println("✅ config encrypted →", EncryptedPath)
			return
		case "decrypt":
			pw := os.Getenv("TELEGRAM_KEY")
			if pw == "" {
				fmt.Print("🔓 Enter decryption passphrase: ")
				pw = readLine()
				fmt.Println()
			}
			cfg, err := LoadEncrypted(pw)
			if err != nil {
				fmt.Println("decrypt error:", err)
				return
			}
			cfg.Encrypted = false
			if err := SavePlain(cfg); err != nil {
				fmt.Println("save error:", err)
				return
			}
			fmt.Println("✅ decrypted and wrote →", PlainPath)
			return
		}
	}

	// token/chatid set etme
	if token == "" && chatid == "" {
		fmt.Println("usage:")
		fmt.Println("  telegram config show")
		fmt.Println("  telegram config --token <TOKEN> [--profile <name>]")
		fmt.Println("  telegram config --chatid <CHAT_ID> [--profile <name>]")
		fmt.Println("  telegram config use <profile>")
		fmt.Println("  telegram config encrypt")
		fmt.Println("  telegram config decrypt")
		return
	}

	cfg := MustLoadAny()
	if cfg.Profiles == nil {
		cfg.Profiles = map[string]Profile{}
	}
	p := cfg.Profiles[profile]
	if token != "" {
		p.Token = token
	}
	if chatid != "" {
		p.ChatID = chatid
	}
	cfg.Profiles[profile] = p
	if cfg.ActiveProfile == "" {
		cfg.ActiveProfile = profile
	}
	if err := Save(cfg); err != nil {
		fmt.Println("save error:", err)
		return
	}
	fmt.Println("✅ saved for profile", profile)
}

func readLine() string {
	r := bufio.NewReader(os.Stdin)
	s, _ := r.ReadString('\n')
	return trimNL(s)
}

func trimNL(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[len(s)-1] == '\n' || s[len(s)-1] == '\r' {
		return s[:len(s)-1]
	}
	return s
}

// --- Encrypted load/save helpers ---

func LoadEncryptedFromEnvOrPrompt() (Config, error) {
	pw := os.Getenv("TELEGRAM_KEY")
	if pw == "" {
		if _, err := os.Stat(EncryptedPath); err == nil {
			fmt.Print("🔐 Enter decryption passphrase (or set TELEGRAM_KEY): ")
			pw = readLine()
			fmt.Println()
		}
	}
	if pw == "" {
		return Config{}, fmt.Errorf("no TELEGRAM_KEY")
	}
	return LoadEncrypted(pw)
}

func LoadEncrypted(pass string) (Config, error) {
	b, err := os.ReadFile(EncryptedPath)
	if err != nil {
		return Config{}, err
	}
	plain, err := aesgcmDecrypt(deriveKey(pass), b)
	if err != nil {
		return Config{}, err
	}
	var c Config
	if err := json.Unmarshal(plain, &c); err != nil {
		return Config{}, err
	}
	c.Encrypted = true
	return c, nil
}

func SaveEncrypted(c Config, pass string) error {
	c.Encrypted = true
	plain, _ := json.MarshalIndent(c, "", "  ")
	enc, err := aesgcmEncrypt(deriveKey(pass), plain)
	if err != nil {
		return err
	}
	if err := os.WriteFile(EncryptedPath, enc, 0600); err != nil {
		return err
	}
	_ = os.Remove(PlainPath)
	return nil
}

func deriveKey(pass string) []byte {
	sum := sha256.Sum256([]byte(pass))
	return sum[:]
}

func aesgcmEncrypt(key, plain []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	ct := gcm.Seal(nonce, nonce, plain, nil)
	return ct, nil
}

func aesgcmDecrypt(key, ct []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(ct) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce := ct[:gcm.NonceSize()]
	data := ct[gcm.NonceSize():]
	return gcm.Open(nil, nonce, data, nil)
}
