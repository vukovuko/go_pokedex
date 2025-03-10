package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"internal/pokecache"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// cliCommand represents a REPL command.
type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

// config holds global configuration such as pagination URLs and a cache.
type config struct {
	next     string
	previous string
	cache    *pokecache.Cache
}

var commands map[string]cliCommand

// Initialize the commands registry.
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
		"map": {
			name:        "map",
			description: "Displays a list of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas (go back a page)",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location area by name",
			callback:    commandExplore,
		},
	}
}

// commandExit prints the exit message and immediately exits.
func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // Unreachable but required by signature.
}

// commandHelp prints the welcome message and usage information.
func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

// LocationAreaResponse represents a paginated response from the PokeAPI.
type LocationAreaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// fetchLocationAreas fetches a paginated list of location areas from the given URL,
// using the provided cache to avoid duplicate network requests.
func fetchLocationAreas(url string, cache *pokecache.Cache) (*LocationAreaResponse, error) {
	if data, found := cache.Get(url); found {
		var resp LocationAreaResponse
		if err := json.Unmarshal(data, &resp); err == nil {
			return &resp, nil
		}
		// If unmarshaling fails, fall back to HTTP request.
	}
	respHTTP, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer respHTTP.Body.Close()
	data, err := ioutil.ReadAll(respHTTP.Body)
	if err != nil {
		return nil, err
	}
	cache.Add(url, data)
	var respData LocationAreaResponse
	if err := json.Unmarshal(data, &respData); err != nil {
		return nil, err
	}
	return &respData, nil
}

// commandMap fetches and displays the next 20 location areas.
func commandMap(cfg *config, args []string) error {
	if cfg.next == "" {
		cfg.next = "https://pokeapi.co/api/v2/location-area?limit=20"
	}
	apiResp, err := fetchLocationAreas(cfg.next, cfg.cache)
	if err != nil {
		return err
	}
	for _, area := range apiResp.Results {
		fmt.Println(area.Name)
	}
	cfg.next = apiResp.Next
	cfg.previous = apiResp.Previous
	return nil
}

// commandMapB fetches and displays the previous 20 location areas.
func commandMapB(cfg *config, args []string) error {
	if cfg.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	apiResp, err := fetchLocationAreas(cfg.previous, cfg.cache)
	if err != nil {
		return err
	}
	for _, area := range apiResp.Results {
		fmt.Println(area.Name)
	}
	cfg.next = apiResp.Next
	cfg.previous = apiResp.Previous
	return nil
}

// LocationAreaDetail represents detailed information for a specific location area.
type LocationAreaDetail struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

// fetchLocationAreaDetail fetches detailed information for a given location area name,
// using caching to speed up subsequent requests.
func fetchLocationAreaDetail(areaName string, cache *pokecache.Cache) (*LocationAreaDetail, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + areaName
	if data, found := cache.Get(url); found {
		var detail LocationAreaDetail
		if err := json.Unmarshal(data, &detail); err == nil {
			return &detail, nil
		}
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	cache.Add(url, data)
	var detail LocationAreaDetail
	if err := json.Unmarshal(data, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}

// commandExplore takes a location area name as an argument,
// fetches detailed data from the PokeAPI, and displays the names of the Pokemon encountered there.
func commandExplore(cfg *config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Please provide a location area to explore")
		return nil
	}
	areaName := args[0]
	fmt.Printf("Exploring %s...\n", areaName)
	detail, err := fetchLocationAreaDetail(areaName, cfg.cache)
	if err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, encounter := range detail.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func main() {
	// Initialize a cache with an expiration interval (e.g., 5 seconds).
	cache := pokecache.NewCache(5 * time.Second)
	cfg := &config{
		next:     "https://pokeapi.co/api/v2/location-area?limit=20",
		previous: "",
		cache:    cache,
	}
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
		args := fields[1:]
		if cmd, exists := commands[commandName]; exists {
			if err := cmd.callback(cfg, args); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
