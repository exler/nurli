package cmd

import (
	"bufio"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func GetUserInput() (text string) {
	text, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	text = strings.Trim(text, "\r\n")
	return
}

func GetUserSecureInput() (text string) {
	bytePassword, _ := term.ReadPassword(int(syscall.Stdin))
	return string(bytePassword)
}
