package utils

import (
	"fmt"
	"os"
)

func LogInfo(format string, a ...any) {
	fmt.Printf("• "+format+"\n", a...)
}

func LogFail(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "✖ "+format+"\n", a...)
	os.Exit(1)
}

func LogOK(quiet, raw bool, body []byte, msg string) {
	if quiet {
		return
	}
	if raw && body != nil {
		fmt.Println(string(body))
		return
	}
	fmt.Println(msg)
}
