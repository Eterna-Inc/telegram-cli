package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"telegram-cli/config"
	"telegram-cli/utils"
)

const (
	maxTextLen = 4096
)

type apiMessageResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
	} `json:"result"`
}

// Send: medya/konum/metin/exec içerir; body ve son message_id döner
func Send(p config.Profile, cfg *config.Config, opts Options, msg string) ([]byte, int, error) {
	// proxy
	client := &http.Client{}
	if opts.Proxy != "" {
		if u, err := url.Parse(opts.Proxy); err == nil {
			client.Transport = &http.Transport{Proxy: http.ProxyURL(u)}
		}
	}

	// exec destekli
	if opts.ExecLine != "" {
		out, err := utils.RunCmd(opts.ExecLine)

		// 🔹 Kod blokları ekle
		var codeWrapped string
		if opts.Format == "markdown" {
			// Markdown için üçlü backtick ile
			codeWrapped = fmt.Sprintf("```\n%s\n```", out)
		} else {
			// HTML için <pre><code> ile
			codeWrapped = fmt.Sprintf("<pre><code>%s</code></pre>", html.EscapeString(out))
		}

		if err != nil {
			msg = fmt.Sprintf("❌ Exec failed: %s\n\n%s", opts.ExecLine, codeWrapped)
		} else {
			msg = fmt.Sprintf("🖥️ %s\n\n%s", opts.ExecLine, codeWrapped)
		}

		// Eğer çıktı çok uzunsa (4096+ karakter) dosya olarak gönder
		if len([]rune(msg)) > maxTextLen {
			tmp := filepath.Join(os.TempDir(), "telegram_exec_output.txt")
			_ = os.WriteFile(tmp, []byte(out), 0644)
			resp, body, err := sendDocument(client, p, tmp, msg, opts)
			if err != nil {
				return body, 0, err
			}
			return body, resp.Result.MessageID, nil
		}
	}

	parseMode := parseMode(opts.Format)
	replyID := opts.ReplyID
	if replyID == "" && opts.Thread != "" {
		if cfg.Threads != nil {
			if t := cfg.Threads[p.Name]; t != nil {
				if last := t[opts.Thread]; last > 0 {
					replyID = fmt.Sprintf("%d", last)
				}
			}
		}
	}

	// medya önceliği: video → photo → audio → voice → file → location → text
	switch {
	case opts.Video != "":
		resp, body, err := sendMultipart(client, p, "sendVideo", "video", opts.Video, map[string]string{
			"caption": msg, "parse_mode": parseMode,
			"disable_notification":        utils.BoolStr(opts.Silent),
			"reply_to_message_id":         replyID,
			"allow_sending_without_reply": "true",
		})
		if err != nil {
			return body, 0, err
		}
		return body, resp.Result.MessageID, nil

	case opts.Photo != "":
		resp, body, err := sendMultipart(client, p, "sendPhoto", "photo", opts.Photo, map[string]string{
			"caption": msg, "parse_mode": parseMode,
			"disable_notification":        utils.BoolStr(opts.Silent),
			"reply_to_message_id":         replyID,
			"allow_sending_without_reply": "true",
		})
		if err != nil {
			return body, 0, err
		}
		return body, resp.Result.MessageID, nil

	case opts.Audio != "":
		resp, body, err := sendMultipart(client, p, "sendAudio", "audio", opts.Audio, map[string]string{
			"caption": msg, "parse_mode": parseMode,
			"disable_notification":        utils.BoolStr(opts.Silent),
			"reply_to_message_id":         replyID,
			"allow_sending_without_reply": "true",
		})
		if err != nil {
			return body, 0, err
		}
		return body, resp.Result.MessageID, nil

	case opts.Voice != "":
		resp, body, err := sendMultipart(client, p, "sendVoice", "voice", opts.Voice, map[string]string{
			"caption": msg, "parse_mode": parseMode,
			"disable_notification":        utils.BoolStr(opts.Silent),
			"reply_to_message_id":         replyID,
			"allow_sending_without_reply": "true",
		})
		if err != nil {
			return body, 0, err
		}
		return body, resp.Result.MessageID, nil

	case opts.File != "":
		resp, body, err := sendDocument(client, p, opts.File, msg, opts)
		if err != nil {
			return body, 0, err
		}
		return body, resp.Result.MessageID, nil

	default:
		if opts.Location != "" {
			lat, lon, err := utils.ParseLatLon(opts.Location)
			if err != nil {
				return nil, 0, err
			}
			payload := map[string]string{
				"latitude":                    lat,
				"longitude":                   lon,
				"disable_notification":        utils.BoolStr(opts.Silent),
				"reply_to_message_id":         replyID,
				"allow_sending_without_reply": "true",
			}
			resp, body, err := callJSON(client, p, "sendLocation", payload)
			if err != nil {
				return body, 0, err
			}
			return body, resp.Result.MessageID, nil
		}
	}

	// text (gerekirse bölerek gönder)
	if msg == "" {
		return nil, 0, fmt.Errorf("message is empty")
	}
	chunks := utils.SplitRunes(msg, maxTextLen)
	var lastBody []byte
	lastMsgID := 0
	for i, ch := range chunks {
		payload := map[string]string{
			"text": ch, "parse_mode": parseMode,
			"disable_notification":        utils.BoolStr(opts.Silent),
			"reply_to_message_id":         replyID,
			"allow_sending_without_reply": "true",
		}
		resp, body, err := callJSON(client, p, "sendMessage", payload)
		if err != nil {
			return body, 0, err
		}
		lastBody = body
		if i == 0 && resp.Result.MessageID > 0 {
			replyID = fmt.Sprintf("%d", resp.Result.MessageID)
		}
		lastMsgID = resp.Result.MessageID
	}
	return lastBody, lastMsgID, nil
}

// --- HTTP helpers ---

func apiURL(token, method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, method)
}

func callJSON(client *http.Client, p config.Profile, method string, payload map[string]string) (apiMessageResponse, []byte, error) {
	var out apiMessageResponse
	payload["chat_id"] = p.ChatID
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", apiURL(p.Token, method), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return out, nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(b, &out)
	if !out.OK {
		return out, b, fmt.Errorf("telegram api error: %s", string(b))
	}
	return out, b, nil
}

func sendMultipart(client *http.Client, p config.Profile, method, field, path string, fields map[string]string) (apiMessageResponse, []byte, error) {
	var out apiMessageResponse

	f, err := os.Open(path)
	if err != nil {
		return out, nil, err
	}
	defer f.Close()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	_ = w.WriteField("chat_id", p.ChatID)
	for k, v := range fields {
		if v != "" {
			_ = w.WriteField(k, v)
		}
	}
	part, err := w.CreateFormFile(field, filepath.Base(path))
	if err != nil {
		return out, nil, err
	}
	if _, err := io.Copy(part, f); err != nil {
		return out, nil, err
	}
	_ = w.Close()

	req, _ := http.NewRequest("POST", apiURL(p.Token, method), &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return out, nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &out)
	if !out.OK {
		return out, body, fmt.Errorf("telegram api error: %s", string(body))
	}
	return out, body, nil
}

func sendDocument(client *http.Client, p config.Profile, path, caption string, opts Options) (apiMessageResponse, []byte, error) {
	return sendMultipart(client, p, "sendDocument", "document", path, map[string]string{
		"caption":                     caption,
		"parse_mode":                  parseMode(opts.Format),
		"disable_notification":        utils.BoolStr(opts.Silent),
		"reply_to_message_id":         opts.ReplyID,
		"allow_sending_without_reply": "true",
	})
}
