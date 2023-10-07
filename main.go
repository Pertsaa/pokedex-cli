package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func main() {
	commands := getCommands()

	for {
		fmt.Print("pokedex > ")

		scanner := bufio.NewScanner(os.Stdin)

		scanner.Scan()
		command := scanner.Text()

		if _, ok := commands[command]; !ok {
			fmt.Printf("Command not found: %s\n", command)

			continue
		}

		commands[command].callback()
	}
}

func commandHelp() error {
	fmt.Print("\nWelcome to the Pokedex!\n")
	fmt.Print("\nUsage:\n")

	commands := getCommands()

	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	fmt.Println()

	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
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
	}
}
