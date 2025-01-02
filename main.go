package main

import (
	"ARC-Tech/Nmap"
	"fmt"
	"regexp"
	"strings"
)

func main() {
	DisplayIntro()
	ipAddress := StartProgram()
	if strings.HasSuffix(ipAddress, "FOR THE REPUBLIC!!!") {
		ip := extractIPAddress(ipAddress)
		fmt.Println("Executing default Nmap series for:", ip)
		Nmap.ExecuteNmap(ip, true)
	} else {
		Nmap.ExecuteNmap(ipAddress, false)
	}
}

// Extract the IP address from the input if "FOR THE REPUBLIC!!!" is present
func extractIPAddress(input string) string {
	re := regexp.MustCompile(`^(\d{1,3}(\.\d{1,3}){3})\s+FOR THE REPUBLIC!!!$`)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 0 {
		return matches[1] // Return the IP address part
	}
	return input // Return input as is if not in the expected format
}

//Scan me Nmap: 45.33.32.156
