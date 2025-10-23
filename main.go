package main

import (
	"os"
	"time"

	"telegram-cli/config"
	"telegram-cli/telegram"
	"telegram-cli/utils"
)

func main() {
	// "config" alt komutlarını ayrı ele al
	if len(os.Args) > 1 && os.Args[1] == "config" {
		config.Handle(os.Args[2:])
		return
	}

	// bayrakları (flags) oku
	opts := telegram.ParseFlags()

	// zamanlama (schedule)
	if d := telegram.WaitUntil(opts.ScheduleIn, opts.ScheduleAt); d > 0 {
		if !opts.Quiet {
			utils.LogInfo("⏳ waiting %s...", d)
		}
		time.Sleep(d)
	}

	// config yükle + profil çöz
	cfg := config.MustLoadAny()
	profile := cfg.ResolveProfile(opts.Profile)
	if profile.Token == "" || profile.ChatID == "" {
		utils.LogFail("profile %q is not configured. Use: telegram config --token XXX --chatid YYY --profile %s",
			profile.Name, profile.Name)
	}

	// mesajı derle (status/mention/format/exec dahil)
	msg := telegram.BuildMessage(opts)

	// gönder (medya/konum/metin)
	body, lastMsgID, err := telegram.Send(profile, &cfg, opts, msg)
	if err != nil {
		utils.LogFail("send error: %v", err)
	}

	// thread kaydı
	if opts.Thread != "" && lastMsgID > 0 {
		config.SaveThread(cfg, profile.Name, opts.Thread, lastMsgID)
	}

	utils.LogOK(opts.Quiet, opts.Raw, body, "✅ sent")
}
