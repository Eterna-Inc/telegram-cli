package telegram

import (
	"flag"
	"strings"
)

// seçenekler
type Options struct {
	Photo      string
	Video      string
	Audio      string
	Voice      string
	File       string
	Location   string
	Format     string
	Silent     bool
	ReplyID    string
	Thread     string
	Status     string
	Mentions   string
	Profile    string
	ScheduleIn string // --in "10m"
	ScheduleAt string // --at "15:30"
	ExecLine   string
	Quiet      bool
	Raw        bool
	Proxy      string
	CodeLang   string
	Args       []string // kalan argümanlar (mesaj)
}

func ParseFlags() Options {
	photo := flag.String("photo", "", "Photo file path")
	video := flag.String("video", "", "Video file path")
	audio := flag.String("audio", "", "Audio file path")
	voice := flag.String("voice", "", "Voice OGG/OPUS file path")
	filef := flag.String("file", "", "Document file path")
	location := flag.String("location", "", "lat,lon")
	format := flag.String("format", "html", "Message format: html|markdown")
	silent := flag.Bool("silent", false, "Disable notification")
	reply := flag.String("reply", "", "Reply to message ID")
	thread := flag.String("thread", "", "Thread name to reply under")
	status := flag.String("status", "", "Status: success|fail|info")
	mentions := flag.String("mention", "", "Space-separated mentions: @user or numeric IDs")
	profile := flag.String("profile", "", "Profile name to use")
	inDelay := flag.String("in", "", "Delay duration (e.g. 10m, 2h)")
	atTime := flag.String("at", "", "Specific time today HH:MM (24h)")
	execCmd := flag.String("exec", "", "Execute shell command and send its output")
	quiet := flag.Bool("quiet", false, "Quiet mode")
	raw := flag.Bool("raw", false, "Print raw API response")
	proxy := flag.String("proxy", "", "HTTP proxy URL")
	codeLang := flag.String("code", "", "Send message as code block (optional language)")
	flag.Parse()

	return Options{
		Photo: *photo, Video: *video, Audio: *audio, Voice: *voice, File: *filef,
		Location: *location, Format: strings.ToLower(*format), Silent: *silent,
		ReplyID: *reply, Thread: *thread, Status: strings.ToLower(*status),
		Mentions: *mentions, Profile: *profile, ScheduleIn: *inDelay, ScheduleAt: *atTime,
		ExecLine: *execCmd, Quiet: *quiet, Raw: *raw, Proxy: *proxy, CodeLang: *codeLang,
		Args: flag.Args(),
	}
}
