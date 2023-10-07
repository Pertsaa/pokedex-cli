package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/Pertsaa/pokedex-cli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *Config) error
}

type Config struct {
	Next     *string
	Previous *string
}

func main() {
	commands := getCommands()

	locationAreaURL := "https://pokeapi.co/api/v2/location-area/"

	config := &Config{
		Next:     &locationAreaURL,
		Previous: nil,
	}

	for {
		fmt.Print("pokedex > ")

		scanner := bufio.NewScanner(os.Stdin)

		scanner.Scan()
		command := scanner.Text()

		if _, ok := commands[command]; !ok {
			fmt.Printf("Command not found: %s\n", command)

			continue
		}

		err := commands[command].callback(config)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 location areas",
			callback:    commandMapB,
		},
	}
}

func commandHelp(c *Config) error {
	fmt.Print("\nWelcome to the Pokedex!\n")
	fmt.Print("\nUsage:\n\n")

	commands := getCommands()

	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	fmt.Println()

	return nil
}

func commandExit(c *Config) error {
	os.Exit(0)
	return nil
}

func commandMap(c *Config) error {
	if c.Next == nil {
		return errors.New("No next areas")
	}

	resp, err := pokeapi.GetLocationAreas(*c.Next)
	if err != nil {
		return err
	}

	c.Next = resp.Next
	c.Previous = resp.Previous

	for _, area := range resp.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandMapB(c *Config) error {
	if c.Previous == nil {
		return errors.New("No previous areas")
	}

	resp, err := pokeapi.GetLocationAreas(*c.Previous)
	if err != nil {
		return err
	}

	c.Next = resp.Next
	c.Previous = resp.Previous

	for _, area := range resp.Results {
		fmt.Println(area.Name)
	}

	return nil
}
