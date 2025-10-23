package telegram

import (
	"time"

	"telegram-cli/utils"
)

// --in / --at için bekleme süresini hesaplar
func WaitUntil(inDur, atHHMM string) time.Duration {
	if inDur != "" {
		if d, err := time.ParseDuration(inDur); err == nil && d > 0 {
			return d
		}
	}
	if atHHMM != "" {
		if t, err := utils.ParseAt(atHHMM, time.Now()); err == nil {
			return time.Until(t)
		}
	}
	return 0
}
