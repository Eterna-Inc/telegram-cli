package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func IsDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func SplitRunes(s string, n int) []string {
	r := []rune(s)
	if len(r) <= n {
		return []string{s}
	}
	var out []string
	for len(r) > 0 {
		if len(r) <= n {
			out = append(out, string(r))
			break
		}
		out = append(out, string(r[:n]))
		r = r[n:]
	}
	return out
}

func ParseAt(hhmm string, now time.Time) (time.Time, error) {
	p := strings.Split(hhmm, ":")
	if len(p) != 2 {
		return time.Time{}, fmt.Errorf("use HH:MM")
	}
	hh, err := strconv.Atoi(p[0])
	if err != nil {
		return time.Time{}, err
	}
	mm, err := strconv.Atoi(p[1])
	if err != nil {
		return time.Time{}, err
	}
	t := time.Date(now.Year(), now.Month(), now.Day(), hh, mm, 0, 0, now.Location())
	if t.Before(now) {
		t = t.Add(24 * time.Hour)
	}
	return t, nil
}

func ParseLatLon(s string) (string, string, error) {
	parts := strings.Split(strings.TrimSpace(s), ",")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("expected lat,lon")
	}
	lat := strings.TrimSpace(parts[0])
	lon := strings.TrimSpace(parts[1])
	if _, err := strconv.ParseFloat(lat, 64); err != nil {
		return "", "", err
	}
	if _, err := strconv.ParseFloat(lon, 64); err != nil {
		return "", "", err
	}
	return lat, lon, nil
}

func BoolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
