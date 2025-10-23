package config

type Profile struct {
	Name   string `json:"-"`
	Token  string `json:"token"`
	ChatID string `json:"chat_id"`
}

type Threads map[string]int // threadName -> last message_id

type Config struct {
	ActiveProfile string             `json:"active_profile"`
	Profiles      map[string]Profile `json:"profiles"`
	Threads       map[string]Threads `json:"threads"`
	Encrypted     bool               `json:"encrypted"`
}
