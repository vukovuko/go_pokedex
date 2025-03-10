# Pokedex CLI

A simple command-line Pokedex application built in Go. This project allows you to explore the Pokemon world by catching Pokemon, inspecting them for detailed information, and exploring various location areas.

## Features

- **Catch Pokemon:** Attempt to catch a Pokemon by name and store its details.
- **Inspect Pokemon:** View detailed information (name, height, weight, stats, and types) of a Pokemon you've caught.
- **Explore Locations:** List and explore location areas, and view encountered Pokemon.
- **Pagination:** Navigate location areas with `map` (next page) and `mapb` (previous page).
- **Help & Exit:** Get usage instructions and exit the application when done.

## Prerequisites

- [Go](https://golang.org/dl/) (version 1.16 or higher recommended)

## Installation

1. **Clone the Repository:**
    ```bash
    git clone https://github.com/vukovuko/go_pokedex.git
    cd go_pokedex
    ```

2. **Download Dependencies:**
    ```bash
    go mod download
    ```

## Running the Application

You have two options to run the application:

### Option 1: Run Directly with `go run`
    go run main.go
### Option 2: Build and Run the Executable

1. Build the application:
    ```bash
    go build -o pokedex
    ```

2. Run the generated executable:
    ```bash
    ./pokedex
    ```
    
## How to Play

Once the application starts, you'll see a prompt:

    Pokedex >

You can interact with the Pokedex using these commands:

- **help:**  
  Displays all available commands and their descriptions.

- **catch `<pokemon-name>`:**  
  Attempts to catch the specified Pokemon.  
  Example:
  
      Pokedex > catch pidgey

- **inspect `<pokemon-name>`:**  
  Shows detailed information about the specified Pokemon if it has been caught.  
  If not caught, it prints:
  
      you have not caught that pokemon

- **map:**  
  Displays a list of 20 location areas in the Pokemon world.

- **mapb:**  
  Displays the previous 20 location areas (go back a page).

- **explore `<location-area>`:**  
  Explores a location area by name and lists the encountered Pokemon.

- **exit:**  
  Exits the application.

## Testing the Application

Interact with the CLI by entering commands at the `Pokedex >` prompt. Here are a few scenarios to try:

1. **Catch and Inspect a Pokemon:**
   - Try `catch pidgey` and then `inspect pidgey` to see its details.
   - If you try to inspect a Pokemon that hasn't been caught, you'll see a message indicating that.

2. **Explore Location Areas:**
   - Use `map` to view a list of locations.
   - Use `explore <location-name>` to see which Pokemon are encountered in that area.
   - Use `mapb` to go back a page if needed.
