package Nmap

import (
	"ARC-Tech/Utilities"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

// ExecuteNmap initiates a Nmap scan for the given target IP, displays themed messages during execution,
// and processes the output with a Python script.
// Parameters:
// - targetIP: The target IP address to scan.
func ExecuteNmap(targetIP string) {
	greenBold := color.New(color.FgGreen, color.Bold)
	Utilities.ErrorCheckedColourPrint(greenBold, "Starting Nmap scan...")

	// Start themed messages in a goroutine
	stopChan := make(chan bool)
	go showCloneWarsThemedMessages(stopChan)

	// Execute Nmap scan and save the output to a .txt file
	nmapOutput, err := runNmap(targetIP)

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

// Function to build the Nmap command from the configuration
func buildNmapCommand(target string, configs []map[string]interface{}) ([]string, error) {
	// Initialize the nmap command with the target
	cmdArgs := []string{target}

	// Loop through the flags from the config file and append them to the command
	for _, config := range configs {
		// Access the "nmap_flags" directly as it is no longer under "Nmap"
		if nmapFlags, ok := config["nmap_flags"]; ok {
			// nmapFlags is a slice of maps, where each map has a "flag" and "value"
			if flags, ok := nmapFlags.([]interface{}); ok {
				for _, flagEntry := range flags {
					// Each entry is a map containing "flag" and "value"
					if flagMap, ok := flagEntry.(map[string]interface{}); ok {
						flag, flagOk := flagMap["flag"].(string)
						value, valueOk := flagMap["value"].(string)

						// If flag exists, append it to the command arguments
						if flagOk {
							if valueOk && value != "" {
								// If value is present, add flag and value
								cmdArgs = append(cmdArgs, flag, value)
							} else {
								// If value is empty, just add the flag
								cmdArgs = append(cmdArgs, flag)
							}
						}
					}
				}
			}
		}
	}

	return cmdArgs, nil
}

// runNmap executes a Nmap scan with the specified parameters and captures the output.
// Parameters:
// - target: The target IP or hostname to scan.
// Returns:
// - string: The captured output from the Nmap scan.
// - error: Any errors encountered during execution.
func runNmap(target string) (string, error) {
	// Load the nmap config file using the SearchFileNames method
	configFile, err := Utilities.SearchFileNames("Nmap", "config")
	if err != nil {
		return "", fmt.Errorf("could not find config file: %v", err)
	}

	// Load the actual config file
	configs, err := Utilities.LoadFilenamesConfig(configFile)
	if err != nil {
		return "", fmt.Errorf("could not load nmap config: %v", err)
	}

	// Build the Nmap command using the helper method
	cmdArgs, err := buildNmapCommand(target, configs)
	if err != nil {
		return "", fmt.Errorf("could not build Nmap command: %v", err)
	}

	// Execute the Nmap command
	cmd := exec.Command("nmap", cmdArgs...)

	// Capture output
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
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

	// Capture output
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Python script execution failed: %v\nOutput: %s", err, out.String())
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
