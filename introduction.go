package main

import (
	"fmt"
	"github.com/fatih/color"
)

func DisplayIntro() {
	// Define color styles
	blueBold := color.New(color.FgBlue, color.Bold)
	greenBold := color.New(color.FgGreen, color.Bold)
	redBold := color.New(color.FgRed, color.Bold)
	cyanBold := color.New(color.FgCyan, color.Bold)
	magentaBold := color.New(color.FgMagenta, color.Bold)

	// Print title with colors
	fmt.Println("==========================================================================================")
	errorCheckedColourPrint(blueBold, "                                        ARC-Tech")
	errorCheckedColourPrint(blueBold, "                                Advanced Recon Clone Tools")
	fmt.Println("==========================================================================================\n")

	// Print Commander introduction
	errorCheckedColourPrint(cyanBold, "Commander,")
	fmt.Println("ARC-Tech is your ultimate reconnaissance and intelligence-gathering system,")
	fmt.Println("modeled after the precision and resourcefulness of the Republic’s elite ARC Troopers.")
	fmt.Println("Its mission is clear: analyze and expose the digital battlefield to prepare you for")
	fmt.Println("strategic operations.\n")

	// Print Capabilities section with bullet points
	errorCheckedColourPrint(greenBold, "--- Capabilities ---\n")

	// Survey the Landscape
	errorCheckedColourPrint(greenBold, "1. Survey the Landscape:")
	fmt.Println("   ARC-Tech maps the target’s infrastructure, identifying key details like:")
	fmt.Println("   - Open ports (e.g., SSH, HTTP, POP3, SMTP, FTP).")
	fmt.Println("   - Operating system fingerprints and service configurations.")

	// Dive Deeper into Web Targets
	errorCheckedColourPrint(greenBold, "2. Dive Deeper into Web Targets:")
	fmt.Println("   - Perform directory busting to uncover hidden paths and resources on HTTP servers.")
	fmt.Println("   - Extract valuable insights from web pages, including metadata and HTML comments.")
	fmt.Println("   - Search for misconfigured endpoints and overlooked files.")

	// Port-Specific Reconnaissance
	errorCheckedColourPrint(greenBold, "3. Port-Specific Reconnaissance:")
	fmt.Println("   Depending on the services detected, ARC-Tech adapts its approach:")
	fmt.Println("   - SSH: Assess the potential for brute-force attacks or misconfigurations.")
	fmt.Println("   - POP3/SMTP: Analyze for open relays, exposed credentials, or weak policies.")
	fmt.Println("   - FTP: Search for unsecured files and folders, scanning for sensitive data.")

	// Enhanced Intelligence Gathering
	errorCheckedColourPrint(greenBold, "4. Enhanced Intelligence Gathering:")
	fmt.Println("   - Leverage basic OSINT techniques to collect external information about the target.")
	fmt.Println("   - Combine findings to generate a complete and actionable overview.")

	// Autonomous Reporting
	errorCheckedColourPrint(greenBold, "5. Autonomous Reporting:")
	fmt.Println("   ARC-Tech operates independently, requiring no input once deployed.")
	fmt.Println("   It compiles its findings into a comprehensive report, detailing potential")
	fmt.Println("   vulnerabilities and pathways for further operations.")

	// Print Mission section with color
	errorCheckedColourPrint(magentaBold, "\n--- Your Mission ---\n")
	fmt.Println("The digital battleground is vast, and the enemy’s defenses grow more complex by the day.")
	fmt.Println("ARC-Tech was designed to cut through this complexity, operating with the same")
	fmt.Println("autonomy and focus as the ARC Troopers it was inspired by.")

	// Print action points in bold red
	errorCheckedColourPrint(redBold, "Activate ARC-Tech and let it carry out its mission:")
	fmt.Println("   - Discover the unseen.")
	fmt.Println("   - Expose the hidden.")
	fmt.Println("   - Arm you with the intelligence to lead decisively.")

	// Print closing message
	errorCheckedColourPrint(magentaBold, "\nMay the Force guide your efforts, Commander. The Republic depends on you.")
}

// Helper function to simplify error handling
func errorCheckedColourPrint(colour *color.Color, message string) {
	_, err := colour.Println(message)
	if err != nil {
		return
	}
}
