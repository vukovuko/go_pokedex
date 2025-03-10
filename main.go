package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // Unreachable but required by function signature.
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if input == "" {
			continue
		}
		fields := strings.Fields(input)
		commandName := fields[0]
		if cmd, exists := commands[commandName]; exists {
			if err := cmd.callback(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
