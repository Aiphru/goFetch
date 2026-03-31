package gofetch

import (
	"fmt"
	"strings"
)

var red string = "\033[38;5;161m"
var reset string = "\033[0m"
var green string = "\033[0;32m"
var orange string = "\033[38;5;208m"

func displayLine(title string, value string) {
	spaces := 10 - len(title)
	printSpace := ""
	for i := 0; i < spaces; i++ {
		printSpace = printSpace + " "
	}
	fmt.Print(" ")
	fmt.Print(red)
	fmt.Print(title)
	fmt.Print(printSpace)
	fmt.Print(reset)
	fmt.Print(" : ")
	fmt.Print(value + "\n")
}

func formatLine(title string, value string) string {
	spaces := 10 - len(title)
	printSpace := ""
	for i := 0; i < spaces; i++ {
		printSpace = printSpace + " "
	}
	return " " + red + title + printSpace + reset + " : " + value
}

func Run() {
	ascii := getAsciiArt()
	title := getNameHostName()
	info := []string{
		" " + title,
		" " + strings.Repeat("-", visibleLength(title)),
		formatLine("OS", getDistroName()),
		formatLine("Kernel", getKernelName()),
		formatLine("Host", getDeviceName()),
		formatLine("Shell", getShell()),
		formatLine("Packages", getPackages()),
		formatLine("CPU", getCPU()),
		formatLine("Memory", getRam()),
		formatLine("Uptime", getUptime()),
		formatLine("Network", getIp()),
		formatLine("Locale", getLocale()),
	}

	maxLengthAsciiLine := 0
	for _, line := range ascii {
		if visibleLength(line) > maxLengthAsciiLine {
			maxLengthAsciiLine = visibleLength(line)
		}
	}
	maxLines := len(info)
	if len(ascii) > maxLines {
		maxLines = len(ascii)
	}
	for i := 0; i < maxLines; i++ {
		left := ""
		if i < len(ascii) {
			left = ascii[i]
		}
		right := ""
		if i < len(info) {
			right = info[i]
		}
		paddingAmount := (maxLengthAsciiLine - visibleLength(left)) + 2
		padding := strings.Repeat(" ", paddingAmount)
		fmt.Printf("%s%s%s\n", left, padding, right)
	}
}
