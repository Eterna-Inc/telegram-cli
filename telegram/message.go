package telegram

import (
	"fmt"
	"html"
	"strings"

	"telegram-cli/utils"
)

// mesajı hazırla (status, mention, exec vb. ile)
func BuildMessage(opts Options) string {
	// exec çıktısı sonrası mesaj boş olsa bile üretilecek, api tarafı halleder
	msg := strings.TrimSpace(strings.Join(opts.Args, " "))

	// Eğer --code aktifse, özel biçim uygulayalım
	if opts.CodeLang != "" {
		lang := opts.CodeLang
		if opts.Format == "markdown" {
			if lang != "" {
				msg = fmt.Sprintf("```%s\n%s\n```", lang, msg)
			} else {
				msg = fmt.Sprintf("```\n%s\n```", msg)
			}
		} else { // HTML format
			safe := html.EscapeString(msg)
			if lang != "" {
				msg = fmt.Sprintf("<pre><code class=\"language-%s\">%s</code></pre>", lang, safe)
			} else {
				msg = fmt.Sprintf("<pre><code>%s</code></pre>", safe)
			}
		}
		// code block kullanıldığında status ve mentions eklenmez
		return msg
	}

	// status
	switch opts.Status {
	case "success", "ok", "done":
		msg = "✅ " + msg
	case "fail", "error":
		msg = "❌ " + msg
	case "info":
		msg = "ℹ️ " + msg
	}

	// mentions
	if m := strings.TrimSpace(opts.Mentions); m != "" {
		toks := strings.Fields(m)
		var built []string
		for _, t := range toks {
			if t == "" {
				continue
			}
			if strings.HasPrefix(t, "@") {
				built = append(built, t)
			} else if utils.IsDigits(t) {
				built = append(built, fmt.Sprintf(`<a href="tg://user?id=%s">mention</a>`, t))
			} else {
				built = append(built, t)
			}
		}
		if len(built) > 0 {
			if msg == "" {
				msg = strings.Join(built, " ")
			} else {
				msg = msg + "\n" + strings.Join(built, " ")
			}
		}
	}
	return msg
}

func parseMode(format string) string {
	switch strings.ToLower(format) {
	case "markdown", "md", "markdownv2":
		return "Markdown"
	default:
		return "HTML"
	}
}
