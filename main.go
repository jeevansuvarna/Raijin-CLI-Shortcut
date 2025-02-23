package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

var configFile = getDesktopConfigPath()

type Config struct {
	Shortcuts map[string]string `json:"shortcuts"`
}

func getDesktopConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("\033[31mError: Unable to find home directory\033[0m")
		os.Exit(1)
	}
	return homeDir + "/Desktop/shortcuts.json"
}

func getConfig() (*Config, error) {
	file, err := os.ReadFile(configFile)

	if err != nil {
		return &Config{Shortcuts: make(map[string]string)}, nil
	}
	var config Config

	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return &config, nil

}

func saveConfig(config *Config) error {
	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config: %v", err)
	}
	err = os.WriteFile(configFile, file, 0644)

	if err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}

func addShortcut(shortcut string, command string) error {
	config, err := getConfig()

	if err != nil {
		return err
	}
	config.Shortcuts[shortcut] = command
	return saveConfig(config)
}

func executeShortcut(shortcut string) error {
	config, err := getConfig()
	if err != nil {
		return err
	}

	command, exists := config.Shortcuts[shortcut]
	if !exists {
		return fmt.Errorf("shortcut '%s' not found", shortcut)
	}

	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("\033[34mUsage: raijin <add|run> [shortcut] [command]\033[0m")
		return
	}

	action := os.Args[1]

	if action == "add" {

		if len(os.Args) != 4 {
			fmt.Println("\033[34mUsage: raijin add <shortcut> <command>\033[0m")
			return
		}

		shortcut := os.Args[2]
		command := os.Args[3]

		err := addShortcut(shortcut, command)
		if err != nil {
			fmt.Println("\033[31mError:", err, "\033[0m")
		} else {
			fmt.Printf("\033[32mShortcut '%s' added successfully.\033[0m\n", shortcut)
		}
	} else {
		if len(os.Args) < 2 {
			fmt.Println("Usage: raijin run <shortcut>")
			return
		}
		shortcut := ""

		if len(os.Args) == 3 {
			shortcut = os.Args[2]
		} else {
			shortcut = os.Args[1]
		}

		err := executeShortcut(shortcut)
		if err != nil {
			fmt.Println("\033[31mError:", err, "\033[0m")
		}
	}
}
