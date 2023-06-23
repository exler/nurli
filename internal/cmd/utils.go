package cmd

import (
	"bufio"
	"fmt"
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

func PrintListHeader(columnNames ...string) {
	var header string
	var lines string
	for _, columnName := range columnNames {
		header += columnName + "\t"
		lines += strings.Repeat("-", len(columnName)) + "\t"
	}
	header = strings.TrimRight(header, "\t")
	lines = strings.TrimRight(lines, "\t")
	fmt.Println(header)
	fmt.Println(lines)

}
