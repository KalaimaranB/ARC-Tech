package main

import "ARC-Tech/Nmap"

func main() {
	DisplayIntro()
	ipAddress := StartProgram()
	Nmap.ExecuteNmap(ipAddress)
	//Scan me Nmap: 45.33.32.156
}
