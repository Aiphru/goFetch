package gofetch

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
	return ""
}

func getRam() string {
	var memTotal, memFree, memUsed float64
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
			break
		}
	}
	return fmt.Sprintf("%.2f GiB / %.2f GiB", memUsed, memTotal)

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
	if uptimeMinutes >= 60 {
		uptimeHours := uptimeMinutes / 60
		remainder := uptimeMinutes % 60
		return fmt.Sprint(uptimeHours) + "h " + fmt.Sprint(remainder) + "mns"
	}
	return fmt.Sprint(uptimeMinutes) + "mns"
}

func Run() {
	fmt.Println("OS : " + getDistroName())
	fmt.Println("Kernel : " + getKernelName())
	fmt.Println("Uptime : " + getUptime())
	fmt.Println("Memory : " + getRam())
}
