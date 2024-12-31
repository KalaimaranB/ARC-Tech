package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// ExecuteNmap /**
func ExecuteNmap(targetIP string) {
	fmt.Println("Starting Nmap scan...")

	// Execute Nmap scan
	nmapOutput, err := runNmap(targetIP)
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

	fmt.Println("Processing complete. Check logs for details.")
}

func runNmap(target string) (string, error) {
	// Example Nmap command
	cmd := exec.Command("nmap", "-sV", target)

	// Capture output
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Nmap execution failed: %v\nOutput: %s", err, out.String())
	}

	return out.String(), nil
}

func runPythonScript(input string) error {
	// Example Python script invocation
	cmd := exec.Command("python3", "process_nmap.py")

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

	fmt.Println("Python script output:\n", out.String())
	return nil
}
