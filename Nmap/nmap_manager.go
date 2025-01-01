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

// ExecuteNmap initiates a Nmap scan for the given target IP, displays themed messages during execution,
// and processes the output with a Python script.
// Parameters:
// - targetIP: The target IP address to scan.
func ExecuteNmap(targetIP string) {
	greenBold := color.New(color.FgGreen, color.Bold)
	var userInput string

	fmt.Println("First task is to run an nmap scan with Commander Cody. Would you like to:")
	fmt.Println("(1) Use default options")
	fmt.Println("(2) Type your own nmap command line flags")
	fmt.Println("(3) Have Commander Cody assist in building a command")

	_, err := fmt.Scanln(&userInput)

	// Check user input and handle accordingly
	var cmdArgs []string
	switch userInput {
	case "1":
		// Search for the config file path
		filePath, err := Utilities.SearchFileNames("Nmap", "config")
		if err != nil {
			log.Fatalf("Failed to locate config file: %v", err)
		}

		// Load default Nmap options
		Utilities.ErrorCheckedColourPrint(greenBold, "Using default Nmap options.")
		cmdArgs, err = loadDefaultNmapFlags(filePath)
		if err != nil {
			log.Fatalf("Failed to load default options: %v", err)
		}
		cmdArgs = append(cmdArgs, targetIP)

	case "2":
		fmt.Println("Enter your custom nmap flags:")
		reader := bufio.NewReader(os.Stdin)
		customFlags, err := reader.ReadString('\n') // Read the full line of input
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		customFlags = strings.TrimSpace(customFlags)
		// Proceed to use the flags provided
		fmt.Println("You entered:", customFlags)
		// Add customFlags to the nmap command
		customFlagsList := strings.Fields(customFlags) // Split by spaces
		cmdArgs = append(cmdArgs, customFlagsList...)  // Add custom flags to the nmap command
		cmdArgs = append(cmdArgs, targetIP)

	case "3":
		// Integrate Commander Cody's assistance (interactive flag selection)
		Utilities.ErrorCheckedColourPrint(greenBold, "Commander Cody is assisting you with flag selection...")
		cmdArgs, err = getCommanderCodyFlags(targetIP)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Println("Invalid option selected, exiting.")
		os.Exit(1)
	}

	// Start themed messages in a goroutine
	stopChan := make(chan bool)
	go showCloneWarsThemedMessages(stopChan)

	// Execute Nmap scan with the constructed command
	fmt.Println("Executing nmap with :" + strings.Join(cmdArgs, " ") + " ...")
	nmapOutput, err := runNmap(cmdArgs)

	// Stop themed messages
	stopChan <- true

	if err != nil {
		fmt.Printf("Error running Nmap: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Nmap scan complete. Passing data to Python script...")

	// Execute Python script
	err = runPythonScript(nmapOutput)
	if err != nil {
		fmt.Printf("Error running Python script: %v\n", err)
		os.Exit(1)
	}

	Utilities.ErrorCheckedColourPrint(greenBold, "Nmap processing complete.")
}

// Function to load default Nmap flags from the JSON file
// Function to load default Nmap flags from the JSON file
func loadDefaultNmapFlags(filename string) ([]string, error) {
	configs, err := Utilities.LoadFilenamesConfig(filename)
	if err != nil {
		return nil, fmt.Errorf("error loading Nmap flags from %s: %v", filename, err)
	}

	// Assume the first config contains the "nmap_flags" key
	if len(configs) == 0 {
		return nil, fmt.Errorf("no configurations found in file")
	}

	nmapFlags, ok := configs[0]["nmap_flags"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("'nmap_flags' key is missing or malformed")
	}

	// Convert the flags into a slice of strings
	var cmdArgs []string
	for _, item := range nmapFlags {
		flag, isMap := item.(map[string]interface{})
		if !isMap {
			continue
		}

		// Extract flag and value
		flagName, _ := flag["flag"].(string)
		flagValue, _ := flag["value"].(string)

		cmdArgs = append(cmdArgs, flagName)
		if flagValue != "" {
			cmdArgs = append(cmdArgs, flagValue)
		}
	}

	return cmdArgs, nil
}

// getCommanderCodyFlags interacts with Commander Cody to select Nmap flags
func getCommanderCodyFlags(targetIP string) ([]string, error) {
	// Create the Python command
	cmd := exec.Command("python3", "Nmap/nmap_wizard.py")

	// Attach Stdin, Stdout, and Stderr to allow interaction
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the Python script
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error invoking Commander Cody: %v", err)
	}

	// Read JSON file
	jsonFile, err := os.Open("intermediate_data/selected_nmap_flags.json")
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %v", err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)

	// Parse JSON file into Go structure
	var selectedFlags []map[string]string
	if err := json.NewDecoder(jsonFile).Decode(&selectedFlags); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	// Convert the parsed flags to command-line arguments
	var cmdArgs []string
	for _, flag := range selectedFlags {
		if value, exists := flag["value"]; exists {
			cmdArgs = append(cmdArgs, fmt.Sprintf("%s=%s", flag["flag"], value))
		} else {
			cmdArgs = append(cmdArgs, flag["flag"])
		}
	}

	// Add target IP as the first argument
	return append([]string{targetIP}, cmdArgs...), nil
}

// runNmap executes a Nmap scan with the specified parameters and captures the output.
// Parameters:
// - target: The target IP or hostname to scan.
// - cmdArgs: The arguments for the Nmap command.
// Returns:
// - string: The captured output from the Nmap scan.
// - error: Any errors encountered during execution.
func runNmap(cmdArgs []string) (string, error) {
	// Execute the Nmap command
	cmd := exec.Command("nmap", cmdArgs...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Nmap execution failed: %v\nOutput: %s", err, out.String())
	}
	return out.String(), nil
}

// runPythonScript invokes a Python script and passes the Nmap output via stdin.
// Parameters:
// - input: The output from the Nmap scan to be processed by the Python script.
// Returns:
// - error: Any errors encountered during execution.
func runPythonScript(input string) error {
	cmd := exec.Command("python3", "Nmap/process_nmap.py")

	// Pass Nmap output to Python script via stdin
	cmd.Stdin = bytes.NewReader([]byte(input))

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Python script execution failed: %v\n", err)
	}
	return nil
}

// showCloneWarsThemedMessages prints themed messages at regular intervals during the scan process.
// It stops when a signal is received on the stopChan channel.
// Parameters:
// - stopChan: A channel used to signal when the messages should stop.
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

	// Function to shuffle the messages
	shuffleMessages := func() []string {
		shuffled := make([]string, len(messages))
		copy(shuffled, messages)
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		return shuffled
	}

	shuffledMessages := shuffleMessages()
	index := 0

	// Record the start time
	startTime := time.Now()

	for {
		select {
		case <-ticker.C:
			// Calculate elapsed time
			elapsedTime := time.Since(startTime).Truncate(time.Second)
			// Print the next message and the elapsed time
			fmt.Printf("%s (Elapsed: %s)\n", color.YellowString(shuffledMessages[index]), elapsedTime)
			index++

			// Reshuffle messages when all have been used
			if index >= len(shuffledMessages) {
				shuffledMessages = shuffleMessages()
				index = 0
			}
		case <-stopChan:
			// Stop the ticker when the channel receives a signal
			return
		}
	}
}
