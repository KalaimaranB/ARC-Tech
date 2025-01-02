package Nmap

import (
	"ARC-Tech/Utilities"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ExecuteNmap initiates a Nmap scan for the given target IP.
//
// Parameters:
//   - targetIP: The IP address of the target for the Nmap scan.
//   - defaultState: A boolean flag to determine whether to use default options.
//
// Objective:
//
//	Depending on user input, this function either runs a default scan,
//	allows custom flags, or invokes an interactive wizard for flag selection.
//
// Error Handling:
//   - Handles errors in reading user input and exits with a log message if invalid.
func ExecuteNmap(targetIP string, defaultState bool) {
	if defaultState {
		useDefaultOptions(targetIP)
		return
	}

	fmt.Println("First task is to run an Nmap scan with Commander Cody. Would you like to:")
	fmt.Println("(1) Use default options")
	fmt.Println("(2) Type your own Nmap command line flags")
	fmt.Println("(3) Have Commander Cody assist in building a command")

	var userInput string
	_, err := fmt.Scanln(&userInput)
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	switch userInput {
	case "1":
		useDefaultOptions(targetIP)
	case "2":
		useCustomOptions(targetIP)
	case "3":
		useCommanderCody(targetIP)
	default:
		fmt.Println("Invalid option selected, exiting.")
		os.Exit(1)
	}
}

// useDefaultOptions runs a Nmap scan using default configurations.
//
// Parameters:
//   - targetIP: The IP address of the target for the Nmap scan.
//
// Objective:
//
//	Reads default configurations from a JSON file and executes the Nmap scan.
//
// Error Handling:
//   - Logs and exits if the config file is not found or cannot be parsed.
func useDefaultOptions(targetIP string) {
	greenBold := color.New(color.FgGreen, color.Bold)
	filePath, err := Utilities.SearchFileNames("Nmap", "config")
	if err != nil {
		log.Fatalf("Failed to locate config file: %v", err)
	}

	Utilities.ErrorCheckedColourPrint(greenBold, "Using default Nmap options.")
	cmdArgs, err := loadDefaultNmapFlags(filePath)
	if err != nil {
		log.Fatalf("Failed to load default options: %v", err)
	}
	executeNmapScan(targetIP, cmdArgs)
}

// useCustomOptions prompts the user to input custom Nmap flags.
//
// Parameters:
//   - targetIP: The IP address of the target for the Nmap scan.
//
// Objective:
//
//	Allows users to manually specify Nmap flags and executes the scan with them.
//
// Error Handling:
//   - Logs and exits if there is an error reading user input.
func useCustomOptions(targetIP string) {
	fmt.Println("Enter your custom Nmap flags:")
	reader := bufio.NewReader(os.Stdin)
	customFlags, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	customFlags = strings.TrimSpace(customFlags)
	cmdArgs := strings.Fields(customFlags)
	executeNmapScan(targetIP, cmdArgs)
}

// useCommanderCody uses an interactive script to assist in building Nmap flags.
//
// Parameters:
//   - targetIP: The IP address of the target for the Nmap scan.
//
// Objective:
//
//	Invokes a Python-based wizard to generate Nmap flags and executes the scan.
//
// Error Handling:
//   - Logs and exits if there is an error invoking the wizard or processing its output.
func useCommanderCody(targetIP string) {
	greenBold := color.New(color.FgGreen, color.Bold)
	Utilities.ErrorCheckedColourPrint(greenBold, "Commander Cody is assisting you with flag selection...")
	cmdArgs, err := getCommanderCodyFlags()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	executeNmapScan(targetIP, cmdArgs)
}

// executeNmapScan runs the Nmap scan and handles its output.
//
// Parameters:
//   - targetIP: The IP address of the target for the Nmap scan.
//   - cmdArgs: A slice of strings representing the Nmap command-line arguments.
//
// Objective:
//
//	Executes the Nmap command with the specified arguments, processes its output, and
//	invokes a Python script for further analysis.
//
// Error Handling:
//   - Logs and exits if Nmap execution or Python script processing fails.
func executeNmapScan(targetIP string, cmdArgs []string) {
	stopChan := make(chan bool)
	go showCloneWarsThemedMessages(stopChan)

	cmdArgs = append(cmdArgs, targetIP, "-oX", "output/nmap_xml.xml")
	fmt.Println("Executing Nmap with:", strings.Join(cmdArgs, " "))

	nmapOutput, err := runNmap(cmdArgs)
	stopChan <- true
	if err != nil {
		log.Fatalf("Error running Nmap: %v", err)
	}

	fmt.Println("Nmap scan complete. Passing data to Python script...")
	if err := runPythonScript(nmapOutput); err != nil {
		log.Fatalf("Error running Python script: %v", err)
	}

	color.New(color.FgGreen, color.Bold).Println("Nmap processing complete.")
}

// loadDefaultNmapFlags loads default Nmap flags from a JSON file.
//
// Parameters:
//   - filename: Path to the JSON configuration file.
//
// Returns:
//   - A slice of strings containing Nmap command-line arguments.
//   - An error if the file cannot be opened or parsed.
func loadDefaultNmapFlags(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error loading Nmap flags from %s: %v", filename, err)
	}
	defer file.Close()

	var configs []map[string]string
	if err := json.NewDecoder(file).Decode(&configs); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	var cmdArgs []string
	for _, flag := range configs {
		cmdArgs = append(cmdArgs, flag["flag"])
		if value := flag["value"]; value != "" {
			cmdArgs = append(cmdArgs, value)
		}
	}

	return cmdArgs, nil
}

// getCommanderCodyFlags invokes the wizard to select Nmap flags.
//
// Returns:
//   - A slice of strings containing Nmap command-line arguments.
//   - An error if the wizard or its output processing fails.
func getCommanderCodyFlags() ([]string, error) {
	cmd := exec.Command("python3", "Nmap/nmap_wizard.py")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error invoking Commander Cody: %v", err)
	}

	jsonFile, err := os.Open("intermediate_data/selected_nmap_flags.json")
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %v", err)
	}
	defer jsonFile.Close()

	var selectedFlags []map[string]string
	if err := json.NewDecoder(jsonFile).Decode(&selectedFlags); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	var cmdArgs []string
	for _, flag := range selectedFlags {
		cmdArgs = append(cmdArgs, flag["flag"])
		if value, exists := flag["value"]; exists && value != "" {
			cmdArgs = append(cmdArgs, value)
		}
	}

	return cmdArgs, nil
}

// runNmap executes a Nmap scan and captures the output.
//
// Parameters:
//   - cmdArgs: A slice of strings representing the Nmap command-line arguments.
//
// Returns:
//   - A string containing the output from the Nmap scan.
//   - An error if the Nmap command fails.
func runNmap(cmdArgs []string) (string, error) {
	cmd := exec.Command("nmap", cmdArgs...)
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Nmap execution failed: %v\nOutput: %s", err, out.String())
	}
	return out.String(), nil
}

// runPythonScript processes Nmap output using a Python script.
//
// Parameters:
//   - input: A string containing the output of the Nmap scan.
//
// Returns:
//   - An error if the Python script execution fails.
func runPythonScript(input string) error {
	cmd := exec.Command("python3", "Nmap/process_nmap.py")
	cmd.Stdin = bytes.NewReader([]byte(input))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("python script execution failed: %v", err)
	}
	return nil
}

// showCloneWarsThemedMessages prints themed messages during the scan process.
//
// Parameters:
//   - stopChan: A channel to signal when to stop printing messages.
//
// Objective:
//
//	Displays periodic themed messages to entertain the user during the scan.
func showCloneWarsThemedMessages(stopChan chan bool) {
	messages := []string{
		"The battlefield is quiet, but the Force is with us.",
		"Scanning enemy territory for vital intelligence...",
		"Advanced recon units have eyes on the objective.",
		"Trust in your training, Commander. We'll uncover the truth.",
		"Synchronizing Republic data feeds for optimal accuracy...",
		"Patience is a virtue, even in the heat of battle.",
		"ARC Troopers stand ready to adapt and overcome.",
		"Gathering intel is the first step to victory.",
		"Stealth and precision, the hallmarks of our success.",
		"The Jedi Council awaits our findings. Proceed with caution.",
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	shuffleMessages := func() []string {
		shuffled := make([]string, len(messages))
		copy(shuffled, messages)
		rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		return shuffled
	}

	shuffledMessages := shuffleMessages()
	index := 0
	startTime := time.Now()

	for {
		select {
		case <-ticker.C:
			elapsedTime := time.Since(startTime).Truncate(time.Second)
			fmt.Printf("%s (Elapsed: %s)\n", color.YellowString(shuffledMessages[index]), elapsedTime)
			index++
			if index >= len(shuffledMessages) {
				shuffledMessages = shuffleMessages()
				index = 0
			}
		case <-stopChan:
			return
		}
	}
}
