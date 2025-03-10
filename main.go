package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Create a new scanner that reads from standard input.
	scanner := bufio.NewScanner(os.Stdin)

	// Infinite loop: one iteration per user command.
	for {
		// Print the prompt without a newline.
		fmt.Print("Pokedex > ")

		// Wait for input from the user.
		if !scanner.Scan() {
			// If scanning fails (e.g., EOF), exit the loop.
			break
		}

		// Get the user's input text.
		input := scanner.Text()

		// Clean the input:
		// 1. Trim any leading or trailing whitespace.
		// 2. Convert the text to lowercase.
		// 3. Split the text into words using whitespace as a delimiter.
		words := strings.Fields(strings.ToLower(strings.TrimSpace(input)))

		// If there are words in the input, print the first word.
		if len(words) > 0 {
			fmt.Printf("Your command was: %s\n", words[0])
		} else {
			// If no valid input was provided, print an empty command response.
			fmt.Println("Your command was:")
		}
	}
}
