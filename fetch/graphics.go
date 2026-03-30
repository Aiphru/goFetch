package gofetch

func getAsciiArt() []string {
	red := "\033[38;5;196m"
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
