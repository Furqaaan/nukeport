package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

func listPorts() {
	conns, _ := net.Connections("inet")

	fmt.Printf("%-8s %-8s %-20s\n", "PORT", "PID", "PROCESS")
	fmt.Println(strings.Repeat("-", 40))

	seen := make(map[string]bool)

	for _, c := range conns {
		if c.Status != "LISTEN" {
			continue
		}

		port := c.Laddr.Port
		pid := c.Pid

		key := fmt.Sprintf("%d-%d", port, pid)
		if seen[key] {
			continue
		}
		seen[key] = true

		name := ""
		if p, err := process.NewProcess(pid); err == nil {
			name, _ = p.Name()
		}

		fmt.Printf("%-8d %-8d %-20s\n", port, pid, name)
	}
}

func killPort(port int) {
	conns, _ := net.Connections("inet")
	killed := false

	for _, c := range conns {
		if c.Laddr.Port == uint32(port) && c.Status == "LISTEN" {
			pid := c.Pid

			p, err := process.NewProcess(pid)
			if err != nil {
				continue
			}

			name, _ := p.Name()
			fmt.Printf("Killing PID %d (%s)\n", pid, name)

			p.Kill()
			killed = true
		}
	}

	if !killed {
		fmt.Println("No process found on port", port)
	}
}

func main() {
	if len(os.Args) == 1 {
		listPorts()
		return
	}

	portStr := os.Args[1]

	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Println("Invalid port:", portStr)
		return
	}

	killPort(port)
}
