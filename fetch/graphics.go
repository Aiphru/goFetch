package gofetch

import "regexp"

func getAsciiArt() []string {
	red := "\033[38;5;161m"
	reset := "\033[0m"

	return []string{
		"  " + red + ",-.       _,---._ __  / \\" + reset,
		" " + red + "/  )    .-'       `./ /   \\" + reset,
		red + "(  (   ,'            `/    /|" + reset,
		" " + red + "\\  `-\"             \\'\\   / |" + reset,
		"  " + red + "`.              ,  \\ \\ /  |" + reset,
		"   " + red + "/`.          ,'-`----Y   |" + reset,
		"  " + red + "(            ;        |   '" + reset,
		"  " + red + "|  ,-.    ,-'         |  /" + reset,
		"  " + red + "|  | (   |            | /" + reset,
		"  " + red + ")  |  \\  `.___________|/" + reset,
		"  " + red + "`--'   `--'" + reset,
	}
}

func visibleLength(s string) int {
	var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return len(ansiRegex.ReplaceAllString(s, ""))
}

func StringArrayMaxLength(s []string) int {
	maxLength := 0
	for _, line := range s {
		if visibleLength(line) > maxLength {
			maxLength = visibleLength(line)
		}
	}
	return maxLength
}
