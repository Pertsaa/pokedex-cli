package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Pertsaa/pokedex-cli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *Config, args ...string) error
}

type Config struct {
	API      *pokeapi.PokeApi
	Next     *string
	Previous *string
	Pokedex  map[string]Pokemon
}

type Pokemon struct {
	Name   string
	Height int
	Weight int
	Stats  map[string]int
	Types  []string
}

func main() {
	commands := getCommands()

	config := &Config{
		API: pokeapi.New(5 * time.Minute),
	}

	for {
		fmt.Print("pokedex > ")

		scanner := bufio.NewScanner(os.Stdin)

		scanner.Scan()
		input := scanner.Text()

		split := strings.Split(input, " ")
		command := split[0]
		args := []string{}

		if len(split) > 1 {
			args = split[1:]
		}

		if _, ok := commands[command]; !ok {
			fmt.Printf("Command not found: %s\n", command)

			continue
		}

		err := commands[command].callback(config, args...)
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
		"explore": {
			name:        "explore <area>",
			description: "Explore area for pokemon",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon>",
			description: "Try to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon>",
			description: "Inspect a caught pokemon",
			callback:    commandInspect,
		},
	}
}

func commandHelp(c *Config, args ...string) error {
	fmt.Print("\nWelcome to the Pokedex!\n")
	fmt.Print("\nUsage:\n\n")

	commands := getCommands()

	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	fmt.Println()

	return nil
}

func commandExit(c *Config, args ...string) error {
	os.Exit(0)
	return nil
}

func commandMap(c *Config, args ...string) error {
	resp, err := c.API.GetLocationAreas(c.Next)
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

func commandMapB(c *Config, args ...string) error {
	if c.Previous == nil {
		return errors.New("No previous areas")
	}

	resp, err := c.API.GetLocationAreas(c.Previous)
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

func commandExplore(c *Config, args ...string) error {
	if len(args) != 1 {
		return errors.New("Invalid args")
	}

	resp, err := c.API.GetAreaEncounters(args[0])
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range resp.Encounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(c *Config, args ...string) error {
	if len(args) != 1 {
		return errors.New("Invalid args")
	}

	resp, err := c.API.GetPokemon(args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])

	if rand.Intn(resp.BaseExperience) < 40 {
		fmt.Printf("%s was caught!\n", args[0])

		stats := make(map[string]int)
		for _, stat := range resp.Stats {
			stats[stat.Stat.Name] = stat.BaseStat
		}

		types := make([]string, len(resp.Types))
		for i, t := range resp.Types {
			types[i] = t.Type.Name
		}

		if c.Pokedex == nil {
			c.Pokedex = make(map[string]Pokemon)
		}

		c.Pokedex[args[0]] = Pokemon{
			Name:   resp.Name,
			Height: resp.Height,
			Weight: resp.Weight,
			Stats:  stats,
			Types:  types,
		}
	} else {
		fmt.Printf("%s escaped!\n", args[0])
	}

	return nil
}

func commandInspect(c *Config, args ...string) error {
	if len(args) != 1 {
		return errors.New("Invalid args")
	}

	pokemon, ok := c.Pokedex[args[0]]
	if !ok {
		return errors.New("Pokemon not found")
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for stat, value := range pokemon.Stats {
		fmt.Printf("-%s: %d\n", stat, value)
	}
	fmt.Printf("Types:\n")
	for _, t := range pokemon.Types {
		fmt.Printf("-%s\n", t)
	}

	return nil
}
