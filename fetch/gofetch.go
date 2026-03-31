package gofetch

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func getNameHostName() string {
	name := os.Getenv("USER")
	hostname_FILE, _ := os.ReadFile("/etc/hostname")
	hostname := strings.Trim(string(hostname_FILE), "\n")
	return red + name + reset + "@" + red + hostname
}

func getDeviceName() string {
	_, err := os.Stat("/sys/class/dmi/id/product_name")
	if err == nil {
		FILE, _ := os.ReadFile("/sys/class/dmi/id/product_name")
		return strings.Trim(string(FILE), "\n")
	}
	_, err = os.Stat("/sys/class/dmi/id/board_name")
	if err == nil {
		FILE, _ := os.ReadFile("sys/class/dmi/id/board_name")
		return strings.Trim(string(FILE), "\n")
	}
	_, err = os.Stat("/proc/device-tree/model")
	if err == nil {
		FILE, _ := os.ReadFile("/proc/device-tree/model")
		return strings.Trim(string(FILE), "\n")
	}
	return "WSL"

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
			if ip.IP.To4() != nil {
				return ip.IP.String()
			} //Come back and check for multiple interfaces
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
	if percentage < 50 {
		return fmt.Sprintf("%.2f GiB / %.2f GiB (%s%d%%%s)", memUsed, memTotal, green, percentage, reset)
	}
	if percentage >= 50 && percentage < 75 {
		return fmt.Sprintf("%.2f GiB / %.2f GiB (%s%d%%%s)", memUsed, memTotal, orange, percentage, reset)
	}
	return fmt.Sprintf("%.2f GiB / %.2f GiB (%s%d%%%s)", memUsed, memTotal, red, percentage, reset)

}

func getDebianPackages() string {
	debian, err := os.Open("/var/lib/dpkg/status")
	defer debian.Close()
	if err != nil {
		fmt.Println("Error opening debian file")
	}
	var count int64 = 0
	scanner := bufio.NewScanner(debian)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Package:") {
			count++
		}
	}
	return strconv.FormatInt(count, 10)
}

func getArchPackages() string {
	dir, err := os.ReadDir("/var/lib/pacman/local")
	if err != nil {
		fmt.Println("Error opening pacman dir")
	}
	var count int64 = 0
	for range dir {
		count++
	}
	//Account for ALPM_DB_VERSION
	return strconv.FormatInt(count-1, 10)
}

// AI placeholder, can't be bothered to test this rn.
func getRPMPackages() string {
	// We check if the RPM database exists first
	if _, err := os.Stat("/var/lib/rpm/rpmdb.sqlite"); err == nil {
		// Since it's a binary DB, we'll use a quick exec here
		// OR return "N/A" if you really want to avoid exec.
		out, _ := exec.Command("rpm", "-qa").Output()
		lines := strings.Split(string(out), "\n")
		return strconv.Itoa(len(lines) - 1)
	}
	return ""
}

func getPackages() string {
	_, err := os.Stat("/var/lib/dpkg/status")
	if err == nil {
		return getDebianPackages()
	}
	_, err = os.Stat("/var/lib/pacman")
	if err == nil {
		return getArchPackages()
	}
	_, err = os.Stat("/var/lib/rpm/rpmdb.sqlite")
	if err == nil {
		return getRPMPackages()
	}

	return "Unknown"
}

func getShell() string {
	shell := os.Getenv("SHELL")
	splitString := strings.Split(shell, "/")
	return splitString[len(splitString)-1]
}

func getLocale() string {
	return os.Getenv("LANG")
}
