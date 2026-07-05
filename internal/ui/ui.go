package ui

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Tomato
const (
	IconStart    = "🌱"
	IconWorking  = "🍅"
	IconSuccess  = "🥫"
	IconInfo     = "💧"
	IconWarn     = "⚠️"
	IconError    = "💥"
	IconQuestion = "🌿"
)

// Emoji
const (
// IconStart    = "🚀"
// IconWorking  = "️⌛"
// IconSuccess  = "✨"
// IconInfo     = "ℹ️"
// IconWarn     = "⚠️"
// IconError    = "❌"
// IconQuestion = "❓"
)

// Mark
const (
// IconStart    = "[>]"
// IconWorking  = "[*]"
// IconSuccess  = "[+]"
// IconInfo     = "[*]"
// IconWarn     = "[!]"
// IconError    = "[-]"
// IconQuestion = "[?]"
)

func Start(format string, a ...interface{}) {
	fmt.Printf("%s %s\n", IconStart, fmt.Sprintf(format, a...))
}
func Working(format string, a ...interface{}) {
	fmt.Printf("%s %s\n", IconWorking, fmt.Sprintf(format, a...))
}
func Success(format string, a ...interface{}) {
	fmt.Printf("%s %s\n", IconSuccess, fmt.Sprintf(format, a...))
}
func Info(format string, a ...interface{}) {
	fmt.Printf("%s %s\n", IconInfo, fmt.Sprintf(format, a...))
}
func Warn(format string, a ...interface{}) {
	fmt.Printf("%s %s\n", IconWarn, fmt.Sprintf(format, a...))
}
func Error(format string, a ...interface{}) {
	fmt.Printf("%s %s\n", IconError, fmt.Sprintf(format, a...))
}
func Question(format string, a ...interface{}) {
	fmt.Printf("%s %s ", IconQuestion, fmt.Sprintf(format, a...))
}

func Progress(percent float64) {
	count := int(percent * 10)

	if percent >= 1.0 {
		fmt.Printf("\r%s 完了！\n", IconSuccess)
		return
	}

	var workingIcons strings.Builder
	for range count {
		workingIcons.WriteString(IconWorking)
	}

	fmt.Printf("\r%s 進行中: %s", IconStart, workingIcons.String())
}

func Prompt(label string, defaultValue string) string {
	Question("%s [%s]:", label, defaultValue)

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			return defaultValue
		}
		return input
	}
	return defaultValue
}

func Confirm(label string, defaultYes bool) bool {
	var prompt string
	if defaultYes {
		prompt = "[Y/n]"
	} else {
		prompt = "[y/N]"
	}

	for {
		Question("%s %s:", label, prompt)

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input := strings.TrimSpace(strings.ToLower(scanner.Text()))

			if input == "" {
				return defaultYes
			}
			if input == "y" {
				return true
			}
			if input == "n" {
				return false
			}
		}
		Info("y または n で入力してください。")
	}
}

func SelectFile(label string, defaultValue string, globPattern string) string {
	files, _ := filepath.Glob(globPattern)

	Info("%s (またはファイル名を直接入力)", label)
	if len(files) > 0 {
		for i, file := range files {
			fmt.Printf("  %d) %s\n", i+1, file)
		}
	}

	input := Prompt("選択または入力", defaultValue)

	if idx, err := strconv.Atoi(input); err == nil {
		if idx >= 1 && idx <= len(files) {
			return files[idx-1]
		}
	}

	return input
}
