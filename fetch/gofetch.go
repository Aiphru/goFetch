package gofetch

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func displayLine(title string, value string) {
	red := "\033[38;5;161m"
	reset := "\033[0m"
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

func getNameHostName() string {
	name := os.Getenv("USER")
	hostname_FILE, _ := os.ReadFile("/etc/hostname")
	hostname := strings.Trim(string(hostname_FILE), "\n")
	return name + "@" + hostname
}

func getDistroName() string {
	osRelease, err := os.Open("/etc/os-release")
	if err != nil {
		return "Error opening /etc/os-release"
	}
	defer osRelease.Close()
	scanner := bufio.NewScanner(osRelease)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			name := strings.TrimPrefix(line, "PRETTY_NAME=")
			name = strings.Trim(name, "\"")
			return name
		}
	}
	return "No distro installed :("
}

func getIp() string {
	addresses, _ := net.InterfaceAddrs()
	for _, address := range addresses {
		ip := address.(*net.IPNet)
		if !ip.IP.IsLoopback() {
			return ip.String() //Come back and check for multiple interfaces and what happens if there's no interface available.
		}
	}
	return "No network(?)"
}

func getKernelName() string {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return "Error opening /proc/version"
	}
	kernelVersion := strings.TrimSpace(string(data))
	splitString := strings.Fields(kernelVersion)
	final := splitString[0] + " " + splitString[1] + " " + splitString[2]
	return final
}

func getUptime() string {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "Error opening /proc/uptime"
	}
	splitString := strings.Fields(string(data))
	secondsOnly := strings.Split(splitString[0], ".")[0]
	uptimeSeconds, err := strconv.ParseInt(secondsOnly, 10, 32)
	uptimeMinutes := uptimeSeconds / 60
	if uptimeMinutes > 1440 {
		uptimeDays := uptimeMinutes / (60 * 24)
		uptimeHours := (uptimeMinutes % (60 * 24)) / 60
		return fmt.Sprintf("%d days, %d hrs", uptimeDays, uptimeHours)
	}
	if uptimeMinutes >= 60 {
		uptimeHours := uptimeMinutes / 60
		remainder := uptimeMinutes % 60
		return fmt.Sprint(uptimeHours) + " hrs " + fmt.Sprint(remainder) + " min "
	}
	return fmt.Sprint(uptimeMinutes) + "min"
}

func getCPU() string {
	cpu_FILE, _ := os.Open("/proc/cpuinfo")
	defer cpu_FILE.Close()
	scanner := bufio.NewScanner(cpu_FILE)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "model name") {
			splitString := strings.SplitN(line, ": ", 2)
			return string(splitString[1])
		}
	}
	return "U don't have a cpu"
}

func getRam() string {
	var memTotal, memFree, memUsed float64
	var percentage int
	replacer := strings.NewReplacer("MemTotal:", "", "MemFree:", "", "kB", "", " ", "")
	memInfo_FILE, _ := os.Open("/proc/meminfo")
	defer memInfo_FILE.Close()
	scanner := bufio.NewScanner(memInfo_FILE)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			memTotalString := replacer.Replace(line)
			memTotal, _ = strconv.ParseFloat(memTotalString, 32)
			memTotal = memTotal / 1024 / 1024
		}
		if strings.HasPrefix(line, "MemFree:") {
			memFreeString := replacer.Replace(line)
			memFree, _ = strconv.ParseFloat(memFreeString, 32)
			memFree = memFree / 1024 / 1024
			memUsed = memTotal - memFree
			percentage = int(memUsed) * 100 / int(memTotal)
			break
		}
	}
	return fmt.Sprintf("%.2f GiB / %.2f GiB (\033[32m%d%%\033[0m)", memUsed, memTotal, percentage)
}

func getShell() string {
	shell := os.Getenv("SHELL")
	splitString := strings.Split(shell, "/")
	return splitString[len(splitString)-1]
}

func Run() {
	title := getNameHostName() + "       goFetch"
	fmt.Printf("\n %s", title)
	fmt.Println()
	for i := 0; i < len(title)+1; i++ {
		fmt.Print("-")
	}
	fmt.Println("\n")
	displayLine("OS", getDistroName())
	displayLine("Kernel", getKernelName())
	displayLine("Shell", getShell())
	displayLine("CPU", getCPU())
	displayLine("Memory", getRam())
	displayLine("Uptime", getUptime())
	displayLine("Network", getIp())
	fmt.Println("")
}
